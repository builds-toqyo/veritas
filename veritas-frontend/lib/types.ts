export interface RiskAssessment {
  risk_score: number;
  liquidity_score: number;
  timestamp: string;
}

export interface ScenarioResponse {
  scenario: 'base' | 'stress' | 'bull';
  prediction: RiskAssessment;
}

export interface HealthCheck {
  status: string;
  model_version: string;
  timestamp: string;
}

export type ScenarioType = 'base' | 'stress' | 'bull';
