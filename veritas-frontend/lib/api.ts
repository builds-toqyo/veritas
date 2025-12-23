import { RiskAssessment, ScenarioResponse, HealthCheck, ScenarioType } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_ML_API_URL || 'http://localhost:5000';

// Mock data for development when API server is not available
const mockRiskAssessment: RiskAssessment = {
  risk_score: 0.75,
  liquidity_score: 0.85,
  timestamp: new Date().toISOString()
};

const mockHealthCheck: HealthCheck = {
  status: 'healthy',
  model_version: 'v2.1.0',
  timestamp: new Date().toISOString()
};

const mockScenarioResponse = (scenario: ScenarioType): ScenarioResponse => ({
  scenario,
  prediction: {
    risk_score: scenario === 'stress' ? 0.95 : scenario === 'bull' ? 0.35 : 0.65,
    liquidity_score: scenario === 'stress' ? 0.45 : scenario === 'bull' ? 0.95 : 0.75,
    timestamp: new Date().toISOString()
  }
});

export async function getRiskAssessment(): Promise<RiskAssessment> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/risk-assessment`);
    if (!response.ok) {
      throw new Error('API not available');
    }
    return response.json();
  } catch (error) {
    // Return mock data when API is not available
    console.warn('Using mock data for risk assessment:', error);
    return new Promise(resolve => {
      setTimeout(() => resolve(mockRiskAssessment), 500);
    });
  }
}

export async function getScenarioAnalysis(scenario: ScenarioType): Promise<ScenarioResponse> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/scenario/${scenario}`);
    if (!response.ok) {
      throw new Error('API not available');
    }
    return response.json();
  } catch (error) {
    // Return mock data when API is not available
    console.warn(`Using mock data for ${scenario} scenario:`, error);
    return new Promise(resolve => {
      setTimeout(() => resolve(mockScenarioResponse(scenario)), 800);
    });
  }
}

export async function getHealthCheck(): Promise<HealthCheck> {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    if (!response.ok) {
      throw new Error('API not available');
    }
    return response.json();
  } catch (error) {
    // Return mock data when API is not available
    console.warn('Using mock data for health check:', error);
    return new Promise(resolve => {
      setTimeout(() => resolve(mockHealthCheck), 300);
    });
  }
}
