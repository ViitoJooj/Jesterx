import { useEffect, useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../styles/pages/Pages.module.scss";
import { CreatePageForm } from "../components/CreatePageForm";
import { ThemeInstallerModal } from "../components/ThemeInstallerModal";
import { DeletePageModal } from "../components/DeletePageModal";
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
  tenant_id: string;
  name: string;
  page_id: string;
  domain?: string;
  theme_id?: string;
  created_at: string;
  updated_at: string;
}

export function Pages() {
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [showSiteModal, setShowSiteModal] = useState(false);
  const [showThemeModal, setShowThemeModal] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [pages, setPages] = useState<Page[]>([]);
  const [loadingPages, setLoadingPages] = useState(false);

  const [siteName, setSiteName] = useState("");
  const [siteSlug, setSiteSlug] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState("");

  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [pageToDelete, setPageToDelete] = useState<Page | null>(null);

  useEffect(() => {
    checkUserAndSite();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function themeKey(tenant: string) {
    return `theme:${tenant}`;
  }

  function pendingThemeKey(tenant: string) {
    return `pendingTheme:${tenant}`;
  }

  function shouldOpenThemeModal(tenant: string) {
    const hasPending = sessionStorage.getItem(pendingThemeKey(tenant)) === "1";
    const hasTheme = !!localStorage.getItem(themeKey(tenant));
    if (hasTheme && hasPending) sessionStorage.removeItem(pendingThemeKey(tenant));
    return hasPending && !hasTheme;
  }

  async function checkUserAndSite() {
    try {
      const response = await get<User>("/v1/auth/me");

      if (!response.data) {
        navigate("/login");
        return;
      }

      setUser(response.data);

      if (!response.data.plan || response.data.plan === "free") {
        navigate("/pricing");
        return;
      }

      const tenant = getCurrentTenant();

      if (!tenant) {
        setShowSiteModal(true);
        return;
      }

      if (shouldOpenThemeModal(tenant)) {
        setShowThemeModal(true);
      }

      await loadPages();
    } catch (err: any) {
      const status = err?.status ?? err?.response?.status;
      if (status === 401) {
        navigate("/login");
        return;
      }
    } finally {
      setLoading(false);
    }
  }

  async function loadPages() {
    if (!getCurrentTenant()) return;

    setLoadingPages(true);
    try {
      const response = await get<Page[]>("/v1/pages");
      if (response.data) setPages(response.data);
    } catch (err) {
      console.error("Erro ao carregar páginas:", err);
    } finally {
      setLoadingPages(false);
    }
  }

  async function handleCreateSite(e: FormEvent) {
    e.preventDefault();
    setError("");
    setCreating(true);

    try {
      await post("/v1/sites", { name: siteName, slug: siteSlug });

      setCurrentTenant(siteSlug);
      setShowSiteModal(false);
      setSiteName("");
      setSiteSlug("");

      setIsOpen(true);
    } catch (err: any) {
      if (err?.message?.includes("plan") || err?.message?.includes("limit")) {
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

    const tenant = getCurrentTenant();
    if (tenant && !localStorage.getItem(themeKey(tenant))) {
      sessionStorage.setItem(pendingThemeKey(tenant), "1");
      setShowThemeModal(true);
      return;
    }

    await loadPages();
  }

  function handleViewPage(pageId: string) {
    navigate(`/pages/${pageId}`);
  }

  function handleEditPage(pageId: string) {
    navigate(`/pages/${pageId}/edit`);
  }

  function openDeleteModal(page: Page) {
    setPageToDelete(page);
    setShowDeleteModal(true);
  }

  function closeDeleteModal() {
    setShowDeleteModal(false);
    setPageToDelete(null);
  }

  if (loading) {
    return (
      <main className={styles.main}>
        <div className={styles.loading}>
          <p>Carregando...</p>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <div className={styles.header}>
        <div>
          <h1>Minhas Páginas 🎉</h1>
          <p className={styles.subtitle}>Crie um site e depois adicione páginas e temas.</p>
        </div>
        <div className={styles.headerActions}>
          <button className={styles.primaryBtn} onClick={() => setShowSiteModal(true)}>
            Criar site
          </button>
          <button
            className={styles.secondaryBtn}
            onClick={() => {
              if (!getCurrentTenant()) setShowSiteModal(true);
              else setIsOpen(true);
            }}
          >
            Criar página
          </button>
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
              if (!getCurrentTenant()) setShowSiteModal(true);
              else setIsOpen(true);
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
                <button className={styles.deleteButton} onClick={() => openDeleteModal(page)}>
                  Excluir
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showSiteModal && (
        <div className={styles.modalOverlay} onClick={() => { }}>
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
                <input
                  type="text"
                  placeholder="Ex:  minha-loja"
                  value={siteSlug}
                  onChange={(e) => setSiteSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ""))}
                  required
                  disabled={creating}
                />
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
            <CreatePageForm onClose={() => setIsOpen(false)} onSuccess={handlePageCreated} />
          </div>
        </div>
      )}

      {pageToDelete && (
        <DeletePageModal
          open={showDeleteModal}
          pageName={pageToDelete.name}
          pageId={pageToDelete.page_id}
          onClose={closeDeleteModal}
          onDeleted={async () => {
            await loadPages();
          }}
        />
      )}

      <ThemeInstallerModal
        open={showThemeModal}
        onClose={async () => {
          setShowThemeModal(false);
          await loadPages();
        }}
        onApplied={async () => {
          const tenant = getCurrentTenant();
          if (tenant) sessionStorage.removeItem(pendingThemeKey(tenant));
          await loadPages();
        }}
        onSkip={async () => {
          const tenant = getCurrentTenant();
          if (tenant) sessionStorage.removeItem(pendingThemeKey(tenant));
          await loadPages();
        }}
      />
    </main>
  );
}
