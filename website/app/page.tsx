import Link from "next/link";

export default function Home() {
  return (
    <section className="flex flex-col md:flex-row items-center justify-center w-full min-h-[calc(100vh-64px)] px-8 md:px-20">
      <div className="flex flex-col gap-6 w-full max-w-sm md:max-w-md">
        <div className="h-16 bg-gray-500/30 rounded-xl" />
        <div className="h-16 bg-gray-500/30 rounded-xl" />
        <div className="h-16 bg-gray-500/30 rounded-xl" />
        <div className="h-16 bg-gray-500/30 rounded-xl" />
      </div>
      <div className="hidden md:block w-32" />
      <div className="flex flex-col items-center md:items-start max-w-2xl text-center md:text-left pt-12 md:pt-0">
        <h1 className="text-5xl md:text-6xl font-extrabold leading-tight mb-4">
          Looking like the best bot for
          <span className="text-secondary"> polls and suggestions</span>
        </h1>
        <p className="text-lg mb-8 max-w-md md:max-w-lg">
          Mr Poll is a question and feedback focused Discord bot powering 30k+
          servers.
        </p>
        <div className="flex flex-row gap-4">
          <Link
            href="/invite"
            className="px-6 py-3 bg-primary hover:bg-secondary rounded-xl font-semibold text-white transition"
          >
            Add Mr Poll
          </Link>
          <Link
            href="/documentaion"
            className="px-6 py-3 bg-primary hover:bg-secondary rounded-xl font-semibold text-white transition"
          >
            Documentation
          </Link>
          <Link
            href="/dashboard"
            className="px-6 py-3 bg-secondary hover:bg-primary rounded-xl font-semibold text-white transition"
          >
            Dashboard
          </Link>
        </div>
      </div>
    </section>
  );
}
