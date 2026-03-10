import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Button from "../../components/button/Button";
import { useAuthContext } from "../../hooks/AuthContext";
import { apiFetch, API_URL } from "../../hooks/api";
import styles from "./Pages.module.scss";

type WebsiteType = "ECOMMERCE" | "LANDING_PAGE" | "SOFTWARE_SELL" | "COURSE" | "VIDEO";
type EditorMode = "ELEMENTOR" | "REACT" | "SVELTE";

type CustomRoute = {
  id: string;
  path: string;
  title: string;
  requiresAuth: boolean;
};

type CreatedWebsite = {
  id: string;
  type: string;
  name: string;
  short_description?: string;
  description?: string;
  created_at?: string;
};

type WebsiteResponse = {
  success: boolean;
  message: string;
  data: CreatedWebsite;
};

type WebsitesResponse = {
  success: boolean;
  message: string;
  data: CreatedWebsite[];
};

type VersionResponse = {
  success: boolean;
  message: string;
  data: {
    id: string;
    version: number;
    scan_status: "clean" | "warning" | "blocked";
    scan_score: number;
    source_type: string;
  };
};

const TYPE_META: Record<WebsiteType, { label: string; icon: string; color: string; description: string; starterRoutes: string[] }> = {
  ECOMMERCE:     { label: "E-commerce",       icon: "🛍️", color: "#ff6029", description: "Catálogo, carrinho e checkout.",       starterRoutes: ["/", "/produtos", "/produto/:slug", "/carrinho", "/checkout"] },
  LANDING_PAGE:  { label: "Landing Page",     icon: "🎯", color: "#2c7ef5", description: "Alta conversão com CTA e formulário.", starterRoutes: ["/", "/obrigado"] },
  SOFTWARE_SELL: { label: "Venda de Software", icon: "💻", color: "#7c3aed", description: "Trial, planos e download seguro.",     starterRoutes: ["/", "/precos", "/download", "/faq"] },
  COURSE:        { label: "Curso",             icon: "🎓", color: "#059669", description: "Vendas, módulos e aulas.",             starterRoutes: ["/", "/inscricao", "/modulos", "/aula/:id"] },
  VIDEO:         { label: "Canal de Vídeo",   icon: "📹", color: "#e11d48", description: "Vitrine de vídeos e séries.",          starterRoutes: ["/", "/videos", "/video/:slug"] },
};

const EDITOR_OPTIONS: { value: EditorMode; label: string; desc: string }[] = [
  { value: "ELEMENTOR", label: "Elementor (Visual)", desc: "Drag-and-drop, sem código" },
  { value: "REACT",     label: "React",               desc: "Componentes e hooks" },
  { value: "SVELTE",    label: "Svelte",               desc: "Sintaxe enxuta e rápida" },
];

function normalizeRoutePath(path: string) {
  const trimmed = path.trim();
  if (!trimmed) return "";
  return trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
}

function toSourceType(editor: EditorMode): "ELEMENTOR_JSON" | "REACT" | "SVELTE" {
  if (editor === "ELEMENTOR") return "ELEMENTOR_JSON";
  if (editor === "REACT") return "REACT";
  return "SVELTE";
}

function buildSource(editor: EditorMode, name: string, desc: string, routes: CustomRoute[]) {
  if (editor === "ELEMENTOR") {
    return JSON.stringify({
      name,
      description: desc,
      blocks: [
        { type: "hero", title: name },
        { type: "text", content: desc || "Nova pagina criada com Jesterx Builder." },
      ],
      routes: routes.map((r) => ({ path: normalizeRoutePath(r.path), title: r.title.trim(), private: r.requiresAuth })),
    }, null, 2);
  }
  if (editor === "REACT") {
    return JSON.stringify({ component: `function App(){return (<main><h1>${name}</h1><p>${desc || "Projeto criado no builder."}</p></main>)};`, css: "main{padding:24px}" });
  }
  return JSON.stringify({ component: `<script>\n  const title = "${name}";\n</script>\n<main><h1>{title}</h1><p>${desc}</p></main>`, css: "main{padding:24px}" });
}

function getPlanRouteLimit(plan?: string) {
  const p = (plan ?? "").toLowerCase();
  if (!p) return 0;
  if (p.includes("enterprise") || p.includes("ultra")) return 100;
  if (p.includes("pro") || p.includes("business")) return 30;
  if (p.includes("starter") || p.includes("basic") || p.includes("essencial")) return 8;
  return 15;
}

export const Pages: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, websiteId, loading, me } = useAuthContext();

  const [type, setType] = useState<WebsiteType>("LANDING_PAGE");
  const [editor, setEditor] = useState<EditorMode>("ELEMENTOR");
  const [name, setName] = useState("");
  const [shortDesc, setShortDesc] = useState("");
  const [desc, setDesc] = useState("");
  const [routes, setRoutes] = useState<CustomRoute[]>([{ id: "route-1", path: "/", title: "Home", requiresAuth: false }]);
  const [routeSeed, setRouteSeed] = useState(2);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [created, setCreated] = useState<CreatedWebsite | null>(null);
  const [websites, setWebsites] = useState<CreatedWebsite[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [deleting, setDeleting] = useState<string | null>(null);

  const activePlan = (me?.user_plan ?? "").trim();
  const hasActivePlan = activePlan.length > 0;
  const routeLimit = useMemo(() => getPlanRouteLimit(me?.user_plan), [me?.user_plan]);

  const canSubmit = useMemo(() => (
    isAuthenticated && hasActivePlan && name.trim().length >= 3 && routes.length > 0 && routes.length <= routeLimit && !submitting
  ), [hasActivePlan, isAuthenticated, name, routeLimit, routes.length, submitting]);

  function applyTemplate(nextType: WebsiteType) {
    const meta = TYPE_META[nextType];
    setType(nextType);
    const seeded = meta.starterRoutes.slice(0, Math.max(routeLimit, 1)).map((path, i) => ({
      id: `route-${i + 1}`, path, title: i === 0 ? "Home" : `Rota ${i + 1}`, requiresAuth: false,
    }));
    setRoutes(seeded.length ? seeded : [{ id: "route-1", path: "/", title: "Home", requiresAuth: false }]);
    setRouteSeed(Math.max(seeded.length + 1, 2));
  }

  function handleAddRoute() {
    if (routes.length >= routeLimit) return;
    setRoutes((prev) => [...prev, { id: `route-${routeSeed}`, path: `/rota-${routeSeed}`, title: `Rota ${routeSeed}`, requiresAuth: false }]);
    setRouteSeed((prev) => prev + 1);
  }

  function handleUpdateRoute(id: string, field: "path" | "title" | "requiresAuth", value: string | boolean) {
    setRoutes((prev) => prev.map((r) => r.id === id ? { ...r, [field]: value } : r));
  }

  function handleRemoveRoute(id: string) {
    setRoutes((prev) => prev.filter((r) => r.id !== id));
  }

  const [siteEditorMap, setSiteEditorMap] = useState<Record<string, string>>({});

  async function detectEditorTypes(siteIds: string[]) {
    const results = await Promise.allSettled(
      siteIds.map((id) =>
        apiFetch<{ success: boolean; data: { source_type: string }[] }>(
          `/api/v1/sites/${id}/versions`, { method: "GET", websiteId }
        ).then((r) => ({ id, type: r.data[0]?.source_type ?? "" }))
      )
    );
    const map: Record<string, string> = {};
    results.forEach((r) => { if (r.status === "fulfilled") map[r.value.id] = r.value.type; });
    setSiteEditorMap((prev) => ({ ...prev, ...map }));
  }

  async function loadWebsites() {
    if (!isAuthenticated) return;
    try {
      const resp = await apiFetch<WebsitesResponse>("/api/v1/websites", { method: "GET", websiteId });
      const list = resp.data ?? [];
      setWebsites(list);
      if (list.length > 0) detectEditorTypes(list.map((s) => s.id));
    } catch { /* silencioso */ }
  }

  useEffect(() => { loadWebsites(); }, [isAuthenticated, websiteId]); // eslint-disable-line react-hooks/exhaustive-deps

  async function handleCreate(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!canSubmit) return;
    setError(null);
    setSubmitting(true);
    try {
      const invalid = routes.find((r) => {
        const p = normalizeRoutePath(r.path);
        return p.length < 1 || r.title.trim().length < 2 || p.includes(" ");
      });
      if (invalid) throw new Error("Revise as rotas: caminho precisa começar com '/' e não pode ter espaços.");

      const resp = await apiFetch<WebsiteResponse>("/api/v1/websites", {
        method: "POST", websiteId,
        body: JSON.stringify({ type, name: name.trim(), short_description: shortDesc.trim(), description: desc.trim() }),
      });

      await apiFetch(`/api/v1/sites/${resp.data.id}/routes`, {
        method: "POST", websiteId,
        body: JSON.stringify({ routes: routes.map((r) => ({ path: normalizeRoutePath(r.path), title: r.title.trim(), requires_auth: r.requiresAuth })) }),
      });

      const vResp = await apiFetch<VersionResponse>(`/api/v1/sites/${resp.data.id}/versions`, {
        method: "POST", websiteId,
        body: JSON.stringify({ source_type: toSourceType(editor), source: buildSource(editor, name.trim(), desc.trim(), routes) }),
      });

      if (vResp.data.scan_status !== "blocked") {
        await apiFetch(`/api/v1/sites/${resp.data.id}/publish/${vResp.data.version}`, { method: "POST", websiteId });
      }

      setCreated(resp.data);
      setName(""); setShortDesc(""); setDesc("");
      applyTemplate(type);
      setShowCreate(false);
      await loadWebsites();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao criar página");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id: string) {
    if (!window.confirm("Tem certeza que deseja excluir este website? Esta ação não pode ser desfeita.")) return;
    setDeleting(id);
    try {
      await apiFetch(`/api/v1/sites/${id}`, { method: "DELETE", websiteId });
      await loadWebsites();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao excluir website");
    } finally {
      setDeleting(null);
    }
  }

  function openEditor(siteId: string, mode: EditorMode) {
    if (mode === "ELEMENTOR") navigate(`/pages/${siteId}/elementor`);
    else if (mode === "REACT") navigate(`/pages/${siteId}/react`);
    else navigate(`/pages/${siteId}/svelte`);
  }

  if (!loading && !isAuthenticated) {
    return (
      <main className={styles.main}>
        <section className={styles.centerBox}>
          <h1>Área de páginas</h1>
          <p>Faça login para criar e gerenciar suas páginas.</p>
          <div className={styles.actions}>
            <Button to="/login" variant="primary">Entrar</Button>
            <Button to="/register" variant="secondary">Criar conta</Button>
          </div>
        </section>
      </main>
    );
  }

  const websitesByType = useMemo(() => {
    const map: Partial<Record<WebsiteType, CreatedWebsite[]>> = {};
    websites.forEach((site) => {
      const key = site.type as WebsiteType;
      if (!map[key]) map[key] = [];
      map[key]!.push(site);
    });
    return map;
  }, [websites]);

  const orderedTypes = (Object.keys(TYPE_META) as WebsiteType[]).filter((t) => (websitesByType[t]?.length ?? 0) > 0);

  return (
    <main className={styles.main}>
      <div className={styles.topBar}>
        <div>
          <h1>Minhas Páginas</h1>
          <p>{hasActivePlan ? `Plano: ${activePlan} · ${routeLimit} rotas por site` : "Sem plano ativo"}</p>
        </div>
        <div className={styles.topActions}>
          {!hasActivePlan && <Button type="button" variant="secondary" onClick={() => navigate("/plans")}>Ativar plano</Button>}
          <Button type="button" variant="primary" onClick={() => setShowCreate((v) => !v)}>
            {showCreate ? "✕ Fechar" : "+ Novo Site"}
          </Button>
        </div>
      </div>

      {error && <p className={styles.error}>{error}</p>}
      {created && (
        <div className={styles.successBanner}>
          ✅ Site <strong>{created.name}</strong> criado com sucesso!{" "}
          <button type="button" onClick={() => setCreated(null)}>✕</button>
        </div>
      )}

      {showCreate && (
        <section className={styles.createForm}>
          <h2>Criar novo site</h2>

          <div className={styles.typeGrid}>
            {(Object.keys(TYPE_META) as WebsiteType[]).map((t) => {
              const meta = TYPE_META[t];
              return (
                <button key={t} type="button"
                  className={`${styles.typeCard} ${type === t ? styles.typeCardActive : ""}`}
                  style={{ "--type-color": meta.color } as React.CSSProperties}
                  onClick={() => applyTemplate(t)}
                >
                  <span className={styles.typeIcon}>{meta.icon}</span>
                  <strong>{meta.label}</strong>
                  <span>{meta.description}</span>
                </button>
              );
            })}
          </div>

          <form onSubmit={handleCreate} className={styles.form} noValidate>
            <div className={styles.fieldGroup}>
              <label htmlFor="websiteName">Nome do site *</label>
              <input
                id="websiteName" value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Ex: Loja Aurora Fitness"
                minLength={3} required
                className={styles.textInput}
              />
            </div>

            <div className={styles.fieldGroup}>
              <label htmlFor="shortDesc">Descrição curta</label>
              <input id="shortDesc" value={shortDesc} onChange={(e) => setShortDesc(e.target.value)} placeholder="Uma frase objetiva" className={styles.textInput} />
            </div>

            <div className={styles.fieldGroup}>
              <label htmlFor="desc">Descrição completa</label>
              <textarea id="desc" value={desc} onChange={(e) => setDesc(e.target.value)} placeholder="Contexto, público e diferenciais." rows={3} />
            </div>

            <div className={styles.fieldGroup}>
              <label>Editor</label>
              <div className={styles.editorRow}>
                {EDITOR_OPTIONS.map((opt) => (
                  <button key={opt.value} type="button"
                    className={`${styles.editorChip} ${editor === opt.value ? styles.editorChipActive : ""}`}
                    onClick={() => setEditor(opt.value)}
                  >
                    <strong>{opt.label}</strong>
                    <span>{opt.desc}</span>
                  </button>
                ))}
              </div>
            </div>

            <div className={styles.routeSection}>
              <div className={styles.routeHeader}>
                <h3>Rotas do projeto</h3>
                <span>{routes.length}/{routeLimit} usadas</span>
              </div>
              <div className={styles.routeList}>
                {routes.map((route) => (
                  <article key={route.id} className={styles.routeItem}>
                    <input value={route.path} onChange={(e) => handleUpdateRoute(route.id, "path", e.target.value)} placeholder="/minha-rota" className={styles.textInput} />
                    <input value={route.title} onChange={(e) => handleUpdateRoute(route.id, "title", e.target.value)} placeholder="Título" className={styles.textInput} />
                    <label className={styles.routeCheck}>
                      <input type="checkbox" checked={route.requiresAuth} onChange={(e) => handleUpdateRoute(route.id, "requiresAuth", e.target.checked)} />
                      🔒 Privada
                    </label>
                    <button type="button" className={styles.removeRoute} onClick={() => handleRemoveRoute(route.id)} disabled={routes.length === 1}>✕</button>
                  </article>
                ))}
              </div>
              <Button type="button" variant="secondary" onClick={handleAddRoute} disabled={routes.length >= routeLimit}>+ Adicionar rota</Button>
            </div>

            {error && <p className={styles.error}>{error}</p>}

            <div className={styles.actions}>
              <Button type="submit" variant="primary" disabled={!canSubmit}>{submitting ? "Criando..." : "✓ Criar site"}</Button>
            </div>
          </form>
        </section>
      )}

      {!websites.length && !showCreate && (
        <div className={styles.empty}>
          <p>Nenhum site criado ainda. Clique em <strong>+ Novo Site</strong> para começar.</p>
        </div>
      )}

      {orderedTypes.map((t) => {
        const meta = TYPE_META[t];
        const list = websitesByType[t] ?? [];
        return (
          <section key={t} className={styles.typeSection}>
            <div className={styles.typeSectionHeader}>
              <span className={styles.typeSectionIcon} style={{ background: meta.color }}>{meta.icon}</span>
              <h2>{meta.label}</h2>
              <span className={styles.typeCount}>{list.length} site{list.length !== 1 ? "s" : ""}</span>
            </div>
            <div className={styles.siteGrid}>
              {list.map((site) => (
                <article key={site.id} className={styles.siteCard}>
                  <div className={styles.siteCardTop} style={{ borderColor: meta.color }}>
                    <div className={styles.siteInfo}>
                      <Link to={`/store/${site.id}`} className={styles.siteName}>{site.name}</Link>
                      {site.short_description && <span>{site.short_description}</span>}
                    </div>
                    <button
                      type="button"
                      className={styles.deleteBtn}
                      disabled={deleting === site.id}
                      onClick={() => handleDelete(site.id)}
                      title="Excluir website"
                    >
                      {deleting === site.id ? "..." : "🗑"}
                    </button>
                  </div>
                  <div className={styles.siteActions}>
                    {(() => {
                      const srcType = siteEditorMap[site.id];
                      const isElementor = !srcType || srcType === "ELEMENTOR_JSON";
                      const isReact = srcType === "REACT";
                      const isSvelte = srcType === "SVELTE";
                      return (
                        <>
                          {(isElementor || !srcType) && (
                            <button type="button" className={styles.editorBtn} onClick={() => openEditor(site.id, "ELEMENTOR")}>
                              🎨 Visual
                            </button>
                          )}
                          {isReact && (
                            <button type="button" className={styles.editorBtn} onClick={() => openEditor(site.id, "REACT")}>
                              ⚛ React
                            </button>
                          )}
                          {isSvelte && (
                            <button type="button" className={styles.editorBtn} onClick={() => openEditor(site.id, "SVELTE")}>
                              🔥 Svelte
                            </button>
                          )}
                          {!srcType && (
                            <>
                              <button type="button" className={`${styles.editorBtn} ${styles.editorBtnDim}`} onClick={() => openEditor(site.id, "REACT")}>⚛ React</button>
                              <button type="button" className={`${styles.editorBtn} ${styles.editorBtnDim}`} onClick={() => openEditor(site.id, "SVELTE")}>🔥 Svelte</button>
                            </>
                          )}
                        </>
                      );
                    })()}
                    <a href={`${API_URL}/p/${site.id}`} target="_blank" rel="noreferrer" className={styles.openBtn}>
                      ↗ Abrir
                    </a>
                  </div>
                </article>
              ))}
            </div>
          </section>
        );
      })}
    </main>
  );
};