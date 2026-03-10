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

// ── Jesterx Components Library (injected into every preview) ────────────────
const JX_REACT_LIB = `
// ── Jesterx Components Library ──────────────────────────────────────────────
// Disponíveis automaticamente em todos os componentes:
// JxHeading, JxParagraph, JxButton, JxImage, JxCarousel, JxInputVar,
// JxVariableText, JxProfileCard, JxProductCard, JxProductList,
// JxCartItems, JxDivider, JxVideo, JxUserAvatar, JxAdminAddBtn
// ────────────────────────────────────────────────────────────────────────────

function JxHeading({ text, fontSize, fontWeight, color, textAlign, style: s }) {
  var st = { fontSize: (fontSize||28)+'px', fontWeight: fontWeight||700, color: color||'#20283a', textAlign: textAlign||'left', margin: 0, ...(s||{}) };
  return <h2 style={st}>{text||'Título'}</h2>;
}

function JxParagraph({ text, fontSize, color, textAlign, style: s }) {
  var st = { fontSize: (fontSize||16)+'px', color: color||'#5a6379', textAlign: textAlign||'left', margin: 0, lineHeight: 1.6, ...(s||{}) };
  return <p style={st}>{text||'Parágrafo'}</p>;
}

function JxButton({ label, href, onClick, bgColor, textColor, borderRadius, fontSize, style: s }) {
  var st = { background: bgColor||'#ff5d1f', color: textColor||'#fff', border: 0, borderRadius: (borderRadius||8)+'px', padding: '10px 22px', fontSize: (fontSize||15)+'px', fontWeight: 600, cursor: 'pointer', fontFamily: 'inherit', display:'inline-block', ...(s||{}) };
  return <button style={st} onClick={onClick||(href?()=>window.open(href,'_self'):undefined)}>{label||'Botão'}</button>;
}

function JxImage({ src, alt, objectFit, borderRadius, style: s }) {
  var st = { width:'100%', height:'100%', objectFit: objectFit||'cover', borderRadius: (borderRadius||0)+'px', display:'block', ...(s||{}) };
  if (!src) return <div style={{background:'#f0f4ff',width:'100%',height:200,display:'flex',alignItems:'center',justifyContent:'center',borderRadius:8,color:'#9aa5bc',fontSize:13}}>🖼 Imagem</div>;
  return <img src={src} alt={alt||''} style={st} />;
}

function JxDivider({ color, thickness, marginTop, marginBottom, style: s }) {
  var st = { border:0, borderTop: (thickness||1)+'px solid '+(color||'#dde3f0'), margin: (marginTop||16)+'px 0 '+(marginBottom||16)+'px 0', ...(s||{}) };
  return <hr style={st} />;
}

function JxVideo({ url, height, style: s }) {
  function getEmbed(u) {
    if (!u) return '';
    var yt = u.match(/(?:youtube\\.com\\/watch\\?v=|youtu\\.be\\/)([^&?/]+)/);
    if (yt) return 'https://www.youtube.com/embed/' + yt[1];
    return u;
  }
  var embed = getEmbed(url);
  if (!embed) return <div style={{background:'#111',color:'#fff',width:'100%',height:(height||240)+'px',display:'flex',alignItems:'center',justifyContent:'center',borderRadius:8,...(s||{})}}>▶ Vídeo</div>;
  return <iframe src={embed} width="100%" height={height||240} frameBorder="0" allowFullScreen style={{display:'block',borderRadius:8,...(s||{})}} />;
}

function JxProfileCard({ name, subtitle, image, bgColor, style: s }) {
  var st = { background: bgColor||'#fff', border:'1px solid #dde3f0', borderRadius:12, padding:20, display:'flex', flexDirection:'column', alignItems:'center', gap:8, textAlign:'center', ...(s||{}) };
  return (
    <div style={st}>
      <img src={image||'https://placehold.co/80x80'} alt={name||''} style={{width:80,height:80,borderRadius:'50%',objectFit:'cover'}} />
      <div style={{fontWeight:700,fontSize:16,color:'#1a2740'}}>{name||'Nome'}</div>
      <div style={{fontSize:13,color:'#6b7a99'}}>{subtitle||'Subtítulo'}</div>
    </div>
  );
}

function JxInputVar({ placeholder, inputType, value, onChange, style: s }) {
  var st = { width:'100%', padding:'10px 14px', border:'1.5px solid #dde3f0', borderRadius:8, fontSize:14, fontFamily:'inherit', outline:'none', boxSizing:'border-box', ...(s||{}) };
  return <input type={inputType||'text'} placeholder={placeholder||'Digite aqui...'} value={value} onChange={onChange} style={st} />;
}

function JxVariableText({ text, vars, fontSize, color, style: s }) {
  var resolved = (text||'').replace(/\\{\\{(\\w+)\\}\\}/g, function(_,k){ return (vars&&vars[k])||''; });
  var st = { fontSize:(fontSize||14)+'px', color:color||'#20283a', ...(s||{}) };
  return <span style={st}>{resolved||'{{variavel}}'}</span>;
}

function JxUserAvatar({ src, size, objectFit, style: s }) {
  var sz = size||48;
  var st = { width:sz+'px', height:sz+'px', borderRadius:'50%', objectFit:objectFit||'cover', display:'block', ...(s||{}) };
  return <img src={src||'https://placehold.co/48x48'} alt="avatar" style={st} />;
}

function JxCarousel({ images, height, objectFit, style: s }) {
  var [idx, setIdx] = React.useState(0);
  var imgs = images||[];
  if (!imgs.length) return <div style={{background:'#f0f4ff',borderRadius:8,height:(height||220)+'px',display:'flex',alignItems:'center',justifyContent:'center',color:'#9aa5bc',fontSize:13}}>🖼 Carrossel</div>;
  return (
    <div style={{position:'relative',width:'100%',...(s||{})}}>
      <img src={imgs[idx]} alt="" style={{width:'100%',height:(height||220)+'px',objectFit:objectFit||'cover',borderRadius:8,display:'block'}} />
      {imgs.length>1 && (
        <div style={{position:'absolute',bottom:10,left:0,right:0,display:'flex',justifyContent:'center',gap:6}}>
          {imgs.map((_,i)=><button key={i} onClick={()=>setIdx(i)} style={{width:8,height:8,borderRadius:'50%',border:0,background:i===idx?'#ff5d1f':'rgba(255,255,255,.6)',padding:0,cursor:'pointer'}} />)}
        </div>
      )}
      {imgs.length>1 && <button onClick={()=>setIdx(i=>Math.max(0,i-1))} style={{position:'absolute',top:'50%',left:8,transform:'translateY(-50%)',background:'rgba(0,0,0,.35)',color:'#fff',border:0,borderRadius:6,width:28,height:28,cursor:'pointer',fontSize:14}}>‹</button>}
      {imgs.length>1 && <button onClick={()=>setIdx(i=>Math.min(imgs.length-1,i+1))} style={{position:'absolute',top:'50%',right:8,transform:'translateY(-50%)',background:'rgba(0,0,0,.35)',color:'#fff',border:0,borderRadius:6,width:28,height:28,cursor:'pointer',fontSize:14}}>›</button>}
    </div>
  );
}

function JxProductCard({ productId, siteId, style: s }) {
  var [prod, setProd] = React.useState(null);
  var [loading, setLoading] = React.useState(!!productId);
  React.useEffect(function(){
    if (!productId||!siteId){setLoading(false);return;}
    fetch('/api/store/'+siteId+'/products/'+productId).then(function(r){return r.ok?r.json():null;}).then(function(d){setProd(d&&d.data||null);setLoading(false);}).catch(function(){setLoading(false);});
  },[productId,siteId]);
  var box = { background:'#fff', border:'1px solid #dde3f0', borderRadius:12, overflow:'hidden', display:'flex', flexDirection:'column', ...(s||{}) };
  if (loading) return <div style={{...box,alignItems:'center',justifyContent:'center',padding:24,color:'#9aa5bc',fontSize:13}}>Carregando...</div>;
  if (!prod) return <div style={{...box,alignItems:'center',justifyContent:'center',padding:24,color:'#9aa5bc',fontSize:13}}>🛍 product_card</div>;
  return (
    <div style={box}>
      {prod.images&&prod.images[0]&&<img src={prod.images[0]} alt={prod.name} style={{width:'100%',height:180,objectFit:'cover'}}/>}
      <div style={{padding:'12px 16px',flex:1,display:'flex',flexDirection:'column',gap:6}}>
        <div style={{fontWeight:700,fontSize:15,color:'#1a2740'}}>{prod.name}</div>
        <div style={{fontSize:16,fontWeight:700,color:'#ff5d1f'}}>{'R$ '+Number(prod.price||0).toFixed(2)}</div>
        {prod.description&&<div style={{fontSize:12,color:'#6b7a99',lineHeight:1.5}}>{prod.description}</div>}
      </div>
    </div>
  );
}

function JxProductList({ siteId, pageSize, style: s }) {
  var [prods, setProds] = React.useState([]);
  var [page, setPage] = React.useState(0);
  var [loading, setLoading] = React.useState(true);
  var ps = Math.max(1, pageSize||6);
  React.useEffect(function(){
    if (!siteId){setLoading(false);return;}
    fetch('/api/store/'+siteId+'/products').then(function(r){return r.ok?r.json():null;}).then(function(d){setProds(d&&d.data||[]);setLoading(false);}).catch(function(){setLoading(false);});
  },[siteId]);
  if (loading) return <div style={{color:'#9aa5bc',padding:20,fontSize:13}}>Carregando produtos...</div>;
  if (!prods.length) return <div style={{color:'#9aa5bc',padding:20,fontSize:13}}>Nenhum produto encontrado.</div>;
  var total = Math.ceil(prods.length/ps);
  var items = prods.slice(page*ps,(page+1)*ps);
  return (
    <div style={{display:'flex',flexDirection:'column',gap:12,...(s||{})}}>
      <div style={{display:'grid',gridTemplateColumns:'repeat(auto-fill,minmax(160px,1fr))',gap:12}}>
        {items.map(function(p){return (
          <div key={p.id} style={{background:'#fff',border:'1px solid #dde3f0',borderRadius:10,overflow:'hidden'}}>
            {p.images&&p.images[0]&&<img src={p.images[0]} alt={p.name} style={{width:'100%',height:120,objectFit:'cover'}}/>}
            <div style={{padding:'8px 10px'}}>
              <div style={{fontWeight:600,fontSize:13,color:'#1a2740'}}>{p.name}</div>
              <div style={{fontSize:13,fontWeight:700,color:'#ff5d1f'}}>{'R$ '+Number(p.price||0).toFixed(2)}</div>
            </div>
          </div>
        );})}
      </div>
      {total>1&&(
        <div style={{display:'flex',gap:8,alignItems:'center',justifyContent:'center'}}>
          <button onClick={function(){setPage(function(pg){return Math.max(0,pg-1);});}} disabled={page===0} style={{padding:'4px 12px',cursor:'pointer',background:'#f4f6fb',border:'1px solid #dde3f0',borderRadius:6,fontFamily:'inherit'}}>← Anterior</button>
          <span style={{fontSize:12,color:'#6b7a99'}}>{(page+1)+' / '+total}</span>
          <button onClick={function(){setPage(function(pg){return Math.min(total-1,pg+1);});}} disabled={page>=total-1} style={{padding:'4px 12px',cursor:'pointer',background:'#f4f6fb',border:'1px solid #dde3f0',borderRadius:6,fontFamily:'inherit'}}>Próximo →</button>
        </div>
      )}
    </div>
  );
}

function JxCartItems({ storageKey, style: s }) {
  var sk = storageKey||('jx_cart_'+window.location.hostname);
  var [items, setItems] = React.useState(function(){ try{return JSON.parse(localStorage.getItem(sk)||'[]');}catch(e){return[];} });
  function save(c){ localStorage.setItem(sk,JSON.stringify(c)); setItems(c); document.dispatchEvent(new CustomEvent('jx:cartupdate')); }
  React.useEffect(function(){
    function h(){ try{setItems(JSON.parse(localStorage.getItem(sk)||'[]'));}catch(e){} }
    document.addEventListener('jx:cartupdate',h);
    return function(){ document.removeEventListener('jx:cartupdate',h); };
  },[sk]);
  if (!items.length) return <div style={{color:'#9aa5bc',padding:'20px',textAlign:'center',fontSize:13}}>Carrinho vazio</div>;
  return (
    <div style={{display:'flex',flexDirection:'column',gap:8,...(s||{})}}>
      {items.map(function(it){ return (
        <div key={it.id} style={{display:'flex',alignItems:'center',gap:8,background:'#fff',border:'1px solid #dbe2f3',borderRadius:8,padding:'8px 10px'}}>
          <div style={{flex:1,fontWeight:600,fontSize:13,color:'#1a2740',overflow:'hidden'}}>{it.name}</div>
          <div style={{display:'flex',alignItems:'center',gap:4,flexShrink:0}}>
            <button onClick={function(){save(items.map(function(x){return x.id===it.id?{...x,qty:Math.max(0,(x.qty||1)-1)}:x;}).filter(function(x){return (x.qty||0)>0;}));}} style={{width:26,height:26,border:'1px solid #ccd5e8',background:'#f4f6fb',borderRadius:6,cursor:'pointer',fontWeight:700,fontFamily:'inherit',fontSize:15}}>−</button>
            <span style={{minWidth:22,textAlign:'center',fontWeight:700,fontSize:13}}>{it.qty||1}</span>
            <button onClick={function(){save(items.map(function(x){return x.id===it.id?{...x,qty:(x.qty||1)+1}:x;}));}} style={{width:26,height:26,border:'1px solid #ccd5e8',background:'#f4f6fb',borderRadius:6,cursor:'pointer',fontWeight:700,fontFamily:'inherit',fontSize:15}}>+</button>
          </div>
          <div style={{fontWeight:700,color:'#ff5d1f',fontSize:13,minWidth:64,textAlign:'right',flexShrink:0}}>{'R$ '+(Number(it.price||0)*(it.qty||1)).toFixed(2)}</div>
          <button onClick={function(){save(items.filter(function(x){return x.id!==it.id;}));}} style={{background:'none',border:0,color:'#9aa5bc',cursor:'pointer',fontSize:18,lineHeight:1,padding:'0 4px',flexShrink:0}}>×</button>
        </div>
      );})}
    </div>
  );
}

function JxAdminAddBtn({ label, onAdd, style: s }) {
  var st = { background:'#1a2740', color:'#fff', border:0, borderRadius:8, padding:'8px 18px', fontSize:13, fontWeight:600, cursor:'pointer', fontFamily:'inherit', display:'inline-flex', alignItems:'center', gap:6, ...(s||{}) };
  return <button style={st} onClick={onAdd}>＋ {label||'Adicionar'}</button>;
}
// ── End Jesterx Components Library ──────────────────────────────────────────
`;

// ── Preview builder ────────────────────────────────────────────────────────
function buildPreview(bundle: Bundle, activePath: string) {
  const routeComp = bundle.routes[activePath]?.component ?? DEFAULT_COMPONENT;
  const safeCode = [JX_REACT_LIB, bundle.header, routeComp, bundle.footer]
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

const JX_REACT_COMPONENTS: { name: string; snippet: string }[] = [
  { name: "JxHeading",      snippet: '<JxHeading text="Título" fontSize={32} color="#20283a" />' },
  { name: "JxParagraph",    snippet: '<JxParagraph text="Seu texto aqui." fontSize={16} />' },
  { name: "JxButton",       snippet: '<JxButton label="Clique aqui" href="/pagina" bgColor="#ff5d1f" />' },
  { name: "JxImage",        snippet: '<JxImage src="https://..." objectFit="cover" borderRadius={8} />' },
  { name: "JxCarousel",     snippet: '<JxCarousel images={["url1","url2","url3"]} height={240} />' },
  { name: "JxDivider",      snippet: '<JxDivider color="#dde3f0" thickness={1} />' },
  { name: "JxVideo",        snippet: '<JxVideo url="https://youtube.com/watch?v=..." height={320} />' },
  { name: "JxInputVar",     snippet: '<JxInputVar placeholder="Digite..." inputType="text" />' },
  { name: "JxVariableText", snippet: '<JxVariableText text="Olá {{user_name}}!" vars={{ user_name: "Visitante" }} />' },
  { name: "JxProfileCard",  snippet: '<JxProfileCard name="João" subtitle="Dev" image="https://..." />' },
  { name: "JxUserAvatar",   snippet: '<JxUserAvatar src="https://..." size={48} />' },
  { name: "JxProductCard",  snippet: '<JxProductCard siteId="SEU_SITE_ID" productId="PRODUTO_ID" />' },
  { name: "JxProductList",  snippet: '<JxProductList siteId="SEU_SITE_ID" pageSize={6} />' },
  { name: "JxCartItems",    snippet: '<JxCartItems storageKey="jx_cart_SEU_SITE_ID" />' },
  { name: "JxAdminAddBtn",  snippet: '<JxAdminAddBtn label="Produto" onAdd={() => alert("add")} />' },
];

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
      const msg = err instanceof Error ? err.message : "Erro ao carregar";
      if (msg.includes("403") || msg.toLowerCase().includes("forbidden") || msg.toLowerCase().includes("unauthorized")) {
        navigate("/pages", { replace: true });
        return;
      }
      setError(msg);
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
            <p className={styles.sideHead}>Biblioteca Jx</p>
            <button className={styles.addBtn} onClick={() => setJxOpen(v => !v)} title={jxOpen ? "Fechar" : "Abrir"}>
              {jxOpen ? "▲" : "▼"}
            </button>
          </div>
          {jxOpen && JX_REACT_COMPONENTS.map((c) => (
            <button
              key={c.name}
              className={styles.fileItem}
              title={`Copiar: ${c.snippet}`}
              onClick={() => navigator.clipboard.writeText(c.snippet).then(() => {}).catch(() => {})}
              style={{ fontFamily: "monospace", fontSize: 12 }}
            >
              {"<" + c.name + " />"}
            </button>
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
