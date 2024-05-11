import styles from "./page.module.css";
import { ClientSideComponent } from "@/app/client";
import { ServerSideComponent } from "@/app/server";

export default function Home() {
  return (
    <main className={styles.main}>
      <ClientSideComponent />
      <ServerSideComponent />
    </main>
  );
}
