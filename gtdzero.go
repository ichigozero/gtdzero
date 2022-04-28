package gtdzero

import (
	"os"
	"strconv"
)

var (
	AccessSecret   = getEnv("ACCESS_SECRET", "access-secret")
	RefreshSecret  = getEnv("REFRESH_SECRET", "refresh-secret")
	CookieHashKey  = getEnv("COOKIE_HASH_KEY", "very-secret")
	CookieBlockKey = getEnv("COOKIE_BLOCK_KEY", "a-lot-secret")
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return fallback
}
