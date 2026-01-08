import { useState } from "react";
import styles from "../styles/components/BuilderCanvas.module.scss";
import { HeroBlock } from "./blocks/HeroBlock";
import { ProductsGridBlock } from "./blocks/ProductsGridBlock";
import { CtaBlock } from "./blocks/CtaBlock";

type PreviewMode = "desktop" | "tablet" | "mobile";

export function BuilderCanvas({ blocks, selected, onSelect }: any) {
    const [previewMode, setPreviewMode] = useState<PreviewMode>("desktop");

    const canvasWidth = 
        previewMode === "mobile" ? "375px" :
        previewMode === "tablet" ? "768px" :
        "100%";

    return (
        <main className={styles.canvas}>
            <div className={styles.toolbar}>
                <span className={styles.toolbarTitle}>
                    {blocks.length === 0 ? "Canvas vazio" : `${blocks.length} bloco(s)`}
                </span>
                <div className={styles.previewModes}>
                    <button
                        className={`${styles.modeButton} ${previewMode === "desktop" ? styles.active : ""}`}
                        onClick={() => setPreviewMode("desktop")}
                        title="Desktop"
                    >
                        🖥️
                    </button>
                    <button
                        className={`${styles.modeButton} ${previewMode === "tablet" ? styles.active : ""}`}
                        onClick={() => setPreviewMode("tablet")}
                        title="Tablet"
                    >
                        📟
                    </button>
                    <button
                        className={`${styles.modeButton} ${previewMode === "mobile" ? styles.active : ""}`}
                        onClick={() => setPreviewMode("mobile")}
                        title="Mobile"
                    >
                        📱
                    </button>
                </div>
            </div>

            <div className={styles.canvasInner} style={{ width: canvasWidth }}>
                {blocks.length === 0 ? (
                    <div className={styles.empty}>
                        <div className={styles.emptyIcon}>📄</div>
                        <h3 className={styles.emptyTitle}>Canvas vazio</h3>
                        <p className={styles.emptyText}>
                            Adicione blocos da barra lateral para começar a construir sua página
                        </p>
                    </div>
                ) : (
                    blocks.map((block: any) => (
                        <div
                            key={block.id}
                            className={`${styles.block} ${selected === block.id ? styles.active : ""}`}
                            onClick={() => onSelect(block.id)}
                        >
                            <span className={styles.blockLabel}>
                                {block.type === "hero" && "Hero"}
                                {block.type === "products" && "Produtos"}
                                {block.type === "cta" && "CTA"}
                            </span>
                            {block.type === "hero" && <HeroBlock {...block.props} />}
                            {block.type === "products" && <ProductsGridBlock />}
                            {block.type === "cta" && <CtaBlock />}
                        </div>
                    ))
                )}
            </div>
        </main>
    );
}
