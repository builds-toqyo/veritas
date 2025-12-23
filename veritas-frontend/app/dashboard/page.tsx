import { RiskDashboard } from "@/components/risk-dashboard";

export default function DashboardPage() {
  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 items-center">
          <div className="mr-4 hidden md:flex">
            <a className="mr-6 flex items-center space-x-2" href="/">
              <span className="hidden font-bold sm:inline-block">Veritas Risk Dashboard</span>
            </a>
          </div>
        </div>
      </header>

      <main className="flex-1">
        <div className="container grid items-center gap-6 pb-8 pt-6 md:py-10">
          <div className="flex max-w-[980px] flex-col items-start gap-2">
            <h1 className="text-3xl font-extrabold leading-tight tracking-tighter md:text-4xl">
              Risk Assessment Dashboard
            </h1>
            <p className="text-lg text-muted-foreground">
              Real-time monitoring and analysis of Veritas Vault risk metrics.
            </p>
          </div>

          <RiskDashboard />
        </div>
      </main>
    </div>
  );
}
