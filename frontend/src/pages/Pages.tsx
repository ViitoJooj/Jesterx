import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../styles/pages/Pages.module.scss";
import { CreatePageForm } from "../components/CreatePageForm";
import { post, get, getCurrentTenant, setCurrentTenant } from "../utils/api";

interface User {
  id: string;
  profile_img: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  plan: string;
}

interface Page {
  id: string;
  name: string;
  page_id: string;
  domain: string;
  created_at: string;
  updated_at: string;
}

export function Pages() {
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [showSiteModal, setShowSiteModal] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [pages, setPages] = useState<Page[]>([]);
  const [loadingPages, setLoadingPages] = useState(false);

  const [siteName, setSiteName] = useState("");
  const [siteSlug, setSiteSlug] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    checkUserAndSite();
  }, []);

  async function checkUserAndSite() {
    try {
      const response = await get("/v1/auth/me");

      if (response.data) {
        setUser(response.data);

        if (!response.data.plan || response.data.plan === "free") {
          navigate("/pricing");
          return;
        }

        if (!getCurrentTenant()) {
          setShowSiteModal(true);
        } else {
          await loadPages();
        }
      } else {
        navigate("/login");
      }
    } catch (err: any) {
      if (err?.status === 401) {
        navigate("/login");
      } else {
        setLoading(false);
      }
    } finally {
      setLoading(false);
    }
  }

  async function loadPages() {
    if (!getCurrentTenant()) return;

    setLoadingPages(true);
    try {
      const response = await get("/v1/pages");
      if (response.data) {
        setPages(response.data);
      }
    } catch (err) {
      console.error("Erro ao carregar páginas:", err);
    } finally {
      setLoadingPages(false);
    }
  }

  async function handleCreateSite(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setCreating(true);

    try {
      await post("/v1/sites", {
        name: siteName,
        slug: siteSlug,
      });

      setCurrentTenant(siteSlug);
      setShowSiteModal(false);
      setSiteName("");
      setSiteSlug("");

      await loadPages();
    } catch (err: any) {
      if (err?.message?.includes("plan") || err?.message?.includes("limit")) {
        alert("Você atingiu o limite do seu plano. Faça upgrade!");
        navigate("/pricing");
        return;
      }
      setError(err?.message || "Erro ao criar site");
    } finally {
      setCreating(false);
    }
  }

  async function handlePageCreated() {
    setIsOpen(false);
    await loadPages();
  }

  function handleViewPage(pageId: string) {
    window.open(`/pages/${pageId}`, "_blank");
  }

  function handleEditPage(pageId: string) {
    navigate(`/pages/${pageId}/edit`);
  }

  if (loading) {
    return (
      <main className={styles.main}>
        <div className={styles.loading}>
          <p>Carregando... </p>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <div className={styles.header}>
        <div>
          <h1>Minhas Páginas 🎉</h1>
        </div>
      </div>

      {loadingPages ? (
        <div className={styles.loading}>
          <p>Carregando páginas...</p>
        </div>
      ) : (
        <div className={styles.projectsContainer}>
          <button
            className={styles.createNewStore}
            onClick={() => {
              if (!getCurrentTenant()) {
                setShowSiteModal(true);
              } else {
                setIsOpen(true);
              }
            }}
          >
            <span className={styles.plus}>+</span>
            <p className={styles.title}>Criar nova página</p>
            <span className={styles.subtitle}>Comece uma nova experiência</span>
          </button>

          {pages.map((page) => (
            <div key={page.id} className={styles.pageCard}>
              <div className={styles.pageHeader}>
                <h3>{page.name}</h3>
                <span className={styles.pageType}>{page.domain || page.page_id}</span>
              </div>

              <div className={styles.pageInfo}>
                <p>
                  <strong>URL:</strong> {page.page_id}
                </p>
                <p className={styles.pageDate}>Criada em {new Date(page.created_at).toLocaleDateString("pt-BR")}</p>
              </div>

              <div className={styles.pageActions}>
                <button className={styles.viewButton} onClick={() => handleViewPage(page.page_id)}>
                  Visualizar
                </button>
                <button className={styles.editButton} onClick={() => handleEditPage(page.page_id)}>
                  Editar
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showSiteModal && (
        <div className={styles.modalOverlay} onClick={() => {}}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <div className={styles.modalHeader}>
              <h2>Primeiro, crie seu site</h2>
              <p>Você precisa ter um site antes de criar páginas</p>
            </div>

            <form onSubmit={handleCreateSite} className={styles.siteForm}>
              {error && <div className={styles.errorBox}>{error}</div>}

              <label>
                Nome do site
                <input type="text" placeholder="Ex:  Minha Loja" value={siteName} onChange={(e) => setSiteName(e.target.value)} required disabled={creating} />
              </label>

              <label>
                URL do site
                <input type="text" placeholder="Ex:  minha-loja" value={siteSlug} onChange={(e) => setSiteSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ""))} required disabled={creating} />
                <small>Seu site: {siteSlug || "sua-url"}. jesterx.com</small>
              </label>

              <div className={styles.modalActions}>
                <button type="submit" className={styles.primaryBtn} disabled={creating}>
                  {creating ? "Criando..." : "Criar site"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isOpen && (
        <div className={styles.modalOverlay} onClick={() => setIsOpen(false)}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <CreatePageForm onClose={handlePageCreated} />
          </div>
        </div>
      )}
    </main>
  );
}
