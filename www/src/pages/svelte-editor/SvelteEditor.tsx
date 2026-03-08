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

// ── Jesterx Svelte custom-element sources ────────────────────────────────────
// Each entry: [customElementName, svelteSource]
const JX_SVELTE_COMPONENTS: [string, string][] = [
  ["jx-heading", `<svelte:options customElement="jx-heading"/>
<script>
  export let text = 'Título';
  export let fontsize = '28';
  export let fontweight = '700';
  export let color = '#20283a';
  export let textalign = 'left';
</script>
<h2 style="font-size:{fontsize}px;font-weight:{fontweight};color:{color};text-align:{textalign};margin:0">{text}</h2>`],

  ["jx-paragraph", `<svelte:options customElement="jx-paragraph"/>
<script>
  export let text = 'Parágrafo';
  export let fontsize = '16';
  export let color = '#5a6379';
  export let textalign = 'left';
</script>
<p style="font-size:{fontsize}px;color:{color};text-align:{textalign};margin:0;line-height:1.6">{text}</p>`],

  ["jx-button", `<svelte:options customElement="jx-button"/>
<script>
  export let label = 'Botão';
  export let href = '';
  export let bgcolor = '#ff5d1f';
  export let textcolor = '#fff';
  export let borderradius = '8';
  export let fontsize = '15';
  function handleClick() { if (href) window.location.href = href; }
</script>
<button on:click={handleClick} style="background:{bgcolor};color:{textcolor};border:0;border-radius:{borderradius}px;padding:10px 22px;font-size:{fontsize}px;font-weight:600;cursor:pointer;font-family:inherit">{label}</button>`],

  ["jx-image", `<svelte:options customElement="jx-image"/>
<script>
  export let src = '';
  export let alt = '';
  export let objectfit = 'cover';
  export let borderradius = '0';
</script>
{#if src}
  <img {src} {alt} style="width:100%;height:100%;object-fit:{objectfit};border-radius:{borderradius}px;display:block" />
{:else}
  <div style="width:100%;height:200px;background:#f0f4ff;display:flex;align-items:center;justify-content:center;color:#9aa5bc;font-size:13px;border-radius:8px">🖼 Imagem</div>
{/if}`],

  ["jx-divider", `<svelte:options customElement="jx-divider"/>
<script>
  export let color = '#dde3f0';
  export let thickness = '1';
  export let margintop = '16';
  export let marginbottom = '16';
</script>
<hr style="border:0;border-top:{thickness}px solid {color};margin:{margintop}px 0 {marginbottom}px 0" />`],

  ["jx-video", `<svelte:options customElement="jx-video"/>
<script>
  export let url = '';
  export let height = '240';
  function getEmbed(u) {
    if (!u) return '';
    const yt = u.match(/(?:youtube\\.com\\/watch\\?v=|youtu\\.be\\/)([^&?/]+)/);
    if (yt) return 'https://www.youtube.com/embed/' + yt[1];
    return u;
  }
  $: embed = getEmbed(url);
</script>
{#if embed}
  <iframe src={embed} width="100%" {height} frameborder="0" allowfullscreen style="display:block;border-radius:8px"></iframe>
{:else}
  <div style="background:#111;color:#fff;width:100%;height:{height}px;display:flex;align-items:center;justify-content:center;border-radius:8px;font-size:15px">▶ Vídeo</div>
{/if}`],

  ["jx-profile-card", `<svelte:options customElement="jx-profile-card"/>
<script>
  export let name = 'Nome';
  export let subtitle = 'Subtítulo';
  export let image = '';
  export let bgcolor = '#fff';
</script>
<div style="background:{bgcolor};border:1px solid #dde3f0;border-radius:12px;padding:20px;display:flex;flex-direction:column;align-items:center;gap:8px;text-align:center">
  <img src={image||'https://placehold.co/80x80'} alt={name} style="width:80px;height:80px;border-radius:50%;object-fit:cover" />
  <div style="font-weight:700;font-size:16px;color:#1a2740">{name}</div>
  <div style="font-size:13px;color:#6b7a99">{subtitle}</div>
</div>`],

  ["jx-input-var", `<svelte:options customElement="jx-input-var"/>
<script>
  export let placeholder = 'Digite aqui...';
  export let inputtype = 'text';
  export let value = '';
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>
<input type={inputtype} {placeholder} {value} on:input={e => dispatch('change', e.target.value)}
  style="width:100%;padding:10px 14px;border:1.5px solid #dde3f0;border-radius:8px;font-size:14px;font-family:inherit;outline:none;box-sizing:border-box" />`],

  ["jx-variable-text", `<svelte:options customElement="jx-variable-text"/>
<script>
  export let text = '{{variavel}}';
  export let fontsize = '14';
  export let color = '#20283a';
  $: resolved = text.replace(/\\{\\{(\\w+)\\}\\}/g, '');
</script>
<span style="font-size:{fontsize}px;color:{color}">{resolved}</span>`],

  ["jx-user-avatar", `<svelte:options customElement="jx-user-avatar"/>
<script>
  export let src = '';
  export let size = '48';
  export let objectfit = 'cover';
</script>
<img src={src||'https://placehold.co/'+size+'x'+size} alt="avatar"
  style="width:{size}px;height:{size}px;border-radius:50%;object-fit:{objectfit};display:block" />`],

  ["jx-product-card", `<svelte:options customElement="jx-product-card"/>
<script>
  export let siteid = '';
  export let productid = '';
  let prod = null, loading = true;
  $: if (siteid && productid) {
    fetch('/api/store/'+siteid+'/products/'+productid)
      .then(r => r.ok ? r.json() : null)
      .then(d => { prod = d?.data || null; loading = false; })
      .catch(() => { loading = false; });
  } else { loading = false; }
</script>
<div style="background:#fff;border:1px solid #dde3f0;border-radius:12px;overflow:hidden;display:flex;flex-direction:column">
  {#if loading}
    <div style="padding:24px;text-align:center;color:#9aa5bc;font-size:13px">Carregando...</div>
  {:else if prod}
    {#if prod.images?.[0]}<img src={prod.images[0]} alt={prod.name} style="width:100%;height:180px;object-fit:cover" />{/if}
    <div style="padding:12px 16px;display:flex;flex-direction:column;gap:6px">
      <div style="font-weight:700;font-size:15px;color:#1a2740">{prod.name}</div>
      <div style="font-size:16px;font-weight:700;color:#ff5d1f">R$ {Number(prod.price||0).toFixed(2)}</div>
    </div>
  {:else}
    <div style="padding:24px;text-align:center;color:#9aa5bc;font-size:13px">🛍 product-card</div>
  {/if}
</div>`],

  ["jx-product-list", `<svelte:options customElement="jx-product-list"/>
<script>
  export let siteid = '';
  export let pagesize = '6';
  let prods = [], page = 0, loading = true;
  $: ps = Math.max(1, parseInt(pagesize)||6);
  $: total = Math.ceil(prods.length / ps);
  $: items = prods.slice(page*ps, (page+1)*ps);
  $: if (siteid) {
    loading = true;
    fetch('/api/store/'+siteid+'/products')
      .then(r => r.ok ? r.json() : null)
      .then(d => { prods = d?.data || []; loading = false; })
      .catch(() => { loading = false; });
  }
</script>
<div style="display:flex;flex-direction:column;gap:12px">
  {#if loading}<div style="color:#9aa5bc;padding:20px;font-size:13px">Carregando...</div>
  {:else if !prods.length}<div style="color:#9aa5bc;padding:20px;font-size:13px">Nenhum produto.</div>
  {:else}
    <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:12px">
      {#each items as p}
        <div style="background:#fff;border:1px solid #dde3f0;border-radius:10px;overflow:hidden">
          {#if p.images?.[0]}<img src={p.images[0]} alt={p.name} style="width:100%;height:120px;object-fit:cover" />{/if}
          <div style="padding:8px 10px">
            <div style="font-weight:600;font-size:13px;color:#1a2740">{p.name}</div>
            <div style="font-size:13px;font-weight:700;color:#ff5d1f">R$ {Number(p.price||0).toFixed(2)}</div>
          </div>
        </div>
      {/each}
    </div>
    {#if total>1}
      <div style="display:flex;gap:8px;align-items:center;justify-content:center">
        <button on:click={()=>page=Math.max(0,page-1)} disabled={page===0} style="padding:4px 12px;background:#f4f6fb;border:1px solid #dde3f0;border-radius:6px;cursor:pointer;font-family:inherit">← Anterior</button>
        <span style="font-size:12px;color:#6b7a99">{page+1} / {total}</span>
        <button on:click={()=>page=Math.min(total-1,page+1)} disabled={page>=total-1} style="padding:4px 12px;background:#f4f6fb;border:1px solid #dde3f0;border-radius:6px;cursor:pointer;font-family:inherit">Próximo →</button>
      </div>
    {/if}
  {/if}
</div>`],

  ["jx-cart-items", `<svelte:options customElement="jx-cart-items"/>
<script>
  export let storagekey = 'jx_cart';
  let items = [];
  function load() { try { items = JSON.parse(localStorage.getItem(storagekey)||'[]'); } catch(e) { items = []; } }
  function save(c) { localStorage.setItem(storagekey, JSON.stringify(c)); items = c; document.dispatchEvent(new CustomEvent('jx:cartupdate')); }
  function dec(id) { save(items.map(x=>x.id===id?{...x,qty:Math.max(0,(x.qty||1)-1)}:x).filter(x=>(x.qty||0)>0)); }
  function inc(id) { save(items.map(x=>x.id===id?{...x,qty:(x.qty||1)+1}:x)); }
  function rm(id)  { save(items.filter(x=>x.id!==id)); }
  load();
  document.addEventListener('jx:cartupdate', load);
</script>
{#if !items.length}
  <div style="color:#9aa5bc;padding:20px;text-align:center;font-size:13px">Carrinho vazio</div>
{:else}
  <div style="display:flex;flex-direction:column;gap:8px">
    {#each items as it}
      <div style="display:flex;align-items:center;gap:8px;background:#fff;border:1px solid #dbe2f3;border-radius:8px;padding:8px 10px">
        <div style="flex:1;font-weight:600;font-size:13px;color:#1a2740;overflow:hidden">{it.name}</div>
        <div style="display:flex;align-items:center;gap:4px">
          <button on:click={()=>dec(it.id)} style="width:26px;height:26px;border:1px solid #ccd5e8;background:#f4f6fb;border-radius:6px;cursor:pointer;font-weight:700;font-family:inherit">−</button>
          <span style="min-width:22px;text-align:center;font-weight:700;font-size:13px">{it.qty||1}</span>
          <button on:click={()=>inc(it.id)} style="width:26px;height:26px;border:1px solid #ccd5e8;background:#f4f6fb;border-radius:6px;cursor:pointer;font-weight:700;font-family:inherit">+</button>
        </div>
        <div style="font-weight:700;color:#ff5d1f;font-size:13px;min-width:64px;text-align:right">R$ {(Number(it.price||0)*(it.qty||1)).toFixed(2)}</div>
        <button on:click={()=>rm(it.id)} style="background:none;border:0;color:#9aa5bc;cursor:pointer;font-size:18px;line-height:1;padding:0 4px">×</button>
      </div>
    {/each}
  </div>
{/if}`],

  ["jx-admin-add-btn", `<svelte:options customElement="jx-admin-add-btn"/>
<script>
  export let label = 'Adicionar';
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>
<button on:click={()=>dispatch('add')} style="background:#1a2740;color:#fff;border:0;border-radius:8px;padding:8px 18px;font-size:13px;font-weight:600;cursor:pointer;font-family:inherit;display:inline-flex;align-items:center;gap:6px">＋ {label}</button>`],
];

function escTpl(value: string) {
  return value.replace(/\\/g, "\\\\").replace(/`/g, "\\`").replace(/\$\{/g, "\\${").replace(/<\/script>/gi, "<\\/script>");
}

function buildPreview(bundle: Bundle, activePath: string): string {
  const routeComp = bundle.routes[activePath]?.component ?? DEFAULT_COMPONENT;
  const mainSfc   = routeComp.includes("<style") ? routeComp : `<style>\n${bundle.globalCss}\n</style>\n${routeComp}`;
  const headerSfc = bundle.header || `<slot/>`;
  const footerSfc = bundle.footer || ``;

  const svelteInternalUrl = "https://unpkg.com/svelte@4.2.19/src/runtime/internal/index.mjs";

  // Serialize custom element sources as template literals for injection
  const jxEntriesJs = JX_SVELTE_COMPONENTS.map(([, src]) =>
    `\`${escTpl(src)}\``
  ).join(",\n    ");

  return `<!doctype html><html><head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/>
<style>body{margin:0;font-family:Inter,system-ui,sans-serif}</style>
</head><body><div id="app"></div>
<script type="module">
import { compile } from "https://unpkg.com/svelte@4.2.19/compiler.mjs";

function fixImports(code) {
  const INT = "${svelteInternalUrl}";
  return code
    .replaceAll('from "svelte/internal"', 'from "' + INT + '"')
    .replaceAll("from 'svelte/internal'", "from '" + INT + "'");
}

async function compileAndLoad(source, name) {
  try {
    const result = compile(source, { generate:"dom", format:"esm", name });
    const code = fixImports(result.js.code);
    return (await import(URL.createObjectURL(new Blob([code],{type:"text/javascript"})))).default;
  } catch(e) { console.warn("Compile error "+name, e); return null; }
}

async function compileCustomEl(source) {
  try {
    const result = compile(source, { generate:"dom", format:"esm", customElement: true });
    const code = fixImports(result.js.code);
    await import(URL.createObjectURL(new Blob([code],{type:"text/javascript"})));
  } catch(e) { console.warn("Custom element error", e); }
}

// Register Jx custom elements
const jxSources = [
    ${jxEntriesJs}
];
await Promise.all(jxSources.map(s => compileCustomEl(s)));

// Compile and mount Header, App, Footer
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

const JX_SVELTE_SNIPPETS: { name: string; snippet: string }[] = [
  { name: "jx-heading",       snippet: '<jx-heading text="Título" fontsize="32"></jx-heading>' },
  { name: "jx-paragraph",     snippet: '<jx-paragraph text="Seu texto." fontsize="16"></jx-paragraph>' },
  { name: "jx-button",        snippet: '<jx-button label="Clique" href="/pagina" bgcolor="#ff5d1f"></jx-button>' },
  { name: "jx-image",         snippet: '<jx-image src="https://..." objectfit="cover"></jx-image>' },
  { name: "jx-carousel",      snippet: '<!-- jx-carousel não disponível em Svelte (use store nativo) -->' },
  { name: "jx-divider",       snippet: '<jx-divider color="#dde3f0" thickness="1"></jx-divider>' },
  { name: "jx-video",         snippet: '<jx-video url="https://youtube.com/watch?v=..." height="320"></jx-video>' },
  { name: "jx-input-var",     snippet: '<jx-input-var placeholder="Digite..." inputtype="text"></jx-input-var>' },
  { name: "jx-variable-text", snippet: '<jx-variable-text text="Olá {{user_name}}!" fontsize="14"></jx-variable-text>' },
  { name: "jx-profile-card",  snippet: '<jx-profile-card name="João" subtitle="Dev" image="https://..."></jx-profile-card>' },
  { name: "jx-user-avatar",   snippet: '<jx-user-avatar src="https://..." size="48"></jx-user-avatar>' },
  { name: "jx-product-card",  snippet: '<jx-product-card siteid="SEU_SITE_ID" productid="PRODUTO_ID"></jx-product-card>' },
  { name: "jx-product-list",  snippet: '<jx-product-list siteid="SEU_SITE_ID" pagesize="6"></jx-product-list>' },
  { name: "jx-cart-items",    snippet: '<jx-cart-items storagekey="jx_cart_SEU_SITE_ID"></jx-cart-items>' },
  { name: "jx-admin-add-btn", snippet: '<jx-admin-add-btn label="Produto"></jx-admin-add-btn>' },
];

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
  const [jxOpen, setJxOpen]             = useState(false);

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
            <p className={styles.sideHead}>Biblioteca Jx</p>
            <button className={styles.addBtn} onClick={() => setJxOpen(v => !v)} title={jxOpen ? "Fechar" : "Abrir"}>
              {jxOpen ? "▲" : "▼"}
            </button>
          </div>
          {jxOpen && JX_SVELTE_SNIPPETS.map((c) => (
            <button
              key={c.name}
              className={styles.fileItem}
              title={`Copiar: ${c.snippet}`}
              onClick={() => navigator.clipboard.writeText(c.snippet).then(() => {}).catch(() => {})}
              style={{ fontFamily: "monospace", fontSize: 12 }}
            >
              {"<" + c.name + ">"}
            </button>
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
