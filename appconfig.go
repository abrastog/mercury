package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// AppConfiguration struct
type AppConfiguration struct {
	ProgramLog struct {
		LogFile      string `yaml:"LogFile"`
		LogLimitInMB int    `yaml:"LogLimitInMB"`
	} `yaml:"ProgramLog"`

	DataManagement struct {
		DataPath         string `yaml:"DataPath"`
		OMSPath          string `yaml:"OMSPath"`
		InstrumentCache  string `yaml:"InstrumentCache"`
		StorageLimitInGB int    `yaml:"StorageLimitInGB"`
		StoragePurgeTime string `yaml:"StoragePurgeTime"`
		PublishSocket    string `yaml:"PublishSocket"`
	} `yaml:"DataManagement"`

	Ticker struct {
		Enable     bool     `yaml:"Enable"`
		NseIndices []string `yaml:"NseIndices"`
		NseSymbols []string `yaml:"NseSymbols"`
	} `yaml:"Ticker"`

	TickerNseFutures struct {
		Enable     bool     `yaml:"Enable"`
		NseSymbols []string `yaml:"NseSymbols"`
	} `yaml:"TickerNseFutures"`

	TickerNiftyWeeklyOptions struct {
		Enable          bool `yaml:"Enable"`
		LimitITMStrikes int  `yaml:"LimitITMStrikes"`
		LimitOTMStrikes int  `yaml:"LimitOTMStrikes"`
	} `yaml:"TickerNiftyWeeklyOptions"`

	TickerBankNiftyWeeklyOptions struct {
		Enable          bool `yaml:"Enable"`
		LimitITMStrikes int  `yaml:"LimitITMStrikes"`
		LimitOTMStrikes int  `yaml:"LimitOTMStrikes"`
	} `yaml:"TickerBankNiftyWeeklyOptions"`

	TickerNseStockMonthlyOptions struct {
		Enable          bool `yaml:"Enable"`
		LimitITMStrikes int  `yaml:"LimitITMStrikes"`
		LimitOTMStrikes int  `yaml:"LimitOTMStrikes"`
	} `yaml:"TickerNseStockMonthlyOptions"`

	KiteConnect struct {
		Enable            bool   `yaml:"Enable"`
		Key               string `yaml:"Key"`
		Secret            string `yaml:"Secret"`
		HTTPSPort         int    `yaml:"HTTPSPort"`
		CertificateFile   string `yaml:"CertificateFile"`
		KeyFile           string `yaml:"KeyFile"`
		TokenRedirectPath string `yaml:"TokenRedirectPath"`
		PostbackPath      string `yaml:"PostbackPath"`
		TimeToReconnect   string `yaml:"TimeToReconnect"`
		MarketBeginTime   string `yaml:"MarketBeginTime"`
		MarketEndTime     string `yaml:"MarketEndTime"`
	} `yaml:"KiteConnect"`

	KiteConnectAutoLogin struct {
		ClientID string `yaml:"ClientID"`
		Password string `yaml:"Password"`
		PIN      string `yaml:"PIN"`
	} `yaml:"KiteConnectAutoLogin"`
}

// AppConfig - Config for application
var AppConfig AppConfiguration

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// NewConfig returns a new decoded Config struct
func parseApplicatioConfigFile(configPath string) error {
	// Create config structure
	AppConfig = AppConfiguration{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&AppConfig); err != nil {
		return err
	}

	return nil
}

func loadAppConfiguration() error {
	var configFilePath string
	flag.StringVar(&configFilePath, "conf", "~/.mercury/mercury.yml", "a string")
	flag.Parse()
	err := validateConfigPath(configFilePath)

	if err == nil {
		go watchAppConfig(configFilePath)
		return parseApplicatioConfigFile(configFilePath)
	}
	return errors.New("Application config missing")
}

func watchAppConfig(configFilePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	err = watcher.Add(configFilePath)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal("Failed to watch config changes")
				os.Exit(1)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Fatal("App config modified")
				os.Exit(1)
			}
		case <-watcher.Errors:
			log.Fatal("Error watching App config changes.")
			os.Exit(1)
		}
	}
}
