package keeper

import (
	"crypto/ecdsa"
	"math/big"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MantleRPC string
	ChainID   int64
	LeveragedStrategyAddr string
	InvoiceTokenAddr      string
	KYCVerifierAddr       string

	MLAPIEndpoint string

	MaxGasPrice *big.Int
	GasLimit    uint64

	PrivateKey string

	CriticalRisk    float64
	HighRisk        float64
	MaxLTV          float64
	MinHealthFactor float64
	MinLiquidity    float64
}

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
	mutex         sync.Mutex

	leveragedStrategy common.Address
	invoiceToken      common.Address
	kycVerifier       common.Address
}

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
