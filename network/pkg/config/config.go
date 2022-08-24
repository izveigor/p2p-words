package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ConfigType struct {
	LemmatizerSvcUrl string `mapstructure:"LEMMATIZER_SVC_URL"`
	P2PSvcUrl        string `mapstructure:"P2P_SVC_URL"`
	PathToData       string `mapstructure:"PATH_TO_DATA"`
	P2PAddress       string `mapstructure:"P2P_ADDRESS"`
	InitialPort      int    `mapstructure:"INITIAL_PORT"`
}

func LoadConfig() (c *ConfigType) {
	if !strings.HasSuffix(os.Args[0], ".test") {
		viper.AddConfigPath("./pkg/config/envs")
	} else {
		viper.AddConfigPath("../config/envs")
	}

	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return
}

var Config *ConfigType = LoadConfig()
