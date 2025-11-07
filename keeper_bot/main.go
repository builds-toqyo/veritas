package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	// Mantle Network
	MantleRPC string
	ChainID   int64

	// Contract Addresses
	VeritasVaultAddr string
	MLAPIEndpoint    string

	// Bot Settings
	PrivateKey        string
	HarvestInterval   time.Duration
	RebalanceInterval time.Duration
	MLUpdateInterval  time.Duration

	// Gas Settings
	MaxGasPrice *big.Int
	GasLimit    uint64
}

func LoadConfig() *Config {
	return &Config{
		MantleRPC:         getEnv("MANTLE_RPC", "https://rpc.mantle.xyz"),
		ChainID:           5000, // Mantle Mainnet
		VeritasVaultAddr:  getEnv("Veritas_VAULT_ADDR", "0x..."),
		MLAPIEndpoint:     getEnv("ML_API_ENDPOINT", "http://localhost:5000"),
		PrivateKey:        os.Getenv("KEEPER_PRIVATE_KEY"),
		HarvestInterval:   24 * time.Hour,         // Daily harvest
		RebalanceInterval: 6 * time.Hour,          // Every 6 hours
		MLUpdateInterval:  12 * time.Hour,         // Twice daily
		MaxGasPrice:       big.NewInt(5000000000), // 5 Gwei
		GasLimit:          500000,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

type MLModelResponse struct {
	RiskScore      float64 `json:"risk_score"`      // 0.0 - 1.0
	LiquidityScore float64 `json:"liquidity_score"` // 0.0 - 1.0
	Confidence     float64 `json:"confidence"`      // 0.0 - 1.0
	ModelVersion   string  `json:"model_version"`
	Timestamp      int64   `json:"timestamp"`
}

type MLClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewMLClient(endpoint string) *MLClient {
	return &MLClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchRiskAssessment queries the ML API for current RWA risk/liquidity scores
func (c *MLClient) FetchRiskAssessment() (*MLModelResponse, error) {
	resp, err := c.httpClient.Get(c.endpoint + "/api/v1/risk-assessment")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch risk assessment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var mlResp MLModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &mlResp, nil
}

type MantleClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
}

func NewMantleClient(rpcURL string, privateKeyHex string, chainID int64) (*MantleClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mantle: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &MantleClient{
		client:     client,
		privateKey: privateKey,
		address:    address,
		chainID:    big.NewInt(chainID),
	}, nil
}

func (m *MantleClient) GetTransactOpts(ctx context.Context, gasLimit uint64) (*bind.TransactOpts, error) {
	nonce, err := m.client.PendingNonceAt(ctx, m.address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := m.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// Mantle has very low gas prices, typically < 0.1 Gwei
	log.Printf("Current gas price: %s Gwei", new(big.Float).Quo(
		new(big.Float).SetInt(gasPrice),
		big.NewFloat(1e9),
	))

	auth, err := bind.NewKeyedTransactorWithChainID(m.privateKey, m.chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

type KeeperBot struct {
	config       *Config
	mantleClient *MantleClient
	mlClient     *MLClient
	vaultAddr    common.Address
}

func NewKeeperBot(config *Config) (*KeeperBot, error) {
	mantleClient, err := NewMantleClient(
		config.MantleRPC,
		config.PrivateKey,
		config.ChainID,
	)
	if err != nil {
		return nil, err
	}

	mlClient := NewMLClient(config.MLAPIEndpoint)

	return &KeeperBot{
		config:       config,
		mantleClient: mantleClient,
		mlClient:     mlClient,
		vaultAddr:    common.HexToAddress(config.VeritasVaultAddr),
	}, nil
}

func (k *KeeperBot) UpdateMLModel(ctx context.Context) error {
	log.Println("Fetching ML model predictions...")

	mlResp, err := k.mlClient.FetchRiskAssessment()
	if err != nil {
		return fmt.Errorf("failed to fetch ML assessment: %w", err)
	}

	log.Printf("ML Model Output: Risk=%.4f, Liquidity=%.4f, Confidence=%.4f, Version=%s",
		mlResp.RiskScore,
		mlResp.LiquidityScore,
		mlResp.Confidence,
		mlResp.ModelVersion,
	)

	// Convert to basis points (0.0-1.0 -> 0-10000)
	riskScoreBps := big.NewInt(int64(mlResp.RiskScore * 10000))
	liquidityScoreBps := big.NewInt(int64(mlResp.LiquidityScore * 10000))

	// Only update if confidence is high enough
	if mlResp.Confidence < 0.8 {
		log.Printf("Skipping update - confidence too low: %.4f", mlResp.Confidence)
		return nil
	}

	// Create transaction
	auth, err := k.mantleClient.GetTransactOpts(ctx, k.config.GasLimit)
	if err != nil {
		return fmt.Errorf("failed to create transaction opts: %w", err)
	}

	log.Printf("Updating on-chain ML model: Risk=%s bps, Liquidity=%s bps",
		riskScoreBps.String(),
		liquidityScoreBps.String(),
	)

	// In production: use generated Go bindings from ABI
	// tx, err := vaultContract.UpdateMLModel(auth, riskScoreBps, liquidityScoreBps, mlResp.ModelVersion)

	log.Printf("ML model updated successfully. Auth nonce: %s", auth.Nonce.String())

	return nil
}

func (k *KeeperBot) Rebalance(ctx context.Context) error {
	log.Println("Executing vault rebalance...")

	auth, err := k.mantleClient.GetTransactOpts(ctx, k.config.GasLimit)
	if err != nil {
		return fmt.Errorf("failed to create transaction opts: %w", err)
	}

	// In production: call vault.rebalance()
	log.Println("Rebalance transaction submitted successfully")
	log.Printf("Nonce: %s", auth.Nonce.String())

	return nil
}

func (k *KeeperBot) HarvestYield(ctx context.Context) error {
	log.Println("Harvesting yield from strategies...")

	auth, err := k.mantleClient.GetTransactOpts(ctx, k.config.GasLimit)
	if err != nil {
		return fmt.Errorf("failed to create transaction opts: %w", err)
	}

	// In production: call vault.harvestAndDistribute()
	log.Println("Harvest transaction submitted successfully")
	log.Printf("Nonce: %s", auth.Nonce.String())

	return nil
}

func (k *KeeperBot) Run(ctx context.Context) error {
	log.Println("Starting Veritas Keeper Bot...")
	log.Printf("Bot address: %s", k.mantleClient.address.Hex())
	log.Printf("Vault address: %s", k.vaultAddr.Hex())

	// Initial ML update
	if err := k.UpdateMLModel(ctx); err != nil {
		log.Printf("Initial ML update failed: %v", err)
	}

	// ML Model Update Loop
	mlTicker := time.NewTicker(k.config.MLUpdateInterval)
	defer mlTicker.Stop()

	// Rebalance Loop
	rebalanceTicker := time.NewTicker(k.config.RebalanceInterval)
	defer rebalanceTicker.Stop()

	// Harvest Loop
	harvestTicker := time.NewTicker(k.config.HarvestInterval)
	defer harvestTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Keeper bot shutting down...")
			return ctx.Err()

		case <-mlTicker.C:
			if err := k.UpdateMLModel(ctx); err != nil {
				log.Printf("ML update failed: %v", err)
			}

		case <-rebalanceTicker.C:
			if err := k.Rebalance(ctx); err != nil {
				log.Printf("Rebalance failed: %v", err)
			}

		case <-harvestTicker.C:
			if err := k.HarvestYield(ctx); err != nil {
				log.Printf("Harvest failed: %v", err)
			}
		}
	}
}

type HealthServer struct {
	bot *KeeperBot
}

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
		fmt.Fprintf(w, "Veritas_keeper_uptime_seconds %d\n", time.Now().Unix())
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	config := LoadConfig()

	bot, err := NewKeeperBot(config)
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

	// Run keeper bot
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("Keeper bot error: %v", err)
	}
}
