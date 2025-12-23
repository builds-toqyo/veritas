import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Nav() {
  return (
    <header className="sticky top-0 z-50 w-full" style={{borderBottom: '1px solid #635a5e', backgroundColor: '#030303'}}>
      <div className="container max-w-6xl mx-auto px-4 flex h-16 items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/" className="text-xl font-bold transition-colors duration-300" style={{color: '#f7e9e9'}}>
            Veritas
          </Link>
        </div>
        <nav className="flex items-center space-x-6">
          <Link 
            href="/dashboard" 
            className="text-sm font-medium transition-colors duration-300"
            style={{color: '#aaa4a5'}}
          >
            Dashboard
          </Link>
          <Button 
            asChild 
            className="transition-colors duration-300"
            style={{backgroundColor: '#aaa4a5', color: '#f7e9e9'}}
          >
            <Link href="/dashboard">Launch App</Link>
          </Button>
        </nav>
      </div>
    </header>
  );
}
