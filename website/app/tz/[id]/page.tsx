import axios from "axios";
import { use } from "react";

export default async function Home({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params;
  let isErr = false
  const res = await axios.post("https://localhost:3001", {
    body: {
      offset: -(new Date().getTimezoneOffset())
    }
  }).catch(err => {
    isErr = true
    return err
  })

  if (!isErr)
    return (
      <section className="text-center">
        <h1 className="text-4xl font-bold mb-4">Timezone Updated</h1>
        <p className="mt-10">
          Your timezone has been updated on Mr Poll. You can close this page.
        </p>
      </section>
    );
  else
    return (
      <section className="text-center">
        <h1 className="text-4xl font-bold mb-4">Timezone Update Failed</h1>
        <p className="mt-10">
          Failed to update your timezone, please try the command on the bot again.
        </p>
      </section>
    );
}