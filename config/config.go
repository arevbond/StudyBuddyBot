package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env           string `yaml:"env"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	AdminsID      []int
	PostgresSettings
	PgAdminSettings
}

type PostgresSettings struct {
	PostgresDBName   string `env:"POSTGRES_DB"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
}

type PgAdminSettings struct {
	PgAdminEmail    string `env:"PGADMIN_DEFAULT_EMAIL"`
	PgAdminPassword string `env:"PGADMIN_DEFAULT_PASSWORD"`
}

// New returns a new Config struct
func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	adminsIdStrings := getEnv("ADMINS_ID", "")
	adminsID := make([]int, 0)
	for _, str := range strings.Split(adminsIdStrings, ",") {
		if str == "" {
			continue
		}
		s := strings.TrimSpace(str)
		id, err := strconv.Atoi(s)
		if err == nil {
			adminsID = append(adminsID, id)
		} else {
			log.Printf("[ERROR] can't convert %s to int", s)
		}
	}
	configPath := getEnv("CONFIG_PATH", "")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("configs file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read configs from %s: %v", configPath, err)
	}

	cfg.AdminsID = adminsID

	return &cfg
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
