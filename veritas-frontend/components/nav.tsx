import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Nav() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-zinc-800 bg-black">
      <div className="container max-w-6xl mx-auto px-4 flex h-16 items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/" className="text-xl font-bold text-white hover:text-zinc-200 transition-colors duration-300">
            Veritas
          </Link>
        </div>
        <nav className="flex items-center space-x-6">
          <Link 
            href="/dashboard" 
            className="text-sm font-medium text-zinc-400 hover:text-white transition-colors duration-300"
          >
            Dashboard
          </Link>
          <Button 
            asChild 
            className="bg-zinc-800 text-white hover:bg-zinc-700 transition-colors duration-300"
          >
            <Link href="/dashboard">Launch App</Link>
          </Button>
        </nav>
      </div>
    </header>
  );
}
