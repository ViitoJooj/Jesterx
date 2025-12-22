import styles from "../styles/components/CreatePageForm.module.scss";

import laddingPage from "../imgs/lading-Page.png";
import ecommerceImg from "../imgs/ecommerce.png";
import softwareSellImg from "../imgs/softwareSell.png";
import videoPage from "../imgs/videos.png";

type PageType = {
  id: string;
  title: string;
  image: string;
};

const pageTypes: PageType[] = [
  {
    id: "landing",
    title: "Landing Page",
    image: laddingPage,
  },
  {
    id: "ecommerce",
    title: "E-commerce",
    image: ecommerceImg,
  },
  {
    id: "software",
    title: "Software Sell",
    image: softwareSellImg,
  },
  {
    id: "video",
    title: "Video page",
    image: videoPage,
  },
];

type CreatePageFormProps = {
  onClose: () => void;
  onSelectType?: (type: string) => void;
};

export function CreatePageForm({ onClose, onSelectType }: CreatePageFormProps) {
  return (
    <>
      <div className={styles.shadow} onClick={onClose} />

      <div className={styles.main} role="dialog" aria-modal="true">
        <button type="button" onClick={onClose} className={styles.closeButton} aria-label="Fechar modal"></button>

        <h1 className={styles.title}>Qual o tipo de página?</h1>

        <div className={styles.typeSiteContainer}>
          {pageTypes.map((page) => (
            <button key={page.id} type="button" className={styles.pageTypeCard} onClick={() => onSelectType?.(page.id)}>
              <img src={page.image} alt={`${page.title} icon`} />
              <h2>{page.title}</h2>
            </button>
          ))}
        </div>

        <p className={styles.RecommendTypePage}>
          Não encontrou seu tipo ? <a href="suporte">Nos recomende a sua!</a>
        </p>
      </div>
    </>
  );
}
