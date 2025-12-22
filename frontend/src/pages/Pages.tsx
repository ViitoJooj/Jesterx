import { useState } from "react";
import styles from "../styles/pages/Store.module.scss";
import { CreatePageForm } from "../components/CreatePageForm";

export function Pages() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <main className={styles.main}>
      <h1>Crie sua página🎉</h1>

      <div className={styles.projectsContainer}>
        <button className={styles.createNewStore} onClick={() => setIsOpen(true)}>
          <span className={styles.plus}>+</span>
          <p className={styles.title}>Crie um novo negocio!</p>
          <span className={styles.subtitle}>Comece uma nova experiência</span>
        </button>
      </div>

      {isOpen && (
        <div className={styles.modalOverlay} onClick={() => setIsOpen(false)}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <CreatePageForm onClose={() => setIsOpen(false)} />
          </div>
        </div>
      )}
    </main>
  );
}
