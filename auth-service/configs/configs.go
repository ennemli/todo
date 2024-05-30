package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig dbConfig
	Service  serviceConfig
}

type dbConfig struct {
	DB_NAME     string
	DB_PASSWORD string
	DB_USER     string
	DB_HOST     string
	DB_PORT     int
}
type serviceConfig struct {
	APP_PORT     int
	APP_DEBUG    bool
	API_ENDPOINT string
	USER_URI     string
	JWT_SK       string
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
		USER_URI:     viper.GetString("USER_URI"),
		JWT_SK:       viper.GetString("JWT_SK"),
	}

	return &Config{
		Service:  APP,
		DBConfig: DB,
	}
}
