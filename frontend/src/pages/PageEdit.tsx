import { useEffect, useState, FormEvent } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../styles/pages/PageEdit.module.scss";
import { get, put } from "../utils/api";

type PageMeta = {
    id: string;
    tenant_id: string;
    name: string;
    page_id: string;
    domain?: string;
    theme_id?: string;
    created_at: string;
    updated_at: string;
};

type RawPage = {
    id: string;
    tenant_id: string;
    page_id: string;
    svelte: string;
    header?: string;
    footer?: string;
    show_header?: boolean;
    show_footer?: boolean;
    components?: string[];
    created_at: string;
    updated_at: string;
};

const elementorBlocks = [
    {
        id: "hero-ecommerce",
        label: "Hero para loja",
        description: "Banner inicial com destaque e CTA.",
        snippet: `<section style="padding:48px;background:linear-gradient(135deg,#111827,#1f2937);color:#fff;font-family:Arial,sans-serif;">
  <div style="max-width:960px;margin:0 auto;display:flex;gap:28px;align-items:center;flex-wrap:wrap;">
    <div style="flex:1 1 320px;">
      <p style="letter-spacing:1px;text-transform:uppercase;font-size:12px;opacity:.8;">Sua marca</p>
      <h1 style="font-size:36px;margin:12px 0;">Apresente seu melhor produto</h1>
      <p style="font-size:16px;opacity:.9;">Adicione textos, botões e imagens direto pelo editor.</p>
      <div style="margin-top:18px;display:flex;gap:12px;flex-wrap:wrap;">
        <a href="#produtos" style="padding:12px 18px;background:#f97316;color:#fff;border-radius:10px;text-decoration:none;">Ver produtos</a>
        <a href="#detalhes" style="padding:12px 18px;border:1px solid #fff;border-radius:10px;color:#fff;text-decoration:none;">Saiba mais</a>
      </div>
    </div>
    <div style="flex:1 1 280px;display:flex;justify-content:center;">
      <div style="width:260px;height:260px;border-radius:20px;background:#fff;display:grid;place-items:center;color:#111;font-weight:600;">Imagem do produto</div>
    </div>
  </div>
</section>`
    },
    {
        id: "grade-produtos",
        label: "Grade de produtos",
        description: "Cards para listar itens da loja.",
        snippet: `<section id="produtos" style="padding:42px;font-family:Arial,sans-serif;">
  <div style="max-width:1024px;margin:0 auto;">
    <div style="display:flex;justify-content:space-between;align-items:center;gap:12px;flex-wrap:wrap;">
      <h2 style="font-size:28px;margin:0;">Produtos em destaque</h2>
      <a href="#" style="color:#f97316;text-decoration:none;font-weight:600;">Ver todos</a>
    </div>
    <div style="margin-top:22px;display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:16px;">
      <article style="border:1px solid #e5e7eb;border-radius:12px;padding:14px;background:#fff;">
        <div style="height:140px;border-radius:10px;background:#f3f4f6;display:grid;place-items:center;color:#6b7280;">Imagem</div>
        <h3 style="margin:12px 0 6px;font-size:18px;">Produto exemplo</h3>
        <p style="margin:0;color:#6b7280;">Descrição curta.</p>
        <strong style="display:block;margin-top:10px;">R$ 199,00</strong>
        <button style="margin-top:12px;width:100%;padding:10px;border:none;background:#111827;color:#fff;border-radius:8px;cursor:pointer;">Adicionar ao carrinho</button>
      </article>
    </div>
  </div>
</section>`
    },
    {
        id: "cta-leads",
        label: "Captura de leads",
        description: "Formulário simples para landing pages.",
        snippet: `<section style="padding:38px;background:#0f172a;color:#fff;font-family:Arial,sans-serif;">
  <div style="max-width:760px;margin:0 auto;text-align:center;">
    <h2 style="margin:0 0 12px;font-size:30px;">Receba novidades e ofertas</h2>
    <p style="margin:0 0 18px;opacity:.9;">Adapte o formulário para seu CRM ou ferramenta favorita.</p>
    <form style="display:flex;gap:12px;flex-wrap:wrap;justify-content:center;">
      <input type="text" placeholder="Seu nome" style="padding:12px 14px;border-radius:10px;border:none;min-width:220px;" />
      <input type="email" placeholder="Seu email" style="padding:12px 14px;border-radius:10px;border:none;min-width:240px;" />
      <button type="submit" style="padding:12px 18px;border:none;background:#f97316;color:#fff;border-radius:10px;cursor:pointer;">Quero receber</button>
    </form>
  </div>
</section>`
    },
    {
        id: "rodape-simples",
        label: "Rodapé com links",
        description: "Rodapé enxuto para todas as páginas.",
        snippet: `<footer style="padding:28px;background:#111827;color:#e5e7eb;font-family:Arial,sans-serif;">
  <div style="max-width:1024px;margin:0 auto;display:flex;flex-wrap:wrap;gap:20px;justify-content:space-between;align-items:center;">
    <div>
      <strong>Minha marca</strong>
      <p style="margin:6px 0 0;font-size:14px;opacity:.8;">Texto sobre sua empresa.</p>
    </div>
    <nav style="display:flex;gap:14px;flex-wrap:wrap;">
      <a href="#produtos" style="color:#e5e7eb;text-decoration:none;">Produtos</a>
      <a href="#suporte" style="color:#e5e7eb;text-decoration:none;">Suporte</a>
      <a href="#contato" style="color:#e5e7eb;text-decoration:none;">Contato</a>
    </nav>
  </div>
</footer>`
    }
];

export function PageEdit() {
    const { pageId } = useParams<{ pageId: string }>();
    const navigate = useNavigate();

    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const [pageMeta, setPageMeta] = useState<PageMeta | null>(null);
    const [name, setName] = useState("");
    const [slug, setSlug] = useState("");
    const [domain, setDomain] = useState("");
    const [svelteCode, setSvelteCode] = useState("");
    const [headerContent, setHeaderContent] = useState("");
    const [footerContent, setFooterContent] = useState("");
    const [showHeader, setShowHeader] = useState(true);
    const [showFooter, setShowFooter] = useState(true);
    const [componentIds, setComponentIds] = useState<string[]>([]);
    const [showPreview, setShowPreview] = useState(false);

    useEffect(() => {
        if (!pageId) {
            setError("Página inválida.");
            setLoading(false);
            return;
        }

        (async () => {
            try {
                const [metaRes, rawRes] = await Promise.all([
                    get<PageMeta>(`/v1/pages/${pageId}`),
                    get<RawPage>(`/v1/pages/${pageId}/raw`),
                ]);

                if (metaRes.data) {
                    setPageMeta(metaRes.data);
                    setName(metaRes.data.name);
                    setSlug(metaRes.data.page_id);
                    setDomain(metaRes.data.domain || "");
                }

                if (rawRes.data) {
                    setSvelteCode(rawRes.data.svelte);
                    setHeaderContent(rawRes.data.header || "");
                    setFooterContent(rawRes.data.footer || "");
                    setShowHeader(rawRes.data.show_header ?? true);
                    setShowFooter(rawRes.data.show_footer ?? true);
                    setComponentIds(rawRes.data.components || []);
                }
            } catch (err: any) {
                const status = err?.status ?? err?.response?.status;
                if (status === 404) {
                    setError("Página não encontrada.");
                    return;
                }
                setError(err?.message || "Erro ao carregar página.");
            } finally {
                setLoading(false);
            }
        })();
    }, [pageId]);

    async function handleSubmit(e: FormEvent) {
        e.preventDefault();

        if (!pageId) return;

        setError("");
        setSuccess("");
        setSaving(true);

        try {
            await put(`/v1/pages/${pageId}`, {
                name,
                page_id: slug || undefined,
                domain: domain || undefined,
                svelte: svelteCode,
                header: headerContent,
                footer: footerContent,
                show_header: showHeader,
                show_footer: showFooter,
                components: componentIds,
            });

            setSuccess("Página atualizada com sucesso.");
        } catch (err: any) {
            setError(err?.message || "Erro ao salvar página.");
        } finally {
            setSaving(false);
        }
    }

    if (loading) {
        return (
            <main className={styles.main}>
                <div className={styles.loading}>Carregando editor...</div>
            </main>
        );
    }

    if (error && !pageMeta) {
        return (
            <main className={styles.main}>
                <p className={styles.error}>{error}</p>
                <button type="button" className={styles.secondaryButton} onClick={() => navigate("/pages")}>
                    Voltar para minhas páginas
                </button>
            </main>
        );
    }

    return (
        <main className={styles.main}>
            <header className={styles.header}>
                <h1 className={styles.title}>Editar página</h1>
                {pageMeta && (
                    <p className={styles.meta}>
                        ID: {pageMeta.id} • Slug: {pageMeta.page_id}
                    </p>
                )}
            </header>

            <form className={styles.form} onSubmit={handleSubmit}>
                {error && <div className={styles.errorBox}>{error}</div>}
                {success && <div className={styles.successBox}>{success}</div>}

                <label className={styles.field}>
                    <span className={styles.label}>Nome da página</span>
                    <input
                        className={styles.input}
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        disabled={saving}
                        required
                    />
                </label>

                <label className={styles.field}>
                    <span className={styles.label}>Rota / slug</span>
                    <input
                        className={styles.input}
                        type="text"
                        value={slug}
                        onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ""))}
                        disabled={saving}
                        placeholder="ex: minha-pagina"
                    />
                    <small>Essa rota será usada na URL pública.</small>
                </label>

                <label className={styles.field}>
                    <span className={styles.label}>Domínio (opcional)</span>
                    <input
                        className={styles.input}
                        type="text"
                        value={domain}
                        onChange={(e) => setDomain(e.target.value)}
                        disabled={saving}
                        placeholder="ex: meusite.com.br"
                    />
                </label>

                <div className={styles.inlineRow}>
                    <label className={styles.checkboxLabel}>
                        <input type="checkbox" checked={showHeader} onChange={(e) => setShowHeader(e.target.checked)} disabled={saving} />
                        <span>Exibir cabeçalho nesta página</span>
                    </label>
                    <label className={styles.checkboxLabel}>
                        <input type="checkbox" checked={showFooter} onChange={(e) => setShowFooter(e.target.checked)} disabled={saving} />
                        <span>Exibir rodapé nesta página</span>
                    </label>
                </div>

                <label className={styles.field}>
                    <span className={styles.label}>Conteúdo do cabeçalho</span>
                    <textarea
                        className={styles.textarea}
                        value={headerContent}
                        onChange={(e) => setHeaderContent(e.target.value)}
                        disabled={saving}
                        rows={6}
                    />
                </label>

                <label className={styles.field}>
                    <span className={styles.label}>Conteúdo do rodapé</span>
                    <textarea
                        className={styles.textarea}
                        value={footerContent}
                        onChange={(e) => setFooterContent(e.target.value)}
                        disabled={saving}
                        rows={6}
                    />
                </label>

                <section className={styles.blocksSection}>
                    <div className={styles.blocksHeader}>
                        <div>
                            <h2 className={styles.blocksTitle}>Componentes estilo Elementor</h2>
                            <p className={styles.blocksSubtitle}>Adicione blocos prontos para e-commerce, landing pages e capturas.</p>
                        </div>
                        {componentIds.length > 0 && (
                            <div className={styles.pillList}>
                                {componentIds.map((id) => (
                                    <span key={id} className={styles.pill}>
                                        {id}
                                    </span>
                                ))}
                            </div>
                        )}
                    </div>

                    <div className={styles.blocksGrid}>
                        {elementorBlocks.map((block) => (
                            <div key={block.id} className={styles.blockCard}>
                                <div className={styles.blockHeader}>
                                    <strong>{block.label}</strong>
                                    <span className={styles.blockBadge}>{block.id}</span>
                                </div>
                                <p className={styles.blockDescription}>{block.description}</p>
                                <div className={styles.blockActions}>
                                    <button
                                        type="button"
                                        className={styles.secondaryButton}
                                        onClick={() => {
                                            setSvelteCode((prev) => `${prev}\n\n${block.snippet}`);
                                            setComponentIds((prev) => Array.from(new Set([...prev, block.id])));
                                        }}
                                        disabled={saving}
                                    >
                                        Inserir no conteúdo
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                <label className={styles.field}>
                    <span className={styles.label}>Código Svelte</span>
                    <textarea
                        className={styles.textarea}
                        value={svelteCode}
                        onChange={(e) => setSvelteCode(e.target.value)}
                        disabled={saving}
                        rows={18}
                    />
                </label>

                <div className={styles.previewActions}>
                    <label className={styles.checkboxLabel}>
                        <input type="checkbox" checked={showPreview} onChange={(e) => setShowPreview(e.target.checked)} />
                        <span>Visualizar página renderizada</span>
                    </label>
                </div>

                {showPreview && (
                    <div className={styles.previewBox}>
                        <iframe
                            title="Pré-visualização"
                            className={styles.previewFrame}
                            sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
                            srcDoc={`${showHeader ? headerContent : ""}\n${svelteCode}\n${showFooter ? footerContent : ""}`}
                        />
                    </div>
                )}

                <div className={styles.actions}>
                    <button
                        type="button"
                        className={styles.secondaryButton}
                        onClick={() => navigate("/pages")}
                        disabled={saving}
                    >
                        Cancelar
                    </button>
                    <button type="submit" className={styles.primaryButton} disabled={saving}>
                        {saving ? "Salvando..." : "Salvar alterações"}
                    </button>
                </div>
            </form>
        </main>
    );
}
