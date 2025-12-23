// Veritas Keeper Bot - Automated Risk Monitoring and Strategy Management
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/veritas/keeper-bot/keeper"
)

// HealthServer handles HTTP health check endpoints
type HealthServer struct {
	bot *keeper.Bot
}

// ServeHTTP implements http.Handler interface
func (h *HealthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
		return
	}

	if r.URL.Path == "/metrics" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "# Veritas Keeper Bot Metrics\n")
		fmt.Fprintf(w, "veritas_keeper_uptime_seconds %d\n", time.Now().Unix())
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	// Load configuration
	config := &keeper.Config{
		MantleRPC:             getEnv("MANTLE_RPC", "https://rpc.mantle.xyz"),
		ChainID:               5000, // Mantle Mainnet
		LeveragedStrategyAddr: os.Getenv("LEVERAGED_STRATEGY_ADDR"),
		InvoiceTokenAddr:      os.Getenv("INVOICE_TOKEN_ADDR"),
		KYCVerifierAddr:       os.Getenv("KYC_VERIFIER_ADDR"),
		MLAPIEndpoint:         getEnv("ML_API_ENDPOINT", "http://localhost:5000"),
		PrivateKey:            os.Getenv("KEEPER_PRIVATE_KEY"),
		MaxGasPrice:           big.NewInt(5000000000), // 5 Gwei
		GasLimit:              500000,

		// Risk thresholds
		CriticalRisk:    0.8,
		HighRisk:        0.6,
		MaxLTV:          0.65,
		MinHealthFactor: 1.3,
		MinLiquidity:    0.3,
	}

	// Validate required config
	if config.PrivateKey == "" {
		log.Fatal("KEEPER_PRIVATE_KEY environment variable required")
	}

	// Initialize keeper bot
	bot, err := keeper.New(config)
	if err != nil {
		log.Fatalf("Failed to initialize keeper bot: %v", err)
	}

	ctx := context.Background()

	// Start health check server
	healthServer := &HealthServer{bot: bot}
	go func() {
		log.Println("Starting health check server on :8080")
		if err := http.ListenAndServe(":8080", healthServer); err != nil {
			log.Printf("Health server error: %v", err)
		}
	}()

	// Start keeper bot
	if err := bot.Start(ctx); err != nil {
		log.Fatalf("Keeper bot error: %v", err)
	}
}

// getEnv gets environment variable with default fallback
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
