import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Button from "../../components/button/Button";
import { useAuthContext } from "../../hooks/AuthContext";
import { apiFetch, API_URL } from "../../hooks/api";
import styles from "./Pages.module.scss";

const TYPE_META = {
  ECOMMERCE:     { label: "E-commerce",        icon: "🛍️", color: "var(--jx-color-primary)", description: "Catálogo, carrinho e checkout completo.", starter: ["/", "/produtos", "/produto/:slug", "/carrinho", "/checkout"] },
  LANDING_PAGE:  { label: "Landing Page",      icon: "🎯", color: "#2c7ef5", description: "Alta conversão com CTA e formulário.",    starter: ["/", "/obrigado"] },
  SOFTWARE_SELL: { label: "Venda de Software", icon: "💻", color: "#7c3aed", description: "Trial, planos e download seguro.",         starter: ["/", "/precos", "/download", "/faq"] },
  COURSE:        { label: "Curso",             icon: "🎓", color: "#059669", description: "Vendas, módulos e aulas.",                 starter: ["/", "/inscricao", "/modulos", "/aula/:id"] },
  VIDEO:         { label: "Canal de Vídeo",    icon: "📹", color: "#e11d48", description: "Vitrine de vídeos e séries.",              starter: ["/", "/videos", "/video/:slug"] },
};

const EDITOR_OPTIONS = [
  { value: "ELEMENTOR", label: "Visual", icon: "🎨", desc: "Drag-and-drop, sem código" },
  { value: "REACT",     label: "React",  icon: "⚛️", desc: "Componentes e hooks" },
  { value: "SVELTE",    label: "Svelte", icon: "🔥", desc: "Sintaxe enxuta" },
];

const EDITOR_LABELS = {
  ELEMENTOR_JSON: { label: "Visual", color: "#e44c65" },
  REACT:          { label: "React",  color: "#61dafb" },
  SVELTE:         { label: "Svelte", color: "var(--jx-color-primary)" },
};

function normPath(path) {
  const t = path.trim();
  if (!t) return "";
  return t.startsWith("/") ? t : `/${t}`;
}

function toSourceType(editor) {
  return editor === "ELEMENTOR" ? "ELEMENTOR_JSON" : editor;
}

function buildSource(editor, name, desc, routes) {
  if (editor === "ELEMENTOR") {
    return JSON.stringify({
      name, description: desc,
      blocks: [{ type: "hero", title: name }, { type: "text", content: desc || "Nova página criada com Jesterx." }],
      routes: routes.map(r => ({ path: normPath(r.path), title: r.title.trim(), private: r.requiresAuth })),
    }, null, 2);
  }
  if (editor === "REACT") {
    return JSON.stringify({ component: `function App(){return(<main><h1>${name}</h1><p>${desc||"Projeto criado."}</p></main>)};`, css: "main{padding:24px}" });
  }
  return JSON.stringify({ component: `<script>\n  const title="${name}";\n</script>\n<main><h1>{title}</h1><p>${desc}</p></main>`, css: "main{padding:24px}" });
}

export const Pages = () => {
  const navigate = useNavigate();
  const { isAuthenticated, websiteId, loading, me, loadMe } = useAuthContext();

  const [websites, setWebsites]       = useState([]);
  const [editorMap, setEditorMap]     = useState({});
  const [deleting, setDeleting]       = useState(null);
  const [listError, setListError]     = useState(null);
  const [created, setCreated]         = useState(null);
  const [deleteTarget, setDeleteTarget] = useState(null); // { id, name }

  // Modal
  const [showModal, setShowModal] = useState(false);
  const [step, setStep]           = useState(1);
  const [type, setType]           = useState("LANDING_PAGE");
  const [editor, setEditor]       = useState("ELEMENTOR");
  const [name, setName]           = useState("");
  const [shortDesc, setShortDesc] = useState("");
  const [routes, setRoutes]       = useState([{ id: "r1", path: "/", title: "Home", requiresAuth: false }]);
  const [routeSeed, setRouteSeed] = useState(2);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError]   = useState(null);

  const activePlan    = (me?.user_plan ?? "").trim();
  const hasActivePlan = activePlan.length > 0;
  const maxRoutes     = me?.plan_max_routes ?? 0;
  const maxSites      = me?.plan_max_sites  ?? 0;

  // Always fetch fresh me so plan info is current
  useEffect(() => {
    loadMe().catch(() => {});
  }, []); // eslint-disable-line

  async function loadWebsites() {
    if (!isAuthenticated) return;
    try {
      const resp = await apiFetch("/api/v1/websites", { method: "GET", websiteId });
      const list = resp.data ?? [];
      setWebsites(list);
      if (!list.length) return;
      const results = await Promise.allSettled(
        list.map(s =>
          apiFetch(`/api/v1/sites/${s.id}/versions`, { method: "GET", websiteId })
            .then(r => ({ id: s.id, type: r.data?.[0]?.source_type ?? "" }))
        )
      );
      const map = {};
      results.forEach(r => { if (r.status === "fulfilled") map[r.value.id] = r.value.type; });
      setEditorMap(map);
    } catch { /* silent */ }
  }

  useEffect(() => { loadWebsites(); }, [isAuthenticated, websiteId]); // eslint-disable-line

  function applyStarter(nextType) {
    const starter = (TYPE_META[nextType]?.starter ?? ["/"]).slice(0, Math.max(maxRoutes, 1));
    const seeded = starter.map((path, i) => ({ id: `r${i + 1}`, path, title: i === 0 ? "Home" : `Rota ${i + 1}`, requiresAuth: false }));
    setRoutes(seeded.length ? seeded : [{ id: "r1", path: "/", title: "Home", requiresAuth: false }]);
    setRouteSeed(seeded.length + 1);
  }

  function openModal() {
    setStep(1); setType("LANDING_PAGE"); setEditor("ELEMENTOR");
    setName(""); setShortDesc(""); setFormError(null);
    applyStarter("LANDING_PAGE");
    setShowModal(true);
  }

  function closeModal() { if (!submitting) setShowModal(false); }

  function selectType(t) { setType(t); applyStarter(t); }

  function addRoute() {
    if (routes.length >= maxRoutes) return;
    setRoutes(prev => [...prev, { id: `r${routeSeed}`, path: `/rota-${routeSeed}`, title: `Rota ${routeSeed}`, requiresAuth: false }]);
    setRouteSeed(n => n + 1);
  }

  function updateRoute(id, field, value) {
    setRoutes(prev => prev.map(r => r.id === id ? { ...r, [field]: value } : r));
  }

  function removeRoute(id) { setRoutes(prev => prev.filter(r => r.id !== id)); }

  async function handleCreate() {
    if (submitting) return;
    setFormError(null);
    const invalid = routes.find(r => { const p = normPath(r.path); return !p || r.title.trim().length < 2 || p.includes(" "); });
    if (invalid) { setFormError("Revise as rotas: caminho deve começar com '/' e sem espaços."); return; }
    setSubmitting(true);
    try {
      const resp = await apiFetch("/api/v1/websites", {
        method: "POST", websiteId,
        body: JSON.stringify({ type, name: name.trim(), short_description: shortDesc.trim() }),
      });
      await apiFetch(`/api/v1/sites/${resp.data.id}/routes`, {
        method: "POST", websiteId,
        body: JSON.stringify({ routes: routes.map(r => ({ path: normPath(r.path), title: r.title.trim(), requires_auth: r.requiresAuth })) }),
      });
      const vResp = await apiFetch(`/api/v1/sites/${resp.data.id}/versions`, {
        method: "POST", websiteId,
        body: JSON.stringify({ source_type: toSourceType(editor), source: buildSource(editor, name.trim(), shortDesc.trim(), routes) }),
      });
      if (vResp.data.scan_status !== "blocked") {
        await apiFetch(`/api/v1/sites/${resp.data.id}/publish/${vResp.data.version}`, { method: "POST", websiteId });
      }
      setCreated(resp.data);
      setShowModal(false);
      await loadWebsites();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "Falha ao criar site");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete() {
    if (!deleteTarget) return;
    const id = deleteTarget.id;
    setDeleting(id);
    setDeleteTarget(null);
    try {
      await apiFetch(`/api/v1/sites/${id}`, { method: "DELETE", websiteId });
      await loadWebsites();
    } catch (err) {
      setListError(err instanceof Error ? err.message : "Erro ao excluir");
    } finally {
      setDeleting(null);
    }
  }

  function openEditor(siteId, srcType) {
    if (srcType === "REACT") navigate(`/pages/${siteId}/react`);
    else if (srcType === "SVELTE") navigate(`/pages/${siteId}/svelte`);
    else navigate(`/pages/${siteId}/elementor`);
  }

  const canStep2 = name.trim().length >= 3;
  const canCreate = canStep2 && routes.length > 0 && routes.length <= maxRoutes && !submitting;

  if (!loading && !isAuthenticated) {
    return (
      <main className={styles.main}>
        <div className={styles.inner}>
          <div className={styles.centerBox}>
            <h1>Minhas Páginas</h1>
            <p>Faça login para criar e gerenciar seus sites.</p>
            <div className={styles.centerActions}>
              <Button to="/login" variant="primary">Entrar</Button>
              <Button to="/register" variant="secondary">Criar conta</Button>
            </div>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <div className={styles.inner}>

        <div className={styles.topBar}>
          <div>
            <h1>Minhas Páginas</h1>
            <p>
              {hasActivePlan
                ? `Plano ${activePlan} · ${websites.length}/${maxSites} sites · ${maxRoutes} rotas por site`
                : "Nenhum plano ativo"}
            </p>
          </div>
          <div className={styles.topActions}>
            {!hasActivePlan && <Button variant="secondary" to="/plans">Ver planos</Button>}
            <Button variant="primary" onClick={openModal} disabled={!hasActivePlan}>+ Criar site</Button>
          </div>
        </div>

        {listError && <p className={styles.alert}>{listError}</p>}

        {created && (
          <div className={styles.successBanner}>
            ✅ <strong>{created.name}</strong> criado com sucesso!
            <button onClick={() => setCreated(null)}>✕</button>
          </div>
        )}

        {websites.length === 0 ? (
          <div className={styles.empty}>
            <div className={styles.emptyIcon}>🌐</div>
            <h2>Nenhum site criado ainda</h2>
            <p>
              {hasActivePlan
                ? 'Clique em "+ Criar site" para começar.'
                : "Assine um plano para criar seu primeiro site."}
            </p>
            {!hasActivePlan && <Button variant="primary" to="/plans">Ver planos</Button>}
          </div>
        ) : (
          <div className={styles.siteGrid}>
            {websites.map(site => {
              const meta = TYPE_META[site.type];
              const srcType = editorMap[site.id];
              const edLabel = EDITOR_LABELS[srcType] ?? EDITOR_LABELS["ELEMENTOR_JSON"];
              const cardColor = meta?.color ?? "var(--jx-color-primary)";
              return (
                <article key={site.id} className={styles.siteCard} style={{ "--card-color": cardColor }}>
                  <div className={styles.cardBanner}>
                    <div className={styles.cardLogo}>
                      {meta?.icon ?? site.name.charAt(0).toUpperCase()}
                    </div>
                    <button
                      className={styles.deleteBtn}
                      disabled={deleting === site.id}
                      onClick={() => setDeleteTarget({ id: site.id, name: site.name })}
                      title="Excluir site"
                    >
                      {deleting === site.id ? "…" : "✕"}
                    </button>
                  </div>
                  <div className={styles.cardBody}>
                    <div className={styles.siteInfo}>
                      <Link to={`/store/${site.id}`} className={styles.siteName}>{site.name}</Link>
                      <span className={styles.siteDesc}>
                        {site.short_description || "Sem descrição"}
                      </span>
                      <div className={styles.cardBadges}>
                        {meta && (
                          <span className={styles.typeBadge} style={{ background: `${cardColor}22`, color: cardColor }}>
                            {meta.label}
                          </span>
                        )}
                        {edLabel && (
                          <span className={styles.editorBadge} style={{ background: `${edLabel.color}18`, color: edLabel.color }}>
                            {edLabel.label}
                          </span>
                        )}
                      </div>
                    </div>
                    <div className={styles.siteActions}>
                      <button className={styles.editorBtn} onClick={() => openEditor(site.id, srcType)}>
                        ✏️ Editar
                      </button>
                      <a className={styles.openBtn} href={`${API_URL}/p/${site.id}`} target="_blank" rel="noreferrer">
                        ↗ Abrir
                      </a>
                    </div>
                  </div>
                </article>
              );
            })}
          </div>
        )}
      </div>

      {showModal && (
        <div className={styles.overlay} onClick={closeModal}>
          <div className={styles.modal} onClick={e => e.stopPropagation()}>

            <div className={styles.modalHeader}>
              <div className={styles.stepIndicator}>
                {[1, 2, 3].map(s => (
                  <div key={s} className={`${styles.stepDot} ${step >= s ? styles.stepDotActive : ""}`}>
                    <span>{s}</span>
                  </div>
                ))}
                <div className={styles.stepLine} />
              </div>
              <button className={styles.closeBtn} onClick={closeModal}>✕</button>
            </div>

            {step === 1 && (
              <div className={styles.modalBody}>
                <h2>Tipo de site</h2>
                <p>Escolha o modelo que melhor descreve seu projeto.</p>
                <div className={styles.typeGrid}>
                  {Object.entries(TYPE_META).map(([key, meta]) => (
                    <button
                      key={key}
                      className={`${styles.typeCard} ${type === key ? styles.typeCardActive : ""}`}
                      style={{ "--tc": meta.color }}
                      onClick={() => selectType(key)}
                    >
                      <span className={styles.typeIcon}>{meta.icon}</span>
                      <strong>{meta.label}</strong>
                      <span>{meta.description}</span>
                    </button>
                  ))}
                </div>
              </div>
            )}

            {step === 2 && (
              <div className={styles.modalBody}>
                <h2>Informações</h2>
                <p>Dê um nome ao seu site e escolha o editor.</p>
                <div className={styles.fieldGroup}>
                  <label>Nome *</label>
                  <input
                    className={styles.input}
                    value={name}
                    onChange={e => setName(e.target.value)}
                    placeholder="Ex: Loja Aurora Fitness"
                    maxLength={50}
                    autoFocus
                  />
                </div>
                <div className={styles.fieldGroup}>
                  <label>Descrição curta</label>
                  <input
                    className={styles.input}
                    value={shortDesc}
                    onChange={e => setShortDesc(e.target.value)}
                    placeholder="Uma frase objetiva (opcional)"
                    maxLength={150}
                  />
                </div>
                <div className={styles.fieldGroup}>
                  <label>Editor</label>
                  <div className={styles.editorRow}>
                    {EDITOR_OPTIONS.map(opt => (
                      <button
                        key={opt.value}
                        className={`${styles.editorChip} ${editor === opt.value ? styles.editorChipActive : ""}`}
                        onClick={() => setEditor(opt.value)}
                      >
                        <span className={styles.editorIcon}>{opt.icon}</span>
                        <strong>{opt.label}</strong>
                        <span>{opt.desc}</span>
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {step === 3 && (
              <div className={styles.modalBody}>
                <div className={styles.routeHeader}>
                  <div>
                    <h2>Rotas</h2>
                    <p>Configure as páginas do seu site.</p>
                  </div>
                  <span className={styles.routeCount}>{routes.length}/{maxRoutes}</span>
                </div>
                <div className={styles.routeList}>
                  {routes.map(route => (
                    <div key={route.id} className={styles.routeItem}>
                      <input
                        className={styles.input}
                        value={route.path}
                        onChange={e => updateRoute(route.id, "path", e.target.value)}
                        placeholder="/caminho"
                      />
                      <input
                        className={styles.input}
                        value={route.title}
                        onChange={e => updateRoute(route.id, "title", e.target.value)}
                        placeholder="Título"
                      />
                      <button
                        className={styles.removeRoute}
                        onClick={() => removeRoute(route.id)}
                        disabled={routes.length === 1}
                      >✕</button>
                    </div>
                  ))}
                </div>
                <button className={styles.addRouteBtn} onClick={addRoute} disabled={routes.length >= maxRoutes}>
                  + Adicionar rota
                </button>
                {formError && <p className={styles.formError}>{formError}</p>}
              </div>
            )}

            <div className={styles.modalFooter}>
              {step > 1 ? (
                <button className={styles.backBtn} onClick={() => setStep(s => s - 1)} disabled={submitting}>
                  ← Voltar
                </button>
              ) : <div />}
              {step < 3 ? (
                <Button variant="primary" onClick={() => setStep(s => s + 1)} disabled={step === 2 && !canStep2}>
                  Continuar →
                </Button>
              ) : (
                <Button variant="primary" onClick={handleCreate} disabled={!canCreate}>
                  {submitting ? "Criando…" : "✓ Criar site"}
                </Button>
              )}
            </div>

          </div>
        </div>
      )}

      {deleteTarget && (
        <div className={styles.overlay} onClick={() => !deleting && setDeleteTarget(null)}>
          <div className={styles.deleteModal} onClick={e => e.stopPropagation()}>
            <div className={styles.deleteModalIcon}>⚠️</div>
            <h2 className={styles.deleteModalTitle}>Excluir site</h2>
            <p className={styles.deleteModalText}>
              Tem certeza que deseja excluir <strong>{deleteTarget.name}</strong>? Esta ação é <strong>irreversível</strong> — o site e todos os dados serão apagados permanentemente.
            </p>
            <div className={styles.deleteModalActions}>
              <button
                className={styles.cancelDeleteBtn}
                onClick={() => setDeleteTarget(null)}
                disabled={!!deleting}
              >
                Cancelar
              </button>
              <button
                className={styles.confirmDeleteBtn}
                onClick={handleDelete}
                disabled={!!deleting}
              >
                {deleting ? "Excluindo…" : "Excluir definitivamente"}
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
};
