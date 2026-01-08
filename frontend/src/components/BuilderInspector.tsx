import styles from "../styles/components/BuilderInspector.module.scss";

export function BuilderInspector({ block, onChange, onSave, saving, onDelete }: any) {
    if (!block) {
        return (
            <aside className={styles.inspector}>
                <div className={styles.emptyState}>
                    <div className={styles.emptyIcon}>⚙️</div>
                    <p className={styles.emptyText}>
                        Selecione um bloco no canvas para editar suas propriedades
                    </p>
                </div>
            </aside>
        );
    }

    const getBlockDisplayName = (type: string) => {
        const names: Record<string, string> = {
            hero: "Hero Banner",
            products: "Grade de Produtos",
            cta: "Captura de Leads",
        };
        return names[type] || type;
    };

    return (
        <aside className={styles.inspector}>
            <div className={styles.header}>
                <h3 className={styles.title}>Propriedades do Bloco</h3>
                <p className={styles.blockType}>{getBlockDisplayName(block.type)}</p>
            </div>

            <div className={styles.section}>
                <h4 className={styles.sectionTitle}>Conteúdo</h4>

                {block.type === "hero" && (
                    <>
                        <div className={styles.field}>
                            <label className={styles.label}>
                                Título
                            </label>
                            <input
                                className={styles.input}
                                type="text"
                                value={block.props.title || ""}
                                placeholder="Digite o título do hero"
                                onChange={(e) =>
                                    onChange({
                                        ...block,
                                        props: { ...block.props, title: e.target.value },
                                    })
                                }
                            />
                        </div>

                        <div className={styles.field}>
                            <label className={styles.label}>
                                Subtítulo
                            </label>
                            <textarea
                                className={styles.textarea}
                                value={block.props.subtitle || ""}
                                placeholder="Digite o subtítulo ou descrição"
                                onChange={(e) =>
                                    onChange({
                                        ...block,
                                        props: { ...block.props, subtitle: e.target.value },
                                    })
                                }
                            />
                        </div>

                        <div className={styles.field}>
                            <label className={styles.label}>
                                Texto do botão
                            </label>
                            <input
                                className={styles.input}
                                type="text"
                                value={block.props.buttonText || ""}
                                placeholder="Ex: Começar agora"
                                onChange={(e) =>
                                    onChange({
                                        ...block,
                                        props: { ...block.props, buttonText: e.target.value },
                                    })
                                }
                            />
                        </div>
                    </>
                )}

                {block.type === "cta" && (
                    <>
                        <div className={styles.field}>
                            <label className={styles.label}>
                                Título
                            </label>
                            <input
                                className={styles.input}
                                type="text"
                                value={block.props.title || ""}
                                placeholder="Digite o título do CTA"
                                onChange={(e) =>
                                    onChange({
                                        ...block,
                                        props: { ...block.props, title: e.target.value },
                                    })
                                }
                            />
                        </div>

                        <div className={styles.field}>
                            <label className={styles.label}>
                                Descrição
                            </label>
                            <textarea
                                className={styles.textarea}
                                value={block.props.description || ""}
                                placeholder="Digite a descrição"
                                onChange={(e) =>
                                    onChange({
                                        ...block,
                                        props: { ...block.props, description: e.target.value },
                                    })
                                }
                            />
                        </div>
                    </>
                )}

                {block.type === "products" && (
                    <p className={styles.emptyText}>
                        Este bloco exibe automaticamente os produtos cadastrados.
                    </p>
                )}
            </div>

            <div className={styles.actions}>
                <button 
                    className={styles.saveButton}
                    onClick={onSave} 
                    disabled={saving}
                >
                    {saving ? "Salvando…" : "💾 Salvar página"}
                </button>
                
                {onDelete && (
                    <button 
                        className={styles.deleteButton}
                        onClick={() => onDelete(block.id)}
                    >
                        🗑️ Excluir bloco
                    </button>
                )}
            </div>
        </aside>
    );
}
