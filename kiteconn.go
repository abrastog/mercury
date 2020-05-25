package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sivamgr/gokitelogin"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

var kiteClient *kiteconnect.Client = nil
var mapTokenToSymbol map[uint32]string
var mapSymbolToToken map[string]uint32
var mapSymbolInfo map[string]kiteconnect.Instrument

func setupKiteConnection() {
	if !AppConfig.KiteConnect.Enable {
		return
	}

	kiteClient = kiteconnect.New(AppConfig.KiteConnect.Key)

	setupKiteCallbacks()
	go manageKiteConnection()

}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK\n"))
	log.Println("Token : ", r.URL.String())
	params := r.URL.Query()
	if params["status"][0] == "success" {
		reqToken := params["request_token"][0]
		log.Println("Token : ", reqToken)
		go restartKiteSession(reqToken)
	} else {
		log.Println("Failed to read request token")
	}
}

func defHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(""))
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error processing order update : %+v", err)
	}
	log.Printf("Order Update : %s", reqBody)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK\n"))
	updateOrdersAndPositions()
}

func setupKiteCallbacks() {
	if AppConfig.KiteConnect.HTTPSPort <= 0 {
		log.Panicf("Need https port for setting up kite connection")
		return
	}
	http.HandleFunc("/", defHandler)
	http.HandleFunc(AppConfig.KiteConnect.TokenRedirectPath, tokenHandler)
	http.HandleFunc(AppConfig.KiteConnect.PostbackPath, hookHandler)
	go launchRunWebserver(":" + strconv.Itoa(AppConfig.KiteConnect.HTTPSPort))
}

func launchRunWebserver(addr string) {
	err := http.ListenAndServeTLS(addr, AppConfig.KiteConnect.CertificateFile, AppConfig.KiteConnect.KeyFile, nil)
	if err != nil {
		log.Printf("http.ListenAndServeTLS Failed. %+v", err)
	}

}

func authenticateKiteConnection(loginurl string) error {
	return gokitelogin.Login(loginurl,
		AppConfig.KiteConnectAutoLogin.ClientID,
		AppConfig.KiteConnectAutoLogin.Password,
		AppConfig.KiteConnectAutoLogin.PIN)

}

func isConnectionTime() bool {
	wd := time.Now().Local().Weekday()
	if (wd == time.Sunday) || (wd == time.Saturday) {
		return false
	}

	tm := time.Now().Local().Format("15:04")
	if tm < AppConfig.KiteConnect.TimeToReconnect {
		return false
	}

	if tm > AppConfig.KiteConnect.MarketEndTime {
		return false
	}

	return true
}

// get a quote to check connection
func waitForDisconnection() {
	for {
		_, err := kiteClient.GetLTP("NSE:NIFTY 50")
		if err != nil {
			return
		}
		time.Sleep(1 * time.Minute)
	}
}

func manageKiteConnection() {
	// Give http server a second
	time.Sleep(1 * time.Second)

	for {

		if !isConnectionTime() {
			time.Sleep(15 * time.Minute)
			continue
		}

		err := authenticateKiteConnection(kiteClient.GetLoginURL())

		if err != nil {
			log.Printf("Kite auth error. %+v", err)
			time.Sleep(5 * time.Minute)
			continue
		} else {
			//Wait for 10 seconds for to get access token
			time.Sleep(10 * time.Second)
		}

		waitForDisconnection()
		//Once disconnected, wait for a minute and Try again
		time.Sleep(1 * time.Minute)
	}
}

func restartKiteSession(reqTok string) {
	// Get user details and access token
	data, err := kiteClient.GenerateSession(reqTok, AppConfig.KiteConnect.Secret)
	if err != nil {
		log.Printf("Failed to Get Kite Access Token: %v", err)
		return
	}
	// Set access token
	kiteClient.SetAccessToken(data.AccessToken)
	downloadInstruments()
	updateOrdersAndPositions()
	startKiteTicker(data.AccessToken)
}

func createDirForFile(filepath string) {
	dir := path.Dir(filepath)
	createDir(dir)
}

func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadInstruments() {
	cache := AppConfig.DataManagement.InstrumentCache
	createDirForFile(cache)
	if fileExists(cache) {
		fcache, err := os.Stat(cache)
		if err == nil {
			hoursSince := time.Now().Sub(fcache.ModTime()).Hours()
			// skip downloads if cache is less than 8 hours old
			if hoursSince < 8 {
				loadInstrumentsFromCache(cache)
				return
			}
		}
	}

	log.Println("Downloading Instruments")
	instruments, err := kiteClient.GetInstruments()
	if err == nil {
		storeObjectToFile(cache, instruments)
	}
	loadInstrumentsFromCache(cache)
}

func loadInstrumentsFromCache(cache string) {
	log.Println("Loading Instruments from Cache")
	buildSymbolTokenMaps(cache)
	log.Println("Loading Instruments Completed")
}

func buildSymbolTokenMaps(cache string) {
	kInstruments := new(kiteconnect.Instruments)
	loadObjectFromFile(cache, kInstruments)

	mapTokenToSymbol = make(map[uint32]string)
	mapSymbolToToken = make(map[string]uint32)
	mapSymbolInfo = make(map[string]kiteconnect.Instrument)

	for _, sym := range *kInstruments {
		symKey := sym.Tradingsymbol + "@" + sym.Segment + "@" + sym.Exchange
		mapSymbolInfo[symKey] = sym
		mapTokenToSymbol[uint32(sym.InstrumentToken)] = sym.Tradingsymbol
		mapSymbolToToken[sym.Tradingsymbol] = uint32(sym.InstrumentToken)
	}
}

func updateOrdersAndPositions() {
	createDir(AppConfig.DataManagement.OMSPath)
	ordersFile := path.Join(AppConfig.DataManagement.OMSPath, "orders.json")
	positionFile := path.Join(AppConfig.DataManagement.OMSPath, "positions.json")

	p, err := kiteClient.GetPositions()
	if err == nil {
		log.Printf("Positions : %+v", p)
		storeObjectToJSONFile(positionFile, p)
	}

	o, err := kiteClient.GetOrders()
	if err == nil {
		log.Printf("Orders : %+v", o)
		storeObjectToJSONFile(ordersFile, o)
	}
}
