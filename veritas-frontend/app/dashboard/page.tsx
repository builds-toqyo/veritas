import { RiskDashboard } from "@/components/risk-dashboard";
import { Nav } from "@/components/nav";

export default function DashboardPage() {
  return (
    <div className="flex min-h-screen flex-col bg-black">
      <Nav />

      <main className="flex-1">
        <div className="container max-w-6xl mx-auto px-4 grid items-center gap-6 pb-8 pt-6 md:py-10">
          <div className="flex max-w-[980px] flex-col items-start gap-2">
            <h1 className="text-3xl font-extrabold leading-tight tracking-tighter text-white md:text-4xl">
              Risk Assessment Dashboard
            </h1>
            <p className="text-lg text-[#373b3b]">
              Real-time monitoring and analysis of Veritas Vault risk metrics.
            </p>
          </div>

          <RiskDashboard />
        </div>
      </main>
    </div>
  );
}
