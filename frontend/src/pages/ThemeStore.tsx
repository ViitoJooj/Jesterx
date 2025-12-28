import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import styles from "../styles/pages/ThemeStore.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";
import { communityThemes, type CommunityTheme } from "../config/themes";
import { getCurrentTenant, post } from "../utils/api";

export function ThemeStore() {
  const navigate = useNavigate();
  const [installing, setInstalling] = useState<string | null>(null);
  const [creating, setCreating] = useState<string | null>(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const tenant = getCurrentTenant();

  function makeSlug(prefix: string) {
    return `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`;
  }

  async function handleInstall(themeId: string) {
    if (!tenant) {
      setError("Crie um site antes de instalar um tema. Vá em 'Minhas Páginas'.");
      navigate("/pages");
      return;
    }

    setMessage("");
    setError("");
    setInstalling(themeId);
    try {
      await post("/v1/themes/apply", { theme_id: themeId });
      setMessage("Tema instalado no site atual. Crie páginas para visualizar o tema.");
    } catch (err: any) {
      setError(err?.message || "Não foi possível instalar o tema.");
    } finally {
      setInstalling(null);
    }
  }

  async function handleCreatePage(theme: CommunityTheme) {
    if (!tenant) {
      setError("É preciso ter um site ativo para clonar o tema.");
      navigate("/pages");
      return;
    }

    setMessage("");
    setError("");
    setCreating(theme.id);
    const slug = makeSlug(theme.id);

    try {
      await post("/v1/pages", {
        name: `Página ${theme.name}`,
        page_type: theme.pageType,
        template: theme.template,
        page_id: slug,
      });
      setMessage("Página criada com o tema selecionado.");
      navigate(`/pages/${slug}`);
    } catch (err: any) {
      setError(err?.message || "Não foi possível criar a página com o tema.");
    } finally {
      setCreating(null);
    }
  }

  async function handleCopy(theme: CommunityTheme) {
    setMessage("");
    setError("");
    try {
      await navigator.clipboard.writeText(theme.template);
      setMessage("Código do tema copiado. Cole no seu editor ou na criação de página.");
    } catch {
      setError("Não foi possível copiar o tema, copie manualmente no preview.");
    }
  }

  return (
    <main className={styles.main}>
      <div className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Loja de temas</p>
          <h1>Instale temas criados pela comunidade</h1>
          <p className={styles.subtitle}>Escolha um tema, instale no site atual e crie quantas páginas quiser. Cada tema tem preview acessível por URL.</p>
        </div>
        <div className={styles.headerActions}>
          <Link to="/pages" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
            Minhas Páginas
          </Link>
          <Link to="/pricing" className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
            Planos
          </Link>
        </div>
      </div>

      {(message || error) && (
        <div className={`${styles.alert} ${error ? styles.error : styles.success}`}>
          {error || message}
        </div>
      )}

      <div className={styles.grid}>
        {communityThemes.map((theme) => (
          <article key={theme.id} className={styles.card} style={{ borderColor: theme.accent }}>
            <div className={styles.preview}>
              <iframe title={theme.name} sandbox="allow-same-origin allow-scripts allow-forms allow-popups" srcDoc={theme.previewHtml} />
              <span className={styles.tag}>{theme.author}</span>
            </div>

            <div className={styles.meta}>
              <div>
                <h2>{theme.name}</h2>
                <p>{theme.description}</p>
              </div>
              <div className={styles.tags}>
                {theme.tags.map((tag) => (
                  <span key={tag}>{tag}</span>
                ))}
              </div>
            </div>

            <div className={styles.actions}>
              <Link to={`/themes/${theme.id}`} className={styles.linkButton}>
                Abrir URL do tema
              </Link>
              <button type="button" onClick={() => handleInstall(theme.id)} disabled={installing === theme.id} className={styles.installButton}>
                {installing === theme.id ? "Instalando..." : "Instalar no site"}
              </button>
              <button type="button" onClick={() => handleCreatePage(theme)} disabled={creating === theme.id} className={styles.primaryButton} style={{ background: theme.accent }}>
                {creating === theme.id ? "Clonando..." : "Criar página com tema"}
              </button>
            </div>

            <div className={styles.secondaryActions}>
              <button type="button" onClick={() => handleCopy(theme)} className={styles.copyButton}>
                Copiar código (Svelte/HTML)
              </button>
              <span className={styles.helper}>URL pública: /themes/{theme.id}</span>
            </div>
          </article>
        ))}
      </div>
    </main>
  );
}
