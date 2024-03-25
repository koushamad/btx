package repositories

import (
	"encoding/json"
	"github.com/koushamad/btx/src/domain/model"
	"os"
)

func WriteAddressesToFile(fileName string, addresses []model.BitcoinAddress) error {
	// Open the file in write mode
	file, err := os.Create("storage/address" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new JSON encoder that writes to the file
	encoder := json.NewEncoder(file)

	// Write the addresses slice to the file in JSON format
	err = encoder.Encode(addresses)
	if err != nil {
		return err
	}

	return nil
}
