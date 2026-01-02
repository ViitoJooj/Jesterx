export function ThemeStore() {

    const apiResponse = {
        themes: [
            { id: 1, name: "Tema Moderno", price: "R$ 29,00", description: "Um tema moderno e responsivo para seu site.", purchased: false, type: "landing page", images: ["image1", "image2"] },
            { id: 2, name: "Tema Clássico", price: "R$ 19,00", description: "Um tema clássico e elegante para blogs.", purchased: false, type: "ecommerce", images: ["image3", "image4"] },
            { id: 3, name: "Tema Minimalista", price: "R$ 39,00", description: "Um tema minimalista focado em conteúdo.", purchased: true, type: "video player", images: ["image5", "image6"] },
            { id: 4, name: "Tema Fotografia", price: "R$ 49,00", description: "Um tema perfeito para portfólios de fotografia.", purchased: false, type: "software sell", images: ["image7", "image8"] },
        ]
    }

    return (
        <>
            <h1>Bem vindo a Loja de temas</h1>
            <p>Aqui você pode navegar e adquirir temas para personalizar sua experiência.</p>
            <button>Crie ou poste seu tema !</button>

            <div>
                <label>Filtrar por ordem</label>
                <select>
                    <option value="price-asc">Preço: Menor para Maior</option>
                    <option value="price-desc">Preço: Maior para Menor</option>
                    <option value="name-asc">Nome: A-Z</option>
                    <option value="name-desc">Nome: Z-A</option>
                </select>
            </div>

            <div>
                <label>Quantidade por página</label>
                <select>
                    <option value="10">20</option>
                    <option value="20">50</option>
                    <option value="50">100</option>
                </select>
            </div>

            <div>
                <div>
                    <button>Cards</button>
                    <button>Lista</button>
                </div>
            </div>

            <div>
                {apiResponse.themes.map((theme) => (
                    <div key={theme.id}>
                        <h2>{theme.name}</h2>
                        <p>{theme.description}</p>
                        <p>Preço: {theme.price}</p>
                        <div>
                            {theme.images.map((img, index) => (
                                <img key={index} src={img} alt={theme.name} />
                            ))}
                        </div>
                        <button>
                            {theme.purchased ? "Instalar" : "Comprar"}
                        </button>
                    </div>
                ))}
            </div>
        </>
    )
}