package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig dbConfig
	Service  serviceConfig
}

type dbConfig struct {
	DB_NAME     string `json:"db_name"`
	DB_PASSWORD string `json:"db_password"`
	DB_USER     string `json:"db_user"`
	DB_HOST     string `json:"db_host"`
	DB_PORT     int    `json:"db_port"`
}

type serviceConfig struct {
	APP_PORT     int
	APP_DEBUG    bool
	API_ENDPOINT string
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
	DB := dbConfig{
		DB_NAME:     viper.GetString("DB_NAME"),
		DB_HOST:     viper.GetString("DB_HOST"),
		DB_PASSWORD: viper.GetString("DB_PASSWORD"),
		DB_PORT:     viper.GetInt("DB_PORT"),
		DB_USER:     viper.GetString("DB_USER"),
	}
	APP := serviceConfig{
		APP_DEBUG:    viper.GetBool("APP_DEBUG"),
		APP_PORT:     viper.GetInt("APP_PORT"),
		API_ENDPOINT: viper.GetString("API_ENDPOINT"),
	}

	return &Config{
		DBConfig: DB,
		Service:  APP,
	}
}
