import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ThemeCard } from "../components/ThemeCard";
import styles from "../styles/pages/ThemeStore.module.scss";
import { getPublic } from "../utils/api";

type ThemeEntry = {
  id: string;
  slug: string;
  name: string;
  description?: string;
  thumbnail?: string;
  rating?: number | null;
  installs?: number | null;
};

export function ThemeStore() {
  const [themes, setThemes] = useState<ThemeEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    load();
  }, []);

  async function load() {
    setLoading(true);
    setError(false);

    try {
      const res = await getPublic<ThemeEntry[]>("/v1/themes/store");
      setThemes(Array.isArray(res.data) ? res.data : []);
    } catch {
      setThemes([]);
      setError(true);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className={styles.main}>
      <header className={styles.header}>
        <h1>Loja de temas</h1>
        <p>Templates modernos prontos para produção</p>
      </header>

      {loading && <p>Carregando...</p>}

      {!loading && error && (
        <p className={styles.error}>Não foi possível carregar os temas</p>
      )}

      {!loading && !error && (
        <section className={styles.grid}>
          {themes.map((theme) => (
            <ThemeCard
              key={theme.id}
              name={theme.name}
              description={theme.description || "Tema moderno e otimizado"}
              thumbnail={theme.thumbnail}
              rating={theme.rating}
              installs={theme.installs}
              onClick={() => navigate(`/themes/${theme.slug}`)}
            />
          ))}
        </section>
      )}
    </main>
  );
}
