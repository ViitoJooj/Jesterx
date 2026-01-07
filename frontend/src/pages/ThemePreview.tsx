import { useLocation, useNavigate } from "react-router-dom";
import { useMemo, useState } from "react";
import { post } from "../utils/api";

type ThemePreviewState = {
  id: string;
  name: string;
  pageType: string;
  template: string;
  description?: string;
};

function makeSlug(value: string): string {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export function ThemePreview() {
  const location = useLocation();
  const navigate = useNavigate();
  const theme = useMemo(() => location.state as ThemePreviewState | undefined, [location.state]);

  const [loadingAction, setLoadingAction] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function applyTheme() {
    if (!theme) {
      setError("Tema não encontrado.");
      return;
    }

    setError("");
    setMessage("");
    setLoadingAction(true);
    try {
      await post("/v1/themes/apply", { theme_id: theme.id });
      setMessage("Tema instalado no site atual.");
    } catch (err: any) {
      setError(err?.message || "Não foi possível instalar o tema.");
    } finally {
      setLoadingAction(false);
    }
  }

  async function createPageFromTheme() {
    if (!theme) {
      setError("Tema não encontrado.");
      return;
    }

    setError("");
    setMessage("");
    setLoadingAction(true);

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
      setError(err?.message || "Não foi possível criar a página.");
    } finally {
      setLoadingAction(false);
    }
  }

  if (!theme) {
    return (
      <main style={{ padding: "32px" }}>
        <p>Não encontramos dados do tema selecionado.</p>
      </main>
    );
  }

  return (
    <main style={{ padding: "32px", maxWidth: 960, margin: "0 auto", display: "grid", gap: 16 }}>
      <header>
        <p style={{ margin: "0 0 4px", color: "#6b7280", textTransform: "uppercase", fontSize: 12 }}>Pré-visualização</p>
        <h1 style={{ margin: 0 }}>{theme.name}</h1>
        {theme.description && <p style={{ marginTop: 6, color: "#4b5563" }}>{theme.description}</p>}
      </header>

      {message && <div style={{ padding: 12, borderRadius: 8, background: "#ecfdf3", color: "#166534" }}>{message}</div>}
      {error && <div style={{ padding: 12, borderRadius: 8, background: "#fef2f2", color: "#b91c1c" }}>{error}</div>}

      <section
        style={{
          padding: 16,
          border: "1px solid #e5e7eb",
          borderRadius: 12,
          display: "flex",
          flexDirection: "column",
          gap: 12,
        }}
      >
        <div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
          <button
            type="button"
            onClick={applyTheme}
            disabled={loadingAction}
            style={{
              padding: "12px 16px",
              borderRadius: 10,
              border: "none",
              background: "#111827",
              color: "#fff",
              cursor: loadingAction ? "not-allowed" : "pointer",
            }}
          >
            {loadingAction ? "Processando..." : "Aplicar tema no site"}
          </button>
          <button
            type="button"
            onClick={createPageFromTheme}
            disabled={loadingAction}
            style={{
              padding: "12px 16px",
              borderRadius: 10,
              border: "1px solid #e5e7eb",
              background: "#fff",
              color: "#111827",
              cursor: loadingAction ? "not-allowed" : "pointer",
            }}
          >
            {loadingAction ? "Processando..." : "Criar página com tema"}
          </button>
        </div>

        <div style={{ padding: 16, borderRadius: 10, background: "#f9fafb", color: "#374151", lineHeight: 1.5 }}>
          <p style={{ margin: 0 }}>
            Este tema usa o template <strong>{theme.pageType}</strong>. Ao criar uma página, o código será clonado para sua conta
            e poderá ser editado livremente no editor de páginas.
          </p>
        </div>
      </section>
    </main>
  );
}
