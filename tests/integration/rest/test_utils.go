package rest

import (
	"encoding/json"
	"log"
	"os"
)

func convertJsonFileToType(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&v); err != nil {
		log.Fatal(err)
	}

	return nil
}
