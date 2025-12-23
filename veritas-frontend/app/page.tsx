import { Nav } from "@/components/nav";
import { HeroSection, StatsSection, FeaturesSection } from "./home";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col" style={{backgroundColor: '#030303', color: '#f7e9e9'}}>
      <Nav />

      <main className="flex-1 overflow-hidden">
        <HeroSection />
        <StatsSection />
        <FeaturesSection />

        {/* CTA Section */}
        <section className="py-32 relative" style={{borderTop: '1px solid #635a5e', backgroundColor: '#030303'}}>
          <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-5 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
          <div className="container max-w-6xl mx-auto px-4 flex flex-col items-center text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl" style={{color: '#f7e9e9'}}>
              Ready to enhance your risk management?
            </h2>
            <p className="mt-4 max-w-2xl text-lg" style={{color: '#aaa4a5'}}>
              Get started with Veritas today and access enterprise-grade risk assessment tools.
            </p>
            <Button 
              asChild 
              size="lg" 
              className="mt-8 transition-colors duration-300"
              style={{backgroundColor: '#aaa4a5', color: '#f7e9e9'}}
            >
              <Link href="/dashboard">Access Dashboard</Link>
            </Button>
          </div>
        </section>
      </main>

      <footer className="py-8" style={{borderTop: '1px solid #635a5e'}}>
        <div className="container max-w-6xl mx-auto px-4">
          <p className="text-center text-sm" style={{color: '#aaa4a5'}}>
            © 2025 Veritas. Built with advanced ML models for optimal risk management.
          </p>
        </div>
      </footer>
    </div>
  );
}