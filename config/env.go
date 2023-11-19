package config

import "github.com/spf13/viper"

type EnvVars struct {
	TOKEN  string `mapstructure:"TOKEN"`
	CMC_API_KEY  string `mapstructure:"CMC_API_KEY"`
}

func LoadConfig() (config EnvVars, err error) {
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