package conf

import (
	"github.com/spf13/viper"
	"log"
)

var Config config

type config struct {
	ProcessChunkSize int `yaml:"process_chunk_size"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.btx")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
