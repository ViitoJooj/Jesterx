export function FooterBlock() {
    return (
        <footer
            style={{
                padding: "32px 48px",
                background: "#111827",
                color: "#e5e7eb",
                borderRadius: 20,
            }}
        >
            <div
                style={{
                    maxWidth: 1024,
                    margin: "0 auto",
                    display: "flex",
                    justifyContent: "space-between",
                    gap: 24,
                    flexWrap: "wrap",
                }}
            >
                <div>
                    <strong>Minha marca</strong>
                    <p style={{ marginTop: 6, opacity: 0.8 }}>
                        Construindo experiências digitais
                    </p>
                </div>
                <nav style={{ display: "flex", gap: 16 }}>
                    <a>Produtos</a>
                    <a>Suporte</a>
                    <a>Contato</a>
                </nav>
            </div>
        </footer>
    );
}
