export function ProductsGridBlock() {
    return (
        <section
            style={{
                padding: "56px 48px",
                background: "#fff",
                borderRadius: 20,
                border: "1px solid #e5e7eb",
            }}
        >
            <div
                style={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    marginBottom: 24,
                }}
            >
                <h2 style={{ margin: 0, fontSize: 28 }}>Produtos em destaque</h2>
                <a style={{ color: "#f97316", fontWeight: 600 }}>Ver todos</a>
            </div>

            <div
                style={{
                    display: "grid",
                    gridTemplateColumns: "repeat(auto-fill,minmax(220px,1fr))",
                    gap: 20,
                }}
            >
                {Array.from({ length: 4 }).map((_, i) => (
                    <article
                        key={i}
                        style={{
                            border: "1px solid #e5e7eb",
                            borderRadius: 16,
                            padding: 14,
                            display: "flex",
                            flexDirection: "column",
                            gap: 10,
                        }}
                    >
                        <div
                            style={{
                                height: 140,
                                background: "#f3f4f6",
                                borderRadius: 12,
                                display: "grid",
                                placeItems: "center",
                                color: "#6b7280",
                            }}
                        >
                            Imagem
                        </div>
                        <strong>Produto exemplo</strong>
                        <span style={{ color: "#6b7280" }}>Descrição curta</span>
                        <strong>R$ 199,00</strong>
                        <button
                            style={{
                                marginTop: 8,
                                padding: 10,
                                background: "#111827",
                                color: "#fff",
                                border: "none",
                                borderRadius: 10,
                                cursor: "pointer",
                            }}
                        >
                            Adicionar ao carrinho
                        </button>
                    </article>
                ))}
            </div>
        </section>
    );
}
