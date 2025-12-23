"use client"

import { useEffect, useState } from 'react'
import { RiskAssessment, ScenarioResponse, HealthCheck } from '@/lib/types'
import { getRiskAssessment, getScenarioAnalysis, getHealthCheck } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Alert, AlertDescription } from "@/components/ui/alert"

export function RiskDashboard() {
  const [riskData, setRiskData] = useState<RiskAssessment | null>(null)
  const [scenarioData, setScenarioData] = useState<ScenarioResponse | null>(null)
  const [healthData, setHealthData] = useState<HealthCheck | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const risk = await getRiskAssessment()
        setRiskData(risk)
        const health = await getHealthCheck()
        setHealthData(health)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to fetch data')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
    // Refresh data every 5 minutes
    const interval = setInterval(fetchData, 5 * 60 * 1000)
    return () => clearInterval(interval)
  }, [])

  const handleScenarioAnalysis = async (scenario: 'base' | 'stress' | 'bull') => {
    try {
      setLoading(true)
      const data = await getScenarioAnalysis(scenario)
      setScenarioData(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to run scenario analysis')
    } finally {
      setLoading(false)
    }
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    )
  }

  return (
    <Tabs defaultValue="overview" className="w-full">
      <TabsList>
        <TabsTrigger value="overview">Overview</TabsTrigger>
        <TabsTrigger value="scenarios">Scenario Analysis</TabsTrigger>
        <TabsTrigger value="health">System Health</TabsTrigger>
      </TabsList>
      
      <TabsContent value="overview" className="space-y-4">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Risk Score</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {riskData ? riskData.risk_score.toFixed(2) : '-'}
              </div>
              <Progress 
                value={riskData ? riskData.risk_score * 100 : 0} 
                className="mt-2" 
              />
              <p className="text-xs text-muted-foreground mt-2">
                Updated {new Date(riskData?.timestamp || '').toLocaleTimeString()}
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Liquidity Score</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {riskData ? riskData.liquidity_score.toFixed(2) : '-'}
              </div>
              <Progress 
                value={riskData ? riskData.liquidity_score * 100 : 0} 
                className="mt-2" 
              />
              <p className="text-xs text-muted-foreground mt-2">
                Updated {new Date(riskData?.timestamp || '').toLocaleTimeString()}
              </p>
            </CardContent>
          </Card>
        </div>
      </TabsContent>

      <TabsContent value="scenarios" className="space-y-4">
        <div className="grid gap-4 md:grid-cols-3">
          {['base', 'stress', 'bull'].map((scenario) => (
            <Card key={scenario}>
              <CardHeader>
                <CardTitle className="capitalize">{scenario} Scenario</CardTitle>
                <CardDescription>Analysis of {scenario} market conditions</CardDescription>
              </CardHeader>
              <CardContent>
                {scenarioData && scenarioData.scenario === scenario ? (
                  <div className="space-y-2">
                    <div>Risk Score: {scenarioData.prediction.risk_score.toFixed(2)}</div>
                    <div>Liquidity Score: {scenarioData.prediction.liquidity_score.toFixed(2)}</div>
                  </div>
                ) : (
                  <Button 
                    className="w-full" 
                    onClick={() => handleScenarioAnalysis(scenario as 'base' | 'stress' | 'bull')}
                    disabled={loading}
                  >
                    Run Analysis
                  </Button>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      </TabsContent>

      <TabsContent value="health" className="space-y-4">
        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
            <CardDescription>Current health metrics of the system</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span>ML Engine</span>
                <span className={healthData?.status === 'healthy' ? 'text-green-500' : 'text-red-500'}>
                  {healthData?.status || 'Unknown'}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span>Model Version</span>
                <span>{healthData?.model_version || 'Unknown'}</span>
              </div>
              <div className="flex items-center justify-between">
                <span>Last Update</span>
                <span>
                  {healthData ? new Date(healthData.timestamp).toLocaleString() : 'Unknown'}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  )
}
