package main

import (
	"log"
	"time"

	kiteconnect "github.com/zerodhatech/gokiteconnect"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

var ticker *kiteticker.Ticker

func startKiteTicker(accessToken string) {
	// Create new Kite ticker instance
	ticker = kiteticker.New(AppConfig.KiteConnect.Key, accessToken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	go ticker.Serve()
}

// Triggered when any error is raised
func onError(err error) {
	log.Println("kiteticker Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	log.Println("kiteticker Close: ", code, reason)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	log.Printf("kiteticker Order: %+v", order)
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	log.Printf("kiteticker Maximum no of reconnect attempt reached: %d\n", attempt)
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	log.Printf("kiteticker Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	log.Println("kiteticker Connected")
	subscribeStaticSymbols()
}

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {
	// Process only ticks with good timestamp
	if isConnectionTime() && (tick.Timestamp.Time.Year() > 2019) {
		if len(tickChannel) < cap(tickChannel) {
			tickChannel <- tick
		}
	}
}

func subscribeStaticSymbols() {
	tokensToSubscribe := make([]uint32, 0)

	for _, sym := range AppConfig.Ticker.NseIndices {
		symkey := sym + "@INDICES@NSE"
		if inst, ok := mapSymbolInfo[symkey]; ok {
			tokensToSubscribe = append(tokensToSubscribe, uint32(inst.InstrumentToken))
		}
	}

	for _, sym := range AppConfig.Ticker.NseSymbols {
		symkey := sym + "@NSE@NSE"
		if inst, ok := mapSymbolInfo[symkey]; ok {
			tokensToSubscribe = append(tokensToSubscribe, uint32(inst.InstrumentToken))
		}
	}
	err := ticker.Subscribe(tokensToSubscribe)
	if err != nil {
		log.Println("err: ", err)
	} else {
		ticker.SetMode(kiteticker.ModeFull, tokensToSubscribe)
	}
}
