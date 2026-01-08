import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import styles from "../styles/pages/PageEdit.module.scss";
import { get, put } from "../utils/api";
import { BuilderSidebar } from "../components/BuilderSidebar";
import { BuilderCanvas } from "../components/BuilderCanvas";
import { BuilderInspector } from "../components/BuilderInspector";

export type BlockInstance = {
    id: string;
    type: string;
    props: Record<string, any>;
};

export function PageEdit() {
    const { pageId } = useParams<{ pageId: string }>();

    const [loading, setLoading] = useState(true);
    const [blocks, setBlocks] = useState<BlockInstance[]>([]);
    const [selected, setSelected] = useState<string | null>(null);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (!pageId) return;

        (async () => {
            const res = await get<any>(`/v1/pages/${pageId}/raw`);
            setBlocks(res.data?.components || []);
            setLoading(false);
        })();
    }, [pageId]);

    async function save() {
        if (!pageId) return;
        setSaving(true);

        await put(`/v1/pages/${pageId}`, {
            components: blocks,
        });

        setSaving(false);
    }

    if (loading) return <div className={styles.loading}>Carregando editor…</div>;

    return (
        <div className={styles.builder}>
            <BuilderSidebar
                onAdd={(block: any) => setBlocks((prev) => [...prev, block])}
            />

            <BuilderCanvas
                blocks={blocks}
                selected={selected}
                onSelect={setSelected}
            />

            <BuilderInspector
                block={blocks.find((b) => b.id === selected) || null}
                onChange={(updated: any) =>
                    setBlocks((prev) =>
                        prev.map((b) => (b.id === updated.id ? updated : b))
                    )
                }
                onSave={save}
                saving={saving}
            />
        </div>
    );
}
