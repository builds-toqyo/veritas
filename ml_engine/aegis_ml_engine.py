import numpy as np
import pandas as pd
import torch
import torch.nn as nn
from typing import Dict, List, Tuple
from datetime import datetime, timedelta
from flask import Flask, jsonify, request
import logging

class RWARiskLSTM(nn.Module):
    """
    LSTM model for predicting RWA default risk and liquidity
    
    Inputs:
    - Historical default rates (30-day window)
    - Liquidity depth metrics
    - Macro indicators (interest rates, credit spreads)
    - On-chain TVL metrics
    
    Outputs:
    - Risk score (0.0-1.0): probability of default
    - Liquidity score (0.0-1.0): market liquidity depth
    """
    
    def __init__(self, input_size: int, hidden_size: int, num_layers: int):
        super(RWARiskLSTM, self).__init__()
        
        self.hidden_size = hidden_size
        self.num_layers = num_layers
        
        # LSTM layers
        self.lstm = nn.LSTM(
            input_size, 
            hidden_size, 
            num_layers, 
            batch_first=True,
            dropout=0.2
        )
        
        # Output heads
        self.fc_risk = nn.Sequential(
            nn.Linear(hidden_size, 64),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(64, 1),
            nn.Sigmoid()  # Risk score 0-1
        )
        
        self.fc_liquidity = nn.Sequential(
            nn.Linear(hidden_size, 64),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(64, 1),
            nn.Sigmoid()  # Liquidity score 0-1
        )
        
    def forward(self, x):
        """
        Args:
            x: (batch_size, seq_length, input_size)
        Returns:
            risk_score: (batch_size, 1)
            liquidity_score: (batch_size, 1)
        """
        # LSTM forward pass
        lstm_out, (h_n, c_n) = self.lstm(x)
        
        # Use last hidden state
        last_hidden = h_n[-1]
        
        # Predict risk and liquidity
        risk_score = self.fc_risk(last_hidden)
        liquidity_score = self.fc_liquidity(last_hidden)
        
        return risk_score, liquidity_score

class RWADataProcessor:
    """
    Processes real-time RWA market data for model inference
    """
    
    def __init__(self, window_size: int = 30):
        self.window_size = window_size
        self.feature_names = [
            'default_rate',
            'avg_maturity_days',
            'weighted_credit_score',
            'on_chain_tvl',
            'liquidity_depth_usdc',
            'treasury_yield_10y',
            'credit_spread_bbb',
            'volatility_index',
        ]
        
    def fetch_market_data(self) -> pd.DataFrame:
        """
        Fetch latest RWA market data
        In production: pull from APIs (Dune, DefiLlama, Bloomberg)
        """
        # Simulated data for demonstration
        dates = pd.date_range(
            end=datetime.now(), 
            periods=self.window_size, 
            freq='D'
        )
        
        data = {
            'date': dates,
            'default_rate': np.random.uniform(0.01, 0.05, self.window_size),
            'avg_maturity_days': np.random.uniform(30, 180, self.window_size),
            'weighted_credit_score': np.random.uniform(650, 750, self.window_size),
            'on_chain_tvl': np.random.uniform(5e6, 20e6, self.window_size),
            'liquidity_depth_usdc': np.random.uniform(1e6, 5e6, self.window_size),
            'treasury_yield_10y': np.random.uniform(0.03, 0.05, self.window_size),
            'credit_spread_bbb': np.random.uniform(0.015, 0.04, self.window_size),
            'volatility_index': np.random.uniform(10, 30, self.window_size),
        }
        
        return pd.DataFrame(data)
    
    def normalize_features(self, df: pd.DataFrame) -> np.ndarray:
        """
        Normalize features to [0, 1] range
        """
        feature_data = df[self.feature_names].values
        
        # Min-max normalization
        normalized = (feature_data - feature_data.min(axis=0)) / \
                     (feature_data.max(axis=0) - feature_data.min(axis=0) + 1e-8)
        
        return normalized
    
    def prepare_input_tensor(self, df: pd.DataFrame) -> torch.Tensor:
        """
        Convert DataFrame to model input tensor
        """
        normalized = self.normalize_features(df)
        
        # Shape: (1, seq_length, num_features)
        tensor = torch.FloatTensor(normalized).unsqueeze(0)
        
        return tensor


class VeritasMLEngine:
    """
    Main ML inference engine for risk prediction
    """
    
    MODEL_VERSION = "v1.2.0-lstm"
    
    def __init__(self, model_path: str = None):
        self.logger = logging.getLogger(__name__)
        
        # Initialize model
        self.model = RWARiskLSTM(
            input_size=8,  # 8 features
            hidden_size=128,
            num_layers=2
        )
        
        # Load pre-trained weights if available
        if model_path:
            self.model.load_state_dict(torch.load(model_path))
            self.logger.info(f"Loaded model from {model_path}")
        
        self.model.eval()
        
        # Data processor
        self.data_processor = RWADataProcessor(window_size=30)
        
    def predict(self) -> Dict:
        """
        Run inference and return risk/liquidity predictions
        """
        try:
            # Fetch latest market data
            market_data = self.data_processor.fetch_market_data()
            
            # Prepare input tensor
            input_tensor = self.data_processor.prepare_input_tensor(market_data)
            
            # Model inference
            with torch.no_grad():
                risk_score, liquidity_score = self.model(input_tensor)
            
            # Extract values
            risk = float(risk_score.item())
            liquidity = float(liquidity_score.item())
            
            # Calculate confidence based on data quality
            confidence = self._calculate_confidence(market_data)
            
            result = {
                'risk_score': risk,
                'liquidity_score': liquidity,
                'confidence': confidence,
                'model_version': self.MODEL_VERSION,
                'timestamp': int(datetime.now().timestamp()),
                'metadata': {
                    'avg_default_rate': float(market_data['default_rate'].mean()),
                    'tvl': float(market_data['on_chain_tvl'].iloc[-1]),
                    'liquidity_depth': float(market_data['liquidity_depth_usdc'].iloc[-1]),
                }
            }
            
            self.logger.info(f"Prediction: Risk={risk:.4f}, Liquidity={liquidity:.4f}, Confidence={confidence:.4f}")
            
            return result
            
        except Exception as e:
            self.logger.error(f"Prediction failed: {e}")
            raise
    
    def _calculate_confidence(self, df: pd.DataFrame) -> float:
        """
        Calculate prediction confidence based on data quality
        """
        # Check for missing data
        missing_ratio = df.isnull().sum().sum() / (len(df) * len(df.columns))
        
        # Check data variance (low variance = low confidence)
        variance_score = df[self.data_processor.feature_names].var().mean()
        
        # Confidence score: higher is better
        confidence = 1.0 - missing_ratio
        confidence *= min(1.0, variance_score / 100)  # Normalize variance
        
        return max(0.0, min(1.0, confidence))

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)

ml_engine = VeritasMLEngine()

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        'status': 'healthy',
        'model_version': ml_engine.MODEL_VERSION,
        'timestamp': datetime.now().isoformat()
    })

@app.route('/api/v1/risk-assessment', methods=['GET'])
def risk_assessment():
    """
    Main endpoint for RWA risk assessment
    Called by Keeper Bot every 12 hours
    """
    try:
        prediction = ml_engine.predict()
        return jsonify(prediction)
    except Exception as e:
        return jsonify({
            'error': str(e),
            'timestamp': datetime.now().isoformat()
        }), 500

@app.route('/api/v1/scenario/<scenario>', methods=['GET'])
def scenario_analysis(scenario):
    """
    Scenario analysis endpoint
    Scenarios: base, stress, bull
    """
    scenarios = {
        'base': {'default_rate': 0.03, 'liquidity_depth': 3e6},
        'stress': {'default_rate': 0.08, 'liquidity_depth': 1e6},
        'bull': {'default_rate': 0.01, 'liquidity_depth': 10e6},
    }
    
    if scenario not in scenarios:
        return jsonify({'error': 'Invalid scenario'}), 400
    
    # Run prediction with scenario parameters
    prediction = ml_engine.predict()
    
    # Adjust based on scenario
    if scenario == 'stress':
        prediction['risk_score'] = min(1.0, prediction['risk_score'] * 1.5)
        prediction['liquidity_score'] *= 0.5
    elif scenario == 'bull':
        prediction['risk_score'] *= 0.7
        prediction['liquidity_score'] = min(1.0, prediction['liquidity_score'] * 1.2)
    
    return jsonify({
        'scenario': scenario,
        'prediction': prediction
    })

@app.route('/api/v1/leverage-health', methods=['POST'])
def assess_leverage_health():
    """
    Assess health of leveraged RWA strategy position
    Called by keeper bot to monitor LeveragedRWAStrategy contract
    """
    try:
        data = request.get_json()
        
        # Extract position data from smart contract
        total_collateral = data.get('totalCollateral', 0)
        total_borrowed = data.get('totalBorrowed', 0)
        current_health_factor = data.get('currentHealthFactor', 0)
        ait_value = data.get('aitValue', 0)
        
        # Calculate current LTV
        ltv = (total_borrowed / total_collateral) if total_collateral > 0 else 0
        
        # Get current market risk assessment
        market_prediction = ml_engine.predict()
        
        # Calculate position-specific risk
        position_risk = {
            'ltv_risk': min(1.0, ltv / 0.7),  # Risk increases as LTV approaches 70%
            'health_factor_risk': max(0.0, 1.0 - (current_health_factor / 1.5)),
            'market_risk': market_prediction['risk_score'],
            'liquidity_risk': 1.0 - market_prediction['liquidity_score']
        }
        
        # Weighted composite risk score
        composite_risk = (
            position_risk['ltv_risk'] * 0.3 +
            position_risk['health_factor_risk'] * 0.3 +
            position_risk['market_risk'] * 0.25 +
            position_risk['liquidity_risk'] * 0.15
        )
        
        # Risk thresholds
        risk_level = 'LOW'
        action_required = False
        
        if composite_risk > 0.8:
            risk_level = 'CRITICAL'
            action_required = True
        elif composite_risk > 0.6:
            risk_level = 'HIGH'
            action_required = True
        elif composite_risk > 0.4:
            risk_level = 'MEDIUM'
        
        # Recommendations
        recommendations = []
        if ltv > 0.65:
            recommendations.append('REDUCE_LEVERAGE')
        if current_health_factor < 1.3:
            recommendations.append('EMERGENCY_DELEVERAGE')
        if market_prediction['liquidity_score'] < 0.3:
            recommendations.append('PAUSE_NEW_POSITIONS')
        
        return jsonify({
            'composite_risk_score': composite_risk,
            'risk_level': risk_level,
            'action_required': action_required,
            'position_metrics': {
                'ltv': ltv,
                'health_factor': current_health_factor,
                'ait_value_usd': ait_value,
                'exposure_ratio': (ait_value / total_borrowed) if total_borrowed > 0 else 0
            },
            'risk_breakdown': position_risk,
            'recommendations': recommendations,
            'timestamp': int(datetime.now().timestamp())
        })
        
    except Exception as e:
        return jsonify({
            'error': str(e),
            'timestamp': datetime.now().isoformat()
        }), 500

@app.route('/api/v1/kyc-risk-assessment', methods=['POST'])
def assess_kyc_risk():
    """
    ML-based KYC risk assessment for investor verification
    Analyzes patterns to detect potential fraud or compliance issues
    """
    try:
        data = request.get_json()
        
        # Extract KYC data
        investor_data = {
            'investment_amount': data.get('investmentAmount', 0),
            'tier': data.get('tier', 0),
            'jurisdiction': data.get('jurisdiction', ''),
            'transaction_frequency': data.get('transactionFrequency', 0),
            'wallet_age_days': data.get('walletAgeDays', 0),
            'previous_defi_exposure': data.get('previousDefiExposure', 0)
        }
        
        # Risk scoring based on patterns
        risk_factors = {
            'amount_risk': min(1.0, investor_data['investment_amount'] / 1000000),  # Risk increases with amount
            'velocity_risk': min(1.0, investor_data['transaction_frequency'] / 100),
            'wallet_risk': max(0.0, 1.0 - (investor_data['wallet_age_days'] / 365)),  # Newer wallets = higher risk
            'jurisdiction_risk': 0.1 if investor_data['jurisdiction'] in ['US', 'EU', 'UK'] else 0.3
        }
        
        # Composite KYC risk score
        kyc_risk_score = (
            risk_factors['amount_risk'] * 0.3 +
            risk_factors['velocity_risk'] * 0.2 +
            risk_factors['wallet_risk'] * 0.3 +
            risk_factors['jurisdiction_risk'] * 0.2
        )
        
        # Risk classification
        if kyc_risk_score > 0.7:
            risk_classification = 'HIGH_RISK'
            verification_required = True
        elif kyc_risk_score > 0.4:
            risk_classification = 'MEDIUM_RISK'
            verification_required = True
        else:
            risk_classification = 'LOW_RISK'
            verification_required = False
        
        # Compliance flags
        flags = []
        if investor_data['investment_amount'] > 500000:
            flags.append('LARGE_INVESTMENT')
        if investor_data['wallet_age_days'] < 30:
            flags.append('NEW_WALLET')
        if investor_data['transaction_frequency'] > 50:
            flags.append('HIGH_VELOCITY')
        
        return jsonify({
            'kyc_risk_score': kyc_risk_score,
            'risk_classification': risk_classification,
            'verification_required': verification_required,
            'risk_factors': risk_factors,
            'compliance_flags': flags,
            'recommended_tier': min(investor_data['tier'], 2 if kyc_risk_score > 0.5 else 4),
            'timestamp': int(datetime.now().timestamp())
        })
        
    except Exception as e:
        return jsonify({
            'error': str(e),
            'timestamp': datetime.now().isoformat()
        }), 500

@app.route('/api/v1/invoice-nav-prediction', methods=['POST'])
def predict_invoice_nav():
    """
    Predict future NAV for VeritasInvoiceToken based on invoice pool health
    Used by oracle to update NAV and by strategy for yield forecasting
    """
    try:
        data = request.get_json()
        
        # Extract invoice pool data
        pool_data = {
            'total_face_value': data.get('totalFaceValue', 0),
            'number_of_invoices': data.get('numberOfInvoices', 0),
            'weighted_maturity': data.get('weightedMaturity', 0),
            'expected_yield': data.get('expectedYield', 0),
            'current_default_rate': data.get('defaultRate', 0),
            'realized_yield': data.get('realizedYield', 0)
        }
        
        # Get market risk assessment
        market_risk = ml_engine.predict()
        
        # Calculate expected collection rate
        base_collection_rate = 1.0 - pool_data['current_default_rate']
        market_adjusted_rate = base_collection_rate * (1.0 - market_risk['risk_score'] * 0.3)
        
        # Predict future NAV based on collection expectations
        expected_collections = pool_data['total_face_value'] * market_adjusted_rate
        predicted_nav = expected_collections / data.get('totalSupply', 1)
        
        # Calculate confidence based on pool diversification
        diversification_score = min(1.0, pool_data['number_of_invoices'] / 100)
        maturity_risk = min(1.0, pool_data['weighted_maturity'] / 180)
        
        confidence = (
            diversification_score * 0.4 +
            (1.0 - maturity_risk) * 0.3 +
            market_risk['confidence'] * 0.3
        )
        
        # Risk-adjusted yield prediction
        risk_adjusted_yield = pool_data['expected_yield'] * (1.0 - market_risk['risk_score'] * 0.5)
        
        return jsonify({
            'predicted_nav': predicted_nav,
            'confidence': confidence,
            'expected_collection_rate': market_adjusted_rate,
            'risk_adjusted_yield': risk_adjusted_yield,
            'pool_health_score': 1.0 - market_risk['risk_score'],
            'diversification_score': diversification_score,
            'market_risk_impact': market_risk['risk_score'],
            'timestamp': int(datetime.now().timestamp())
        })
        
    except Exception as e:
        return jsonify({
            'error': str(e),
            'timestamp': datetime.now().isoformat()
        }), 500

def train_model(epochs: int = 50, batch_size: int = 32):
    """
    Train LSTM model on historical RWA data
    Run this periodically to update model with new data
    """
    model = RWARiskLSTM(input_size=8, hidden_size=128, num_layers=2)
    
    # Loss functions
    criterion_risk = nn.BCELoss()
    criterion_liquidity = nn.MSELoss()
    
    # Optimizer
    optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
    
    print(f"Training model for {epochs} epochs...")
    
    for epoch in range(epochs):
        # In production: load real training data
        # For now, use synthetic data
        
        # Generate synthetic training batch
        batch_x = torch.randn(batch_size, 30, 8)  # (batch, seq, features)
        batch_y_risk = torch.rand(batch_size, 1)
        batch_y_liquidity = torch.rand(batch_size, 1)
        
        # Forward pass
        pred_risk, pred_liquidity = model(batch_x)
        
        # Calculate losses
        loss_risk = criterion_risk(pred_risk, batch_y_risk)
        loss_liquidity = criterion_liquidity(pred_liquidity, batch_y_liquidity)
        loss = loss_risk + loss_liquidity
        
        # Backward pass
        optimizer.zero_grad()
        loss.backward()
        optimizer.step()
        
        if (epoch + 1) % 10 == 0:
            print(f"Epoch [{epoch+1}/{epochs}], Loss: {loss.item():.4f}")
    
    # Save model
    torch.save(model.state_dict(), 'Veritas_rwa_model.pth')
    print("Model saved to Veritas_rwa_model.pth")

if __name__ == '__main__':
    import sys
    
    if len(sys.argv) > 1 and sys.argv[1] == 'train':
        # Training mode
        train_model(epochs=100)
    else:
        # API server mode
        print("Starting Veritas ML API Server...")
        print(f"Model version: {ml_engine.MODEL_VERSION}")
        app.run(host='0.0.0.0', port=5000, debug=False)
