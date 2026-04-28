package configs

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr            string
	SpiceAddr           string
	SpiceToken          string
	SpiceInsecure       bool
	SpiceConsistency    string
	SpiceRequestTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:            getEnv("HTTP_ADDR", ":8080"),
		SpiceAddr:           getEnv("SPICEDB_ADDR", "localhost:50051"),
		SpiceToken:          getEnv("SPICEDB_TOKEN", getEnv("SPICEDB_PRESHARED_KEY", "")),
		SpiceInsecure:       getEnvBool("SPICEDB_INSECURE", true),
		SpiceConsistency:    getEnv("SPICEDB_CONSISTENCY", "minimize_latency"),
		SpiceRequestTimeout: getEnvDuration("SPICEDB_TIMEOUT", 3*time.Second),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return parsed
}
