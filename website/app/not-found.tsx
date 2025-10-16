import Image from "next/image";

export default function Home() {
  return (
    <section className="text-center">
      <h1 className="text-4xl font-bold mb-4">Wrong page buddy</h1>
      <div className="flex items-center justify-center"><Image src="/paul.jpg" width={500} height={500} alt="paul" /></div>
      
    </section>
  );
}
