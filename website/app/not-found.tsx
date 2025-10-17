import Link from "next/link";

export default function Home() {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">Not Found</h1>
        <p className="mt-10">
          This page doesn't exist. Go <Link className="text-brand" href="/">home</Link>.
        </p>
      </div>
    </div>
  );
}