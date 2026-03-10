import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { API_URL, apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./ElementorEditor.module.scss";

// ─── Types ─────────────────────────────────────────────────────────────────────
type BlockType =
  | "heading" | "paragraph" | "button" | "image" | "carousel"
  | "input_var" | "variable_text" | "profile_card" | "product_card"
  | "video" | "divider" | "product_list" | "user_avatar";

type StyleMap = Record<string, string>;

type Block = {
  id: string; type: BlockType; text?: string; label?: string; href?: string;
  action_type?: "navigate" | "call_api" | "store_login" | "store_logout" | "add_to_cart" | "open_popup" | "close_popup" | "store_register" | "add_product";
  action_target?: string; api_id?: string;
  src?: string; object_fit?: "cover" | "contain" | "fill"; images?: string[]; var_src?: string;
  var_name?: string; placeholder?: string; profile_name?: string;
  profile_subtitle?: string; profile_image?: string; video_url?: string;
  admin_only?: boolean; btn_action_type?: "add_product" | "add_video";
  product_id?: string; page_size?: number;
  popup_id?: string; email_var?: string; password_var?: string; register_mode?: boolean;
  first_name_var?: string; last_name_var?: string; input_type?: string;
  inner_blocks?: Block[]; style: StyleMap;x: number; y: number; w: number; h: number;
  rotation: number; z: number;
};

type PageDoc = { title: string; blocks: Block[] };
type PopupDoc = { title: string; blocks: Block[]; width: number; height: number; background: string };
type GlobalSection = { enabled: boolean; height: number; background: string; blocks: Block[] };
type BuilderDoc = {
  title: string;
  canvas: { width: number; height: number; background: string };
  header: GlobalSection;
  footer: GlobalSection;
  pages: Record<string, PageDoc>;
  popups: Record<string, PopupDoc>;
};
type LegacyDoc = { title?: string; blocks?: Block[] };

type RouteItem = { id: string; path: string; title: string; requires_auth: boolean; position: number };
type RoutesResponse = { success: boolean; message: string; data: RouteItem[] };
type VersionItem = { id: string; version: number; source_type: "JXML" | "REACT" | "SVELTE" | "ELEMENTOR_JSON"; source?: string; scan_status: "clean" | "warning" | "blocked" };
type VersionsResponse = { success: boolean; message: string; data: VersionItem[] };
type CreateVersionResponse = { success: boolean; message: string; data: VersionItem };
type SiteAPI = { id: string; method: string; path: string; label: string; description: string };
type SiteAPIsResponse = { success: boolean; message: string; data: SiteAPI[] };

type GuideState = { vertical: number | null; horizontal: number | null };
type EditSection = "page" | "header" | "footer";
type LeftTab = "elements" | "layers" | "settings" | "templates" | "popups";

type Interaction =
  | { mode: "drag"; id: string; section: EditSection; dx: number; dy: number }
  | { mode: "resize"; id: string; section: EditSection; edge: "right" | "bottom" | "left" | "top" | "corner"; startX: number; startY: number; start: { x: number; y: number; w: number; h: number }; baseFontSize?: number }
  | { mode: "rotate"; id: string; section: EditSection; cx: number; cy: number; startAngle: number; startRotation: number };

// ─── Constants ─────────────────────────────────────────────────────────────────
const ZOOM_LEVELS = [0.4, 0.5, 0.65, 0.75, 1.0, 1.25];

const BLOCK_CATEGORIES = [
  { label: "Texto", items: [{ type: "heading" as BlockType, icon: "🔤", label: "Titulo" }, { type: "paragraph" as BlockType, icon: "📝", label: "Paragrafo" }, { type: "variable_text" as BlockType, icon: "��", label: "Texto Variavel" }] },
  { label: "Midia", items: [{ type: "image" as BlockType, icon: "🖼", label: "Imagem" }, { type: "carousel" as BlockType, icon: "🎠", label: "Carousel" }, { type: "video" as BlockType, icon: "▶", label: "Video" }] },
  { label: "Interativo", items: [{ type: "button" as BlockType, icon: "🔘", label: "Botao" }, { type: "input_var" as BlockType, icon: "📌", label: "Input Variavel" }] },
  { label: "Loja", items: [{ type: "product_card" as BlockType, icon: "🛍", label: "Product Card" }, { type: "product_list" as BlockType, icon: "📦", label: "Lista de Produtos" }] },
  { label: "Layout", items: [{ type: "profile_card" as BlockType, icon: "👤", label: "Perfil" }, { type: "user_avatar" as BlockType, icon: "🙂", label: "Avatar Usuario" }, { type: "divider" as BlockType, icon: "➖", label: "Divider" }] },
];

// ─── Helpers ───────────────────────────────────────────────────────────────────
function isTextBlock(type: BlockType): boolean {
  return type === "heading" || type === "paragraph" || type === "variable_text";
}

const DEFAULT_STYLE: StyleMap = {
  color: "#1f2b43", background: "transparent", padding: "10px 12px",
  border: "0", "border-radius": "0", "font-size": "16px", "font-weight": "500", "text-align": "left",
};

function nextId(prefix: string): string {
  return `${prefix}-${Date.now()}-${Math.floor(Math.random() * 9999)}`;
}

function newBlock(type: BlockType, z: number): Block {
  if (type === "heading") return { id: nextId("blk"), type, text: "Titulo principal", style: { ...DEFAULT_STYLE, "font-size": "48px", "font-weight": "700" }, x: 120, y: 90, w: 620, h: 90, rotation: 0, z };
  if (type === "paragraph") return { id: nextId("blk"), type, text: "Paragrafo com descricao da sua oferta.", style: { ...DEFAULT_STYLE, color: "#4f5f83" }, x: 120, y: 190, w: 580, h: 80, rotation: 0, z };
  if (type === "button") return { id: nextId("blk"), type, label: "Botao", href: "/", action_type: "navigate", action_target: "/", style: { ...DEFAULT_STYLE, background: "#ff5d1f", color: "#ffffff", padding: "12px 18px", "text-align": "center" }, x: 120, y: 290, w: 190, h: 54, rotation: 0, z };
  if (type === "image") return { id: nextId("blk"), type, src: "https://images.unsplash.com/photo-1512436991641-6745cdb1723f?auto=format&fit=crop&w=1000&q=80", style: { ...DEFAULT_STYLE, padding: "0", background: "#e6ebf5" }, x: 780, y: 100, w: 320, h: 240, rotation: 0, z };
  if (type === "carousel") return { id: nextId("blk"), type, images: ["https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=1200&q=80", "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?auto=format&fit=crop&w=1200&q=80"], style: { ...DEFAULT_STYLE, padding: "0", background: "#dce4f3" }, x: 120, y: 370, w: 680, h: 250, rotation: 0, z };
  if (type === "input_var") return { id: nextId("blk"), type, var_name: "nome", placeholder: "Digite seu nome", style: { ...DEFAULT_STYLE, border: "1px solid #cad3e7", background: "#ffffff" }, x: 120, y: 650, w: 340, h: 56, rotation: 0, z };
  if (type === "variable_text") return { id: nextId("blk"), type, text: "Ola, {{nome}}", style: { ...DEFAULT_STYLE, "font-size": "26px", "font-weight": "700" }, x: 490, y: 648, w: 360, h: 62, rotation: 0, z };
  if (type === "profile_card") return { id: nextId("blk"), type, profile_name: "Nome do usuario", profile_subtitle: "Cliente Premium", profile_image: "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=400&q=80", style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#ffffff", padding: "16px" }, x: 860, y: 370, w: 280, h: 210, rotation: 0, z };
  if (type === "product_card") return { id: nextId("blk"), type, style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#ffffff", padding: "12px" }, x: 860, y: 610, w: 300, h: 280, rotation: 0, z };
  if (type === "video") return { id: nextId("blk"), type, video_url: "", style: { ...DEFAULT_STYLE, padding: "0", background: "#000" }, x: 120, y: 400, w: 640, h: 360, rotation: 0, z };
  if (type === "user_avatar") return { id: nextId("blk"), type, object_fit: "cover", style: { ...DEFAULT_STYLE, padding: "0", background: "#e6ebf5", "border-radius": "50%" }, x: 860, y: 100, w: 100, h: 100, rotation: 0, z };
  if (type === "product_list") return { id: nextId("blk"), type, page_size: 6, style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#f8f9ff", padding: "8px" }, x: 80, y: 400, w: 900, h: 480, rotation: 0, z };
  return { id: nextId("blk"), type: "divider", style: { ...DEFAULT_STYLE, background: "#cad3e7", padding: "0" }, x: 120, y: 740, w: 720, h: 2, rotation: 0, z };
}

function ensurePath(path: string): string {
  const trimmed = path.trim();
  if (!trimmed) return "/";
  if (trimmed.startsWith("/")) return trimmed;
  return `/${trimmed}`;
}

function parseDoc(source: string | undefined, routes: RouteItem[]): BuilderDoc | null {
  if (!source) return null;
  const defaultHeader: GlobalSection = { enabled: false, height: 80, background: "#1a2740", blocks: [] };
  const defaultFooter: GlobalSection = { enabled: false, height: 100, background: "#1a2740", blocks: [] };
  const routePaths = routes.length > 0 ? routes.map(r => ensurePath(r.path)) : ["/"];
  try {
    const parsed = JSON.parse(source) as BuilderDoc & LegacyDoc;
    const pages: Record<string, PageDoc> = {};
    if (parsed.pages && typeof parsed.pages === "object") {
      Object.entries(parsed.pages).forEach(([path, page]) => {
        pages[ensurePath(path)] = { title: page.title || "Pagina", blocks: (page.blocks || []).map((block, idx) => ({ ...block, id: block.id || nextId("blk"), style: { ...DEFAULT_STYLE, ...(block.style || {}) }, z: typeof block.z === "number" ? block.z : idx + 1, w: block.w > 0 ? block.w : 220, h: block.h > 0 ? block.h : 90 })) };
      });
    } else if (Array.isArray(parsed.blocks)) {
      pages["/"] = { title: parsed.title || "Pagina Inicial", blocks: parsed.blocks.map((block, idx) => ({ ...block, id: block.id || nextId("blk"), style: { ...DEFAULT_STYLE, ...(block.style || {}) }, z: typeof block.z === "number" ? block.z : idx + 1, w: block.w > 0 ? block.w : 220, h: block.h > 0 ? block.h : 90 })) };
    }
    routePaths.forEach(path => { if (!pages[path]) pages[path] = { title: path === "/" ? "Pagina Inicial" : `Pagina ${path}`, blocks: [] }; });
    return {
      title: parsed.title || "Site",
      canvas: { width: Math.min(2800, Math.max(900, parsed.canvas?.width || 1400)), height: Math.min(2800, Math.max(700, parsed.canvas?.height || 900)), background: parsed.canvas?.background || "#f8f9ff" },
      header: parsed.header || defaultHeader,
      footer: parsed.footer || defaultFooter,
      pages,
      popups: (parsed as any).popups || {},
    };
  } catch { return null; }
}

function readFileAsDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result || ""));
    reader.onerror = () => reject(new Error("Falha ao ler arquivo"));
    reader.readAsDataURL(file);
  });
}

function computeSnap(x: number, y: number, w: number, h: number, movingID: string, blocks: Block[], canvasWidth: number, canvasHeight: number): { x: number; y: number; guides: GuideState } {
  const threshold = 8;
  const pointsX = [x, x + w / 2, x + w];
  const pointsY = [y, y + h / 2, y + h];
  const targetX: number[] = [canvasWidth / 2];
  const targetY: number[] = [canvasHeight / 2];
  blocks.forEach(block => {
    if (block.id === movingID) return;
    targetX.push(block.x, block.x + block.w / 2, block.x + block.w);
    targetY.push(block.y, block.y + block.h / 2, block.y + block.h);
  });
  let snappedX = x, snappedY = y, guideX: number | null = null, guideY: number | null = null;
  let bestDx = threshold + 1;
  pointsX.forEach(point => { targetX.forEach(target => { const abs = Math.abs(target - point); if (abs < bestDx && abs <= threshold) { bestDx = abs; snappedX = x + (target - point); guideX = target; } }); });
  let bestDy = threshold + 1;
  pointsY.forEach(point => { targetY.forEach(target => { const abs = Math.abs(target - point); if (abs < bestDy && abs <= threshold) { bestDy = abs; snappedY = y + (target - point); guideY = target; } }); });
  return { x: Math.round(snappedX), y: Math.round(snappedY), guides: { vertical: guideX, horizontal: guideY } };
}

function getBlockIcon(type: BlockType): string {
  const icons: Record<string, string> = { heading: "🔤", paragraph: "📝", button: "🔘", image: "🖼", carousel: "🎠", input_var: "📌", variable_text: "💬", profile_card: "👤", product_card: "🛍", product_list: "📦", video: "▶", divider: "➖", user_avatar: "🙂" };
  return icons[type] ?? "▪";
}

function renderBlock(block: Block) {
  if (block.type === "heading") return <h1 style={{ fontSize: block.style["font-size"], fontWeight: block.style["font-weight"], margin: 0, lineHeight: 1.15, overflow: "hidden" }}>{block.text}</h1>;
  if (block.type === "paragraph") return <p style={{ fontSize: block.style["font-size"], margin: 0, lineHeight: 1.45, overflow: "hidden" }}>{block.text}</p>;
  if (block.type === "button") return <span className={styles.inlineButton}>{block.label || "Botao"}</span>;
  if (block.type === "user_avatar") return <div className={styles.imagePlaceholder}>🙂 Avatar do Usuário</div>;
  if (block.type === "image") { const fit = block.object_fit || "cover"; if (block.var_src) return <div className={styles.imagePlaceholder}>🖼 {block.var_src}</div>; return block.src ? <img src={block.src} alt="imagem" style={{ width: "100%", height: "100%", objectFit: fit, display: "block" }} /> : <div className={styles.imagePlaceholder}>🖼 Imagem</div>; }
  if (block.type === "carousel") { const first = block.images?.[0]; return first ? <img src={first} alt="carousel" className={styles.image} /> : <div className={styles.imagePlaceholder}>🖼 Carousel</div>; }
  if (block.type === "input_var") return <div className={styles.inputVarWrap}><input className={styles.previewInput} placeholder={block.placeholder || "Digite"} readOnly />{block.var_name && <span className={styles.varBadge}>📌 {block.var_name}</span>}</div>;
  if (block.type === "variable_text") return <p style={{ fontSize: block.style["font-size"], margin: 0, overflow: "hidden" }}>{block.text || "Texto dinamico {{var}}"}</p>;
  if (block.type === "profile_card") return <div className={styles.profileCard}>{block.profile_image && <img src={block.profile_image} alt="perfil" />}<h4>{block.profile_name || "Nome"}</h4><p>{block.profile_subtitle || "Descricao"}</p></div>;
  if (block.type === "product_card") return <div className={styles.productCard}><strong>🛍 Product Card</strong><small>{block.product_id ? `Produto: ${block.product_id.slice(0,8)}…` : "Exibe 1º produto da loja"}</small></div>;
  if (block.type === "product_list") return <div className={styles.productCard}><strong>📦 Lista de Produtos</strong><small>{`${block.page_size || 6} por página · ordenado por + vendidos`}</small></div>;
  if (block.type === "video") {
    const url = block.video_url || "";
    const ytMatch = url.match(/(?:youtu\.be\/|youtube\.com\/(?:watch\?v=|embed\/))([\w-]+)/);
    const embedSrc = ytMatch ? `https://www.youtube.com/embed/${ytMatch[1]}` : url;
    return url ? <iframe src={embedSrc} allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowFullScreen style={{ width: "100%", height: "100%", border: "none", display: "block" }} title="video" /> : <div className={styles.videoPlaceholder}>▶ Adicione uma URL de video</div>;
  }
  return <div className={styles.dividerPreview} />;
}

// ─── Landing Template ──────────────────────────────────────────────────────────
function makeLandingTemplate(): BuilderDoc {
  return {
    title: "Site",
    canvas: { width: 1400, height: 900, background: "#f8f9ff" },
    popups: {},
    header: {
      enabled: true, height: 80, background: "#1a2740",
      blocks: [
        { id: nextId("blk"), type: "heading", text: "MyBrand", style: { ...DEFAULT_STYLE, "font-size": "24px", "font-weight": "700", color: "#ffffff", background: "transparent", padding: "8px 0" }, x: 30, y: 20, w: 200, h: 44, rotation: 0, z: 1 },
        { id: nextId("blk"), type: "button", label: "Inicio", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "transparent", color: "#ffffff", padding: "10px 16px", border: "0", "border-radius": "6px", "text-align": "center" }, x: 580, y: 18, w: 100, h: 44, rotation: 0, z: 2 },
        { id: nextId("blk"), type: "button", label: "Sobre", action_type: "navigate", action_target: "/sobre", href: "/sobre", style: { ...DEFAULT_STYLE, background: "transparent", color: "#ffffff", padding: "10px 16px", border: "0", "border-radius": "6px", "text-align": "center" }, x: 690, y: 18, w: 100, h: 44, rotation: 0, z: 3 },
        { id: nextId("blk"), type: "button", label: "Comecar gratis", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "#ff5d1f", color: "#ffffff", padding: "10px 20px", "border-radius": "8px", border: "0", "font-weight": "600", "text-align": "center" }, x: 1150, y: 18, w: 170, h: 44, rotation: 0, z: 4 },
      ],
    },
    footer: {
      enabled: true, height: 100, background: "#1a2740",
      blocks: [
        { id: nextId("blk"), type: "divider", style: { ...DEFAULT_STYLE, background: "#2d4060", padding: "0" }, x: 0, y: 0, w: 1400, h: 2, rotation: 0, z: 1 },
        { id: nextId("blk"), type: "paragraph", text: "© 2025 MyBrand. Todos os direitos reservados.", style: { ...DEFAULT_STYLE, color: "#7a8fa6", "font-size": "13px", background: "transparent", padding: "0" }, x: 30, y: 35, w: 500, h: 40, rotation: 0, z: 2 },
      ],
    },
    pages: {
      "/": {
        title: "Pagina Inicial",
        blocks: [
          { id: nextId("blk"), type: "heading", text: "Transforme seu negocio digital", style: { ...DEFAULT_STYLE, "font-size": "58px", "font-weight": "700", color: "#1a2740", background: "transparent", padding: "0" }, x: 80, y: 60, w: 680, h: 120, rotation: 0, z: 1 },
          { id: nextId("blk"), type: "paragraph", text: "Crie paginas incriveis com nosso editor visual. Simples, rapido e poderoso.", style: { ...DEFAULT_STYLE, color: "#4f5f83", "font-size": "18px", background: "transparent", padding: "0" }, x: 80, y: 210, w: 600, h: 80, rotation: 0, z: 2 },
          { id: nextId("blk"), type: "button", label: "Comecar agora", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "#ff5d1f", color: "#ffffff", padding: "14px 28px", "border-radius": "8px", border: "0", "font-weight": "600", "font-size": "16px", "text-align": "center" }, x: 80, y: 320, w: 200, h: 54, rotation: 0, z: 3 },
          { id: nextId("blk"), type: "button", label: "Ver demonstracao", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "transparent", color: "#1a2740", padding: "12px 28px", "border-radius": "8px", border: "2px solid #1a2740", "font-weight": "600", "font-size": "16px", "text-align": "center" }, x: 300, y: 320, w: 230, h: 54, rotation: 0, z: 4 },
          { id: nextId("blk"), type: "image", src: "https://images.unsplash.com/photo-1537432376769-00f5c2f4c8d2?auto=format&fit=crop&w=1000&q=80", object_fit: "cover", style: { ...DEFAULT_STYLE, padding: "0", background: "#e6ebf5", "border-radius": "16px" }, x: 820, y: 40, w: 500, h: 380, rotation: 0, z: 5 },
          { id: nextId("blk"), type: "divider", style: { ...DEFAULT_STYLE, background: "#e2e8f3", padding: "0" }, x: 0, y: 500, w: 1400, h: 2, rotation: 0, z: 6 },
          { id: nextId("blk"), type: "heading", text: "Por que escolher nos?", style: { ...DEFAULT_STYLE, "font-size": "36px", "font-weight": "700", color: "#1a2740", background: "transparent", padding: "0", "text-align": "center" }, x: 80, y: 540, w: 1240, h: 60, rotation: 0, z: 7 },
          { id: nextId("blk"), type: "profile_card", profile_name: "Ana Silva", profile_subtitle: "Designer Senior", profile_image: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?auto=format&fit=crop&w=400&q=80", style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#ffffff", padding: "20px", "border-radius": "12px" }, x: 80, y: 640, w: 280, h: 210, rotation: 0, z: 8 },
          { id: nextId("blk"), type: "profile_card", profile_name: "Carlos Mota", profile_subtitle: "Dev Full Stack", profile_image: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=400&q=80", style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#ffffff", padding: "20px", "border-radius": "12px" }, x: 400, y: 640, w: 280, h: 210, rotation: 0, z: 9 },
          { id: nextId("blk"), type: "profile_card", profile_name: "Mariana Costa", profile_subtitle: "Product Manager", profile_image: "https://images.unsplash.com/photo-1494790108755-2616b612b77c?auto=format&fit=crop&w=400&q=80", style: { ...DEFAULT_STYLE, border: "1px solid #dde4f4", background: "#ffffff", padding: "20px", "border-radius": "12px" }, x: 720, y: 640, w: 280, h: 210, rotation: 0, z: 10 },
        ],
      },
    },
  };
}

// ─── Component ─────────────────────────────────────────────────────────────────
export const ElementorEditor = () => {
  const { siteId = "" } = useParams();
  const { websiteId } = useAuthContext();
  const navigate = useNavigate();

  const pageCanvasRef = useRef<HTMLDivElement | null>(null);
  const headerCanvasRef = useRef<HTMLDivElement | null>(null);
  const footerCanvasRef = useRef<HTMLDivElement | null>(null);

  const [doc, setDoc] = useState<BuilderDoc>({
    title: "Site",
    canvas: { width: 2000, height: 900, background: "#f8f9ff" },
    header: { enabled: false, height: 80, background: "#1a2740", blocks: [] },
    footer: { enabled: false, height: 100, background: "#1a2740", blocks: [] },
    pages: { "/": { title: "Pagina Inicial", blocks: [] } },
    popups: {},
  });
  const docRef = useRef<BuilderDoc>(doc);
  docRef.current = doc;

  const [routes, setRoutes] = useState<RouteItem[]>([{ id: "root", path: "/", title: "Inicio", requires_auth: false, position: 0 }]);
  const activeRouteRef = useRef("/");
  const [activeRoute, setActiveRoute] = useState("/");
  activeRouteRef.current = activeRoute;

  const [selected, setSelected] = useState<string | null>(null);
  const [editingPopup, setEditingPopup] = useState<string | null>(null);
  const editingPopupRef = useRef<string | null>(null);
  editingPopupRef.current = editingPopup;
  const [editingProductCard, setEditingProductCard] = useState<string | null>(null);
  const [editingProductCardSection, setEditingProductCardSection] = useState<EditSection>("page");
  const editingProductCardRef = useRef<string | null>(null);
  const editingProductCardSectionRef = useRef<EditSection>("page");
  editingProductCardRef.current = editingProductCard;
  editingProductCardSectionRef.current = editingProductCardSection;
  const [editSection, setEditSection] = useState<EditSection>("page");
  const [leftTab, setLeftTab] = useState<LeftTab>("elements");
  const [zoom, setZoom] = useState(1.0);
  const [interaction, setInteraction] = useState<Interaction | null>(null);
  const [guides, setGuides] = useState<GuideState>({ vertical: null, horizontal: null });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [apis, setApis] = useState<SiteAPI[]>([]);
  const [menu, setMenu] = useState<{ x: number; y: number; id: string } | null>(null);
  const [apiTestResult, setApiTestResult] = useState<string | null>(null);
  const [apiTesting, setApiTesting] = useState(false);
  const [showTemplateModal, setShowTemplateModal] = useState(false);
  const [canUndo, setCanUndo] = useState(false);
  const [canRedo, setCanRedo] = useState(false);

  const undoStack = useRef<BuilderDoc[]>([]);
  const redoStack = useRef<BuilderDoc[]>([]);

  // ─ Derived ──────────────────────────────────────────────────────────────────
  const currentBlocks = useMemo(() => {
    if (editingPopup) return doc.popups[editingPopup]?.blocks || [];
    if (editingProductCard) {
      const pcSec = editingProductCardSection;
      const src = pcSec === "header" ? doc.header.blocks : pcSec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
      return src.find(b => b.id === editingProductCard)?.inner_blocks || [];
    }
    if (editSection === "header") return doc.header.blocks;
    if (editSection === "footer") return doc.footer.blocks;
    return doc.pages[activeRoute]?.blocks || [];
  }, [editSection, doc, activeRoute, editingPopup, editingProductCard, editingProductCardSection]);

  const pcCardBlock = useMemo(() => {
    if (!editingProductCard) return null;
    const pcSec = editingProductCardSection;
    const src = pcSec === "header" ? doc.header.blocks : pcSec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
    return src.find(b => b.id === editingProductCard) ?? null;
  }, [editingProductCard, editingProductCardSection, doc, activeRoute]);

  const selectedBlock = useMemo(() => currentBlocks.find(b => b.id === selected) || null, [currentBlocks, selected]);

  // ─ History ──────────────────────────────────────────────────────────────────
  function pushHistory() {
    undoStack.current = [...undoStack.current, docRef.current].slice(-30);
    redoStack.current = [];
    setCanUndo(true);
    setCanRedo(false);
  }

  function undo() {
    if (undoStack.current.length === 0) return;
    const prev = undoStack.current[undoStack.current.length - 1];
    redoStack.current = [docRef.current, ...redoStack.current];
    undoStack.current = undoStack.current.slice(0, -1);
    setDoc(prev);
    setCanUndo(undoStack.current.length > 0);
    setCanRedo(true);
  }

  function redo() {
    if (redoStack.current.length === 0) return;
    const next = redoStack.current[0];
    undoStack.current = [...undoStack.current, docRef.current];
    redoStack.current = redoStack.current.slice(1);
    setDoc(next);
    setCanUndo(true);
    setCanRedo(redoStack.current.length > 0);
  }

  // ─ Block mutation ────────────────────────────────────────────────────────────
  function clampBlock(block: Block, cw: number, ch: number): Block {
    const w = Math.min(block.w, cw);
    const h = Math.min(block.h, ch);
    const x = Math.max(0, Math.min(block.x, cw - w));
    const y = Math.max(0, Math.min(block.y, ch - h));
    return { ...block, x, y, w, h };
  }

  function updateBlock(section: EditSection, id: string, patch: Partial<Block>) {
    setDoc(prev => {
      const pid = editingPopupRef.current;
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: (prev.popups[pid]?.blocks || []).map(b => b.id === id ? { ...b, ...patch } : b) } } };
      const pcId = editingProductCardRef.current;
      if (pcId) { const pcSec = editingProductCardSectionRef.current; const mp = (bl: Block[]) => bl.map(b => b.id === pcId ? { ...b, inner_blocks: (b.inner_blocks || []).map(ib => ib.id === id ? { ...ib, ...patch } : ib) } : b); if (pcSec === "header") return { ...prev, header: { ...prev.header, blocks: mp(prev.header.blocks) } }; if (pcSec === "footer") return { ...prev, footer: { ...prev.footer, blocks: mp(prev.footer.blocks) } }; const r = activeRouteRef.current; return { ...prev, pages: { ...prev.pages, [r]: { ...prev.pages[r], blocks: mp(prev.pages[r]?.blocks || []) } } }; }
      if (section === "header") return { ...prev, header: { ...prev.header, blocks: prev.header.blocks.map(b => b.id === id ? { ...b, ...patch } : b) } };
      if (section === "footer") return { ...prev, footer: { ...prev.footer, blocks: prev.footer.blocks.map(b => b.id === id ? { ...b, ...patch } : b) } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: (prev.pages[route]?.blocks || []).map(b => b.id === id ? { ...b, ...patch } : b) } } };
    });
  }

  function addBlock(type: BlockType) {
    const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0);
    let block = newBlock(type, maxZ + 1);
    // Clamp to canvas dimensions so block is always visible on add
    const pid = editingPopupRef.current;
    const pcIdClamp = editingProductCardRef.current;
    if (pcIdClamp) {
      const pcSecC = editingProductCardSectionRef.current;
      const srcC = pcSecC === "header" ? doc.header.blocks : pcSecC === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
      const pcBl = srcC.find(b => b.id === pcIdClamp);
      block = clampBlock(block, pcBl?.w || 300, pcBl?.h || 280);
    } else if (pid) {
      const popup = doc.popups[pid];
      block = clampBlock(block, popup?.width || 480, popup?.height || 560);
    } else if (editSection === "header") {
      block = clampBlock(block, doc.canvas.width, doc.header.height || 80);
    } else if (editSection === "footer") {
      block = clampBlock(block, doc.canvas.width, doc.footer.height || 100);
    }
    pushHistory();
    setDoc(prev => {
      const pid = editingPopupRef.current;
      const pcIdAdd = editingProductCardRef.current;
      if (pcIdAdd) { const pcSecA = editingProductCardSectionRef.current; const mpA = (bl: Block[]) => bl.map(b => b.id === pcIdAdd ? { ...b, inner_blocks: [...(b.inner_blocks || []), block] } : b); if (pcSecA === "header") return { ...prev, header: { ...prev.header, blocks: mpA(prev.header.blocks) } }; if (pcSecA === "footer") return { ...prev, footer: { ...prev.footer, blocks: mpA(prev.footer.blocks) } }; const rA = activeRouteRef.current; return { ...prev, pages: { ...prev.pages, [rA]: { ...prev.pages[rA], blocks: mpA(prev.pages[rA]?.blocks || []) } } }; }
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: [...(prev.popups[pid]?.blocks || []), block] } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: [...prev.header.blocks, block] } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: [...prev.footer.blocks, block] } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: [...(prev.pages[route]?.blocks || []), block] } } };
    });
    setSelected(block.id);
  }

  function removeBlock(id: string) {
    pushHistory();
    setDoc(prev => {
      const pid = editingPopupRef.current;
      const pcIdRm = editingProductCardRef.current;
      if (pcIdRm) { const pcSecR = editingProductCardSectionRef.current; const mpR = (bl: Block[]) => bl.map(b => b.id === pcIdRm ? { ...b, inner_blocks: (b.inner_blocks || []).filter(ib => ib.id !== id) } : b); if (pcSecR === "header") return { ...prev, header: { ...prev.header, blocks: mpR(prev.header.blocks) } }; if (pcSecR === "footer") return { ...prev, footer: { ...prev.footer, blocks: mpR(prev.footer.blocks) } }; const rR = activeRouteRef.current; return { ...prev, pages: { ...prev.pages, [rR]: { ...prev.pages[rR], blocks: mpR(prev.pages[rR]?.blocks || []) } } }; }
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: (prev.popups[pid]?.blocks || []).filter(b => b.id !== id) } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: prev.header.blocks.filter(b => b.id !== id) } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: prev.footer.blocks.filter(b => b.id !== id) } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: (prev.pages[route]?.blocks || []).filter(b => b.id !== id) } } };
    });
    if (selected === id) setSelected(null);
  }

  function duplicateBlock(id: string) {
    const source = currentBlocks.find(b => b.id === id);
    if (!source) return;
    const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0);
    const copy: Block = { ...source, id: nextId("blk"), x: source.x + 24, y: source.y + 24, z: maxZ + 1 };
    pushHistory();
    setDoc(prev => {
      const pid = editingPopupRef.current;
      const pcIdDup = editingProductCardRef.current;
      if (pcIdDup) { const pcSecD = editingProductCardSectionRef.current; const mpD = (bl: Block[]) => bl.map(b => b.id === pcIdDup ? { ...b, inner_blocks: [...(b.inner_blocks || []), copy] } : b); if (pcSecD === "header") return { ...prev, header: { ...prev.header, blocks: mpD(prev.header.blocks) } }; if (pcSecD === "footer") return { ...prev, footer: { ...prev.footer, blocks: mpD(prev.footer.blocks) } }; const rD = activeRouteRef.current; return { ...prev, pages: { ...prev.pages, [rD]: { ...prev.pages[rD], blocks: mpD(prev.pages[rD]?.blocks || []) } } }; }
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: [...(prev.popups[pid]?.blocks || []), copy] } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: [...prev.header.blocks, copy] } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: [...prev.footer.blocks, copy] } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: [...(prev.pages[route]?.blocks || []), copy] } } };
    });
    setSelected(copy.id);
  }

  function contextAction(action: "delete" | "duplicate" | "bringForward" | "sendBackward" | "width100" | "height100") {
    if (!menu) return;
    const block = currentBlocks.find(b => b.id === menu.id);
    if (!block) return;
    if (action === "delete") { removeBlock(block.id); setMenu(null); return; }
    if (action === "duplicate") { duplicateBlock(block.id); setMenu(null); return; }
    if (action === "bringForward") { const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0); updateBlock(editSection, block.id, { z: maxZ + 1 }); setMenu(null); return; }
    if (action === "sendBackward") { const minZ = currentBlocks.reduce((acc, b) => Math.min(acc, b.z || 0), Infinity); updateBlock(editSection, block.id, { z: Math.max(0, minZ - 1) }); setMenu(null); return; }
    if (action === "width100") { const _epid = editingPopupRef.current; const _pcid = editingProductCardRef.current; const w = _pcid ? (pcCardBlock?.w || 300) : _epid ? (doc.popups[_epid]?.width || 480) : doc.canvas.width; updateBlock(editSection, block.id, { w }); setMenu(null); return; }
    const _epid2 = editingPopupRef.current;
    const _pcid2 = editingProductCardRef.current;
    const canvasH = _pcid2 ? (pcCardBlock?.h || 280) : _epid2 ? (doc.popups[_epid2]?.height || 560) : (editSection === "header" ? doc.header.height : editSection === "footer" ? doc.footer.height : doc.canvas.height);
    updateBlock(editSection, block.id, { h: canvasH }); setMenu(null);
  }

  function addPage() {
    const newPath = `/pagina-${Date.now()}`;
    const newRoute: RouteItem = { id: nextId("rt"), path: newPath, title: `Pagina ${routes.length + 1}`, requires_auth: false, position: routes.length };
    setRoutes(prev => [...prev, newRoute]);
    setDoc(prev => ({ ...prev, pages: { ...prev.pages, [newPath]: { title: newRoute.title, blocks: [] } } }));
    setActiveRoute(newPath);
  }

  function removePage(path: string) {
    const safePath = ensurePath(path);
    if (routes.length <= 1) return;
    setRoutes(prev => prev.filter(r => ensurePath(r.path) !== safePath));
    setDoc(prev => {
      const newPages = { ...prev.pages };
      delete newPages[safePath];
      return { ...prev, pages: newPages };
    });
    if (activeRoute === safePath) {
      const remaining = routes.filter(r => ensurePath(r.path) !== safePath);
      setActiveRoute(ensurePath(remaining[0]?.path || "/"));
    }
  }

  function addPopup() {
    const id = nextId("popup");
    const title = `Popup ${Object.keys(doc.popups).length + 1}`;
    pushHistory();
    setDoc(prev => ({ ...prev, popups: { ...prev.popups, [id]: { title, blocks: [], width: 480, height: 560, background: "#ffffff" } } }));
    setEditingPopup(id);
    setEditSection("page");
  }

  function removePopup(id: string) {
    pushHistory();
    setDoc(prev => { const p = { ...prev.popups }; delete p[id]; return { ...prev, popups: p }; });
    if (editingPopup === id) setEditingPopup(null);
  }

  function applyTemplate(type: "landing" | "blank") {
    if (type === "landing") {
      const t = makeLandingTemplate();
      setDoc(t);
      const routeKeys = Object.keys(t.pages);
      setRoutes(routeKeys.map((path, i) => ({ id: nextId("rt"), path, title: t.pages[path].title, requires_auth: false, position: i })));
      setActiveRoute("/");
    } else {
      const blankDoc: BuilderDoc = { title: "Site", canvas: { width: 1400, height: 900, background: "#f8f9ff" }, header: { enabled: false, height: 80, background: "#1a2740", blocks: [] }, footer: { enabled: false, height: 100, background: "#1a2740", blocks: [] }, pages: { "/": { title: "Pagina Inicial", blocks: [] } }, popups: {} };
      setDoc(blankDoc);
      if (routes.length === 0) setRoutes([{ id: nextId("rt"), path: "/", title: "Pagina Inicial", requires_auth: false, position: 0 }]);
      setActiveRoute("/");
    }
    setShowTemplateModal(false);
  }

  // ─ API test ──────────────────────────────────────────────────────────────────
  async function testApi(api: SiteAPI) {
    setApiTesting(true); setApiTestResult(null);
    try {
      const resp = await fetch(`${API_URL}${api.path}`, { method: api.method, credentials: "include", headers: { "X-Website-Id": websiteId || "" } });
      const text = await resp.text();
      setApiTestResult(`${resp.status} ${resp.statusText}\n${text.slice(0, 400)}`);
    } catch (err) {
      setApiTestResult(err instanceof Error ? err.message : "Erro na chamada");
    } finally { setApiTesting(false); }
  }

  // ─ Save ──────────────────────────────────────────────────────────────────────
  async function save(publish: boolean) {
    if (!siteId) return;
    setSaving(true); setError(null); setSuccess(null);
    try {
      const routeInputs = routes.map((r, i) => ({ path: ensurePath(r.path), title: r.title || ensurePath(r.path), requires_auth: r.requires_auth, position: i }));
      await apiFetch(`/api/v1/sites/${siteId}/routes`, { method: "POST", websiteId, body: JSON.stringify({ routes: routeInputs }) });
      const source = JSON.stringify(doc);
      const resp = await apiFetch<CreateVersionResponse>(`/api/v1/sites/${siteId}/versions`, { method: "POST", websiteId, body: JSON.stringify({ source_type: "ELEMENTOR_JSON", source }) });
      if (publish) {
        if (resp.data.scan_status === "blocked") { setError("Versao bloqueada pelo scan."); return; }
        await apiFetch(`/api/v1/sites/${siteId}/publish/${resp.data.version}`, { method: "POST", websiteId });
        setSuccess("✅ Publicado com sucesso!");
      } else { setSuccess("💾 Salvo!"); }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao salvar");
    } finally { setSaving(false); }
  }

  // ─ Effects ───────────────────────────────────────────────────────────────────
  useEffect(() => {
    async function load() {
      if (!siteId) return;
      try {
        const [routeResp, versionResp, apiResp] = await Promise.all([
          apiFetch<RoutesResponse>(`/api/v1/sites/${siteId}/routes`, { method: "GET", websiteId }),
          apiFetch<VersionsResponse>(`/api/v1/sites/${siteId}/versions`, { method: "GET", websiteId }),
          apiFetch<SiteAPIsResponse>("/api/v1/site-apis", { method: "GET", websiteId }),
        ]);
        const availableRoutes = routeResp.data.length > 0 ? routeResp.data : [{ id: "root", path: "/", title: "Inicio", requires_auth: false, position: 0 }];
        setRoutes(availableRoutes);
        const latestAny = versionResp.data[0];
        const latestElementor = versionResp.data.find(item => item.source_type === "ELEMENTOR_JSON");
        if (latestAny && latestAny.source_type !== "ELEMENTOR_JSON") { setError(`Este site foi criado em ${latestAny.source_type}. Edite apenas nesse modo.`); return; }
        const parsed = parseDoc(latestElementor?.source, availableRoutes);
        if (parsed === null) {
          const blankDoc: BuilderDoc = {
            title: "Site", canvas: { width: 1400, height: 900, background: "#f8f9ff" },
            header: { enabled: false, height: 80, background: "#1a2740", blocks: [] },
            footer: { enabled: false, height: 100, background: "#1a2740", blocks: [] },
            pages: Object.fromEntries(availableRoutes.map(r => [ensurePath(r.path), { title: r.title || ensurePath(r.path), blocks: [] }])),
            popups: {},
          };
          setDoc(blankDoc);
          setShowTemplateModal(true);
        } else { setDoc(parsed); }
        setActiveRoute(ensurePath(availableRoutes[0]?.path || "/"));
        setApis(apiResp.data);
      } catch (err) {
        const msg = err instanceof Error ? err.message : "Erro ao carregar editor";
        if (msg.includes("403") || msg.toLowerCase().includes("forbidden") || msg.toLowerCase().includes("unauthorized")) {
          navigate("/pages", { replace: true });
          return;
        }
        setError(msg);
      }
    }
    load();
  }, [siteId, websiteId]); // eslint-disable-line react-hooks/exhaustive-deps

  const closeMenu = useCallback(() => setMenu(null), []);
  useEffect(() => { window.addEventListener("click", closeMenu); return () => window.removeEventListener("click", closeMenu); }, [closeMenu]);

  useEffect(() => {
    function onMove(e: MouseEvent) {
      if (!interaction) return;
      const ref = interaction.section === "header" ? headerCanvasRef : interaction.section === "footer" ? footerCanvasRef : pageCanvasRef;
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      const _epid = editingPopupRef.current;
      const _pcidMov = editingProductCardRef.current;
      let blocks: Block[];
      let canvasH: number;
      let canvasW: number;
      if (interaction.section === "header") { blocks = doc.header.blocks; canvasH = doc.header.height; canvasW = doc.canvas.width; }
      else if (interaction.section === "footer") { blocks = doc.footer.blocks; canvasH = doc.footer.height; canvasW = doc.canvas.width; }
      else if (_pcidMov) { const pcSecMov = editingProductCardSectionRef.current; const srcMov = pcSecMov === "header" ? doc.header.blocks : pcSecMov === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || []; const pcBl = srcMov.find(b => b.id === _pcidMov); blocks = pcBl?.inner_blocks || []; canvasH = pcBl?.h || 280; canvasW = pcBl?.w || 300; }
      else if (_epid) { blocks = doc.popups[_epid]?.blocks || []; canvasH = doc.popups[_epid]?.height || 560; canvasW = doc.popups[_epid]?.width || 480; }
      else { blocks = doc.pages[activeRoute]?.blocks || []; canvasH = doc.canvas.height; canvasW = doc.canvas.width; }

      if (interaction.mode === "drag") {
        const block = blocks.find(b => b.id === interaction.id);
        if (!block) return;
        const nextX = (e.clientX - rect.left - interaction.dx) / zoom;
        const nextY = (e.clientY - rect.top - interaction.dy) / zoom;
        const snapped = computeSnap(nextX, nextY, block.w, block.h, block.id, blocks, canvasW, canvasH);
        setGuides(snapped.guides);
        updateBlock(interaction.section, interaction.id, { x: Math.max(0, Math.min(canvasW - block.w, snapped.x)), y: Math.max(0, Math.min(canvasH - block.h, snapped.y)) });
        return;
      }
      if (interaction.mode === "resize") {
        const dx = (e.clientX - interaction.startX) / zoom;
        const dy = (e.clientY - interaction.startY) / zoom;
        let { x, y, w, h } = interaction.start;
        if (interaction.edge === "right" || interaction.edge === "corner") w = interaction.start.w + dx;
        if (interaction.edge === "bottom" || interaction.edge === "corner") h = interaction.start.h + dy;
        if (interaction.edge === "left") { x = interaction.start.x + dx; w = interaction.start.w - dx; }
        if (interaction.edge === "top") { y = interaction.start.y + dy; h = interaction.start.h - dy; }
        w = Math.max(48, Math.min(canvasW, Math.round(w)));
        h = Math.max(24, Math.min(canvasH, Math.round(h)));
        x = Math.max(0, Math.min(canvasW - w, Math.round(x)));
        y = Math.max(0, Math.min(canvasH - h, Math.round(y)));
        const block = blocks.find(b => b.id === interaction.id);
        if (block && isTextBlock(block.type) && interaction.baseFontSize && interaction.start.h > 0) {
          const scaledFont = Math.max(8, Math.round(interaction.baseFontSize * (h / interaction.start.h)));
          updateBlock(interaction.section, interaction.id, { x, y, w, h, style: { ...block.style, "font-size": `${scaledFont}px` } });
        } else { updateBlock(interaction.section, interaction.id, { x, y, w, h }); }
        return;
      }
      const block = blocks.find(b => b.id === interaction.id);
      if (!block) return;
      const angle = (Math.atan2(e.clientY - interaction.cy, e.clientX - interaction.cx) * 180) / Math.PI;
      updateBlock(interaction.section, interaction.id, { rotation: Math.round(interaction.startRotation + (angle - interaction.startAngle)) });
    }
    function onUp() { setInteraction(null); setGuides({ vertical: null, horizontal: null }); }
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => { window.removeEventListener("mousemove", onMove); window.removeEventListener("mouseup", onUp); };
  }, [interaction, doc, activeRoute, zoom]);

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      const target = e.target as HTMLElement;
      if (target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.tagName === "SELECT") return;
      if ((e.key === "Delete" || e.key === "Backspace") && selected) { removeBlock(selected); }
      else if (e.key === "Escape") { setSelected(null); }
      else if ((e.ctrlKey || e.metaKey) && e.key === "d") { e.preventDefault(); if (selected) duplicateBlock(selected); }
      else if ((e.ctrlKey || e.metaKey) && e.key === "z") { e.preventDefault(); undo(); }
      else if ((e.ctrlKey || e.metaKey) && e.key === "y") { e.preventDefault(); redo(); }
    }
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [selected, editSection]);

  // ─ Block renderer helper ─────────────────────────────────────────────────────
  function renderSectionBlocks(blocks: Block[], section: EditSection) {
    const isActive = editSection === section || (section === "page" && (!!editingPopupRef.current || !!editingProductCardRef.current));
    return [...blocks].sort((a, b) => (a.z || 0) - (b.z || 0)).map(block => (
      <article
        key={block.id}
        className={`${styles.block} ${selected === block.id && isActive ? styles.blockSelected : ""} ${!isActive ? styles.blockFrozen : ""}`}
        style={{
          left: block.x, top: block.y, width: block.w, height: block.h, zIndex: block.z,
          transform: `rotate(${block.rotation}deg)`,
          color: block.style.color, background: block.style.background,
          border: block.style.border, borderRadius: block.style["border-radius"],
          padding: block.style.padding, fontSize: block.style["font-size"],
          fontWeight: block.style["font-weight"],
          textAlign: block.style["text-align"] as "left" | "center" | "right" | "justify" | undefined,
          cursor: isActive ? "move" : "pointer",
        }}
        onMouseDown={e => {
          if (!isActive) { e.stopPropagation(); setEditSection(section); setSelected(null); return; }
          const blockRect = e.currentTarget.getBoundingClientRect();
          setInteraction({ mode: "drag", id: block.id, section, dx: e.clientX - blockRect.left, dy: e.clientY - blockRect.top });
          setSelected(block.id); setMenu(null); e.stopPropagation();
        }}
        onClick={e => e.stopPropagation()}
        onContextMenu={e => {
          if (!isActive) return;
          e.preventDefault(); setSelected(block.id); setMenu({ x: e.clientX, y: e.clientY, id: block.id });
        }}
      >
        {selected === block.id && isActive && <div className={styles.blockTag}>{block.type}</div>}
        {renderBlock(block)}
        {selected === block.id && isActive && (
          <>
            <button type="button" className={`${styles.handle} ${styles.rotateHandle}`} onMouseDown={e => { e.stopPropagation(); const r = (e.currentTarget.parentElement as HTMLElement).getBoundingClientRect(); const cx = r.left + r.width / 2; const cy = r.top + r.height / 2; setInteraction({ mode: "rotate", id: block.id, section, cx, cy, startAngle: (Math.atan2(e.clientY - cy, e.clientX - cx) * 180) / Math.PI, startRotation: block.rotation }); }} />
            <button type="button" className={`${styles.handle} ${styles.rightHandle}`} onMouseDown={e => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge: "right", startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
            <button type="button" className={`${styles.handle} ${styles.leftHandle}`} onMouseDown={e => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge: "left", startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
            <button type="button" className={`${styles.handle} ${styles.topHandle}`} onMouseDown={e => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge: "top", startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
            <button type="button" className={`${styles.handle} ${styles.bottomHandle}`} onMouseDown={e => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge: "bottom", startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
            <button type="button" className={`${styles.handle} ${styles.cornerHandle}`} onMouseDown={e => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge: "corner", startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
          </>
        )}
      </article>
    ));
  }

  // ─ Render ────────────────────────────────────────────────────────────────────
  return (
    <main className={styles.main}>
      {showTemplateModal && (
        <div className={styles.modalOverlay} onClick={() => setShowTemplateModal(false)}>
          <div className={styles.modal} onClick={e => e.stopPropagation()}>
            <button type="button" className={styles.modalClose} onClick={() => setShowTemplateModal(false)}>×</button>
            <h2 style={{ margin: "0 0 .4rem" }}>Escolha um template</h2>
            <p style={{ margin: "0 0 1rem", color: "#6a7387", fontSize: ".9rem" }}>Selecione um ponto de partida para o seu site.</p>
            <div className={styles.templateGrid}>
              <div className={styles.templateCard} onClick={() => applyTemplate("landing")}>
                <div className={styles.templateEmoji}>🚀</div>
                <strong>Landing Page</strong>
                <p>Template profissional com header, hero, features e footer</p>
              </div>
              <div className={styles.templateCard} onClick={() => applyTemplate("blank")}>
                <div className={styles.templateEmoji}>📄</div>
                <strong>Em branco</strong>
                <p>Canvas vazio para comecar do zero</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <button type="button" onClick={() => navigate("/pages")} style={{ background: "none", border: "1px solid #d2daeb", borderRadius: 8, padding: "6px 12px", cursor: "pointer", fontSize: ".82rem" }}>← Voltar</button>
          <input className={styles.siteTitleInput} value={doc.title} onChange={e => setDoc(prev => ({ ...prev, title: e.target.value }))} />
          <span className={styles.siteIdLabel}>{siteId}</span>
        </div>
        <div className={styles.headerCenter}>
          <div className={styles.editSectionTabs}>
            <button type="button" className={editSection === "page" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("page"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); }}>📄 Pagina</button>
            <button type="button" className={editSection === "header" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("header"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); if (!doc.header.enabled) setDoc(prev => ({ ...prev, header: { ...prev.header, enabled: true } })); }}>⬆ Header</button>
            <button type="button" className={editSection === "footer" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("footer"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); if (!doc.footer.enabled) setDoc(prev => ({ ...prev, footer: { ...prev.footer, enabled: true } })); }}>⬇ Footer</button>
          </div>
          <div className={styles.zoomControl}>
            {ZOOM_LEVELS.map(z => (
              <button key={z} type="button" className={zoom === z ? styles.zoomActive : styles.zoomBtn} onClick={() => setZoom(z)}>{Math.round(z * 100)}%</button>
            ))}
          </div>
        </div>
        <div className={styles.headerRight}>
          <button type="button" className={styles.undoBtn} disabled={!canUndo} onClick={undo} title="Desfazer (Ctrl+Z)">↩</button>
          <button type="button" className={styles.undoBtn} disabled={!canRedo} onClick={redo} title="Refazer (Ctrl+Y)">↪</button>
          <button type="button" disabled={saving} onClick={() => save(false)} style={{ background: "#fff", border: "1px solid #d2daeb", borderRadius: 8, padding: "6px 14px", cursor: "pointer", fontSize: ".82rem" }}>💾 Salvar</button>
          <button type="button" disabled={saving} onClick={() => save(true)} style={{ background: "#1a2740", color: "#fff", border: "none", borderRadius: 8, padding: "6px 14px", cursor: "pointer", fontSize: ".82rem" }}>🚀 Publicar</button>
          <a href={`${API_URL}/p/${siteId}`} target="_blank" rel="noreferrer" className={styles.openLive}>Abrir ↗</a>
        </div>
      </div>

      {error && <p className={styles.error}>{error}</p>}
      {success && <div className={styles.successBanner}>{success}</div>}

      <div className={styles.workspace}>
        {/* ─ Left Panel ─ */}
        <aside className={styles.left}>
          <div className={styles.leftTabs}>
            <button type="button" className={leftTab === "elements" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("elements")} title="Elementos">🧩</button>
            <button type="button" className={leftTab === "layers" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("layers")} title="Camadas">📐</button>
            <button type="button" className={leftTab === "settings" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("settings")} title="Configuracoes">⚙</button>
            <button type="button" className={leftTab === "templates" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("templates")} title="Templates">🎨</button>
            <button type="button" className={leftTab === "popups" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("popups")} title="Popups">📌</button>
          </div>

          {leftTab === "elements" && (
            <div className={styles.blockList}>
              {editingProductCard && <div className={styles.varHint} style={{ padding: "6px 10px", margin: "0 0 6px", background: "#e0f2fe", borderRadius: 6, fontSize: ".75rem", color: "#0369a1" }}>🛍 <strong>Vars produto:</strong> {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{product_sku}}"}</div>}
              {!editingProductCard && <div className={styles.varHint} style={{ padding: "6px 10px", margin: "0 0 6px", background: "#f1f5f9", borderRadius: 6, fontSize: ".75rem" }}>👤 <strong>Vars usuário:</strong> {"{{user_name}}"} {"{{user_email}}"} {"{{cart_count}}"}</div>}
              {BLOCK_CATEGORIES.map(cat => (
                <div key={cat.label}>
                  <div className={styles.catLabel}>{cat.label}</div>
                  <div className={styles.catItems}>
                    {cat.items.map(item => (
                      <button key={item.type} type="button" className={styles.blockItem} onClick={() => addBlock(item.type)}>
                        <span className={styles.blockIcon}>{item.icon}</span>
                        {item.label}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}

          {leftTab === "layers" && (
            <div className={styles.layersList}>
              <div className={styles.layersHeader}>Camadas · {editSection === "page" ? activeRoute : editSection}</div>
              {currentBlocks.length === 0 && <div className={styles.emptyLayers}>Sem blocos nesta secao</div>}
              {[...currentBlocks].sort((a, b) => (b.z || 0) - (a.z || 0)).map(block => (
                <div key={block.id} className={`${styles.layerItem} ${selected === block.id ? styles.layerSelected : ""}`} onClick={() => setSelected(block.id)}>
                  <span className={styles.layerIcon}>{getBlockIcon(block.type)}</span>
                  <span className={styles.layerName}>{block.text || block.label || block.profile_name || block.type}</span>
                  <button type="button" className={styles.layerDelete} onClick={e => { e.stopPropagation(); removeBlock(block.id); }}>🗑</button>
                </div>
              ))}
            </div>
          )}

          {leftTab === "settings" && (
            <div className={styles.settingsPanel}>
              <div className={styles.settingsGroup}>
                <div className={styles.settingsGroupTitle}>Pagina</div>
                <select value={activeRoute} onChange={e => { setActiveRoute(e.target.value); setSelected(null); }}>
                  {routes.map(r => <option key={r.id} value={ensurePath(r.path)}>{r.title} ({ensurePath(r.path)})</option>)}
                </select>
                <label>Titulo</label>
                <input value={doc.pages[activeRoute]?.title || ""} onChange={e => {
                  const newTitle = e.target.value;
                  setDoc(prev => ({ ...prev, pages: { ...prev.pages, [activeRoute]: { ...prev.pages[activeRoute], title: newTitle } } }));
                  setRoutes(prev => prev.map(r => ensurePath(r.path) === activeRoute ? { ...r, title: newTitle } : r));
                }} />
              </div>
              <div className={styles.settingsGroup}>
                <div className={styles.settingsGroupTitle}>Canvas</div>
                <div className={styles.colorPickerRow}>
                  <label>Fundo</label>
                  <input type="color" value={doc.canvas.background || "#f8f9ff"} onChange={e => setDoc(prev => ({ ...prev, canvas: { ...prev.canvas, background: e.target.value } }))} />
                </div>
                <label>Largura</label>
                <input type="number" value={doc.canvas.width} min={900} max={2800} onChange={e => setDoc(prev => ({ ...prev, canvas: { ...prev.canvas, width: Math.max(900, Math.min(2800, Number(e.target.value) || 900)) } }))} />
                <label>Altura</label>
                <input type="number" value={doc.canvas.height} min={700} max={2800} onChange={e => setDoc(prev => ({ ...prev, canvas: { ...prev.canvas, height: Math.max(700, Math.min(2800, Number(e.target.value) || 700)) } }))} />
              </div>
              <div className={styles.settingsGroup}>
                <div className={styles.settingsGroupTitle}>Header Global</div>
                <div className={styles.toggleRow}>
                  <label>Ativo</label>
                  <input type="checkbox" checked={doc.header.enabled} onChange={e => setDoc(prev => ({ ...prev, header: { ...prev.header, enabled: e.target.checked } }))} />
                </div>
                <label>Altura</label>
                <input type="number" value={doc.header.height} min={40} max={400} onChange={e => setDoc(prev => ({ ...prev, header: { ...prev.header, height: Math.max(40, Number(e.target.value) || 80) } }))} />
                <div className={styles.colorPickerRow}>
                  <label>Fundo</label>
                  <input type="color" value={doc.header.background || "#1a2740"} onChange={e => setDoc(prev => ({ ...prev, header: { ...prev.header, background: e.target.value } }))} />
                </div>
                <button type="button" className={editSection === "header" ? styles.editSectionBtnActive : styles.editSectionBtn} onClick={() => { setEditSection("header"); setSelected(null); }}>✏️ Editar Header</button>
              </div>
              <div className={styles.settingsGroup}>
                <div className={styles.settingsGroupTitle}>Footer Global</div>
                <div className={styles.toggleRow}>
                  <label>Ativo</label>
                  <input type="checkbox" checked={doc.footer.enabled} onChange={e => setDoc(prev => ({ ...prev, footer: { ...prev.footer, enabled: e.target.checked } }))} />
                </div>
                <label>Altura</label>
                <input type="number" value={doc.footer.height} min={40} max={400} onChange={e => setDoc(prev => ({ ...prev, footer: { ...prev.footer, height: Math.max(40, Number(e.target.value) || 100) } }))} />
                <div className={styles.colorPickerRow}>
                  <label>Fundo</label>
                  <input type="color" value={doc.footer.background || "#1a2740"} onChange={e => setDoc(prev => ({ ...prev, footer: { ...prev.footer, background: e.target.value } }))} />
                </div>
                <button type="button" className={editSection === "footer" ? styles.editSectionBtnActive : styles.editSectionBtn} onClick={() => { setEditSection("footer"); setSelected(null); }}>✏️ Editar Footer</button>
              </div>
              <div className={styles.settingsGroup}>
                <div className={styles.settingsGroupTitle}>Paginas</div>
                {routes.map(route => (
                  <div key={route.id} className={styles.pageRow}>
                    <span style={{ fontSize: ".78rem", flex: 1 }}>{ensurePath(route.path)}</span>
                    <button type="button" className={styles.removePageBtn} disabled={routes.length <= 1} onClick={() => removePage(route.path)}>✕</button>
                  </div>
                ))}
                <button type="button" className={styles.addPageBtn} onClick={addPage}>+ Nova Pagina</button>
              </div>
            </div>
          )}

          {leftTab === "templates" && (
            <div className={styles.templatesPanel}>
              <div className={styles.templateCardSmall} onClick={() => applyTemplate("landing")}>
                <span className={styles.templateEmoji}>🚀</span>
                <span>Landing Page</span>
              </div>
              <div className={styles.templateCardSmall} onClick={() => applyTemplate("blank")}>
                <span className={styles.templateEmoji}>📄</span>
                <span>Em branco</span>
              </div>
            </div>
          )}

          {leftTab === "popups" && (
            <div className={styles.layersList}>
              <div className={styles.layersHeader}>Popups</div>
              {Object.entries(doc.popups).map(([id, popup]) => (
                <div key={id} className={`${styles.layerItem} ${editingPopup === id ? styles.layerSelected : ""}`}
                  onClick={() => { setEditingPopup(id); setEditingProductCard(null); setEditSection("page"); setSelected(null); }}>
                  <span className={styles.layerIcon}>📌</span>
                  <span className={styles.layerName}>{popup.title}</span>
                  <button type="button" className={styles.layerDelete} onClick={e => { e.stopPropagation(); removePopup(id); }}>🗑</button>
                </div>
              ))}
              {Object.keys(doc.popups).length === 0 && <div className={styles.emptyLayers}>Nenhum popup criado</div>}
              <button type="button" className={styles.addPageBtn} onClick={addPopup}>+ Novo Popup</button>
            </div>
          )}
        </aside>

        {/* ─ Canvas ─ */}
        <section className={styles.canvasWrap}>
          {editingProductCard && !editingPopup && (
            <div style={{ background: "#0ea5e9", color: "#fff", padding: "8px 16px", fontSize: ".82rem", display: "flex", alignItems: "center", gap: 10, flexWrap: "wrap" }}>
              <strong>🛍 Editando Layout do Product Card</strong>
              <span style={{ opacity: .7, fontSize: ".75rem" }}>Vars: {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{cart_count}}"}</span>
              <button type="button" onClick={() => { setEditingProductCard(null); setSelected(null); }} style={{ background: "none", border: "1px solid rgba(255,255,255,.5)", color: "#fff", borderRadius: 6, padding: "2px 10px", cursor: "pointer", fontSize: ".8rem", marginLeft: "auto" }}>← Sair do Product Card</button>
            </div>
          )}
          {editingPopup && (
            <div style={{ background: "#7c3aed", color: "#fff", padding: "8px 16px", fontSize: ".82rem", display: "flex", alignItems: "center", gap: 10 }}>
              <strong>📌 Editando Popup: {doc.popups[editingPopup]?.title}</strong>
              <button type="button" onClick={() => setEditingPopup(null)} style={{ background: "none", border: "1px solid rgba(255,255,255,.5)", color: "#fff", borderRadius: 6, padding: "2px 10px", cursor: "pointer", fontSize: ".8rem" }}>← Sair do Popup</button>
            </div>
          )}
          <div className={styles.canvasScroll}>
            <div className={styles.canvasZoomHost} style={{ zoom }}>
              <div className={styles.canvasStack}>
                {editingPopup ? (
                  <div
                    ref={pageCanvasRef}
                    className={`${styles.sectionCanvas} ${styles.sectionActive}`}
                    style={{ width: doc.popups[editingPopup]?.width || 480, height: doc.popups[editingPopup]?.height || 560, backgroundColor: doc.popups[editingPopup]?.background || "#ffffff", backgroundImage: "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)", backgroundSize: "20px 20px" }}
                    onClick={() => setSelected(null)}
                  >
                    <div className={styles.sectionLabel}>POPUP · {doc.popups[editingPopup]?.title}</div>
                    {renderSectionBlocks(doc.popups[editingPopup]?.blocks || [], "page")}
                    {guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                ) : editingProductCard ? (
                  <div
                    ref={pageCanvasRef}
                    className={`${styles.sectionCanvas} ${styles.sectionActive}`}
                    style={{ width: pcCardBlock?.w || 300, height: pcCardBlock?.h || 280, backgroundColor: pcCardBlock?.style.background || "#ffffff", backgroundImage: "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)", backgroundSize: "20px 20px" }}
                    onClick={() => setSelected(null)}
                  >
                    <div className={styles.sectionLabel}>PRODUCT CARD · Layout</div>
                    {renderSectionBlocks(currentBlocks, "page")}
                    {guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                ) : (
                  <>
                {doc.header.enabled && (
                  <div
                    ref={headerCanvasRef}
                    className={`${styles.sectionCanvas} ${editSection === "header" ? styles.sectionActive : styles.sectionInactive}`}
                    style={{ width: doc.canvas.width, height: doc.header.height, background: doc.header.background }}
                    onClick={() => { if (editSection !== "header") { setEditSection("header"); setSelected(null); } }}
                  >
                    <div className={styles.sectionLabel}>HEADER</div>
                    {renderSectionBlocks(doc.header.blocks, "header")}
                    {editSection === "header" && guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {editSection === "header" && guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                )}
                <div
                  ref={pageCanvasRef}
                  className={`${styles.sectionCanvas} ${editSection === "page" ? styles.sectionActive : styles.sectionInactive}`}
                  style={{ width: doc.canvas.width, height: doc.canvas.height, backgroundColor: doc.canvas.background, backgroundImage: editSection === "page" ? "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)" : undefined, backgroundSize: editSection === "page" ? "20px 20px" : undefined }}
                  onClick={() => { if (editSection !== "page") { setEditSection("page"); setSelected(null); } }}
                >
                  <div className={styles.sectionLabel}>PAGINA</div>
                  {renderSectionBlocks(doc.pages[activeRoute]?.blocks || [], "page")}
                  {editSection === "page" && guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                  {editSection === "page" && guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                </div>
                {doc.footer.enabled && (
                  <div
                    ref={footerCanvasRef}
                    className={`${styles.sectionCanvas} ${editSection === "footer" ? styles.sectionActive : styles.sectionInactive}`}
                    style={{ width: doc.canvas.width, height: doc.footer.height, background: doc.footer.background }}
                    onClick={() => { if (editSection !== "footer") { setEditSection("footer"); setSelected(null); } }}
                  >
                    <div className={styles.sectionLabel}>FOOTER</div>
                    {renderSectionBlocks(doc.footer.blocks, "footer")}
                    {editSection === "footer" && guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {editSection === "footer" && guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                )}
                  </>
                )}
              </div>
            </div>
          </div>
        </section>

        {/* ─ Inspector ─ */}
        <aside className={styles.right}>
          <h3 style={{ margin: "0 0 .4rem", fontSize: ".8rem", textTransform: "uppercase", letterSpacing: ".04em", color: "#4b5774" }}>🔍 Inspector</h3>
          {!selectedBlock && <p className={styles.emptyInspector}>Clique em um bloco no canvas para editar.</p>}
          {selectedBlock && (
            <div className={styles.form}>
              <div className={styles.inspectorSection}>
                <div className={styles.sectionLabelTag}>CONTEUDO · {selectedBlock.type.toUpperCase()}</div>
                {isTextBlock(selectedBlock.type) && (
                  <>
                    <label>Texto</label>
                    <textarea value={selectedBlock.text || ""} rows={3} onChange={e => updateBlock(editSection, selectedBlock.id, { text: e.target.value })} />
                    <label>Tamanho fonte</label>
                    <input value={selectedBlock.style["font-size"] || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, "font-size": e.target.value } })} placeholder="ex: 48px" />
                    <label>Peso fonte</label>
                    <select value={selectedBlock.style["font-weight"] || "400"} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, "font-weight": e.target.value } })}>
                      <option value="300">Light 300</option><option value="400">Normal 400</option><option value="500">Medium 500</option><option value="600">SemiBold 600</option><option value="700">Bold 700</option><option value="800">ExtraBold 800</option>
                    </select>
                  </>
                )}
                {selectedBlock.type === "button" && (
                  <>
                    <label>Texto do botao</label>
                    <input value={selectedBlock.label || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { label: e.target.value })} />
                    <label>Acao</label>
                    <select value={selectedBlock.action_type || "navigate"} onChange={e => updateBlock(editSection, selectedBlock.id, { action_type: e.target.value as Block["action_type"] })}>
                      <option value="navigate">Navegar para rota</option>
                      <option value="call_api">Chamar API</option>
                      <option value="store_login">Login na loja</option>
                      <option value="store_logout">Logout da loja</option>
                      <option value="store_register">Registrar na loja</option>
                      <option value="add_to_cart">Adicionar ao carrinho</option>
                      <option value="add_product">Adicionar produto (admin)</option>
                      <option value="open_popup">Abrir Popup</option>
                      <option value="close_popup">Fechar Popup</option>
                    </select>
                    {selectedBlock.action_type === "call_api" ? (
                      <>
                        <label>API</label>
                        <select value={selectedBlock.api_id || ""} onChange={e => { const api = apis.find(a => a.id === e.target.value); updateBlock(editSection, selectedBlock.id, { api_id: e.target.value, action_target: api?.path || "", href: "#" }); }}>
                          <option value="">— Selecione API —</option>
                          {apis.map(api => <option key={api.id} value={api.id}>{api.method} {api.path} — {api.label}</option>)}
                        </select>
                        {selectedBlock.api_id && (() => { const api = apis.find(a => a.id === selectedBlock.api_id); return api ? (<div className={styles.apiTestBox}><p className={styles.apiDesc}>{api.description}</p><button type="button" className={styles.apiTestBtn} disabled={apiTesting} onClick={() => testApi(api)}>{apiTesting ? "Testando..." : "▶ Testar API"}</button>{apiTestResult && <pre className={styles.apiResult}>{apiTestResult}</pre>}</div>) : null; })()}
                      </>
                    ) : selectedBlock.action_type === "add_to_cart" ? (
                      <>
                        <label>ID do produto</label>
                        <input value={selectedBlock.product_id || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { product_id: e.target.value })} placeholder="UUID do produto" />
                        <p className={styles.varHint}>Deixe vazio para usar o 1º produto da loja.</p>
                      </>
                    ) : selectedBlock.action_type === "open_popup" ? (
                       <>
                         <label>Popup</label>
                         <select value={selectedBlock.popup_id || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { popup_id: e.target.value })}>
                           <option value="">— Selecione Popup —</option>
                           {Object.entries(doc.popups).map(([id, popup]) => <option key={id} value={id}>{popup.title}</option>)}
                         </select>
                       </>
                     ) : selectedBlock.action_type === "store_login" || selectedBlock.action_type === "store_register" ? (
                       <>
                         <label>Var email</label>
                         <input value={selectedBlock.email_var || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { email_var: e.target.value })} placeholder="nome da variavel de email" />
                         <label>Var senha</label>
                         <input value={selectedBlock.password_var || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { password_var: e.target.value })} placeholder="nome da variavel de senha" />
                         {selectedBlock.action_type === "store_register" && (<>
                           <label>Var nome</label>
                           <input value={selectedBlock.first_name_var || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { first_name_var: e.target.value })} placeholder="nome da variavel de nome" />
                           <label>Var sobrenome</label>
                           <input value={selectedBlock.last_name_var || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { last_name_var: e.target.value })} placeholder="nome da variavel de sobrenome" />
                         </>)}
                         <p className={styles.varHint}>Deixe vazio para abrir o modal de login/registro.</p>
                       </>
                     ) : selectedBlock.action_type === "navigate" || !selectedBlock.action_type ? (
                       <>
                         <label>Rota destino</label>
                         <select value={selectedBlock.action_target || "/"} onChange={e => updateBlock(editSection, selectedBlock.id, { action_target: e.target.value, href: e.target.value })}>
                           {routes.map(r => <option key={r.id} value={ensurePath(r.path)}>{r.title} ({ensurePath(r.path)})</option>)}
                         </select>
                       </>
                     ) : null}
                     <div className={styles.toggleRow} style={{ marginTop: 8 }}>
                       <label>Visivel somente para admins</label>
                       <input type="checkbox" checked={!!selectedBlock.admin_only} onChange={e => updateBlock(editSection, selectedBlock.id, { admin_only: e.target.checked })} />
                     </div>
                     {selectedBlock.admin_only && <p className={styles.varHint}>Oculto para visitantes; visivel apenas para o dono da loja.</p>}
                  </>
                )}
                {selectedBlock.type === "product_card" && (
                  <>
                    <label>ID do produto</label>
                    <input value={selectedBlock.product_id || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { product_id: e.target.value })} placeholder="UUID do produto (opcional)" />
                    <p className={styles.varHint}>Deixe vazio para exibir o 1º produto ativo da loja.</p>
                    {!editingPopup && <button type="button" className={styles.editSectionBtn} style={{ marginTop: 6 }} onClick={() => { setEditingProductCard(selectedBlock.id); setEditingProductCardSection(editSection); setSelected(null); }}>✏️ Editar Layout Interno ({selectedBlock.inner_blocks?.length || 0} blocos)</button>}
                    {(selectedBlock.inner_blocks?.length || 0) > 0 && <button type="button" className={styles.apiTestBtn} style={{ marginTop: 4, background: "#fee2e2", color: "#b91c1c", border: "1px solid #fecaca" }} onClick={() => updateBlock(editSection, selectedBlock.id, { inner_blocks: [] })}>🗑 Limpar layout</button>}
                    <p className={styles.varHint}>Vars: {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{product_sku}}"}</p>
                  </>
                )}
                {selectedBlock.type === "product_list" && (
                  <>
                    <label>Produtos por pagina</label>
                    <input type="number" min={1} max={50} value={selectedBlock.page_size || 6} onChange={e => updateBlock(editSection, selectedBlock.id, { page_size: Math.max(1, Number(e.target.value) || 6) })} />
                    <p className={styles.varHint}>Ordenados por mais vendidos. Botoes de pagina aparecem automaticamente.</p>
                  </>
                )}
                {selectedBlock.type === "image" && (
                  <>
                    <label>URL da imagem</label>
                    <input value={selectedBlock.src || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { src: e.target.value })} placeholder="https://..." />
                    <label>Var source (opcional)</label>
                    <input value={selectedBlock.var_src || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { var_src: e.target.value })} placeholder="ex: {{product_image}}" />
                    <label>Upload</label>
                    <input type="file" accept="image/*" onChange={async e => { const file = e.target.files?.[0]; if (!file) return; try { updateBlock(editSection, selectedBlock.id, { src: await readFileAsDataURL(file) }); } catch { setError("Falha ao carregar imagem"); } }} />
                    <label>Ajuste</label>
                    <select value={selectedBlock.object_fit || "cover"} onChange={e => updateBlock(editSection, selectedBlock.id, { object_fit: e.target.value as "cover" | "contain" | "fill" })}>
                      <option value="cover">Cover</option><option value="contain">Contain</option><option value="fill">Fill</option>
                    </select>
                  </>
                )}
                {selectedBlock.type === "user_avatar" && (
                  <>
                    <p className={styles.varHint}>Exibe o avatar do usuário logado. Atualiza após login/logout.</p>
                    <label>Ajuste</label>
                    <select value={selectedBlock.object_fit || "cover"} onChange={e => updateBlock(editSection, selectedBlock.id, { object_fit: e.target.value as "cover" | "contain" | "fill" })}>
                      <option value="cover">Cover</option><option value="contain">Contain</option><option value="fill">Fill</option>
                    </select>
                  </>
                )}
                {selectedBlock.type === "carousel" && (
                  <>
                    <label>Imagens (uma URL por linha)</label>
                    <textarea value={(selectedBlock.images || []).join("\n")} rows={4} onChange={e => updateBlock(editSection, selectedBlock.id, { images: e.target.value.split("\n").map(i => i.trim()).filter(Boolean) })} />
                  </>
                )}
                {selectedBlock.type === "video" && (
                  <>
                    <label>URL do video</label>
                    <input value={selectedBlock.video_url || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { video_url: e.target.value })} placeholder="https://youtube.com/watch?v=..." />
                  </>
                )}
                {selectedBlock.type === "input_var" && (
                  <>
                    <label>Nome da variavel</label>
                    <input value={selectedBlock.var_name || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { var_name: e.target.value })} placeholder="ex: nome" />
                    <p className={styles.varHint}>Use {"{{" + (selectedBlock.var_name || "nome") + "}}"} em blocos Texto Variavel.</p>
                    <label>Placeholder</label>
                    <input value={selectedBlock.placeholder || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { placeholder: e.target.value })} />
                    <label>Tipo do campo</label>
                    <select value={selectedBlock.input_type || "text"} onChange={e => updateBlock(editSection, selectedBlock.id, { input_type: e.target.value })}>
                      <option value="text">Texto</option>
                      <option value="email">Email</option>
                      <option value="password">Senha (oculta)</option>
                      <option value="number">Numero</option>
                      <option value="tel">Telefone</option>
                    </select>
                  </>
                )}
                {selectedBlock.type === "profile_card" && (
                  <>
                    <label>Nome</label>
                    <input value={selectedBlock.profile_name || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { profile_name: e.target.value })} />
                    <label>Subtitulo</label>
                    <input value={selectedBlock.profile_subtitle || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { profile_subtitle: e.target.value })} />
                    <label>URL da imagem</label>
                    <input value={selectedBlock.profile_image || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { profile_image: e.target.value })} />
                  </>
                )}
              </div>

              <div className={styles.inspectorSection}>
                <div className={styles.sectionLabelTag}>POSICAO &amp; TAMANHO</div>
                <div className={styles.inlineFields}>
                  <span><label>X</label><input type="number" value={selectedBlock.x} onChange={e => updateBlock(editSection, selectedBlock.id, { x: Number(e.target.value) || 0 })} /></span>
                  <span><label>Y</label><input type="number" value={selectedBlock.y} onChange={e => updateBlock(editSection, selectedBlock.id, { y: Number(e.target.value) || 0 })} /></span>
                  <span><label>W</label><input type="number" value={selectedBlock.w} onChange={e => updateBlock(editSection, selectedBlock.id, { w: Number(e.target.value) || 0 })} /></span>
                  <span><label>H</label><input type="number" value={selectedBlock.h} onChange={e => updateBlock(editSection, selectedBlock.id, { h: Number(e.target.value) || 0 })} /></span>
                </div>
                <div className={styles.inlineFields}>
                  <span><label>Rot°</label><input type="number" value={selectedBlock.rotation} onChange={e => updateBlock(editSection, selectedBlock.id, { rotation: Number(e.target.value) || 0 })} /></span>
                  <span><label>Z</label><input type="number" value={selectedBlock.z} onChange={e => updateBlock(editSection, selectedBlock.id, { z: Number(e.target.value) || 1 })} /></span>
                </div>
              </div>

              <div className={styles.inspectorSection}>
                <div className={styles.sectionLabelTag}>ESTILO</div>
                <div className={styles.colorRow}>
                  <span><label>Texto</label><input type="color" value={selectedBlock.style.color?.startsWith("#") ? selectedBlock.style.color : "#1f2b43"} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, color: e.target.value } })} /></span>
                  <span><label>Fundo</label><input type="color" value={selectedBlock.style.background?.startsWith("#") ? selectedBlock.style.background : "#ffffff"} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, background: e.target.value } })} /></span>
                </div>
                <label>Cor texto (hex/rgba)</label>
                <input value={selectedBlock.style.color || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, color: e.target.value } })} placeholder="#1f2b43" />
                <label>Cor fundo</label>
                <input value={selectedBlock.style.background || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, background: e.target.value } })} placeholder="transparent" />
                <label>Alinhamento</label>
                <select value={selectedBlock.style["text-align"] || "left"} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, "text-align": e.target.value } })}>
                  <option value="left">Esquerda</option><option value="center">Centro</option><option value="right">Direita</option><option value="justify">Justificado</option>
                </select>
                <label>Padding</label>
                <input value={selectedBlock.style.padding || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, padding: e.target.value } })} placeholder="10px 12px" />
                <label>Borda</label>
                <input value={selectedBlock.style.border || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, border: e.target.value } })} placeholder="1px solid #ccc" />
                <label>Raio borda</label>
                <input value={selectedBlock.style["border-radius"] || ""} onChange={e => updateBlock(editSection, selectedBlock.id, { style: { ...selectedBlock.style, "border-radius": e.target.value } })} placeholder="8px" />
              </div>

              <button type="button" className={styles.remove} onClick={() => removeBlock(selectedBlock.id)}>🗑 Excluir bloco</button>
            </div>
          )}
        </aside>
      </div>

      {menu && (
        <div className={styles.contextMenu} style={{ left: menu.x, top: menu.y }} onClick={e => e.stopPropagation()}>
          <button type="button" onClick={() => contextAction("duplicate")}>📋 Duplicar</button>
          <button type="button" onClick={() => contextAction("bringForward")}>⬆ Trazer para frente</button>
          <button type="button" onClick={() => contextAction("sendBackward")}>⬇ Enviar para tras</button>
          <button type="button" onClick={() => contextAction("width100")}>↔ Largura total</button>
          <button type="button" onClick={() => contextAction("height100")}>↕ Altura total</button>
          <button type="button" onClick={() => contextAction("delete")} style={{ color: "#d13a2f" }}>🗑 Excluir</button>
        </div>
      )}
    </main>
  );
};
