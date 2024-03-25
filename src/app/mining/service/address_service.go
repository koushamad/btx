package service

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/koushamad/btx/src/domain/model"
	"github.com/koushamad/btx/src/is/repositories"
	"log"
	"strconv"
	"sync"
	"time"
)

func Generator(core int, chunkSize int) error {
	// get computer's number of cores
	addresses := make([]model.BitcoinAddress, 0, chunkSize*core)
	address := make(chan model.BitcoinAddress, 10)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(core)

	for i := 0; i < core; i++ {
		go func() {
			for i := 0; i < chunkSize; i++ {
				address <- createBitcoinAddress()
			}

			wg.Done()
		}()
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
			if err := repositories.WriteAddressesToFile(fileName, addresses); err != nil {
				return err
			}
		}
	}
}

func createBitcoinAddress() model.BitcoinAddress {
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
	return model.BitcoinAddress{
		Address:    address.EncodeAddress(),
		PublicKey:  hex.EncodeToString(publicKey.SerializeCompressed()),
		PrivateKey: hex.EncodeToString(privateKey.Serialize()),
	}
}
