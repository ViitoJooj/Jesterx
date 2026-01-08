import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import styles from "../styles/pages/ThemeProduct.module.scss";
import { get } from "../utils/api";

type ThemeProductData = {
    id: string;
    name: string;
    description: string;
    images: string[];
    rating: number;
    installs: number;
    long_description: string;
    page_id: string;
    domain: string;
};

export function ThemeProduct() {
    const { slug } = useParams();
    const navigate = useNavigate();
    const [theme, setTheme] = useState<ThemeProductData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    const load = useCallback(async () => {
        try {
            setLoading(true);
            setError("");
            const res = await get<ThemeProductData>(`/v1/themes/store/${slug}`);
            if (res.success && res.data) {
                setTheme(res.data);
            } else {
                setError("Tema não encontrado");
            }
        } catch (err: any) {
            setError(err?.message || "Erro ao carregar tema");
        } finally {
            setLoading(false);
        }
    }, [slug]);

    useEffect(() => {
        load();
    }, [load]);

    if (loading) {
        return (
            <div className={styles.loading}>
                <p>Carregando tema...</p>
            </div>
        );
    }

    if (error || !theme) {
        return (
            <div className={styles.error}>
                <p className={styles.errorText}>{error || "Tema não encontrado"}</p>
                <button className={styles.backButton} onClick={() => navigate("/temas")}>
                    ← Voltar para loja de temas
                </button>
            </div>
        );
    }

    return (
        <main className={styles.main}>
            <div className={styles.content}>
                <div className={styles.gallery}>
                    {theme.images.map((img) => (
                        <img 
                            key={img} 
                            src={img} 
                            alt={`${theme.name} - Preview`}
                            className={styles.galleryImage}
                        />
                    ))}
                </div>

                <div className={styles.info}>
                    <h1 className={styles.title}>{theme.name}</h1>
                    
                    <p className={styles.description}>{theme.description}</p>

                    <div className={styles.meta}>
                        <div className={styles.metaItem}>
                            <span className={styles.metaIcon}>⭐</span>
                            <span>{theme.rating.toFixed(1)}</span>
                        </div>
                        <div className={styles.metaItem}>
                            <span className={styles.metaIcon}>👥</span>
                            <span>{theme.installs} instalações</span>
                        </div>
                    </div>

                    <button className={styles.primaryButton}>
                        🎨 Usar este tema
                    </button>
                </div>
            </div>

            <section className={styles.details}>
                <h2 className={styles.detailsTitle}>Sobre o tema</h2>
                <p className={styles.detailsText}>{theme.long_description}</p>

                <div className={styles.features}>
                    <h3 className={styles.detailsTitle}>Recursos inclusos</h3>
                    <div className={styles.featuresGrid}>
                        <div className={styles.featureItem}>
                            <span className={styles.featureIcon}>📱</span>
                            <div className={styles.featureContent}>
                                <h4 className={styles.featureTitle}>Design Responsivo</h4>
                                <p className={styles.featureDescription}>
                                    Funciona perfeitamente em todos os dispositivos
                                </p>
                            </div>
                        </div>
                        <div className={styles.featureItem}>
                            <span className={styles.featureIcon}>⚡</span>
                            <div className={styles.featureContent}>
                                <h4 className={styles.featureTitle}>Performance Otimizada</h4>
                                <p className={styles.featureDescription}>
                                    Carregamento rápido e otimizado para SEO
                                </p>
                            </div>
                        </div>
                        <div className={styles.featureItem}>
                            <span className={styles.featureIcon}>🎨</span>
                            <div className={styles.featureContent}>
                                <h4 className={styles.featureTitle}>Personalizável</h4>
                                <p className={styles.featureDescription}>
                                    Fácil de customizar cores e conteúdo
                                </p>
                            </div>
                        </div>
                        <div className={styles.featureItem}>
                            <span className={styles.featureIcon}>🔒</span>
                            <div className={styles.featureContent}>
                                <h4 className={styles.featureTitle}>Seguro</h4>
                                <p className={styles.featureDescription}>
                                    Código limpo e seguindo as melhores práticas
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}
