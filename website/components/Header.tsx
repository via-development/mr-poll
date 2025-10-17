"use client";

import Link from "next/link";
import Image from "next/image";

export default function Header() {
  return (
    <header className="bg-primary">
      <nav className="w-full px-6 md:px-12 lg:px-20">
        <div className="flex items-center h-16">
          <div className="flex items-center space-x-2">
            <Image
              src="/images/mrpollclear.png"
              alt="Mr Poll Logo"
              width={24}
              height={24}
            />
            <Link href="/" className="text-brand font-bold text-lg">
              Mr Poll
            </Link>
          </div>
          <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
            <Link
              href="/documentation"
              className="px-2 text-sm hover:text-brand transition duration-200 ease-in-out"
            >
              Documentation
            </Link>
            <Link
              href="/dashboard"
              className=" text-sm hover:text-brand transition duration-200 ease-in-out"
            >
              Dashboard
            </Link>
          </div>
        </div>
      </nav>
    </header>
  );
}
