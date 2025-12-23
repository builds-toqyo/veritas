// Package keeper implements the Veritas keeper bot core functionality
package keeper

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
	// Network settings
	MantleRPC string
	ChainID   int64

	// Contract addresses
	LeveragedStrategyAddr string
	InvoiceTokenAddr      string
	KYCVerifierAddr       string

	// ML Engine settings
	MLAPIEndpoint string

	// Gas settings
	MaxGasPrice *big.Int
	GasLimit    uint64

	// Private key for transactions
	PrivateKey string

	// Risk thresholds
	CriticalRisk    float64
	HighRisk        float64
	MaxLTV          float64
	MinHealthFactor float64
	MinLiquidity    float64
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

// Bot represents the Veritas keeper bot
type Bot struct {
	config        *Config
	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	address       common.Address
	chainID       *big.Int
	logger        *logrus.Logger
	httpClient    *http.Client
	cron          *cron.Cron
	emergencyMode bool

	// Contract instances (would use generated bindings)
	leveragedStrategy common.Address
	invoiceToken      common.Address
	kycVerifier       common.Address
}

// New creates a new keeper bot instance
func New(config *Config) (*Bot, error) {
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

	return &Bot{
		config:        config,
		client:        client,
		privateKey:    privateKey,
		address:       address,
		chainID:       big.NewInt(config.ChainID),
		logger:        logger,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		cron:          cron.New(),
		emergencyMode: false,

		// Initialize contract addresses
		leveragedStrategy: common.HexToAddress(config.LeveragedStrategyAddr),
		invoiceToken:      common.HexToAddress(config.InvoiceTokenAddr),
		kycVerifier:       common.HexToAddress(config.KYCVerifierAddr),
	}, nil
}

// MonitorLeverageStrategy monitors the leveraged RWA strategy
func (b *Bot) MonitorLeverageStrategy(ctx context.Context) error {
	b.logger.Info("Monitoring leverage strategy health...")

	// Get position data from contract (simplified - would use contract bindings)
	positionData := map[string]interface{}{
		"totalCollateral":     1000000, // Mock data
		"totalBorrowed":       600000,
		"currentHealthFactor": 1.5,
		"aitValue":            650000,
	}

	// Call ML engine for risk assessment
	response, err := b.callMLAPI("/api/v1/leverage-health", positionData)
	if err != nil {
		return fmt.Errorf("ML API call failed: %w", err)
	}

	var healthResp LeverageHealthResponse
	if err := json.Unmarshal(response, &healthResp); err != nil {
		return fmt.Errorf("failed to parse ML response: %w", err)
	}

	b.logger.WithFields(logrus.Fields{
		"risk_level": healthResp.RiskLevel,
		"risk_score": healthResp.CompositeRiskScore,
	}).Info("Risk assessment completed")

	// Execute actions based on recommendations
	return b.executeRiskActions(ctx, &healthResp)
}

// executeRiskActions performs risk management actions
func (b *Bot) executeRiskActions(ctx context.Context, assessment *LeverageHealthResponse) error {
	for _, recommendation := range assessment.Recommendations {
		switch recommendation {
		case "EMERGENCY_DELEVERAGE":
			b.logger.Warn("EMERGENCY DELEVERAGING TRIGGERED")
			return b.emergencyDeleverage(ctx)
		case "REDUCE_LEVERAGE":
			b.logger.Info("Reducing leverage position")
			return b.reduceLeverage(ctx)
		case "PAUSE_NEW_POSITIONS":
			b.logger.Info("Pausing new positions due to low liquidity")
			// Implementation would pause new borrowing
		}
	}
	return nil
}

// emergencyDeleverage executes emergency deleveraging
func (b *Bot) emergencyDeleverage(ctx context.Context) error {
	auth, err := b.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// In production: call emergencyDeleverage on LeveragedRWAStrategy contract
	b.logger.WithField("nonce", auth.Nonce).Info("Emergency deleverage transaction sent")
	b.emergencyMode = true
	return nil
}

// reduceLeverage gradually reduces leverage
func (b *Bot) reduceLeverage(ctx context.Context) error {
	auth, err := b.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// In production: call harvestRwaYield and repayDebt
	b.logger.WithField("nonce", auth.Nonce).Info("Leverage reduction transaction sent")
	return nil
}

// UpdateInvoiceNAV updates NAV for VeritasInvoiceToken
func (b *Bot) UpdateInvoiceNAV(ctx context.Context) error {
	b.logger.Info("Updating invoice token NAV...")

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

	response, err := b.callMLAPI("/api/v1/invoice-nav-prediction", navData)
	if err != nil {
		return fmt.Errorf("NAV prediction failed: %w", err)
	}

	var navResp NAVPredictionResponse
	if err := json.Unmarshal(response, &navResp); err != nil {
		return fmt.Errorf("failed to parse NAV response: %w", err)
	}

	b.logger.WithFields(logrus.Fields{
		"predicted_nav": navResp.PredictedNAV,
		"confidence":    navResp.Confidence,
	}).Info("NAV prediction completed")

	// Update NAV if confidence is high enough
	if navResp.Confidence > 0.7 {
		return b.updateNAVOnChain(ctx, navResp.PredictedNAV)
	}

	b.logger.Warn("Low confidence NAV prediction, skipping update")
	return nil
}

// updateNAVOnChain updates NAV on the smart contract
func (b *Bot) updateNAVOnChain(ctx context.Context, newNAV float64) error {
	auth, err := b.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	// Convert to wei (assuming 6 decimals for USDC compatibility)
	navWei := big.NewInt(int64(newNAV * 1e6))

	// In production: call updateNav on VeritasInvoiceToken contract
	b.logger.WithFields(logrus.Fields{
		"nav_wei": navWei.String(),
		"nonce":   auth.Nonce,
	}).Info("NAV update transaction sent")

	return nil
}

// MonitorKYCCompliance monitors KYC compliance
func (b *Bot) MonitorKYCCompliance(ctx context.Context) error {
	b.logger.Info("Monitoring KYC compliance...")

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
		response, err := b.callMLAPI("/api/v1/kyc-risk-assessment", investment)
		if err != nil {
			b.logger.WithError(err).Error("KYC risk assessment failed")
			continue
		}

		var kycResp KYCRiskResponse
		if err := json.Unmarshal(response, &kycResp); err != nil {
			b.logger.WithError(err).Error("Failed to parse KYC response")
			continue
		}

		if kycResp.RiskClassification == "HIGH_RISK" {
			b.logger.WithFields(logrus.Fields{
				"risk_score":     kycResp.KYCRiskScore,
				"classification": kycResp.RiskClassification,
				"flags":          kycResp.ComplianceFlags,
			}).Warn("HIGH RISK INVESTMENT DETECTED")
		}
	}

	return nil
}

// callMLAPI makes HTTP calls to the ML engine
func (b *Bot) callMLAPI(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := b.httpClient.Post(
		b.config.MLAPIEndpoint+endpoint,
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
func (b *Bot) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := b.client.PendingNonceAt(ctx, b.address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(b.privateKey, b.chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = b.config.GasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// HealthCheck performs system health check
func (b *Bot) HealthCheck(ctx context.Context) error {
	// Check ML engine health
	resp, err := b.httpClient.Get(b.config.MLAPIEndpoint + "/health")
	if err != nil {
		b.logger.WithError(err).Error("ML engine health check failed")
	} else {
		resp.Body.Close()
		b.logger.Info("ML engine health check: OK")
	}

	// Check blockchain connection
	latestBlock, err := b.client.BlockNumber(ctx)
	if err != nil {
		b.logger.WithError(err).Error("Blockchain connection failed")
	} else {
		b.logger.WithField("block", latestBlock).Info("Blockchain connection: OK")
	}

	// Check account balance
	balance, err := b.client.BalanceAt(ctx, b.address, nil)
	if err != nil {
		b.logger.WithError(err).Error("Failed to get account balance")
	} else {
		ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		b.logger.WithField("balance", ethBalance).Info("Account balance checked")

		if balance.Cmp(big.NewInt(1e17)) < 0 { // Less than 0.1 ETH
			b.logger.Warn("LOW KEEPER ACCOUNT BALANCE - REFILL NEEDED")
		}
	}

	return nil
}

// Start starts the keeper bot with scheduled tasks
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Starting Veritas Keeper Bot...")
	b.logger.WithField("address", b.address.Hex()).Info("Keeper address")

	// Schedule tasks
	b.cron.AddFunc("*/5 * * * *", func() { // Every 5 minutes
		if err := b.MonitorLeverageStrategy(ctx); err != nil {
			b.logger.WithError(err).Error("Leverage monitoring failed")
		}
	})

	b.cron.AddFunc("*/30 * * * *", func() { // Every 30 minutes
		if err := b.UpdateInvoiceNAV(ctx); err != nil {
			b.logger.WithError(err).Error("NAV update failed")
		}
	})

	b.cron.AddFunc("*/15 * * * *", func() { // Every 15 minutes
		if err := b.MonitorKYCCompliance(ctx); err != nil {
			b.logger.WithError(err).Error("KYC monitoring failed")
		}
	})

	b.cron.AddFunc("0 * * * *", func() { // Every hour
		if err := b.HealthCheck(ctx); err != nil {
			b.logger.WithError(err).Error("Health check failed")
		}
	})

	// Start cron scheduler
	b.cron.Start()

	// Initial health check
	b.HealthCheck(ctx)

	// Keep running
	<-ctx.Done()
	b.logger.Info("Keeper bot shutting down...")
	b.cron.Stop()
	return ctx.Err()
}
