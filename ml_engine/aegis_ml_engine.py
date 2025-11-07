"""
Veritas ML Engine - LSTM-based RWA Risk Prediction
Predicts default risk and liquidity scores for tokenized invoices/bonds
"""

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
