import React from "react";
import styles from "./home.module.scss";
import Button from "../../components/Button/button";
import RotatingWord from "../../components/RotatingWord/rotatingWord";

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

export function Home() {
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
    </>
  )
}