package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ConfigType struct {
	P2PSvcUrl string `mapstructure:"P2P_SVC_URL"`
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
