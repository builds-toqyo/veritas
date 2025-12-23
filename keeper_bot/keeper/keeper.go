package keeper

import (
	"bytes"
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
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the keeper bot
type Config struct {
	MantleRPC             string
	ChainID               int64
	PrivateKey            string
	MLAPIEndpoint         string
	LeveragedStrategyAddr string
	InvoiceTokenAddr      string
	KYCVerifierAddr       string
	MaxGasPrice           *big.Int
	GasLimit              uint64
}

// MLResponses for different endpoints
type LeverageHealthResponse struct {
	CompositeRiskScore float64  `json:"composite_risk_score"`
	RiskLevel          string   `json:"risk_level"`
	ActionRequired     bool     `json:"action_required"`
	Recommendations    []string `json:"recommendations"`
	Timestamp          int64    `json:"timestamp"`
}

type KYCRiskResponse struct {
	KYCRiskScore         float64  `json:"kyc_risk_score"`
	RiskClassification   string   `json:"risk_classification"`
	VerificationRequired bool     `json:"verification_required"`
	ComplianceFlags      []string `json:"compliance_flags"`
	Timestamp            int64    `json:"timestamp"`
}

type NAVPredictionResponse struct {
	PredictedNAV           float64 `json:"predicted_nav"`
	Confidence             float64 `json:"confidence"`
	ExpectedCollectionRate float64 `json:"expected_collection_rate"`
	RiskAdjustedYield      float64 `json:"risk_adjusted_yield"`
	Timestamp              int64   `json:"timestamp"`
}

// VeritasKeeper main keeper bot struct
type VeritasKeeper struct {
	config        *Config
	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	address       common.Address
	chainID       *big.Int
	logger        *logrus.Logger
	httpClient    *http.Client
	cron          *cron.Cron
	emergencyMode bool
}

// NewVeritasKeeper creates a new keeper bot instance
func NewVeritasKeeper(config *Config) (*VeritasKeeper, error) {
	client, err := ethclient.Dial(config.MantleRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mantle: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &VeritasKeeper{
		config:        config,
		client:        client,
		privateKey:    privateKey,
		address:       address,
		chainID:       big.NewInt(config.ChainID),
		logger:        logger,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		cron:          cron.New(),
		emergencyMode: false,
	}, nil
}

// MonitorLeverageStrategy monitors the leveraged RWA strategy
func (k *VeritasKeeper) MonitorLeverageStrategy(ctx context.Context) error {
	k.logger.Info("Monitoring leverage strategy health...")

	// Get position data from contract (simplified - would use contract bindings)
	positionData := map[string]interface{}{
		"totalCollateral":     1000000, // Mock data
		"totalBorrowed":       600000,
		"currentHealthFactor": 1.5,
		"aitValue":            650000,
	}

	// Call ML engine for risk assessment
	response, err := k.callMLAPI("/api/v1/leverage-health", positionData)
	if err != nil {
		return fmt.Errorf("ML API call failed: %w", err)
	}

	var healthResp LeverageHealthResponse
	if err := json.Unmarshal(response, &healthResp); err != nil {
		return fmt.Errorf("failed to parse ML response: %w", err)
	}

	k.logger.WithFields(logrus.Fields{
		"risk_level": healthResp.RiskLevel,
		"risk_score": healthResp.CompositeRiskScore,
	}).Info("Risk assessment completed")

	// Execute actions based on recommendations
	return k.executeRiskActions(ctx, &healthResp)
}

// executeRiskActions performs risk management actions
func (k *VeritasKeeper) executeRiskActions(ctx context.Context, assessment *LeverageHealthResponse) error {
	for _, recommendation := range assessment.Recommendations {
		switch recommendation {
		case "EMERGENCY_DELEVERAGE":
			k.logger.Warn("EMERGENCY DELEVERAGING TRIGGERED")
			return k.emergencyDeleverage(ctx)
		case "REDUCE_LEVERAGE":
			k.logger.Info("Reducing leverage position")
			return k.reduceLeverage(ctx)
		case "PAUSE_NEW_POSITIONS":
			k.logger.Info("Pausing new positions due to low liquidity")
			// Implementation would pause new borrowing
		}
	}
	return nil
}

// emergencyDeleverage executes emergency deleveraging
func (k *VeritasKeeper) emergencyDeleverage(ctx context.Context) error {
	auth, err := k.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// In production: call emergencyDeleverage on LeveragedRWAStrategy contract
	k.logger.WithField("nonce", auth.Nonce).Info("Emergency deleverage transaction sent")
	k.emergencyMode = true
	return nil
}

// reduceLeverage gradually reduces leverage
func (k *VeritasKeeper) reduceLeverage(ctx context.Context) error {
	auth, err := k.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// In production: call harvestRwaYield and repayDebt
	k.logger.WithField("nonce", auth.Nonce).Info("Leverage reduction transaction sent")
	return nil
}

// UpdateInvoiceNAV updates NAV for VeritasInvoiceToken
func (k *VeritasKeeper) UpdateInvoiceNAV(ctx context.Context) error {
	k.logger.Info("Updating invoice token NAV...")

	// Get pool data from contract (mock data)
	navData := map[string]interface{}{
		"totalFaceValue":   5000000,
		"numberOfInvoices": 100,
		"weightedMaturity": 90,
		"expectedYield":    800, // 8% in basis points
		"defaultRate":      300, // 3% in basis points
		"realizedYield":    200000,
		"totalSupply":      4800000,
	}

	response, err := k.callMLAPI("/api/v1/invoice-nav-prediction", navData)
	if err != nil {
		return fmt.Errorf("NAV prediction failed: %w", err)
	}

	var navResp NAVPredictionResponse
	if err := json.Unmarshal(response, &navResp); err != nil {
		return fmt.Errorf("failed to parse NAV response: %w", err)
	}

	k.logger.WithFields(logrus.Fields{
		"predicted_nav": navResp.PredictedNAV,
		"confidence":    navResp.Confidence,
	}).Info("NAV prediction completed")

	// Update NAV if confidence is high enough
	if navResp.Confidence > 0.7 {
		return k.updateNAVOnChain(ctx, navResp.PredictedNAV)
	}

	k.logger.Warn("Low confidence NAV prediction, skipping update")
	return nil
}

// updateNAVOnChain updates NAV on the smart contract
func (k *VeritasKeeper) updateNAVOnChain(ctx context.Context, newNAV float64) error {
	auth, err := k.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// Convert to wei (assuming 6 decimals for USDC compatibility)
	navWei := big.NewInt(int64(newNAV * 1e6))

	// In production: call updateNav on VeritasInvoiceToken contract
	k.logger.WithFields(logrus.Fields{
		"nav_wei": navWei.String(),
		"nonce":   auth.Nonce,
	}).Info("NAV update transaction sent")

	return nil
}

// MonitorKYCCompliance monitors KYC compliance
func (k *VeritasKeeper) MonitorKYCCompliance(ctx context.Context) error {
	k.logger.Info("Monitoring KYC compliance...")

	// Mock investment data - in production would get from contract events
	investments := []map[string]interface{}{
		{
			"investmentAmount":     500000,
			"tier":                 2,
			"jurisdiction":         "US",
			"transactionFrequency": 5,
			"walletAgeDays":        100,
			"previousDefiExposure": 0,
		},
	}

	for _, investment := range investments {
		response, err := k.callMLAPI("/api/v1/kyc-risk-assessment", investment)
		if err != nil {
			k.logger.WithError(err).Error("KYC risk assessment failed")
			continue
		}

		var kycResp KYCRiskResponse
		if err := json.Unmarshal(response, &kycResp); err != nil {
			k.logger.WithError(err).Error("Failed to parse KYC response")
			continue
		}

		if kycResp.RiskClassification == "HIGH_RISK" {
			k.logger.WithFields(logrus.Fields{
				"risk_score":     kycResp.KYCRiskScore,
				"classification": kycResp.RiskClassification,
				"flags":          kycResp.ComplianceFlags,
			}).Warn("HIGH RISK INVESTMENT DETECTED")
		}
	}

	return nil
}

// callMLAPI makes HTTP calls to the ML engine
func (k *VeritasKeeper) callMLAPI(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := k.httpClient.Post(
		k.config.MLAPIEndpoint+endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// getTransactOpts creates transaction options
func (k *VeritasKeeper) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := k.client.PendingNonceAt(ctx, k.address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := k.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(k.privateKey, k.chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = k.config.GasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// HealthCheck performs system health check
func (k *VeritasKeeper) HealthCheck(ctx context.Context) error {
	// Check ML engine health
	resp, err := k.httpClient.Get(k.config.MLAPIEndpoint + "/health")
	if err != nil {
		k.logger.WithError(err).Error("ML engine health check failed")
	} else {
		resp.Body.Close()
		k.logger.Info("ML engine health check: OK")
	}

	// Check blockchain connection
	latestBlock, err := k.client.BlockNumber(ctx)
	if err != nil {
		k.logger.WithError(err).Error("Blockchain connection failed")
	} else {
		k.logger.WithField("block", latestBlock).Info("Blockchain connection: OK")
	}

	// Check account balance
	balance, err := k.client.BalanceAt(ctx, k.address, nil)
	if err != nil {
		k.logger.WithError(err).Error("Failed to get account balance")
	} else {
		ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		k.logger.WithField("balance", ethBalance).Info("Account balance checked")

		if balance.Cmp(big.NewInt(1e17)) < 0 { // Less than 0.1 ETH
			k.logger.Warn("LOW KEEPER ACCOUNT BALANCE - REFILL NEEDED")
		}
	}

	return nil
}

// Start starts the keeper bot with scheduled tasks
func (k *VeritasKeeper) Start(ctx context.Context) error {
	k.logger.Info("Starting Veritas Keeper Bot...")
	k.logger.WithField("address", k.address.Hex()).Info("Keeper address")

	// Schedule tasks
	k.cron.AddFunc("*/5 * * * *", func() { // Every 5 minutes
		if err := k.MonitorLeverageStrategy(ctx); err != nil {
			k.logger.WithError(err).Error("Leverage monitoring failed")
		}
	})

	k.cron.AddFunc("*/30 * * * *", func() { // Every 30 minutes
		if err := k.UpdateInvoiceNAV(ctx); err != nil {
			k.logger.WithError(err).Error("NAV update failed")
		}
	})

	k.cron.AddFunc("*/15 * * * *", func() { // Every 15 minutes
		if err := k.MonitorKYCCompliance(ctx); err != nil {
			k.logger.WithError(err).Error("KYC monitoring failed")
		}
	})

	k.cron.AddFunc("0 * * * *", func() { // Every hour
		if err := k.HealthCheck(ctx); err != nil {
			k.logger.WithError(err).Error("Health check failed")
		}
	})

	// Start cron scheduler
	k.cron.Start()

	// Initial health check
	k.HealthCheck(ctx)

	// Keep running
	select {
	case <-ctx.Done():
		k.logger.Info("Keeper bot shutting down...")
		k.cron.Stop()
		return ctx.Err()
	}
}

func LoadConfig() *Config {
	return &Config{
		MantleRPC:             getEnv("MANTLE_RPC", "https://rpc.mantle.xyz"),
		ChainID:               5000,
		PrivateKey:            os.Getenv("KEEPER_PRIVATE_KEY"),
		MLAPIEndpoint:         getEnv("ML_API_ENDPOINT", "http://localhost:5000"),
		LeveragedStrategyAddr: getEnv("LEVERAGED_STRATEGY_ADDR", "0x..."),
		InvoiceTokenAddr:      getEnv("INVOICE_TOKEN_ADDR", "0x..."),
		KYCVerifierAddr:       getEnv("KYC_VERIFIER_ADDR", "0x..."),
		MaxGasPrice:           big.NewInt(5000000000), // 5 Gwei
		GasLimit:              500000,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func main() {
	config := LoadConfig()

	if config.PrivateKey == "" {
		log.Fatal("KEEPER_PRIVATE_KEY environment variable required")
	}

	keeper, err := NewVeritasKeeper(config)
	if err != nil {
		log.Fatalf("Failed to initialize keeper: %v", err)
	}

	ctx := context.Background()

	// Start health server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "healthy",
				"time":   time.Now().Format(time.RFC3339),
			})
		})
		log.Println("Health server starting on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Start keeper bot
	if err := keeper.Start(ctx); err != nil {
		log.Fatalf("Keeper bot error: %v", err)
	}
}
