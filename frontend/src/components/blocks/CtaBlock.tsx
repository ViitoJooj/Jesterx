type Props = {
    title?: string;
    subtitle?: string;
    buttonText?: string;
};

export function CtaBlock({ title, subtitle, buttonText }: Props) {
    return (
        <section
            style={{
                padding: "64px 48px",
                background: "#0f172a",
                color: "#fff",
                borderRadius: 20,
                textAlign: "center",
            }}
        >
            <h2 style={{ fontSize: 32, margin: "0 0 12px" }}>
                {title || "Receba novidades e ofertas"}
            </h2>
            <p style={{ opacity: 0.9, marginBottom: 24 }}>
                {subtitle || "Cadastre seu email para receber promoções"}
            </p>
            <div
                style={{
                    display: "flex",
                    gap: 12,
                    justifyContent: "center",
                    flexWrap: "wrap",
                }}
            >
                <input
                    placeholder="Seu email"
                    style={{
                        padding: "12px 14px",
                        borderRadius: 10,
                        border: "none",
                        minWidth: 240,
                    }}
                />
                <button
                    style={{
                        padding: "12px 20px",
                        background: "#f97316",
                        color: "#fff",
                        border: "none",
                        borderRadius: 10,
                        cursor: "pointer",
                    }}
                >
                    {buttonText || "Quero receber"}
                </button>
            </div>
        </section>
    );
}
