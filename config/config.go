package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TelegramToken string
	AdminsID      []int
}

// New returns a new Config struct
func New() *Config {
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

	return &Config{
		getEnv("TELEGRAM_TOKEN", ""),
		adminsID,
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
