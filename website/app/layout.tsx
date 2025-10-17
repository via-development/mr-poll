import type { Metadata } from "next";
import "./globals.css";
import Header from "@/components/Header";

export const metadata: Metadata = {
  title: "Mr Poll",
  description: "kevin/kefin/muffin/kyle/dexter morgan",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <Header />
        <main className="flex items-center justify-center min-h-screen px-6 max-w-7xl mx-auto">
          {children}
        </main>
      </body>
    </html>
  );
}
