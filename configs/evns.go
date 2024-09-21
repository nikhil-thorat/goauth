package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBConnectionString string

	JWTSecret     string
	JWTExpiration int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:         getEnv("PUBLIC_HOST", "http://localhost"),
		Port:               getEnv("PORT", "8080"),
		DBConnectionString: getEnv("DB_CONNECTION", "db_connection_url"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		JWTExpiration:      getEvnAsInt("JWT_EXPIRATION", 3600*24*7),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

func getEvnAsInt(key string, fallback int64) int64 {

	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return intValue
}
