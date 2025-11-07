# 🚀 Aegis RWA Vault - Deployment Checklist

## Pre-Deployment

### Environment Setup

- [ ] Install Foundry (`curl -L https://foundry.paradigm.xyz | bash`)
- [ ] Install Python 3.9+ (`python --version`)
- [ ] Install Go 1.21+ (`go version`)
- [ ] Install Node.js 18+ (for frontend, optional)

### Configuration Files

- [ ] Create `.env` file with:
  ```bash
  MANTLE_RPC=https://rpc.sepolia.mantle.xyz
  PRIVATE_KEY=your_deployer_private_key
  KEEPER_PRIVATE_KEY=your_keeper_private_key
  ETHERSCAN_API_KEY=your_api_key
  ```

- [ ] Update contract addresses in `Deploy.s.sol`:
  - [ ] mETH address (Mantle liquid staking token)
  - [ ] USDC address (Mantle USDC)
  - [ ] INIT Capital / Lendle address (lending protocol)

### Security Review

- [ ] Review all `AccessControl` roles
- [ ] Verify `onlyRole` modifiers on critical functions
- [ ] Check `ReentrancyGuard` on state-changing functions
- [ ] Audit transfer restrictions on AIT token
- [ ] Review emergency pause mechanisms

---

## Smart Contract Deployment

### Step 1: Compile Contracts

```bash
cd Veritas
forge build
```

**Expected Output:**
```
[⠊] Compiling...
[⠒] Compiling 10 files with 0.8.20
[⠢] Solc 0.8.20 finished in 3.45s
Compiler run successful!
```

- [ ] Compilation successful
- [ ] No warnings or errors

### Step 2: Run Tests

```bash
forge test -vvv
```

**Expected Output:**
```
Running 15 tests for test/AegisVault.t.sol:AegisVaultTest
[PASS] test_KYCTierIssuance() (gas: 245123)
[PASS] test_AITPoolInitialization() (gas: 189456)
...
Test result: ok. 15 passed; 0 failed; finished in 2.34s
```

- [ ] All tests passing
- [ ] Gas usage within acceptable limits (<500k per operation)

### Step 3: Deploy to Testnet

```bash
forge script script/Deploy.s.sol:DeployAegis \
  --rpc-url $MANTLE_RPC \
  --broadcast \
  --verify
```

**Record Deployed Addresses:**
- [ ] AegisInvoiceToken: `0x________________`
- [ ] TieredKYCVerifier: `0x________________`
- [ ] LeveragedRWAStrategy: `0x________________`

### Step 4: Verify Contracts

```bash
# Verify on Mantle Explorer
forge verify-contract <ADDRESS> <CONTRACT_NAME> \
  --chain-id 5003 \
  --etherscan-api-key $ETHERSCAN_API_KEY
```

- [ ] AegisInvoiceToken verified
- [ ] TieredKYCVerifier verified
- [ ] LeveragedRWAStrategy verified

### Step 5: Initial Configuration

```bash
# Run setup script
forge script script/Deploy.s.sol:SetupDemo \
  --rpc-url $MANTLE_RPC \
  --broadcast
```

- [ ] Sample invoice pool initialized
- [ ] Demo investors created
- [ ] Whitelist configured
- [ ] Roles granted

---

## ML Engine Deployment

### Step 1: Install Dependencies

```bash
cd ml_engine
pip install -r requirements.txt
```

- [ ] PyTorch installed
- [ ] Flask installed
- [ ] All dependencies resolved

### Step 2: Train Model

```bash
python aegis_ml_engine.py train
```

**Expected Output:**
```
Training model for 100 epochs...
Epoch [10/100], Loss: 0.4523
Epoch [20/100], Loss: 0.3891
...
Model saved to aegis_rwa_model.pth
```

- [ ] Model trained successfully
- [ ] Model file created (`aegis_rwa_model.pth`)

### Step 3: Start API Server

```bash
python aegis_ml_engine.py
```

**Expected Output:**
```
Starting Aegis ML API Server...
Model version: v1.2.0-lstm
 * Running on http://0.0.0.0:5000
```

- [ ] Server running on port 5000
- [ ] Health check responds: `curl http://localhost:5000/health`

### Step 4: Test API

```bash
# Test risk assessment endpoint
curl http://localhost:5000/api/v1/risk-assessment
```

**Expected Response:**
```json
{
  "risk_score": 0.35,
  "liquidity_score": 0.78,
  "confidence": 0.92,
  "model_version": "v1.2.0-lstm",
  "timestamp": 1704067200
}
```

- [ ] API responding correctly
- [ ] Risk/liquidity scores in valid range (0.0-1.0)
- [ ] Confidence > 0.8

---

## Keeper Bot Deployment

### Step 1: Install Dependencies

```bash
cd keeper_bot
go mod download
```

- [ ] Go modules downloaded
- [ ] `go-ethereum` installed

### Step 2: Configure Environment

```bash
export MANTLE_RPC="https://rpc.sepolia.mantle.xyz"
export AEGIS_VAULT_ADDR="0x..." # From deployment
export ML_API_ENDPOINT="http://localhost:5000"
export KEEPER_PRIVATE_KEY="your_keeper_private_key"
```

- [ ] All environment variables set
- [ ] Keeper wallet funded with MNT for gas

### Step 3: Grant Keeper Roles

```bash
# Grant KEEPER_ROLE to keeper bot address
cast send $LEVERAGE_STRATEGY_ADDR \
  "grantRole(bytes32,address)" \
  $(cast keccak "KEEPER_ROLE") \
  $KEEPER_ADDRESS \
  --rpc-url $MANTLE_RPC \
  --private-key $PRIVATE_KEY

# Grant ORACLE_ROLE to keeper bot
cast send $AIT_ADDR \
  "grantRole(bytes32,address)" \
  $(cast keccak "ORACLE_ROLE") \
  $KEEPER_ADDRESS \
  --rpc-url $MANTLE_RPC \
  --private-key $PRIVATE_KEY
```

- [ ] KEEPER_ROLE granted
- [ ] ORACLE_ROLE granted
- [ ] Roles verified on-chain

### Step 4: Start Keeper Bot

```bash
go run main.go
```

**Expected Output:**
```
Starting Aegis Keeper Bot...
Bot address: 0x...
Vault address: 0x...
Fetching ML model predictions...
ML Model Output: Risk=0.35, Liquidity=0.78, Confidence=0.92
```

- [ ] Bot started successfully
- [ ] ML API connection working
- [ ] Health check responds: `curl http://localhost:8080/health`

---

## Post-Deployment Verification

### Smart Contracts

- [ ] Verify contract state on Mantle Explorer
- [ ] Check initial pool configuration
- [ ] Verify role assignments
- [ ] Test deposit/withdraw functions (small amounts)

### ML Engine

- [ ] Monitor API logs for errors
- [ ] Verify prediction accuracy
- [ ] Check model confidence levels
- [ ] Test scenario analysis endpoints

### Keeper Bot

- [ ] Monitor bot logs for successful operations
- [ ] Verify ML updates executing every 12 hours
- [ ] Check rebalancing every 6 hours
- [ ] Verify yield harvesting daily

### Integration Tests

- [ ] End-to-end user deposit flow
- [ ] ML-driven rebalancing
- [ ] Yield harvesting and distribution
- [ ] Tax event emission

---

## Monitoring Setup

### Metrics to Track

- [ ] Total Value Locked (TVL)
- [ ] Health Factor (leverage strategy)
- [ ] Current Allocation (RWA vs mETH)
- [ ] Blended APY
- [ ] Gas costs per operation
- [ ] ML model confidence levels
- [ ] Keeper bot uptime

### Alerts to Configure

- [ ] Health factor < 1.5x
- [ ] ML confidence < 0.8
- [ ] Keeper bot downtime > 1 hour
- [ ] Unusual gas price spikes
- [ ] Large withdrawals (>10% TVL)

### Dashboards

- [ ] Grafana dashboard for metrics
- [ ] Mantle Explorer bookmarks
- [ ] ML API monitoring
- [ ] Keeper bot health dashboard

---

## Security Checklist

### Access Control

- [ ] Admin keys in cold storage
- [ ] Keeper keys in secure environment
- [ ] Multi-sig for critical operations
- [ ] Role-based access properly configured

### Operational Security

- [ ] Rate limiting on API endpoints
- [ ] DDoS protection for ML API
- [ ] Backup keeper bot instance
- [ ] Emergency pause procedures documented

### Compliance

- [ ] KYC provider integration tested
- [ ] AML blacklist updated
- [ ] Tax reporting tested
- [ ] Privacy measures verified

---

## Mainnet Migration Checklist

### Pre-Migration

- [ ] Full security audit completed
- [ ] Testnet running for 30+ days
- [ ] All edge cases tested
- [ ] Insurance coverage obtained

### Migration Steps

- [ ] Update RPC URLs to mainnet
- [ ] Update contract addresses (mETH, USDC, lending)
- [ ] Deploy contracts to mainnet
- [ ] Verify all contracts
- [ ] Transfer ownership to multi-sig
- [ ] Start ML engine on production server
- [ ] Start keeper bot on production server

### Post-Migration

- [ ] Monitor for 24 hours
- [ ] Gradual TVL ramp (start with $100k cap)
- [ ] Community announcement
- [ ] Documentation published

---

## Rollback Plan

### If Critical Issue Detected

1. [ ] Pause all deposits immediately
2. [ ] Stop keeper bot operations
3. [ ] Assess issue severity
4. [ ] Communicate with users
5. [ ] Execute emergency withdrawal if needed
6. [ ] Fix issue and redeploy
7. [ ] Resume operations after verification

### Emergency Contacts

- [ ] Lead developer: _______________
- [ ] Security auditor: _______________
- [ ] Mantle team contact: _______________
- [ ] Legal counsel: _______________

---

## Success Criteria

### Technical

- [x] All contracts deployed and verified
- [x] ML engine running with >90% uptime
- [x] Keeper bot executing all operations
- [x] Gas costs <$5/month
- [x] No security vulnerabilities

### Business

- [ ] $1M+ TVL within 30 days
- [ ] 100+ unique depositors
- [ ] 6%+ blended APY achieved
- [ ] Zero compliance incidents
- [ ] Positive user feedback

---

## Notes

**Deployment Date**: _______________  
**Deployed By**: _______________  
**Network**: Mantle Sepolia / Mainnet  
**Version**: v1.0.0

**Additional Notes**:
```
[Add any deployment-specific notes here]
```

---

**✅ Deployment Complete!**

Next steps:
1. Monitor all systems for 24 hours
2. Announce to community
3. Begin user onboarding
4. Collect feedback and iterate

---

*For support, contact: team@aegisvault.xyz*
