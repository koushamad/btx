package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

const chunkSize = 1000000

type BitcoinAddress struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func main() {
	totalCore := runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("Total core:", totalCore)
	log.Println("Start generating Bitcoin addresses...")
	for i := 0; true; i++ {

		if err := generator(totalCore, i); err != nil {
			log.Println("Error:", err)
			i--
			continue
		}

		log.Println("Generated", totalCore*(i+1), "million addresses and saved to file.")
	}
}

func generator(core int, round int) error {
	// get computer's number of cores
	addresses := make([]BitcoinAddress, 0, chunkSize*core)
	address := make(chan BitcoinAddress, 10)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(core)

	for i := 0; i < core; i++ {
		go createMillionAddress(address, &wg)
	}

	go func() {
		wg.Wait()
		done <- true
		close(address)
		close(done)
	}()

	for {
		select {
		case addr, ok := <-address:
			if !ok {
				continue
			}

			addresses = append(addresses, addr)
		case <-done:
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			fileName := timestamp + "_addresses.json"
			if err := writeAddressesToFile(fileName, addresses); err != nil {
				return err
			}
		}
	}
}

func createMillionAddress(address chan BitcoinAddress, wg *sync.WaitGroup) {
	for i := 0; i < chunkSize; i++ {
		address <- createBitcoinAddress()
	}

	wg.Done()
}

func createBitcoinAddress() BitcoinAddress {
	// Generate a private key
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Generate a public key
	publicKey := privateKey.PubKey()

	// Generate a Bitcoin address
	address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Failed to generate address: %v", err)
	}

	// Create a BitcoinAddress struct and return it
	return BitcoinAddress{
		Address:    address.EncodeAddress(),
		PublicKey:  hex.EncodeToString(publicKey.SerializeCompressed()),
		PrivateKey: hex.EncodeToString(privateKey.Serialize()),
	}
}

func writeAddressesToFile(fileName string, addresses []BitcoinAddress) error {
	// Open the file in write mode
	file, err := os.Create("storage/" + fileName)
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
