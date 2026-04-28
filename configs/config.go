package configs

import "os"

type Config struct {
	SpiceAddr string
}

func Load() Config {
	return Config{
		SpiceAddr: getEnv("SPICEDB_ADDR", "localhost:50051"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
