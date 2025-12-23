import { Nav } from "@/components/nav";
import { HeroSection, StatsSection, FeaturesSection } from "./home";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      <Nav />

      <main className="flex-1 overflow-hidden">
        <HeroSection />
        <StatsSection />
        <FeaturesSection />

        <section className="py-32 relative bg-light border-t border-muted">
          <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-5 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
          <div className="container max-w-6xl mx-auto px-4 flex flex-col items-center text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl text-primary">
              Ready to enhance your risk management?
            </h2>
            <p className="mt-4 max-w-2xl text-lg text-muted">
              Get started with Veritas today and access enterprise-grade risk assessment tools.
            </p>
            <Button 
              asChild 
              size="lg" 
              className="mt-8 btn-primary"
            >
              <Link href="/dashboard">Access Dashboard</Link>
            </Button>
          </div>
        </section>
      </main>

      <footer className="py-8 footer">
        <div className="container max-w-6xl mx-auto px-4">
          <p className="text-center text-sm text-light opacity-70">
            © 2025 Veritas. Built with advanced ML models for optimal risk management.
          </p>
        </div>
      </footer>
    </div>
  );
}