package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yourusername/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultVersion = "dev"
)

var (
	// Version is set at build time via ldflags
	Version = defaultVersion
	// Commit is set at build time via ldflags "unknown"
)

func main() {
	var (
		port    int
		verbose bool
		showVer bool
	)

	flag.IntVar(&port, "port", getEnvInt("PORT", defaultPort), "HTTP server port")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&showVer, "version", false, "Print version information and exit")
	flag.Parse()

	if showVer {
		fmt.Printf("sub2api version %s (commit: %s)\n", Version, Commit)
		os.Exit(0)
	}

	log.Printf("Starting sub2api %s (commit: %s)", Version, Commit)

	router := handler.NewRouter(handler.Config{
		Verbose: verbose,
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s", addr)

	// Increased timeouts to better handle slow subscription sources.
	// Bumped ReadTimeout and WriteTimeout to 90s for particularly sluggish upstream URLs.
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// getEnvInt retrieves an integer environment variable or returns the fallback value.
func getEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Warning: invalid value for %s=%q, using default %d", key, val, fallback)
		return fallback
	}
	return n
}
