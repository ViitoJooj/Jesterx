import styles from "../styles/components/ThemeCard.module.scss";

type Props = {
    name: string;
    description: string;
    thumbnail?: string;
    rating?: number | null;
    installs?: number | null;
    onClick: () => void;
};

export function ThemeCard({
    name,
    description,
    thumbnail,
    rating,
    installs,
    onClick
}: Props) {
    return (
        <article className={styles.card} onClick={onClick}>
            <div className={styles.preview}>
                {thumbnail ? <img src={thumbnail} alt={name} /> : <div className={styles.fallback} />}
                <div className={styles.overlay}>
                    <span>Ver tema</span>
                </div>
            </div>

            <div className={styles.body}>
                <h3>{name}</h3>
                <p>{description}</p>
            </div>

            <footer className={styles.footer}>
                <span>⭐ {typeof rating === "number" ? rating.toFixed(1) : "—"}</span>
                <span>{installs ?? 0} installs</span>
            </footer>
        </article>
    );
}
