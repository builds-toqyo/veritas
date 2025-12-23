import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Nav() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-neutral-200 bg-white">
      <div className="container max-w-6xl mx-auto px-4 flex h-16 items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/" className="text-xl font-bold text-black hover:text-neutral-800">
            Veritas
          </Link>
        </div>
        <nav className="flex items-center space-x-6">
          <Link 
            href="/dashboard" 
            className="text-sm font-medium text-neutral-600 hover:text-black"
          >
            Dashboard
          </Link>
          <Button 
            asChild 
            className="bg-black text-white hover:bg-neutral-800"
          >
            <Link href="/dashboard">Launch App</Link>
          </Button>
        </nav>
      </div>
    </header>
  );
}
