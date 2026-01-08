import styles from "../styles/components/BuilderSidebar.module.scss";
import { nanoid } from "nanoid";

const blocks = [
    { type: "hero", label: "Hero banner", icon: "🎯" },
    { type: "products", label: "Grade de produtos", icon: "🛍️" },
    { type: "cta", label: "Captura de leads", icon: "📧" },
];

export function BuilderSidebar({ onAdd }: any) {
    return (
        <aside className={styles.sidebar}>
            <div className={styles.header}>
                <h3 className={styles.title}>Blocos disponíveis</h3>
                <p className={styles.subtitle}>Arraste ou clique para adicionar</p>
            </div>

            <div className={styles.blockList}>
                {blocks.map((block) => (
                    <button
                        key={block.type}
                        className={styles.blockItem}
                        onClick={() =>
                            onAdd({
                                id: nanoid(),
                                type: block.type,
                                props: {},
                            })
                        }
                    >
                        <span className={styles.blockIcon}>{block.icon}</span>
                        <span>{block.label}</span>
                    </button>
                ))}
            </div>

            <div className={styles.divider}></div>
            
            <div className={styles.header}>
                <p className={styles.subtitle}>💡 Dica: Selecione um bloco no canvas para editar suas propriedades</p>
            </div>
        </aside>
    );
}
