"use client";

import Link from "next/link";

export default function Header() {
  return (
    <header>
      <div className="flex justify-between items-center px-5 py-5">
        <Link href="/" className="text-xl font-bold">
          Mr Poll
        </Link>
        <nav className="space-x-4">
          <Link href="/" className="hover:text-pink-700">
            Placeholder
          </Link>
          <Link href="/" className="hover:text-pink-700">
            Placeholder
          </Link>
          <Link href="/conact" className="hover:text-pink-700">
            Placeholder
          </Link>
        </nav>
      </div>
    </header>
  );
}
