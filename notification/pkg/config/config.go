package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	MongoPort string `mapstructure:"mongoPort"`
	Port      string `mapstructure:"port"`
	Host      string `mapstructure:"host"`

	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type RedisConfig struct {
	RedisURL string `mapstructure:"redisURL"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	UserDB   int    `mapstructure:"userDB"`
}

func Read() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Yapılandırma dosyası okunamadı: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Yapılandırma çözümlenemedi: %v", err)
	}

	return config
}
