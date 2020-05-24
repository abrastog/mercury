package main

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

func storeObjectToFile(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func storeObjectToJSONFile(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := json.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func loadObjectFromFile(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
