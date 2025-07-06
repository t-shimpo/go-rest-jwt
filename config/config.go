package config

import (
	"os"
	"strconv"
)

var (
	JWTSecret        = os.Getenv("JWT_SECRET")
	JWTIssuer        = os.Getenv("JWT_ISSUER")
	JWTExpireMinutes = getEnvAsInt("JWT_EXPIRE_MINUTES", 60)
)

func getEnvAsInt(key string, defaultValue int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultValue
}
