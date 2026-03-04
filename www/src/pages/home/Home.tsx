import React from "react";
import styles from "./Home.module.scss";
import Button from "../../components/button/Button";
import RotatingWord from "../../components/rotatingWord/RotatingWord";

const howItWorks = [
  {
    icon: "01",
    title: "Escolha o tipo de projeto",
    description:
      "Comece por e-commerce, landing page, loja de software ou site de vídeos com estrutura inicial pronta.",
    examples: ["Loja de roupas", "LP de campanha", "SaaS de assinatura"],
  },
  {
    icon: "02",
    title: "Monte com blocos visuais",
    description:
      "Edite páginas com componentes prontos, organização clara e controle de conteúdo sem fluxo técnico complexo.",
    examples: ["Hero + CTA", "Seção de preços", "Catálogo com filtro"],
  },
  {
    icon: "03",
    title: "Publique e evolua",
    description:
      "Conecte catálogo, pagamentos e conteúdo para operar rápido e escalar com consistência.",
    examples: ["Checkout", "Área de membros", "Campanhas sazonais"],
  },
];

const projects = [
  {
    icon: "🛒",
    title: "E-commerce completo",
    description:
      "Catálogo, páginas de produto e checkout em uma estrutura pronta para vender e crescer.",
    examples: ["Moda", "Eletrônicos", "Cosméticos"],
  },
  {
    icon: "🎯",
    title: "Landing pages",
    description:
      "Páginas orientadas a conversão para captação de leads e validação de ofertas.",
    examples: ["Tráfego pago", "Webinar", "Pré-lançamento"],
  },
  {
    icon: "💻",
    title: "Loja de softwares",
    description:
      "Venda de produtos digitais, planos e assinaturas com jornada de compra limpa.",
    examples: ["Plano mensal", "Trial", "Upgrade de plano"],
  },
  {
    icon: "🎬",
    title: "Site de vídeos",
    description:
      "Biblioteca de conteúdo organizada para navegação simples e descoberta rápida.",
    examples: ["Cursos", "Comunidade", "Portal de conteúdo"],
  },
];

export const Home: React.FC = () => {
  return (
    <>
      <main className={styles.main}>
        <div className={styles.header}>
          <h1>
            Construa seu projeto
            <br />
            <RotatingWord items={["mais rápido", "com clareza", "sem código", "do seu jeito"]} />
          </h1>

          <h2>
            Jester é a plataforma low-code para criar desde e-commerces completos até
            landing pages e experiências digitais em um só lugar. Conecte ERPs, gerencie
            produtos físicos e digitais e lance sua operação sem escrever código.
          </h2>

          <div className={styles.cta}>
            <Button to="/plans" variant="primary">
              Começar agora
            </Button>

            <Button to="/register" variant="secondary">
              Criar conta
            </Button>
          </div>
        </div>
      </main>

      <section className={styles.section}>
        <div className={styles.inner}>
          <div className={styles.sectionHeader}>
            <h3>Como funciona</h3>
            <p>Fluxo visual e objetivo para sair da ideia para um projeto publicado sem fricção.</p>
          </div>
          <div className={styles.grid3}>
            {howItWorks.map((item) => (
              <article key={item.title} className={styles.card}>
                <div className={styles.cardTop}>
                  <span className={styles.iconBadge}>{item.icon}</span>
                  <h4>{item.title}</h4>
                </div>
                <p>{item.description}</p>
                <div className={styles.examples}>
                  {item.examples.map((tag) => (
                    <span key={tag}>{tag}</span>
                  ))}
                </div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className={`${styles.section} ${styles.alt}`}>
        <div className={styles.inner}>
          <div className={styles.sectionHeader}>
            <h3>Projetos que você pode criar</h3>
            <p>Do básico ao avançado com a mesma experiência elegante, simples e escalável.</p>
          </div>
          <div className={styles.grid2}>
            {projects.map((item) => (
              <article key={item.title} className={styles.card}>
                <div className={styles.cardTop}>
                  <span className={`${styles.iconBadge} ${styles.iconProject}`}>{item.icon}</span>
                  <h4>{item.title}</h4>
                </div>
                <p>{item.description}</p>
                <div className={styles.examples}>
                  {item.examples.map((tag) => (
                    <span key={tag}>{tag}</span>
                  ))}
                </div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.inner}>
          <div className={styles.sectionHeader}>
            <h3>Flexível para seu stack</h3>
            <p>
              A proposta é simplificar para o usuário, com liberdade para evoluir usando
              Elementor, Svelte ou React quando necessário.
            </p>
          </div>
          <div className={styles.stackRow}>
            <article className={styles.stackCard}>
              <div className={styles.cardTop}>
                <span className={styles.iconBadge}>🧩</span>
                <h4>Modo visual estilo Elementor</h4>
              </div>
              <p>Ideal para times que precisam montar, revisar e publicar sem depender do dev em cada detalhe.</p>
              <div className={styles.examples}>
                <span>Editor drag-and-drop</span>
                <span>Blocos reaproveitáveis</span>
                <span>Publicação rápida</span>
              </div>
            </article>
            <article className={styles.stackCard}>
              <div className={styles.cardTop}>
                <span className={styles.iconBadge}>⚙️</span>
                <h4>Extensível com Svelte ou React</h4>
              </div>
              <p>Quando o projeto pedir customizações avançadas, o front continua evoluindo sem reescrever tudo.</p>
              <div className={styles.examples}>
                <span>Componentes custom</span>
                <span>Integrações de API</span>
                <span>Escala de produto</span>
              </div>
            </article>
          </div>
        </div>
      </section>
    </>
  );
};
