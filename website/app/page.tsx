import Image from "next/image";

export default function Home() {
  return (
    <section className="text-center">
      <h1 className="text-4xl font-bold mb-4">Mr Paul</h1>
      <p className="text-pink-700">
        Here at mr paul we love pauls, such as paul hollywood and paul rudd
      </p>
      <div className="flex items-center justify-center"><Image src="/paul.jpg" width={500} height={500} alt="paul" /></div>
      
    </section>
  );
}
