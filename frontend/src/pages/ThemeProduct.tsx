import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import styles from "../styles/pages/ThemeProduct.module.scss";
import { get } from "../utils/api";

type ThemeProduct = {
    id: string;
    name: string;
    description: string;
    images: string[];
    rating: number;
    installs: number;
    long_description: string;
};

export function ThemeProduct() {
    const { slug } = useParams();
    const [theme, setTheme] = useState<ThemeProduct | null>(null);

    useEffect(() => {
        load();
    }, [slug]);

    async function load() {
        const res = await get<ThemeProduct>(`/v1/themes/store/${slug}`);
        if (res.success && res.data) setTheme(res.data);
    }

    if (!theme) return <p>Carregando tema...</p>;

    return (
        <main className={styles.main}>
            <div className={styles.gallery}>
                {theme.images.map((img) => (
                    <img key={img} src={img} alt={theme.name} />
                ))}
            </div>

            <section className={styles.info}>
                <h1>{theme.name}</h1>

                <div className={styles.meta}>
                    <span>⭐ {theme.rating.toFixed(1)}</span>
                    <span>👥 {theme.installs} lojas usando</span>
                </div>

                <p className={styles.description}>{theme.description}</p>

                <button className={styles.primary}>
                    Usar este tema
                </button>

                <article className={styles.details}>
                    <h2>Sobre o tema</h2>
                    <p>{theme.long_description}</p>
                </article>
            </section>
        </main>
    );
}
