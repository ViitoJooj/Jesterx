import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Button from "../../components/button/Button";
import CodeEditor from "../../components/code-editor/CodeEditor";
import { API_URL, apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./SvelteEditor.module.scss";

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

interface RouteEntry { component: string; }
interface Bundle {
  header: string;
  footer: string;
  globalCss: string;
  routes: Record<string, RouteEntry>;
}

const DEFAULT_COMPONENT = `<script>
  let title = "Página Inicial";
</script>

<main>
  <h1>{title}</h1>
  <p>Edite e veja o preview em tempo real.</p>
</main>`;

const DEFAULT_HEADER = `<script>
  export let siteName = "Meu Site";
</script>

<header>
  <strong>{siteName}</strong>
</header>

<style>
  header { background:#1f2b43; color:#fff; padding:1rem 2rem; display:flex; align-items:center; }
</style>`;

const DEFAULT_FOOTER = `<footer>© 2025 · Feito com Jesterx</footer>

<style>
  footer { background:#f2f5fb; color:#5a6379; padding:1rem 2rem; text-align:center; font-size:.85rem; }
</style>`;

const DEFAULT_CSS = `main { padding: 24px; max-width: 1024px; margin: 0 auto; }
h1 { color: #1f2b43; }
p { color: #5f6980; }`;

function makeDefaultBundle(routePaths: string[]): Bundle {
  const routes: Record<string, RouteEntry> = {};
  (routePaths.length ? routePaths : ["/"]).forEach((p) => { routes[p] = { component: DEFAULT_COMPONENT }; });
  return { header: DEFAULT_HEADER, footer: DEFAULT_FOOTER, globalCss: DEFAULT_CSS, routes };
}

function parseBundle(raw: string | undefined, routePaths: string[]): Bundle {
  const fallback = makeDefaultBundle(routePaths);
  if (!raw) return fallback;
  try {
    const p = JSON.parse(raw) as Record<string, unknown>;
    // Migrate old format { component, css }
    if (typeof p.component === "string" && !p.routes) {
      const routes: Record<string, RouteEntry> = {};
      (routePaths.length ? routePaths : ["/"]).forEach((path) => { routes[path] = { component: p.component as string }; });
      return { header: (typeof p.header === "string" ? p.header : DEFAULT_HEADER), footer: (typeof p.footer === "string" ? p.footer : DEFAULT_FOOTER), globalCss: (typeof p.css === "string" ? p.css : DEFAULT_CSS), routes };
    }
    const routes: Record<string, RouteEntry> = {};
    const rawRoutes = (p.routes ?? {}) as Record<string, { component?: string }>;
    (routePaths.length ? routePaths : ["/"]).forEach((path) => { routes[path] = { component: rawRoutes[path]?.component ?? DEFAULT_COMPONENT }; });
    return { header: (typeof p.header === "string" ? p.header : DEFAULT_HEADER), footer: (typeof p.footer === "string" ? p.footer : DEFAULT_FOOTER), globalCss: (typeof p.globalCss === "string" ? p.globalCss : DEFAULT_CSS), routes };
  } catch { return fallback; }
}

function escTpl(value: string) {
  return value.replace(/\\/g, "\\\\").replace(/`/g, "\\`").replace(/\$\{/g, "\\${").replace(/<\/script>/gi, "<\\/script>");
}

function buildPreview(bundle: Bundle, activePath: string) {
  const routeComp = bundle.routes[activePath]?.component ?? DEFAULT_COMPONENT;
  const mainSfc = routeComp.includes("<style") ? routeComp : `<style>\n${bundle.globalCss}\n</style>\n${routeComp}`;
  const headerSfc = bundle.header || `<slot/>`;
  const footerSfc = bundle.footer || ``;
  return `<!doctype html><html><head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/>
<style>body{margin:0;font-family:Inter,system-ui,sans-serif}</style>
</head><body><div id="app"></div>
<script type="module">
import { compile } from "https://unpkg.com/svelte@4.2.19/compiler.mjs";
async function compileAndLoad(source, name) {
  try {
    const result = compile(source, { generate:"dom", format:"esm", name });
    let code = result.js.code
      .replaceAll('from "svelte/internal"','from "https://unpkg.com/svelte@4.2.19/src/runtime/internal/index.mjs"')
      .replaceAll("from 'svelte/internal'","from 'https://unpkg.com/svelte@4.2.19/src/runtime/internal/index.mjs'");
    return (await import(URL.createObjectURL(new Blob([code],{type:"text/javascript"})))).default;
  } catch(e) { return null; }
}
try {
  const [Header, App, Footer] = await Promise.all([
    compileAndLoad(\`${escTpl(headerSfc)}\`,"Header"),
    compileAndLoad(\`${escTpl(mainSfc)}\`,"App"),
    compileAndLoad(\`${escTpl(footerSfc)}\`,"Footer"),
  ]);
  const root = document.getElementById("app");
  if (Header) new Header({ target: root });
  if (App) { const d=document.createElement("div"); root.appendChild(d); new App({target:d}); }
  else { root.innerHTML += '<pre>Erro ao compilar App.svelte</pre>'; }
  if (Footer) new Footer({ target: root });
} catch(e) { document.getElementById("app").innerHTML="<pre>"+String(e).replace(/</g,"&lt;")+"</pre>"; }
</script></body></html>`;
}

type SharedTab = "header" | "footer" | "globalCss";
type ActiveTab = { kind: "shared"; tab: SharedTab } | { kind: "route"; path: string };

export const SvelteEditor: React.FC = () => {
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
      if (latestAny && latestAny.source_type !== "SVELTE") throw new Error(`Este site foi criado em ${latestAny.source_type}. Use o editor correto.`);
      const routeList = rResp.data ?? [];
      setRoutes(routeList);
      const routePaths = routeList.map((r) => r.path);
      const latest = vResp.data.find((v) => v.source_type === "SVELTE");
      setBundle(parseBundle(latest?.source, routePaths));
      const firstPath = routePaths[0] ?? "/";
      setActive({ kind: "route", path: firstPath });
      setPreviewPath(firstPath);
    } catch (err) { setError(err instanceof Error ? err.message : "Erro ao carregar"); }
  }

  useEffect(() => { loadData(); }, [siteId]); // eslint-disable-line react-hooks/exhaustive-deps

  async function save(publish: boolean) {
    if (!siteId) return;
    setSaving(true); setError(null);
    try {
      const resp = await apiFetch<CreateVersionResponse>(`/api/v1/sites/${siteId}/versions`, {
        method: "POST", websiteId,
        body: JSON.stringify({ source_type: "SVELTE", source: JSON.stringify(bundle) }),
      });
      if (publish && resp.data.scan_status !== "blocked") {
        await apiFetch(`/api/v1/sites/${siteId}/publish/${resp.data.version}`, { method: "POST", websiteId });
      }
    } catch (err) { setError(err instanceof Error ? err.message : "Erro ao salvar"); }
    finally { setSaving(false); }
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
    setActive({ kind: "route", path }); setPreviewPath(path);
  }

  async function removeRoute(path: string) {
    const newList = routes.filter((r) => r.path !== path);
    setRoutes(newList);
    setBundle((b) => { const next = { ...b.routes }; delete next[path]; return { ...b, routes: next }; });
    await apiFetch(`/api/v1/sites/${siteId}/routes`, {
      method: "POST", websiteId,
      body: JSON.stringify({ routes: newList.map((r) => ({ path: r.path, title: r.title, requires_auth: r.requires_auth })) }),
    }).catch(() => {});
    const fp = newList[0]?.path ?? "/";
    setActive({ kind: "route", path: fp }); setPreviewPath(fp);
  }

  function getCurrentValue(): string {
    if (active.kind === "shared") return bundle[active.tab];
    return bundle.routes[active.path]?.component ?? DEFAULT_COMPONENT;
  }
  function handleChange(val: string) {
    if (active.kind === "shared") setShared(active.tab, val);
    else setRouteComponent(active.path, val);
  }

  const currentLang: "svelte" | "css" = active.kind === "shared" && active.tab === "globalCss" ? "css" : "svelte";
  const activeRouteTitle = active.kind === "route" ? (routes.find((r) => r.path === active.path)?.title ?? active.path) : null;

  const SHARED_FILES: { tab: SharedTab; label: string }[] = [
    { tab: "header",    label: "Header.svelte" },
    { tab: "footer",    label: "Footer.svelte" },
    { tab: "globalCss", label: "globals.css" },
  ];

  return (
    <div className={styles.root}>
      <header className={styles.topbar}>
        <div className={styles.topLeft}>
          <button className={styles.backBtn} onClick={() => navigate("/pages")}>← Voltar</button>
          <span className={styles.editorLabel}>Svelte Editor</span>
          <span className={styles.siteId}>{siteId}</span>
        </div>
        <div className={styles.topRight}>
          {error && <span className={styles.errInline}>{error}</span>}
          <Button type="button" variant="secondary" disabled={saving} onClick={() => save(false)}>Salvar</Button>
          <Button type="button" variant="primary"   disabled={saving} onClick={() => save(true)}>Publicar</Button>
          <a href={`${API_URL}/p/${siteId}`} target="_blank" rel="noreferrer" className={styles.liveLink}>↗ Ver publicado</a>
        </div>
      </header>

      <div className={styles.workspace}>
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
            <button className={styles.addBtn} onClick={() => setAddingRoute((v) => !v)}>＋</button>
          </div>

          {addingRoute && (
            <div className={styles.addRouteBox}>
              <input className={styles.addRouteInput} value={newRoutePath} onChange={(e) => setNewRoutePath(e.target.value)} placeholder="/minha-pagina" onKeyDown={(e) => e.key === "Enter" && addRoute()} autoFocus />
              <button className={styles.addRouteConfirm} onClick={addRoute}>OK</button>
              <button className={styles.addRouteCancel} onClick={() => setAddingRoute(false)}>✕</button>
            </div>
          )}

          {routes.map((r) => (
            <div key={r.path} className={`${styles.routeItem} ${active.kind === "route" && active.path === r.path ? styles.routeItemActive : ""}`}>
              <button className={styles.routeBtn} onClick={() => { setActive({ kind: "route", path: r.path }); setPreviewPath(r.path); }}>
                <span className={styles.routePath}>{r.path}</span>
                <span className={styles.routeTitle}>{r.title}</span>
              </button>
              {r.path !== "/" && <button className={styles.routeRemove} onClick={() => removeRoute(r.path)}>✕</button>}
            </div>
          ))}
        </aside>

        <section className={styles.editorPane}>
          <div className={styles.editorLabel2}>
            {active.kind === "shared"
              ? (active.tab === "globalCss" ? "globals.css" : active.tab === "header" ? "Header.svelte" : "Footer.svelte")
              : `${activeRouteTitle ?? active.path} — App.svelte`}
          </div>
          <div className={styles.editorBody}>
            <CodeEditor key={active.kind + (active.kind === "shared" ? active.tab : active.path)} value={getCurrentValue()} onChange={handleChange} language={currentLang} flat />
          </div>
        </section>

        <section className={styles.previewPane}>
          <div className={styles.previewTabs}>
            {routes.map((r) => (
              <button key={r.path}
                className={`${styles.previewTab} ${previewPath === r.path ? styles.previewTabActive : ""}`}
                onClick={() => setPreviewPath(r.path)}
              >{r.path}</button>
            ))}
          </div>
          <iframe title="svelte-preview" srcDoc={preview} className={styles.previewFrame} />
        </section>
      </div>
    </div>
  );
};
