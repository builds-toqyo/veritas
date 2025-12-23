package keeper

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

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
