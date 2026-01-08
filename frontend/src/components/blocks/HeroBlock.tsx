export function HeroBlock({ title, subtitle }: any) {
    return (
        <section
            style={{
                padding: 64,
                background: "linear-gradient(135deg,#111827,#1f2937)",
                color: "#fff",
                borderRadius: 16,
            }}
        >
            <h1>{title || "Título principal"}</h1>
            <p>{subtitle || "Descrição da seção"}</p>
        </section>
    );
}
