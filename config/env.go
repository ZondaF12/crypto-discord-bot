package config

import (
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	TOKEN       string `mapstructure:"TOKEN"`
	CMC_API_KEY string `mapstructure:"CMC_API_KEY"`
	MONGODB_URI string `mapstructure:"MONGODB_URI"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			MONGODB_URI:    os.Getenv("MONGODB_URI"),
			TOKEN:          os.Getenv("TOKEN"),
			CMC_API_KEY:    os.Getenv("CMC_API_KEY"),
		}, nil
	}
	
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
