import { useMemo, useState } from "react";
import styles from "../styles/components/CreatePageForm.module.scss";
import { post } from "../utils/api";

import laddingPage from "../imgs/lading-Page.png";
import ecommerceImg from "../imgs/ecommerce.png";
import softwareSellImg from "../imgs/ecommerce.png";
import videoPage from "../imgs/videos.png";

type PageType = {
  id: string;
  title: string;
  image: string;
};

const pageTypes: PageType[] = [
  { id: "landing", title: "Landing Page", image: laddingPage },
  { id: "ecommerce", title: "E-commerce", image: ecommerceImg },
  { id: "software", title: "Software Sell", image: softwareSellImg },
  { id: "video", title: "Video page", image: videoPage },
];

type CreatePageFormProps = {
  onClose: () => void;
  onSuccess?: () => void;
};

export function CreatePageForm({ onClose, onSuccess }: CreatePageFormProps) {
  const [storeType, setStoreType] = useState<string>("");
  const [step, setStep] = useState<"choose" | "details">("choose");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>("");

  const [pageName, setPageName] = useState("");
  const [domain, setDomain] = useState("");
  const [goal, setGoal] = useState("");

  const selected = useMemo(() => pageTypes.find((p) => p.id === storeType), [storeType]);

  function handleSelectType(type: string) {
    setStoreType(type);
    setStep("details");
  }

  function handleBack() {
    setStep("choose");
    setError("");
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setError("");
    setLoading(true);

    try {
      const payload: any = {
        name: pageName,
        page_type: storeType,
      };

      if (domain) payload.domain = domain;
      if (goal) payload.goal = goal;

      const response = await post("/v1/pages", payload);

      if (response) {
        onSuccess?.();
      }
    } catch (err: any) {
      setError(err?.message || "Erro ao criar página");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <div className={styles.shadow} onClick={onClose} />

      <div className={`${styles.main} ${step === "details" ? styles.aside : styles.center}`} role="dialog" aria-modal="true">
        <div className={styles.headerBar}>
          {step === "details" ? (
            <button type="button" onClick={handleBack} className={styles.iconButton} aria-label="Voltar">
              {"<"}
            </button>
          ) : (
            <span className={styles.iconPlaceholder} />
          )}

          <div className={styles.headerTitleWrap}>
            <h1 className={styles.headerTitle}>{step === "details" ? "Detalhes da página" : "Qual o tipo de página?"}</h1>
            {step === "details" && selected?.title ? <p className={styles.headerSubtitle}>{selected.title}</p> : null}
          </div>

          <button type="button" onClick={onClose} className={styles.iconButtonClose} aria-label="Fechar modal">
            X
          </button>
        </div>

        <div className={styles.body}>
          <div className={`${styles.step} ${step === "choose" ? styles.stepActive : styles.stepHidden}`}>
            <div className={styles.typeSiteContainer}>
              {pageTypes.map((page) => (
                <button key={page.id} type="button" className={`${styles.pageTypeCard} ${storeType === page.id ? styles.active : ""}`} onClick={() => handleSelectType(page.id)}>
                  <img src={page.image} alt={`${page.title} icon`} />
                  <h2>{page.title}</h2>
                </button>
              ))}
            </div>

            <p className={styles.RecommendTypePage}>
              Não encontrou seu tipo? <a href="/suporte">Nos recomende a sua! </a>
            </p>
          </div>

          <div className={`${styles.step} ${step === "details" ? styles.stepActive : styles.stepHidden}`}>
            <form className={styles.detailsForm} onSubmit={handleSubmit}>
              {error && <div style={{ color: "#e74c3c", marginBottom: "1rem", padding: "0.5rem", background: "#fee", borderRadius: "4px" }}>{error}</div>}

              <label className={styles.field}>
                Nome do projeto
                <input type="text" placeholder="Ex:  Minha Landing" value={pageName} onChange={(e) => setPageName(e.target.value)} required disabled={loading} />
              </label>

              <label className={styles.field}>
                Domínio (opcional)
                <input type="text" placeholder="ex: meusite.com.br" value={domain} onChange={(e) => setDomain(e.target.value)} disabled={loading} />
              </label>

              <label className={styles.field}>
                Objetivo
                <select value={goal} onChange={(e) => setGoal(e.target.value)} disabled={loading}>
                  <option value="">Selecione…</option>
                  <option value="leads">Capturar leads</option>
                  <option value="sales">Vender</option>
                  <option value="brand">Apresentar marca</option>
                </select>
              </label>

              <div className={styles.actions}>
                <button type="button" className={styles.secondary} onClick={handleBack} disabled={loading}>
                  Alterar tipo
                </button>
                <button type="submit" className={styles.primary} disabled={loading}>
                  {loading ? "Criando..." : "Criar página"}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </>
  );
}
