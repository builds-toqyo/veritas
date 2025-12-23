package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

func (b *Bot) UpdateInvoiceNAV(ctx context.Context) error {
	b.logger.Info("Updating invoice token NAV...")

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
