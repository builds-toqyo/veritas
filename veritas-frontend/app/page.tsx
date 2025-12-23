import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Nav } from "@/components/nav";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      <Nav />

      <main className="flex-1">
        {/* Hero Section */}
        <section className="space-y-6 pb-8 pt-6 md:pb-12 md:pt-10 lg:py-32">
          <div className="container flex max-w-[64rem] flex-col items-center gap-4 text-center">
            <h1 className="font-bold text-3xl sm:text-5xl md:text-6xl lg:text-7xl">
              Real-Time Risk Assessment for
              <span className="text-primary"> Veritas Vault</span>
            </h1>
            <p className="max-w-[42rem] leading-normal text-muted-foreground sm:text-xl sm:leading-8">
              Advanced ML-powered risk analysis and monitoring system for DeFi operations. 
              Get real-time insights and scenario analysis for optimal risk management.
            </p>
            <div className="space-x-4">
              <Button asChild size="lg" className="px-8">
                <Link href="/dashboard">Launch Dashboard</Link>
              </Button>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="container space-y-6 bg-slate-50 py-8 dark:bg-transparent md:py-12 lg:py-24">
          <div className="mx-auto flex max-w-[58rem] flex-col items-center space-y-4 text-center">
            <h2 className="font-bold text-3xl leading-[1.1] sm:text-3xl md:text-6xl">
              Features
            </h2>
            <p className="max-w-[85%] leading-normal text-muted-foreground sm:text-lg sm:leading-7">
              Comprehensive risk management tools powered by advanced machine learning
            </p>
          </div>

          <div className="mx-auto grid justify-center gap-4 sm:grid-cols-2 md:max-w-[64rem] md:grid-cols-3">
            <Card className="flex flex-col justify-between">
              <CardHeader>
                <CardTitle>Real-Time Monitoring</CardTitle>
                <CardDescription>Continuous risk assessment and liquidity analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  Get instant insights into your vault's risk profile and market conditions
                </p>
              </CardContent>
            </Card>

            <Card className="flex flex-col justify-between">
              <CardHeader>
                <CardTitle>Scenario Analysis</CardTitle>
                <CardDescription>Simulate different market conditions</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  Test your strategies against base, stress, and bull market scenarios
                </p>
              </CardContent>
            </Card>

            <Card className="flex flex-col justify-between">
              <CardHeader>
                <CardTitle>ML-Powered Insights</CardTitle>
                <CardDescription>Advanced predictive analytics</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  Leverage machine learning models for accurate risk prediction
                </p>
              </CardContent>
            </Card>
          </div>
        </section>
      </main>

      <footer className="border-t py-6 md:py-0">
        <div className="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
          <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
            <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
              Built by Veritas Team. Powered by advanced ML models.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
