import { FormEvent, useEffect, useMemo, useState } from "react";
import styles from "../styles/pages/Products.module.scss";
import { get, post, put } from "../utils/api";

type Page = {
  id: string;
  name: string;
  page_id: string;
};

type Product = {
  id: string;
  name: string;
  description?: string;
  price_cents: number;
  images?: string[];
  visible: boolean;
  page_id: string;
};

export function Products() {
  const [pages, setPages] = useState<Page[]>([]);
  const [selectedPage, setSelectedPage] = useState<string>("");
  const [products, setProducts] = useState<Product[]>([]);
  const [loadingPages, setLoadingPages] = useState(true);
  const [loadingProducts, setLoadingProducts] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [visible, setVisible] = useState(true);

  const formattedPageName = useMemo(() => pages.find((p) => p.page_id === selectedPage)?.name || "", [pages, selectedPage]);

  useEffect(() => {
    (async () => {
      try {
        const res = await get<Page[]>("/v1/pages");
        if (res.data) {
          setPages(res.data);
          if (res.data[0]) {
            setSelectedPage(res.data[0].page_id);
          }
        }
      } catch (err: any) {
        setError(err?.message || "Não foi possível carregar páginas.");
      } finally {
        setLoadingPages(false);
      }
    })();
  }, []);

  useEffect(() => {
    if (!selectedPage) {
      setProducts([]);
      return;
    }

    (async () => {
      setLoadingProducts(true);
      try {
        const res = await get<Product[]>(`/v1/pages/${selectedPage}/products`);
        if (res.data) setProducts(res.data);
      } catch (err: any) {
        setError(err?.message || "Erro ao carregar produtos.");
      } finally {
        setLoadingProducts(false);
      }
    })();
  }, [selectedPage]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    if (!selectedPage) {
      setError("Selecione uma página para vincular o produto.");
      return;
    }

    setSaving(true);
    try {
      const normalized = price.trim().replace(",", ".");
      const parsed = Number(normalized);
      if (!Number.isFinite(parsed)) {
        setError("Preço inválido. Use apenas números, vírgula ou ponto.");
        setSaving(false);
        return;
      }

      const payload = {
        name,
        description,
        price_cents: Math.max(0, Math.round(parsed * 100)),
        images: [],
        visible,
      };
      await post(`/v1/pages/${selectedPage}/products`, payload);
      setName("");
      setDescription("");
      setPrice("");
      setVisible(true);
      const res = await get<Product[]>(`/v1/pages/${selectedPage}/products`);
      if (res.data) setProducts(res.data);
    } catch (err: any) {
      setError(err?.message || "Erro ao salvar produto.");
    } finally {
      setSaving(false);
    }
  }

  async function toggleVisibility(product: Product) {
    try {
      await put(`/v1/pages/${product.page_id}/products/${product.id}`, {
        visible: !product.visible,
      });
      const res = await get<Product[]>(`/v1/pages/${product.page_id}/products`);
      if (res.data) setProducts(res.data);
    } catch (err: any) {
      setError(err?.message || "Não foi possível atualizar o produto.");
    }
  }

  return (
    <main className={styles.main}>
      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Produtos das suas páginas</p>
          <h1>Cadastre e exiba itens nas lojas</h1>
          <p className={styles.lead}>Conecte seus produtos às rotas que você criou. Tudo fica disponível publicamente pelo tenant da página.</p>
        </div>
        <div className={styles.selector}>
          <label>
            Página
            <select value={selectedPage} onChange={(e) => setSelectedPage(e.target.value)} disabled={loadingPages || saving}>
              {pages.map((page) => (
                <option key={page.id} value={page.page_id}>
                  {page.name} ({page.page_id})
                </option>
              ))}
            </select>
          </label>
        </div>
      </header>

      {error && <div className={styles.error}>{error}</div>}

      <section className={styles.content}>
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.formHeader}>
            <h2>Novo produto</h2>
            <p>Produtos ficam atrelados à rota selecionada.</p>
          </div>

          <label className={styles.field}>
            Nome do produto
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="Ex: Camiseta premium" required disabled={saving || loadingPages} />
          </label>

          <label className={styles.field}>
            Descrição
            <textarea value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Conte a história do produto" disabled={saving || loadingPages} />
          </label>

          <div className={styles.inlineRow}>
            <label className={styles.field}>
              Preço (R$)
              <input type="number" step="0.01" value={price} onChange={(e) => setPrice(e.target.value)} placeholder="129.90" disabled={saving || loadingPages} />
            </label>

            <label className={styles.checkbox}>
              <input type="checkbox" checked={visible} onChange={(e) => setVisible(e.target.checked)} disabled={saving || loadingPages} />
              <span>Exibir nas páginas públicas</span>
            </label>
          </div>

          <div className={styles.formActions}>
            <button type="submit" disabled={saving || loadingPages || !selectedPage}>
              {saving ? "Salvando..." : "Adicionar produto"}
            </button>
          </div>
        </form>

        <div className={styles.listSection}>
          <div className={styles.listHeader}>
            <div>
              <h2>Produtos em {formattedPageName || selectedPage || "sua página"}</h2>
              <p>Estes itens são expostos no catálogo público da rota.</p>
            </div>
          </div>

          {loadingProducts ? (
            <p>Carregando produtos...</p>
          ) : products.length === 0 ? (
            <p className={styles.empty}>Nenhum produto cadastrado para esta página.</p>
          ) : (
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Nome</th>
                  <th>Preço</th>
                  <th>Visibilidade</th>
                  <th>Ações</th>
                </tr>
              </thead>
              <tbody>
                {products.map((product) => (
                  <tr key={product.id}>
                    <td>
                      <strong>{product.name}</strong>
                      {product.description && <p className={styles.muted}>{product.description}</p>}
                    </td>
                    <td>{(product.price_cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })}</td>
                    <td>
                      <span className={`${styles.badge} ${product.visible ? styles.badgeSuccess : styles.badgeMuted}`}>
                        {product.visible ? "Publicado" : "Oculto"}
                      </span>
                    </td>
                    <td>
                      <button className={styles.linkButton} type="button" onClick={() => toggleVisibility(product)}>
                        {product.visible ? "Ocultar" : "Publicar"}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </section>
    </main>
  );
}
