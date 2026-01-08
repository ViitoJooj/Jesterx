import styles from "../styles/components/BuilderCanvas.module.scss";
import { HeroBlock } from "./blocks/HeroBlock";
import { ProductsGridBlock } from "./blocks/ProductsGridBlock";
import { CtaBlock } from "./blocks/CtaBlock";

export function BuilderCanvas({ blocks, selected, onSelect }: any) {
    return (
        <main className={styles.canvas}>
            {blocks.map((block: any) => (
                <div
                    key={block.id}
                    className={`${styles.block} ${selected === block.id ? styles.active : ""
                        }`}
                    onClick={() => onSelect(block.id)}
                >
                    {block.type === "hero" && <HeroBlock {...block.props} />}
                    {block.type === "products" && <ProductsGridBlock />}
                    {block.type === "cta" && <CtaBlock />}
                </div>
            ))}
        </main>
    );
}
