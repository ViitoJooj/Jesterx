import { useEffect, useState } from "react";
import styles from "../styles/pages/ThemeStore.module.scss";
import { get, put } from "../utils/api";

type ThemeEntry = {
  id: string;
  page_id: string;
  name: string;
  domain?: string;
  for_sale: boolean;
  owned?: boolean;
  updated_at?: string;
};

export function ThemeStore() {
  const [entries, setEntries] = useState<ThemeEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function load() {
    setLoading(true);
    setError("");
    try {
      const res = await get<ThemeEntry[]>("/v1/themes/store");
      if (res.data) setEntries(res.data);
    } catch (err: any) {
      setError(err?.message || "Não foi possível carregar a loja de temas.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function toggleSale(entry: ThemeEntry) {
    try {
      await put(`/v1/themes/store/${entry.page_id}`, { for_sale: !entry.for_sale });
      await load();
    } catch (err: any) {
      setError(err?.message || "Erro ao atualizar tema.");
    }
  }

  return (
    <main className={styles.main}>
      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Marketplace de temas</p>
          <h1>Veja e publique lojas como temas</h1>
          <p className={styles.lead}>Toda página criada entra aqui. Defina se quer vender ou apenas exibir seu tema.</p>
        </div>
      </header>

      {error && <div className={styles.error}>{error}</div>}

      {loading ? (
        <p>Carregando temas...</p>
      ) : (
        <div className={styles.grid}>
          {entries.map((entry) => (
            <article key={entry.id} className={styles.card}>
              <div className={styles.cardHeader}>
                <div>
                  <p className={styles.cardKicker}>{entry.page_id}</p>
                  <h3>{entry.name}</h3>
                  {entry.domain && <span className={styles.domain}>{entry.domain}</span>}
                </div>
                <span className={`${styles.badge} ${entry.for_sale ? styles.badgeSale : styles.badgeView}`}>
                  {entry.for_sale ? "À venda" : "Somente vitrine"}
                </span>
              </div>
              <p className={styles.muted}>Atualizado {entry.updated_at ? new Date(entry.updated_at).toLocaleString("pt-BR") : "recentemente"}</p>
              {entry.owned ? (
                <button className={styles.toggle} type="button" onClick={() => toggleSale(entry)}>
                  {entry.for_sale ? "Marcar como vitrine" : "Colocar à venda"}
                </button>
              ) : (
                <p className={styles.hint}>Entre com o tenant da loja para habilitar o controle de venda.</p>
              )}
            </article>
          ))}
        </div>
      )}
    </main>
  );
}
