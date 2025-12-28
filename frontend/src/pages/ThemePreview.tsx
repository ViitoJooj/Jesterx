import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import styles from "../styles/pages/ThemeStore.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";
import { getThemeById } from "../config/themes";
import { getCurrentTenant, post } from "../utils/api";

export function ThemePreviewPage() {
  const { themeId } = useParams<{ themeId: string }>();
  const navigate = useNavigate();
  const [loadingAction, setLoadingAction] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const theme = themeId ? getThemeById(themeId) : undefined;
  const tenant = getCurrentTenant();

  if (!theme) {
    return (
      <main className={styles.main}>
        <div className={styles.header}>
          <div>
            <p className={styles.eyebrow}>Loja de temas</p>
            <h1>Tema não encontrado</h1>
            <p className={styles.subtitle}>Confira o link ou volte para a loja de temas para escolher outro layout.</p>
          </div>
          <Link to="/themes" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
            Voltar para loja
          </Link>
        </div>
      </main>
    );
  }

  const currentTheme = theme;

  async function install() {
    if (!tenant) {
      setError("Crie um site antes de instalar um tema. Vá em 'Minhas Páginas'.");
      navigate("/pages");
      return;
    }
    setMessage("");
    setError("");
    setLoadingAction(true);
    try {
      await post("/v1/themes/apply", { theme_id: currentTheme.id });
      setMessage("Tema instalado no site atual.");
    } catch (err: any) {
      setError(err?.message || "Não foi possível instalar o tema.");
    } finally {
      setLoadingAction(false);
    }
  }

  async function clonePage() {
    if (!tenant) {
      setError("É preciso ter um site ativo para clonar o tema.");
      navigate("/pages");
      return;
    }
    setMessage("");
    setError("");
    setLoadingAction(true);
    const slug = `${currentTheme.id}-${Date.now().toString().slice(-5)}`;
    try {
      await post("/v1/pages", {
        name: `Página ${currentTheme.name}`,
        page_type: currentTheme.pageType,
        template: currentTheme.template,
        page_id: slug,
      });
      setMessage("Página criada com o tema selecionado.");
      navigate(`/pages/${slug}`);
    } catch (err: any) {
      setError(err?.message || "Não foi possível criar a página com o tema.");
    } finally {
      setLoadingAction(false);
    }
  }

  return (
    <main className={styles.main}>
      <div className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Loja de temas</p>
          <h1>{currentTheme.name}</h1>
          <p className={styles.subtitle}>{currentTheme.description}</p>
          <div className={styles.tags}>
            {currentTheme.tags.map((tag) => (
              <span key={tag}>{tag}</span>
            ))}
          </div>
        </div>
        <div className={styles.headerActions}>
          <Link to="/themes" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
            Ver todos
          </Link>
          <button type="button" onClick={install} disabled={loadingAction} className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
            {loadingAction ? "Instalando..." : "Instalar tema"}
          </button>
        </div>
      </div>

      {(message || error) && (
        <div className={`${styles.alert} ${error ? styles.error : styles.success}`}>
          {error || message}
        </div>
      )}

      <div className={styles.previewArea}>
        <div className={styles.previewHeader}>
          <div>
            <strong>Preview em URL pública</strong>
            <p className={styles.subtitle}>Compartilhe este endereço com sua equipe para validar o tema.</p>
          </div>
          <div className={styles.previewActions}>
            <button type="button" onClick={clonePage} disabled={loadingAction} className={styles.primaryButton} style={{ background: currentTheme.accent }}>
              {loadingAction ? "Clonando..." : "Criar página com tema"}
            </button>
            <a href={`/themes/${currentTheme.id}`} className={styles.linkButton}>
              {typeof window !== "undefined" && window.location?.origin ? `${window.location.origin}/themes/${currentTheme.id}` : `/themes/${currentTheme.id}`}
            </a>
          </div>
        </div>
        <div className={styles.previewFrame}>
          <iframe title={currentTheme.name} sandbox="allow-same-origin allow-scripts allow-forms" srcDoc={currentTheme.previewHtml} />
        </div>
      </div>
    </main>
  );
}
