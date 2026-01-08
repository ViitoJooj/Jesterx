import styles from "../styles/components/BuilderInspector.module.scss";

export function BuilderInspector({ block, onChange, onSave, saving }: any) {
    if (!block) {
        return (
            <aside className={styles.inspector}>
                <p>Selecione uma seção</p>
            </aside>
        );
    }

    return (
        <aside className={styles.inspector}>
            <h3>Configurações</h3>

            {block.type === "hero" && (
                <>
                    <label>
                        Título
                        <input
                            value={block.props.title || ""}
                            onChange={(e) =>
                                onChange({
                                    ...block,
                                    props: { ...block.props, title: e.target.value },
                                })
                            }
                        />
                    </label>

                    <label>
                        Subtítulo
                        <input
                            value={block.props.subtitle || ""}
                            onChange={(e) =>
                                onChange({
                                    ...block,
                                    props: { ...block.props, subtitle: e.target.value },
                                })
                            }
                        />
                    </label>
                </>
            )}

            <button onClick={onSave} disabled={saving}>
                {saving ? "Salvando…" : "Salvar página"}
            </button>
        </aside>
    );
}
