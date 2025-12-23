// Package keeper implements the Veritas keeper bot core functionality
package keeper

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

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
