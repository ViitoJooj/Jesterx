import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import styles from "../styles/pages/PagePreview.module.scss";
import { apiRequest } from "../utils/api";

type PublicPageResponse = {
  data?: {
    meta: {
      name: string;
      page_id: string;
    };
    content: {
      svelte: string;
      header?: string;
      footer?: string;
      show_header?: boolean;
      show_footer?: boolean;
    };
    products: Array<{
      id: string;
      name: string;
      price_cents: number;
      description?: string;
    }>;
  };
};

export function PublicPage() {
  const { tenant, pageId } = useParams<{ tenant: string; pageId: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState<PublicPageResponse["data"] | null>(null);

  const html = useMemo(() => {
    if (!page?.content) return "";
    const parts = [];
    if (page.content.show_header !== false && page.content.header) parts.push(page.content.header);
    parts.push(page.content.svelte);
    if (page.content.show_footer !== false && page.content.footer) parts.push(page.content.footer);
    return parts.join("\n\n");
  }, [page]);

  useEffect(() => {
    if (!tenant || !pageId) {
      setError("Página inválida.");
      setLoading(false);
      return;
    }
    (async () => {
      try {
        const res = await apiRequest<PublicPageResponse["data"]>(`/v1/public/pages/${pageId}`, {
          headers: { "X-Tenant-Page-Id": tenant },
        });
        if (res.data) {
          setPage(res.data);
        } else {
          setError("Página não encontrada.");
        }
      } catch (err: any) {
        setError(err?.message || "Erro ao carregar página pública.");
      } finally {
        setLoading(false);
      }
    })();
  }, [tenant, pageId]);

  if (loading) {
    return (
      <main className={styles.main}>
        <div className={styles.center}>
          <p>Carregando página...</p>
        </div>
      </main>
    );
  }

  if (error || !page) {
    return (
      <main className={styles.main}>
        <div className={styles.center}>
          <p className={styles.error}>{error || "Página não encontrada."}</p>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.fullPage}>
      <iframe title={page.meta?.name || pageId} className={styles.iframe} srcDoc={html} sandbox="allow-same-origin allow-scripts allow-forms allow-popups" />

      {page.products?.length ? (
        <div className={styles.productStrip}>
          <h2>Produtos desta rota</h2>
          <div className={styles.productGrid}>
            {page.products.map((prod) => (
              <article key={prod.id} className={styles.productCard}>
                <strong>{prod.name}</strong>
                {prod.description && <p>{prod.description}</p>}
                <span className={styles.price}>{(prod.price_cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })}</span>
              </article>
            ))}
          </div>
        </div>
      ) : null}
    </main>
  );
}
