package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Config struct {
	DevConfiguration TypeConfiguration `mapstructure:"dev"`
}

type TypeConfiguration struct {
	Server     Server     `mapstructure:"server"`
	Database   Database   `mapstructure:"database"`
	Parameters Parameters `mapstructure:"parameters"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Database struct {
	DriverName string `mapstructure:"driver_name"`
	DNS        string `mapstructure:"dns"`
}

type Parameters struct {
	PwdMinLength int `mapstructure:"pwd_min_length"`
	PwdMaxLength int `mapstructure:"pwd_max_length"`
}

func LoadConfiguration() (Config, error) {
	var config Config

	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./config")
	vp.AddConfigPath(".")

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
