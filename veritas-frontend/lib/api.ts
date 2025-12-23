import { RiskAssessment, ScenarioResponse, HealthCheck, ScenarioType } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_ML_API_URL || 'http://localhost:5000';

export async function getRiskAssessment(): Promise<RiskAssessment> {
  const response = await fetch(`${API_BASE_URL}/api/v1/risk-assessment`);
  if (!response.ok) {
    throw new Error('Failed to fetch risk assessment');
  }
  return response.json();
}

export async function getScenarioAnalysis(scenario: ScenarioType): Promise<ScenarioResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/scenario/${scenario}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch ${scenario} scenario analysis`);
  }
  return response.json();
}

export async function getHealthCheck(): Promise<HealthCheck> {
  const response = await fetch(`${API_BASE_URL}/health`);
  if (!response.ok) {
    throw new Error('Failed to fetch health status');
  }
  return response.json();
}
