import { useMemo, useState } from "react";
import styles from "../styles/components/CreatePageForm.module.scss";
import { CustomSelect } from "./Select";

import laddingPage from "../imgs/lading-Page.png";
import ecommerceImg from "../imgs/ecommerce.png";
import softwareSellImg from "../imgs/softwareSell.png";
import videoPage from "../imgs/videos.png";

type PageType = "landing" | "ecommerce" | "software" | "video";

type PageOption = {
  id: string;
  title: string;
};

type PageTypeMeta = {
  id: PageType;
  title: string;
  image: string;
  defaultPages: PageOption[];
};

const pageTypes: PageTypeMeta[] = [
  {
    id: "landing",
    title: "Landing Page",
    image: laddingPage,
    defaultPages: [
      { id: "main", title: "Página única" },
      { id: "contact", title: "Seção de contato" },
      { id: "cta", title: "Call to action" },
      { id: "faq", title: "FAQ" },
    ],
  },
  {
    id: "ecommerce",
    title: "E-commerce",
    image: ecommerceImg,
    defaultPages: [
      { id: "home", title: "Home" },
      { id: "products", title: "Produtos" },
      { id: "product", title: "Produto" },
      { id: "cart", title: "Carrinho" },
      { id: "checkout", title: "Checkout" },
      { id: "about", title: "Sobre" },
      { id: "search", title: "Busca" },
      { id: "blog", title: "Blog" },
    ],
  },
  {
    id: "software",
    title: "Software",
    image: softwareSellImg,
    defaultPages: [
      { id: "main", title: "Página principal" },
      { id: "download", title: "Download" },
      { id: "pricing", title: "Preços" },
      { id: "about", title: "Sobre" },
      { id: "docs", title: "Documentação" },
    ],
  },
  {
    id: "video",
    title: "Vídeo Page",
    image: videoPage,
    defaultPages: [
      { id: "gallery", title: "Galeria" },
      { id: "player", title: "Player" },
      { id: "about", title: "Sobre" },
    ],
  },
];

type CreatePageFormProps = {
  onClose: () => void;
  onSuccess?: (data: any) => void;
};

export function CreatePageForm({ onClose, onSuccess }: CreatePageFormProps) {
  const [step, setStep] = useState<"choose" | "details">("choose");
  const [selectedType, setSelectedType] = useState<PageType | "">("");
  const [name, setName] = useState("");
  const [domain, setDomain] = useState("");
  const [description, setDescription] = useState("");
  const [goal, setGoal] = useState("");
  const [logoUrl, setLogoUrl] = useState("");
  const [pages, setPages] = useState<string[]>([]);

  const selected = useMemo(() => pageTypes.find((p) => p.id === selectedType), [selectedType]);

  function handleTypeSelect(type: PageType) {
    setSelectedType(type);
    setPages(pageTypes.find((p) => p.id === type)?.defaultPages.map((p) => p.id) || []);
    setStep("details");
  }

  function handleTogglePage(id: string) {
    setPages((prev) =>
      prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]
    );
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    const payload = {
      name,
      domain,
      description,
      goal,
      logoUrl,
      page_type: selectedType,
      pages,
    };

    console.log("Payload para envio:", payload);
    onSuccess?.(payload);
  }

  return (
    <>
      <div className={styles.shadow} onClick={onClose} />

      <div className={`${styles.main} ${step === "details" ? styles.aside : styles.center}`}>
        <div className={styles.headerBar}>
          {step === "details" ? (
            <button onClick={() => setStep("choose")} className={styles.iconButton}>{"<"}</button>
          ) : (
            <span className={styles.iconPlaceholder} />
          )}
          <h1>{step === "choose" ? "Escolha o tipo de site" : "Detalhes do site"}</h1>
          <button onClick={onClose} className={styles.iconButtonClose}>X</button>
        </div>

        <div className={styles.body}>
          {step === "choose" && (
            <div className={styles.typeSiteContainer}>
              {pageTypes.map((page) => (
                <button key={page.id} onClick={() => handleTypeSelect(page.id)} className={styles.pageTypeCard}>
                  <img src={page.image} alt={page.title} />
                  <h2>{page.title}</h2>
                </button>
              ))}
            </div>
          )}

          {step === "details" && (
            <form className={styles.detailsForm} onSubmit={handleSubmit}>
              <label className={styles.field}>
                Nome do site
                <input value={name} onChange={(e) => setName(e.target.value)} required />
              </label>

              <label className={styles.field}>
                Domínio
                <input value={domain} onChange={(e) => setDomain(e.target.value)} />
              </label>

              <label className={styles.field}>
                Descrição
                <input value={description} onChange={(e) => setDescription(e.target.value)} required />
              </label>

              <label className={styles.field}>
                URL do Logo
                <input value={logoUrl} onChange={(e) => setLogoUrl(e.target.value)} />
              </label>

              <CustomSelect
                value={goal}
                onChange={setGoal}
                placeholder="Selecione o objetivo"
                options={[
                  { value: "leads", label: "Capturar leads" },
                  { value: "sales", label: "Vender" },
                  { value: "brand", label: "Apresentar marca" },
                ]}
              />

              <fieldset className={styles.pageCheckboxes}>
                <legend>Páginas do site</legend>
                {selected?.defaultPages.map((p) => (
                  <label key={p.id}>
                    <input type="checkbox" checked={pages.includes(p.id)} onChange={() => handleTogglePage(p.id)} />
                    {p.title}
                  </label>
                ))}
              </fieldset>

              <div className={styles.actions}>
                <button type="submit" className={styles.primary}>Criar site</button>
              </div>
            </form>
          )}
        </div>
      </div>
    </>
  );
}
