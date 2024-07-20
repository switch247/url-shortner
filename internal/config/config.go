package config

import "os"

var (
	AppDomain = getEnv("APP_DOMAIN", "http://localhost:8080")
	AppName   = getEnv("APP_NAME", "tiny.cloud")
	Port      = getEnv("PORT", "8080")
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
