import { RiskDashboard } from "@/components/risk-dashboard";
import { Nav } from "@/components/nav";

export default function DashboardPage() {
  return (
    <div className="flex min-h-screen flex-col" style={{backgroundColor: '#030303'}}>
      <Nav />

      <main className="flex-1">
        <div className="container max-w-6xl mx-auto px-4 grid items-center gap-6 pb-8 pt-6 md:py-10">
          <div className="flex max-w-[980px] flex-col items-start gap-2">
            <h1 className="text-3xl font-extrabold leading-tight tracking-tighter md:text-4xl" style={{color: '#f7e9e9'}}>
              Risk Assessment Dashboard
            </h1>
            <p className="text-lg" style={{color: '#aaa4a5'}}>
              Real-time monitoring and analysis of Veritas Vault risk metrics.
            </p>
          </div>

          <RiskDashboard />
        </div>
      </main>
    </div>
  );
}
