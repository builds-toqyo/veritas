package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

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
