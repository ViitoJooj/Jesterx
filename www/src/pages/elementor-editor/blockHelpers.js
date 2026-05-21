export const ZOOM_LEVELS = [0.4, 0.5, 0.65, 0.75, 1.0, 1.25];

export const BLOCK_CATEGORIES = [
  {
    label: "Texto",
    items: [
      { type: "heading", icon: "🔤", label: "Titulo" },
      { type: "paragraph", icon: "📝", label: "Paragrafo" },
      { type: "variable_text", icon: "💬", label: "Texto Variavel" },
    ],
  },
  {
    label: "Midia",
    items: [
      { type: "image", icon: "🖼", label: "Imagem" },
      { type: "carousel", icon: "🎠", label: "Carousel" },
      { type: "video", icon: "▶", label: "Video" },
    ],
  },
  {
    label: "Interativo",
    items: [
      { type: "button", icon: "🔘", label: "Botao" },
      { type: "input_var", icon: "📌", label: "Input Variavel" },
    ],
  },
  {
    label: "Loja",
    items: [
      { type: "product_card", icon: "🛍", label: "Product Card" },
      { type: "product_list", icon: "📦", label: "Lista de Produtos" },
    ],
  },
  {
    label: "Layout",
    items: [
      { type: "profile_card", icon: "👤", label: "Perfil" },
      { type: "user_avatar", icon: "🙂", label: "Avatar Usuario" },
      { type: "divider", icon: "➖", label: "Divider" },
    ],
  },
];

export const DEFAULT_STYLE = {
  color: "#1f2b43",
  background: "transparent",
  padding: "10px 12px",
  margin: "0",
  border: "0",
  "border-radius": "0",
  "box-shadow": "none",
  "font-size": "16px",
  "font-family": "inherit",
  "font-weight": "500",
  "line-height": "1.5",
  "letter-spacing": "0",
  "text-align": "left",
  "text-transform": "none",
  "text-decoration": "none",
  opacity: "1",
};

export function nextId(prefix) {
  return `${prefix}-${Date.now()}-${Math.floor(Math.random() * 9999)}`;
}

export function isTextBlock(type) {
  return type === "heading" || type === "paragraph" || type === "variable_text";
}

export function getBlockIcon(type) {
  const icons = {
    heading: "🔤", paragraph: "📝", button: "🔘", image: "🖼",
    carousel: "🎠", input_var: "📌", variable_text: "💬",
    profile_card: "👤", product_card: "🛍", product_list: "📦",
    video: "▶", divider: "➖", user_avatar: "🙂",
  };
  return icons[type] ?? "▪";
}

export function ensurePath(path) {
  const trimmed = path.trim();
  if (!trimmed) return "/";
  if (trimmed.startsWith("/")) return trimmed;
  return `/${trimmed}`;
}

export function readFileAsDataURL(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result || ""));
    reader.onerror = () => reject(new Error("Falha ao ler arquivo"));
    reader.readAsDataURL(file);
  });
}

export function computeSnap(x, y, w, h, movingID, blocks, canvasWidth, canvasHeight) {
  const threshold = 8;
  const pointsX = [x, x + w / 2, x + w];
  const pointsY = [y, y + h / 2, y + h];
  const targetX = [canvasWidth / 2];
  const targetY = [canvasHeight / 2];
  blocks.forEach((block) => {
    if (block.id === movingID) return;
    targetX.push(block.x, block.x + block.w / 2, block.x + block.w);
    targetY.push(block.y, block.y + block.h / 2, block.y + block.h);
  });
  let snappedX = x, snappedY = y, guideX = null, guideY = null;
  let bestDx = threshold + 1;
  pointsX.forEach((point) => {
    targetX.forEach((target) => {
      const abs = Math.abs(target - point);
      if (abs < bestDx && abs <= threshold) { bestDx = abs; snappedX = x + (target - point); guideX = target; }
    });
  });
  let bestDy = threshold + 1;
  pointsY.forEach((point) => {
    targetY.forEach((target) => {
      const abs = Math.abs(target - point);
      if (abs < bestDy && abs <= threshold) { bestDy = abs; snappedY = y + (target - point); guideY = target; }
    });
  });
  return { x: Math.round(snappedX), y: Math.round(snappedY), guides: { vertical: guideX, horizontal: guideY } };
}

export function newBlock(type, z) {
  if (type === "heading") return { id: nextId("blk"), type, text: "Titulo principal", heading_level: "h1", style: { ...DEFAULT_STYLE, "font-size": "48px", "font-weight": "700", "line-height": "1.15", "letter-spacing": "-0.02em" }, x: 120, y: 90, w: 620, h: 90, rotation: 0, z };
  if (type === "paragraph") return { id: nextId("blk"), type, text: "Paragrafo com descricao da sua oferta.", style: { ...DEFAULT_STYLE, color: "#4f5f83", "line-height": "1.65" }, x: 120, y: 190, w: 580, h: 80, rotation: 0, z };
  if (type === "button") return { id: nextId("blk"), type, label: "Botao", href: "/", action_type: "navigate", action_target: "/", style: { ...DEFAULT_STYLE, background: "var(--jx-color-primary)", color: "#ffffff", padding: "12px 24px", "border-radius": "8px", "font-weight": "600", "line-height": "1.4", "text-align": "center" }, x: 120, y: 290, w: 190, h: 54, rotation: 0, z };
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

export function parseDoc(source, routes) {
  if (!source) return null;
  const defaultHeader = { enabled: false, height: 80, background: "#1a2740", blocks: [] };
  const defaultFooter = { enabled: false, height: 100, background: "#1a2740", blocks: [] };
  const routePaths = routes.length > 0 ? routes.map((r) => ensurePath(r.path)) : ["/"];
  try {
    const parsed = JSON.parse(source);
    const pages = {};
    if (parsed.pages && typeof parsed.pages === "object") {
      Object.entries(parsed.pages).forEach(([path, page]) => {
        pages[ensurePath(path)] = {
          title: page.title || "Pagina",
          blocks: (page.blocks || []).map((block, idx) => ({
            ...block,
            id: block.id || nextId("blk"),
            style: { ...DEFAULT_STYLE, ...(block.style || {}) },
            z: typeof block.z === "number" ? block.z : idx + 1,
            w: block.w > 0 ? block.w : 220,
            h: block.h > 0 ? block.h : 90,
          })),
        };
      });
    } else if (Array.isArray(parsed.blocks)) {
      pages["/"] = {
        title: parsed.title || "Pagina Inicial",
        blocks: parsed.blocks.map((block, idx) => ({
          ...block,
          id: block.id || nextId("blk"),
          style: { ...DEFAULT_STYLE, ...(block.style || {}) },
          z: typeof block.z === "number" ? block.z : idx + 1,
          w: block.w > 0 ? block.w : 220,
          h: block.h > 0 ? block.h : 90,
        })),
      };
    }
    routePaths.forEach((path) => {
      if (!pages[path]) {
        pages[path] = { title: path === "/" ? "Pagina Inicial" : `Pagina ${path}`, blocks: [] };
      }
    });
    return {
      title: parsed.title || "Site",
      canvas: {
        width: Math.min(2800, Math.max(900, parsed.canvas?.width || 1400)),
        height: Math.min(2800, Math.max(700, parsed.canvas?.height || 900)),
        background: parsed.canvas?.background || "#f8f9ff",
      },
      header: parsed.header || defaultHeader,
      footer: parsed.footer || defaultFooter,
      pages,
      popups: parsed.popups || {},
    };
  } catch {
    return null;
  }
}

export function makeLandingTemplate() {
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
        { id: nextId("blk"), type: "button", label: "Comecar gratis", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "var(--jx-color-primary)", color: "#ffffff", padding: "10px 20px", "border-radius": "8px", border: "0", "font-weight": "600", "text-align": "center" }, x: 1150, y: 18, w: 170, h: 44, rotation: 0, z: 4 },
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
          { id: nextId("blk"), type: "button", label: "Comecar agora", action_type: "navigate", action_target: "/", href: "/", style: { ...DEFAULT_STYLE, background: "var(--jx-color-primary)", color: "#ffffff", padding: "14px 28px", "border-radius": "8px", border: "0", "font-weight": "600", "font-size": "16px", "text-align": "center" }, x: 80, y: 320, w: 200, h: 54, rotation: 0, z: 3 },
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
