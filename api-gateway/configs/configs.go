package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Service serviceConfig
}

type serviceConfig struct {
	USERS_ENDPOINT string
	TODO_ENDPOINT  string
	AUTH_ENDPOINT  string
	APP_PORT       int
	APP_DEBUG      bool
}

func Initialize(filename string, filepath string, filetype string) {
	viper.SetConfigName(filename)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(filepath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	APP := serviceConfig{
		APP_DEBUG:      viper.GetBool("APP_DEBUG"),
		APP_PORT:       viper.GetInt("APP_PORT"),
		USERS_ENDPOINT: viper.GetString("USERS_ENDPOINT"),
		TODO_ENDPOINT:  viper.GetString("TODO_ENDPOINT"),
		AUTH_ENDPOINT:  viper.GetString("AUTH_ENDPOINT"),
	}

	return &Config{
		Service: APP,
	}
}
