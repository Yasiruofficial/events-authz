package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spicedb/spicedb-go/spicedb"
	"github.com/spicedb/spicedb-go/spicedb/types"
)

func main() {
	// Load configuration from environment
	cfg := loadConfig()

	// Create SpiceDB client
	client, err := spicedb.NewClient(spicedb.ClientOptions{
		Address:            cfg.SpiceAddr,
		PreSharedKey:       cfg.SpiceToken,
		TLSEnabled:         !cfg.SpiceInsecure,
		RequestTimeout:     cfg.SpiceRequestTimeout,
		DefaultConsistency: cfg.SpiceConsistency,
		DisableCache:       false,
	})
	if err != nil {
		log.Fatalf("failed to create SpiceDB client: %v", err)
	}
	defer client.Close()

	// Set up HTTP router
	if os.Getenv("GIN_DEBUG") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	handler := &Handler{client: client}

	router.POST("/check", handler.CheckPermission)
	router.GET("/health", handler.Health)

	log.Printf("AuthZ HTTP service running on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// Handler wraps the SpiceDB client for HTTP requests
type Handler struct {
	client *spicedb.Client
}

// CheckPermission handles permission check requests
func (h *Handler) CheckPermission(c *gin.Context) {
	var req types.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.client.CheckPermission(c.Request.Context(), req)
	if err != nil {
		if spicedb.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadGateway, gin.H{"error": "authorization check failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Health returns a simple health check response
func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// Config holds HTTP service configuration
type Config struct {
	HTTPAddr            string
	SpiceAddr           string
	SpiceToken          string
	SpiceInsecure       bool
	SpiceConsistency    string
	SpiceRequestTimeout time.Duration
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
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

	// Try parsing as duration
	if parsed, err := time.ParseDuration(val); err == nil {
		return parsed
	}

	// Try parsing as milliseconds
	if ms, err := strconv.ParseInt(val, 10, 64); err == nil {
		return time.Duration(ms) * time.Millisecond
	}

	return fallback
}
