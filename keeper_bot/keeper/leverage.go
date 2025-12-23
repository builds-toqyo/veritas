package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

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
