import { useMemo, useState, type DragEvent } from "react";
import styles from "../styles/components/ThemeInstallerModal.module.scss";
import { getCurrentTenant, post } from "../utils/api";
import { communityThemes } from "../config/themes";

type Theme = {
  id: string;
  name: string;
  description: string;
  accent: string;
};

type Props = {
  open: boolean;
  onClose: () => void;
  onApplied?: () => void;
  onSkip?: () => void;
};

export function ThemeInstallerModal({ open, onClose, onApplied, onSkip }: Props) {
  const tenant = getCurrentTenant();

  const themes = useMemo<Theme[]>(
    () =>
      communityThemes.map((theme) => ({
        id: theme.id,
        name: theme.name,
        description: theme.description,
        accent: theme.accent,
      })),
    []
  );

  const [hoveringDrop, setHoveringDrop] = useState(false);
  const [selectedThemeId, setSelectedThemeId] = useState<string | null>(null);
  const [applying, setApplying] = useState(false);
  const [error, setError] = useState("");

  if (!open) return null;

  async function applyTheme(themeId: string) {
    if (!tenant) return;

    setError("");
    setApplying(true);

    try {
      await post("/v1/themes/apply", { theme_id: themeId });
      localStorage.setItem(`theme:${tenant}`, themeId);
      onApplied?.();
      onClose();
    } catch (err: any) {
      setError(err?.message || "Não foi possível aplicar o tema.");
    } finally {
      setApplying(false);
    }
  }

  function onDragStart(e: DragEvent<HTMLButtonElement>, themeId: string) {
    e.dataTransfer.setData("text/plain", themeId);
    e.dataTransfer.effectAllowed = "copy";
  }

  function onDrop(e: DragEvent<HTMLDivElement>) {
    e.preventDefault();
    setHoveringDrop(false);
    const themeId = e.dataTransfer.getData("text/plain");
    if (!themeId) return;
    applyTheme(themeId);
  }

  return (
    <div className={styles.overlay} onClick={() => { }}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <h2>Escolha seu tema</h2>
          <p>Arraste um tema para o seu site</p>
        </div>

        {error && <div className={styles.errorBox}>{error}</div>}

        <div className={styles.installerRow}>
          <div className={styles.themesCol}>
            <div className={styles.colTitle}>Temas</div>

            <div className={styles.themeGrid}>
              {themes.map((t) => (
                <button
                  key={t.id}
                  type="button"
                  className={`${styles.themeIcon} ${selectedThemeId === t.id ? styles.selected : ""}`}
                  draggable
                  onDragStart={(e) => onDragStart(e, t.id)}
                  onClick={() => setSelectedThemeId(t.id)}
                  disabled={applying}
                >
                  <div className={styles.iconPreview} style={{ background: t.accent }} />
                  <div className={styles.iconMeta}>
                    <strong>{t.name}</strong>
                    <span>{t.description}</span>
                  </div>
                </button>
              ))}
            </div>

            <div className={styles.hint}>Clique para selecionar e depois aplicar, ou arraste e solte.</div>
          </div>

          <div className={styles.arrowCol} aria-hidden="true">
            <div className={styles.arrow}>➜</div>
          </div>

          <div className={styles.dropCol}>
            <div className={styles.colTitle}>Seu site</div>

            <div
              className={`${styles.dropZone} ${hoveringDrop ? styles.dropHover : ""}`}
              onDragOver={(e) => {
                e.preventDefault();
                setHoveringDrop(true);
              }}
              onDragLeave={() => setHoveringDrop(false)}
              onDrop={onDrop}
            >
              <div className={styles.siteIcon}>
                <div className={styles.siteDot} />
              </div>
              <div className={styles.siteText}>
                <strong>{tenant ? `${tenant}.jesterx.com` : "seu-site.jesterx.com"}</strong>
                <span>Solte aqui para aplicar o tema</span>
              </div>
            </div>

            <div className={styles.actions}>
              <button
                type="button"
                className={styles.secondary}
                onClick={() => {
                  onSkip?.();
                  onClose();
                }}
                disabled={applying}
              >
                Pular
              </button>

              <button
                type="button"
                className={styles.primary}
                onClick={() => selectedThemeId && applyTheme(selectedThemeId)}
                disabled={applying || !selectedThemeId}
              >
                {applying ? "Aplicando..." : "Aplicar tema"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
