import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Button from "../../components/button/Button";
import CodeEditor from "../../components/code-editor/CodeEditor";
import { API_URL, apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./ReactEditor.module.scss";

type VersionItem = {
  id: string; version: number;
  source_type: "JXML" | "REACT" | "SVELTE" | "ELEMENTOR_JSON";
  source?: string; scan_status: "clean" | "warning" | "blocked";
  scan_score: number; published: boolean;
};
type RouteItem = { id: string; path: string; title: string; requires_auth: boolean; };
type VersionsResponse = { success: boolean; message: string; data: VersionItem[] };
type CreateVersionResponse = { success: boolean; message: string; data: VersionItem };
type RoutesResponse = { success: boolean; message: string; data: RouteItem[] };

// ── Bundle format ─────────────────────────────────────────────────────────
interface RouteEntry { component: string; }
interface Bundle {
  header: string;
  footer: string;
  globalCss: string;
  routes: Record<string, RouteEntry>; // path → { component }
}

const DEFAULT_COMPONENT = `function App() {
  return (
    <main>
      <h1>Página Inicial</h1>
      <p>Edite este componente.</p>
    </main>
  );
}`;
const DEFAULT_HEADER = `function Header() {
  return (
    <header style={{ background:"#1f2b43", color:"#fff", padding:"1rem 2rem" }}>
      <strong>Meu Site</strong>
    </header>
  );
}`;
const DEFAULT_FOOTER = `function Footer() {
  return (
    <footer style={{ background:"#f2f5fb", color:"#5a6379", padding:"1rem 2rem", textAlign:"center", fontSize:".85rem" }}>
      © 2025 · Feito com Jesterx
    </footer>
  );
}`;
const DEFAULT_CSS = `main { padding: 24px; max-width: 1024px; margin: 0 auto; }
h1 { color: #20283a; }
p { color: #5a6379; }`;

function makeDefaultBundle(routePaths: string[]): Bundle {
  const routes: Record<string, RouteEntry> = {};
  (routePaths.length ? routePaths : ["/"]).forEach((p) => {
    routes[p] = { component: DEFAULT_COMPONENT };
  });
  return { header: DEFAULT_HEADER, footer: DEFAULT_FOOTER, globalCss: DEFAULT_CSS, routes };
}

function parseBundle(raw: string | undefined, routePaths: string[]): Bundle {
  const fallback = makeDefaultBundle(routePaths);
  if (!raw) return fallback;
  try {
    const p = JSON.parse(raw) as Record<string, unknown>;
    // Migrate old format { component, header, footer, css }
    if (typeof p.component === "string" && !p.routes) {
      const routes: Record<string, RouteEntry> = {};
      (routePaths.length ? routePaths : ["/"]).forEach((path) => {
        routes[path] = { component: p.component as string };
      });
      return {
        header:    (typeof p.header === "string" ? p.header : DEFAULT_HEADER),
        footer:    (typeof p.footer === "string" ? p.footer : DEFAULT_FOOTER),
        globalCss: (typeof p.css === "string" ? p.css : DEFAULT_CSS),
        routes,
      };
    }
    const routes: Record<string, RouteEntry> = {};
    const rawRoutes = (p.routes ?? {}) as Record<string, { component?: string }>;
    // Merge with current route paths
    (routePaths.length ? routePaths : ["/"]).forEach((path) => {
      routes[path] = { component: rawRoutes[path]?.component ?? DEFAULT_COMPONENT };
    });
    return {
      header:    (typeof p.header === "string" ? p.header : DEFAULT_HEADER),
      footer:    (typeof p.footer === "string" ? p.footer : DEFAULT_FOOTER),
      globalCss: (typeof p.globalCss === "string" ? p.globalCss : DEFAULT_CSS),
      routes,
    };
  } catch { return fallback; }
}

// ── Preview builder ────────────────────────────────────────────────────────
function buildPreview(bundle: Bundle, activePath: string) {
  const routeComp = bundle.routes[activePath]?.component ?? DEFAULT_COMPONENT;
  const safeCode = [bundle.header, routeComp, bundle.footer]
    .join("\n\n")
    .replace(/<\/script>/gi, "<\\/script>");
  return `<!doctype html><html><head><meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/>
<script crossorigin src="https://unpkg.com/react@18/umd/react.development.js"></script>
<script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.development.js"></script>
<script src="https://unpkg.com/@babel/standalone/babel.min.js"></script>
<style>body{margin:0;font-family:Inter,system-ui,sans-serif}${bundle.globalCss}</style>
</head><body><div id="app"></div>
<script type="text/babel">
const React=window.React;const ReactDOM=window.ReactDOM;
${safeCode}
const __H=(typeof Header!=="undefined")?Header:()=>null;
const __A=(typeof App!=="undefined")?App:(()=><main><h1>Crie uma função App()</h1></main>);
const __F=(typeof Footer!=="undefined")?Footer:()=>null;
function __Root(){return(<>< __H/><__A/><__F/></>);}
ReactDOM.createRoot(document.getElementById('app')).render(<__Root/>);
</script></body></html>`;
}

// ── Active editing tab ─────────────────────────────────────────────────────
type SharedTab = "header" | "footer" | "globalCss";
type ActiveTab = { kind: "shared"; tab: SharedTab } | { kind: "route"; path: string };

export const ReactEditor: React.FC = () => {
  const { siteId = "" } = useParams();
  const { websiteId } = useAuthContext();
  const navigate = useNavigate();

  const [bundle, setBundle]       = useState<Bundle>({ header: DEFAULT_HEADER, footer: DEFAULT_FOOTER, globalCss: DEFAULT_CSS, routes: { "/": { component: DEFAULT_COMPONENT } } });
  const [routes, setRoutes]       = useState<RouteItem[]>([]);
  const [active, setActive]       = useState<ActiveTab>({ kind: "route", path: "/" });
  const [previewPath, setPreviewPath] = useState("/");
  const [saving, setSaving]       = useState(false);
  const [error, setError]         = useState<string | null>(null);
  const [newRoutePath, setNewRoutePath] = useState("");
  const [addingRoute, setAddingRoute]   = useState(false);

  const preview = useMemo(() => buildPreview(bundle, previewPath), [bundle, previewPath]);

  function setRouteComponent(path: string, component: string) {
    setBundle((b) => ({ ...b, routes: { ...b.routes, [path]: { component } } }));
  }
  function setShared(field: SharedTab, val: string) {
    setBundle((b) => ({ ...b, [field]: val }));
  }

  async function loadData() {
    if (!siteId) return;
    try {
      const [vResp, rResp] = await Promise.all([
        apiFetch<VersionsResponse>(`/api/v1/sites/${siteId}/versions`, { method: "GET", websiteId }),
        apiFetch<RoutesResponse>(`/api/v1/sites/${siteId}/routes`, { method: "GET", websiteId }),
      ]);
      const latestAny = vResp.data[0];
      if (latestAny && latestAny.source_type !== "REACT") {
        throw new Error(`Este site foi criado em ${latestAny.source_type}. Use o editor correto.`);
      }
      const routeList = rResp.data ?? [];
      setRoutes(routeList);
      const routePaths = routeList.map((r) => r.path);
      const latest = vResp.data.find((v) => v.source_type === "REACT");
      const parsed = parseBundle(latest?.source, routePaths);
      setBundle(parsed);
      // Set active to first route
      const firstPath = routePaths[0] ?? "/";
      setActive({ kind: "route", path: firstPath });
      setPreviewPath(firstPath);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao carregar");
    }
  }

  useEffect(() => { loadData(); }, [siteId]); // eslint-disable-line react-hooks/exhaustive-deps

  async function save(publish: boolean) {
    if (!siteId) return;
    setSaving(true); setError(null);
    try {
      const resp = await apiFetch<CreateVersionResponse>(`/api/v1/sites/${siteId}/versions`, {
        method: "POST", websiteId,
        body: JSON.stringify({ source_type: "REACT", source: JSON.stringify(bundle) }),
      });
      if (publish && resp.data.scan_status !== "blocked") {
        await apiFetch(`/api/v1/sites/${siteId}/publish/${resp.data.version}`, { method: "POST", websiteId });
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao salvar");
    } finally { setSaving(false); }
  }

  async function addRoute() {
    const path = newRoutePath.trim().startsWith("/") ? newRoutePath.trim() : `/${newRoutePath.trim()}`;
    if (!path || path.length < 2) return;
    const newRoutes = [...routes, { id: `r-${Date.now()}`, path, title: path.slice(1) || "Nova página", requires_auth: false }];
    setRoutes(newRoutes);
    setBundle((b) => ({ ...b, routes: { ...b.routes, [path]: { component: DEFAULT_COMPONENT } } }));
    await apiFetch(`/api/v1/sites/${siteId}/routes`, {
      method: "POST", websiteId,
      body: JSON.stringify({ routes: newRoutes.map((r) => ({ path: r.path, title: r.title, requires_auth: r.requires_auth })) }),
    }).catch(() => {});
    setNewRoutePath(""); setAddingRoute(false);
    setActive({ kind: "route", path });
    setPreviewPath(path);
  }

  async function removeRoute(path: string) {
    const newList = routes.filter((r) => r.path !== path);
    setRoutes(newList);
    setBundle((b) => {
      const next = { ...b.routes };
      delete next[path];
      return { ...b, routes: next };
    });
    await apiFetch(`/api/v1/sites/${siteId}/routes`, {
      method: "POST", websiteId,
      body: JSON.stringify({ routes: newList.map((r) => ({ path: r.path, title: r.title, requires_auth: r.requires_auth })) }),
    }).catch(() => {});
    const firstPath = newList[0]?.path ?? "/";
    setActive({ kind: "route", path: firstPath });
    setPreviewPath(firstPath);
  }

  // Determine what's currently in the editor
  function getCurrentValue(): string {
    if (active.kind === "shared") return bundle[active.tab];
    return bundle.routes[active.path]?.component ?? DEFAULT_COMPONENT;
  }
  function handleChange(val: string) {
    if (active.kind === "shared") setShared(active.tab, val);
    else setRouteComponent(active.path, val);
  }
  const currentLang: "tsx" | "css" = active.kind === "shared" && active.tab === "globalCss" ? "css" : "tsx";
  const activeRouteTitle = active.kind === "route" ? (routes.find((r) => r.path === active.path)?.title ?? active.path) : null;

  const SHARED_FILES: { tab: SharedTab; label: string }[] = [
    { tab: "header",    label: "Header.tsx" },
    { tab: "footer",    label: "Footer.tsx" },
    { tab: "globalCss", label: "globals.css" },
  ];

  return (
    <div className={styles.root}>
      {/* Top bar */}
      <header className={styles.topbar}>
        <div className={styles.topLeft}>
          <button className={styles.backBtn} onClick={() => navigate("/pages")}>← Voltar</button>
          <span className={styles.editorLabel}>React Editor</span>
          <span className={styles.siteId}>{siteId}</span>
        </div>
        <div className={styles.topRight}>
          {error && <span className={styles.errInline}>{error}</span>}
          <Button type="button" variant="secondary" disabled={saving} onClick={() => save(false)}>Salvar</Button>
          <Button type="button" variant="primary"   disabled={saving} onClick={() => save(true)}>Publicar</Button>
          <a href={`${API_URL}/p/${siteId}`} target="_blank" rel="noreferrer" className={styles.liveLink}>↗ Ver publicado</a>
        </div>
      </header>

      {/* Main 3-column layout */}
      <div className={styles.workspace}>
        {/* Sidebar */}
        <aside className={styles.sidebar}>
          <p className={styles.sideHead}>Arquivos compartilhados</p>
          {SHARED_FILES.map((f) => (
            <button key={f.tab}
              className={`${styles.fileItem} ${active.kind === "shared" && active.tab === f.tab ? styles.fileItemActive : ""}`}
              onClick={() => setActive({ kind: "shared", tab: f.tab })}
            >{f.label}</button>
          ))}

          <div className={styles.sideDivider} />
          <div className={styles.sideHeadRow}>
            <p className={styles.sideHead}>Páginas / Rotas</p>
            <button className={styles.addBtn} onClick={() => setAddingRoute((v) => !v)} title="Nova página">＋</button>
          </div>

          {addingRoute && (
            <div className={styles.addRouteBox}>
              <input
                className={styles.addRouteInput}
                value={newRoutePath}
                onChange={(e) => setNewRoutePath(e.target.value)}
                placeholder="/minha-pagina"
                onKeyDown={(e) => e.key === "Enter" && addRoute()}
                autoFocus
              />
              <button className={styles.addRouteConfirm} onClick={addRoute}>OK</button>
              <button className={styles.addRouteCancel} onClick={() => setAddingRoute(false)}>✕</button>
            </div>
          )}

          {routes.map((r) => (
            <div key={r.path} className={`${styles.routeItem} ${active.kind === "route" && active.path === r.path ? styles.routeItemActive : ""}`}>
              <button className={styles.routeBtn}
                onClick={() => { setActive({ kind: "route", path: r.path }); setPreviewPath(r.path); }}
              >
                <span className={styles.routePath}>{r.path}</span>
                <span className={styles.routeTitle}>{r.title}</span>
              </button>
              {r.path !== "/" && (
                <button className={styles.routeRemove} onClick={() => removeRoute(r.path)} title="Remover">✕</button>
              )}
            </div>
          ))}
        </aside>

        {/* Editor */}
        <section className={styles.editorPane}>
          <div className={styles.editorLabel2}>
            {active.kind === "shared"
              ? (active.tab === "globalCss" ? "globals.css" : active.tab === "header" ? "Header.tsx" : "Footer.tsx")
              : `${activeRouteTitle ?? active.path} — App.tsx`}
          </div>
          <div className={styles.editorBody}>
            <CodeEditor key={active.kind + (active.kind === "shared" ? active.tab : active.path)} value={getCurrentValue()} onChange={handleChange} language={currentLang} flat />
          </div>
        </section>

        {/* Preview */}
        <section className={styles.previewPane}>
          {/* Route tabs for preview */}
          <div className={styles.previewTabs}>
            {routes.map((r) => (
              <button key={r.path}
                className={`${styles.previewTab} ${previewPath === r.path ? styles.previewTabActive : ""}`}
                onClick={() => setPreviewPath(r.path)}
              >{r.path}</button>
            ))}
          </div>
          <iframe title="react-preview" srcDoc={preview} className={styles.previewFrame} />
        </section>
      </div>
    </div>
  );
};
