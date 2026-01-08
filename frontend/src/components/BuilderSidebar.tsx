import styles from "../styles/components/BuilderSidebar.module.scss";
import { nanoid } from "nanoid";

const blocks = [
    { type: "hero", label: "Hero banner" },
    { type: "products", label: "Grade de produtos" },
    { type: "cta", label: "Captura de leads" },
];

export function BuilderSidebar({ onAdd }: any) {
    return (
        <aside className={styles.sidebar}>
            <h3>Seções</h3>

            {blocks.map((block) => (
                <button
                    key={block.type}
                    onClick={() =>
                        onAdd({
                            id: nanoid(),
                            type: block.type,
                            props: {},
                        })
                    }
                >
                    + {block.label}
                </button>
            ))}
        </aside>
    );
}
