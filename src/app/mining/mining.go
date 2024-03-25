package mining

import (
	"github.com/koushamad/btx/src/app/mining/service"
	"log"
	"runtime"
)

func Mine(chunkSize int) {
	totalCore := runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("Total core:", totalCore)
	log.Println("Start generating Bitcoin addresses...")
	for i := 0; true; i++ {

		if err := service.Generator(totalCore, chunkSize); err != nil {
			i--
			continue
		}

		log.Println("Generated", totalCore*(i+1), "million addresses and saved to file.")
	}
}
