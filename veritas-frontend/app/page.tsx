import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Nav } from "@/components/nav";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-white">
      <Nav />

      <main className="flex-1 overflow-hidden">
        {/* Hero Section */}
        <section className="relative min-h-[80vh] flex items-center justify-center py-20">
          <div className="absolute inset-0 bg-[url('/grid.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
          <div className="container max-w-6xl mx-auto px-4">
            <div className="flex flex-col items-center gap-8 text-center mx-auto max-w-3xl">
              <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl text-black">
                Real-Time Risk Assessment for{" "}
                <span className="bg-gradient-to-r from-black to-neutral-600 bg-clip-text text-transparent">
                  Veritas Vault
                </span>
              </h1>
              <p className="mx-auto max-w-[700px] text-lg text-neutral-600 md:text-xl">
                Advanced ML-powered risk analysis and monitoring system for DeFi operations.
              </p>
              <div className="mt-6">
                <Button asChild size="lg" className="bg-black text-white hover:bg-neutral-800">
                  <Link href="/dashboard">Launch Dashboard</Link>
                </Button>
              </div>
            </div>
          </div>
        </section>

        {/* Stats Section */}
        <section className="border-y border-neutral-200 bg-white py-20">
          <div className="container max-w-6xl mx-auto px-4">
            <div className="mx-auto grid gap-16 md:grid-cols-3 max-w-4xl">
              <div className="text-center">
                <h3 className="text-4xl font-bold text-black">99.9%</h3>
                <p className="mt-2 text-neutral-600">Uptime</p>
              </div>
              <div className="text-center">
                <h3 className="text-4xl font-bold text-black">24/7</h3>
                <p className="mt-2 text-neutral-600">Monitoring</p>
              </div>
              <div className="text-center">
                <h3 className="text-4xl font-bold text-black">ML v2.0</h3>
                <p className="mt-2 text-neutral-600">Latest Model</p>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-20">
          <div className="container max-w-6xl mx-auto px-4">
            <div className="mx-auto mb-16 max-w-3xl text-center">
              <h2 className="text-3xl font-bold tracking-tight text-black sm:text-4xl lg:text-5xl">
                Enterprise-Grade Risk Management
              </h2>
              <p className="mt-6 text-lg text-neutral-600 max-w-2xl mx-auto">
                Comprehensive suite of tools powered by advanced machine learning algorithms
              </p>
            </div>
            <div className="mx-auto grid max-w-5xl gap-8 md:grid-cols-3">
              <Card className="border-2 border-neutral-200">
                <CardHeader>
                  <CardTitle className="text-black">Real-Time Analysis</CardTitle>
                  <CardDescription className="text-neutral-600">
                    Continuous monitoring and assessment
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className="list-inside list-disc space-y-2 text-neutral-600">
                    <li>Live risk scoring</li>
                    <li>Market condition tracking</li>
                    <li>Automated alerts</li>
                  </ul>
                </CardContent>
              </Card>

              <Card className="border-2 border-neutral-200">
                <CardHeader>
                  <CardTitle className="text-black">Scenario Testing</CardTitle>
                  <CardDescription className="text-neutral-600">
                    Advanced market simulations
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className="list-inside list-disc space-y-2 text-neutral-600">
                    <li>Stress testing</li>
                    <li>Bull market scenarios</li>
                    <li>Risk forecasting</li>
                  </ul>
                </CardContent>
              </Card>

              <Card className="border-2 border-neutral-200">
                <CardHeader>
                  <CardTitle className="text-black">ML Intelligence</CardTitle>
                  <CardDescription className="text-neutral-600">
                    Smart predictive analytics
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className="list-inside list-disc space-y-2 text-neutral-600">
                    <li>Pattern recognition</li>
                    <li>Anomaly detection</li>
                    <li>Trend analysis</li>
                  </ul>
                </CardContent>
              </Card>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="border-t border-neutral-200 bg-neutral-50 py-32">
          <div className="container max-w-6xl mx-auto px-4 flex flex-col items-center text-center">
            <h2 className="text-3xl font-bold tracking-tight text-black sm:text-4xl">
              Ready to enhance your risk management?
            </h2>
            <p className="mt-4 max-w-2xl text-lg text-neutral-600">
              Get started with Veritas today and access enterprise-grade risk assessment tools.
            </p>
            <Button 
              asChild 
              size="lg" 
              className="mt-8 bg-black text-white hover:bg-neutral-800"
            >
              <Link href="/dashboard">Access Dashboard</Link>
            </Button>
          </div>
        </section>
      </main>

      <footer className="border-t border-neutral-200 py-8">
        <div className="container">
          <p className="text-center text-sm text-neutral-600">
            © 2025 Veritas. Built with advanced ML models for optimal risk management.
          </p>
        </div>
      </footer>
    </div>
  );
}
