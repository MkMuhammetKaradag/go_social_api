package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	SMTP   SMTPConfig   `mapstructure:"smtp"`
}

type SMTPConfig struct {
	Host     string        `mapstructure:"host"`
	Port     string        `mapstructure:"port"`
	Email    string        `mapstructure:"email"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func Read() *Config {
	// 1. Önce .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Println(".env dosyası bulunamadı, environment variables direkt kullanılacak")
	}

	// 2. Viper ile YAML konfigürasyonunu oku
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Yapılandırma dosyası okunamadı: %v", err)
	}

	// 3. Environment variables ayarla
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.AllowEmptyEnv(true)

	// 4. Özel bağlamalar (YAML'daki yapı ile env variable'ları eşleştir)
	viper.BindEnv("smtp.host", "APP_SMTP_HOST")
	viper.BindEnv("smtp.port", "APP_SMTP_PORT")
	viper.BindEnv("smtp.email", "APP_SMTP_EMAIL")
	viper.BindEnv("smtp.password", "APP_SMTP_PASSWORD")
	viper.BindEnv("smtp.timeout", "APP_SMTP_TIMEOUT")

	// 5. Config struct'ını doldur
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Yapılandırma çözümlenemedi: %v", err)
	}

	// 6. Default değerler
	if config.SMTP.Timeout == 0 {
		config.SMTP.Timeout = 10 * time.Second
	}

	// 7. SMTP şifresini environment'dan al (opsiyonel güvenlik)
	if pass := os.Getenv("APP_SMTP_PASSWORD"); pass != "" {
		config.SMTP.Password = pass
	}

	return &config
}