import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Nav() {
  return (
    <header className="sticky top-0 z-50 w-full nav-header">
      <div className="container max-w-6xl mx-auto px-4 flex h-16 items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/" className="nav-brand">
            Veritas
          </Link>
        </div>
        <nav className="flex items-center space-x-6">
          <Link 
            href="/dashboard" 
            className="text-sm font-medium nav-link"
          >
            Dashboard
          </Link>
          <Button 
            asChild 
            className="btn-primary"
          >
            <Link href="/dashboard">Launch App</Link>
          </Button>
        </nav>
      </div>
    </header>
  );
}
