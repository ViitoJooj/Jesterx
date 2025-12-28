export type CommunityTheme = {
  id: string;
  name: string;
  description: string;
  accent: string;
  author: string;
  tags: string[];
  pageType: "landing" | "ecommerce" | "software" | "video";
  template: string;
  previewHtml: string;
};

function buildTemplate(accent: string, title: string, subtitle: string, cta: string, extra?: string) {
  return `
  <style>
    :root {
      --accent: ${accent};
      --bg: #0f172a;
      --muted: #e2e8f0;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: 'Inter', system-ui, -apple-system, sans-serif;
      background:
        linear-gradient(120deg, rgba(15,23,42,0.95), rgba(15,23,42,0.7)),
        radial-gradient(circle at 20% 20%, rgba(255,255,255,0.06), transparent 35%),
        radial-gradient(circle at 80% 0%, rgba(255,255,255,0.04), transparent 25%),
        #0b1220;
      color: #f8fafc;
    }
    .page { min-height: 100vh; padding: 64px 24px; max-width: 1024px; margin: 0 auto; }
    .badge { display: inline-flex; align-items: center; gap: 8px; padding: 8px 12px; border-radius: 999px; background: rgba(255,255,255,0.08); color: #cbd5e1; font-weight: 600; letter-spacing: 0.02em; }
    .badge svg { width: 18px; height: 18px; }
    h1 { font-size: clamp(2.25rem, 5vw, 3.5rem); margin: 16px 0 12px; line-height: 1.05; }
    p.lead { color: #cbd5e1; font-size: 1.125rem; line-height: 1.7; max-width: 720px; }
    .actions { margin: 28px 0 40px; display: flex; flex-wrap: wrap; gap: 12px; }
    .primary { background: var(--accent); border: none; color: #0b1220; padding: 14px 18px; border-radius: 12px; font-weight: 700; cursor: pointer; box-shadow: 0 18px 40px rgba(0,0,0,0.25); }
    .secondary { background: transparent; border: 2px solid rgba(255,255,255,0.14); color: #e2e8f0; padding: 14px 18px; border-radius: 12px; font-weight: 700; cursor: pointer; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 16px; margin-top: 24px; }
    .card { background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.08); border-radius: 16px; padding: 18px; box-shadow: 0 12px 30px rgba(0,0,0,0.2); }
    .card h3 { margin: 0 0 8px; color: #f8fafc; }
    .card p { margin: 0; color: #cbd5e1; line-height: 1.6; }
    .tag { display: inline-block; padding: 6px 10px; border-radius: 999px; background: rgba(255,255,255,0.08); color: #cbd5e1; font-size: 0.85rem; }
  </style>
  <div class="page">
    <span class="badge"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M12 3l2.29 4.64L19 8.35l-3.5 3.41.83 4.84L12 14.77l-4.33 1.83.83-4.84L5 8.35l4.71-.71L12 3z"/></svg> Tema pronto para usar</span>
    <h1>${title}</h1>
    <p class="lead">${subtitle}</p>
    <div class="actions">
      <button class="primary">${cta}</button>
      <button class="secondary">Explorar coleção</button>
    </div>
    <div class="grid">
      <div class="card">
        <h3>Seção hero</h3>
        <p>Apresente sua oferta com um título forte e CTA direto.</p>
      </div>
      <div class="card">
        <h3>Cards destacados</h3>
        <p>Mostre features, serviços ou aulas da comunidade.</p>
      </div>
      <div class="card">
        <h3>Credibilidade</h3>
        <p>Espaço para depoimentos, métricas ou selo de confiança.</p>
      </div>
      <div class="card tag">${extra ?? "Layout flexível para páginas e lojas"}</div>
    </div>
  </div>
  `;
}

export const communityThemes: CommunityTheme[] = [
  {
    id: "minimal",
    name: "Minimal Studio",
    description: "Tema limpo, com tipografia forte e foco em copy curta.",
    accent: "#111827",
    author: "Equipe Jester",
    tags: ["landing", "portfolio", "produtos"],
    pageType: "landing",
    template: buildTemplate("#111827", "Lance sua ideia com clareza", "Hero enxuto, cards objetivos e CTA sempre visível para converter.", "Começar agora"),
    previewHtml: buildTemplate("#111827", "Lance sua ideia com clareza", "Hero enxuto, cards objetivos e CTA sempre visível para converter.", "Começar agora"),
  },
  {
    id: "elegant",
    name: "Elegant Commerce",
    description: "Grid elegante para produtos digitais ou físicos.",
    accent: "#7C3AED",
    author: "Comunidade",
    tags: ["ecommerce", "produtos", "saas"],
    pageType: "ecommerce",
    template: buildTemplate("#7C3AED", "Catálogo premium sem esforço", "Cards com foto, preço e ação rápida para sua loja full digital.", "Ver vitrines"),
    previewHtml: buildTemplate("#7C3AED", "Catálogo premium sem esforço", "Cards com foto, preço e ação rápida para sua loja full digital.", "Ver vitrines"),
  },
  {
    id: "bold",
    name: "Bold Launch",
    description: "Cores fortes, CTA destacado e seção de métricas.",
    accent: "#F97316",
    author: "Comunidade",
    tags: ["startup", "landing", "curso"],
    pageType: "landing",
    template: buildTemplate("#F97316", "Lance sem complicar", "Página direta com métricas, CTA fixo e espaço para depoimentos.", "Publicar página"),
    previewHtml: buildTemplate("#F97316", "Lance sem complicar", "Página direta com métricas, CTA fixo e espaço para depoimentos.", "Publicar página"),
  },
  {
    id: "neo",
    name: "Neo Cards",
    description: "Layout em cards com vibe moderna e animada.",
    accent: "#22C55E",
    author: "Comunidade",
    tags: ["comunidade", "portfolio", "coleções"],
    pageType: "software",
    template: buildTemplate("#22C55E", "Mostre suas criações em blocos", "Ideal para coleções de temas, aulas ou showcases feitos pela comunidade.", "Clonar tema"),
    previewHtml: buildTemplate("#22C55E", "Mostre suas criações em blocos", "Ideal para coleções de temas, aulas ou showcases feitos pela comunidade.", "Clonar tema"),
  },
];

export function getThemeById(id: string) {
  return communityThemes.find((theme) => theme.id === id);
}
