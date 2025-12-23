import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Nav() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-[#181919] bg-[#040404]">
      <div className="container max-w-6xl mx-auto px-4 flex h-16 items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/" className="text-xl font-bold text-[#a6b8b3] hover:text-[#808d8e] transition-colors duration-300">
            Veritas
          </Link>
        </div>
        <nav className="flex items-center space-x-6">
          <Link 
            href="/dashboard" 
            className="text-sm font-medium text-[#808d8e] hover:text-[#a6b8b3] transition-colors duration-300"
          >
            Dashboard
          </Link>
          <Button 
            asChild 
            className="bg-[#373b3b] text-[#a6b8b3] hover:bg-[#565f5f] transition-colors duration-300"
          >
            <Link href="/dashboard">Launch App</Link>
          </Button>
        </nav>
      </div>
    </header>
  );
}
