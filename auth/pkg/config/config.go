package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
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
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

func Read() *Config {
	viper.SetConfigName("config") // Yapılandırma dosyasının adı
	viper.SetConfigType("yaml")   // Yapılandırma dosyasının türü
	viper.AddConfigPath(".")      // Yapılandırma dosyasının bulunduğu dizin

	// Yapılandırma dosyasını oku
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Yapılandırma dosyası okunamadı: %v", err)
	}

	// Yapılandırmayı bir struct'a yerleştir
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Yapılandırma çözümlenemedi: %v", err)
	}

	return &config
}
