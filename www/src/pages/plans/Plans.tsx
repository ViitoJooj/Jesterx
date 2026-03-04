import { useState } from "react";
import styles from "./Plans.module.scss"

type PlanType = string;

interface Plan {
  id: PlanType;
  name: string;
  price_cents: number;
  features: string[];
  popular?: boolean;
  description?: string;
}

const plans: Plan[] = [
  {
    id: "business",
    name: "Business",
    price_cents: 4900,
    features: [
      "1 site",
      "Páginas ilimitadas",
      "Suporte por email",
      "Templates básicos",
      "SSL incluído",
    ],
  },
  {
    id: "pro",
    name: "Pro",
    price_cents: 9900,
    popular: true,
    features: [
      "Até 10 sites",
      "Páginas ilimitadas",
      "Suporte prioritário",
      "Todos os templates",
      "SSL incluído",
      "Analytics avançado",
      "Integrações premium",
    ],
  },
  {
    id: "enterprise",
    name: "Enterprise",
    price_cents: 19900,
    features: [
      "Até 50 sites",
      "Páginas ilimitadas",
      "Suporte 24/7",
      "Templates customizados",
      "SSL incluído",
      "Analytics avançado",
      "Todas as integrações",
      "API dedicada",
      "Suporte a White Label",
    ],
  },
];

export const Plans: React.FC = () => {
  const [selectedPlan, setSelectedPlan] = useState<PlanType | null>(null);

  function formatPrice(cents: number) {
    return (cents / 100).toLocaleString("pt-BR", {
      style: "currency",
      currency: "BRL",
    });
  }

  function handleSelectPlan(planId: PlanType) {
    setSelectedPlan(planId);
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.header}>
          <h1>Escolha seu plano</h1>
          <p>Comece a criar suas páginas e e-commerces hoje mesmo</p>
        </div>

        <div className={styles.plansGrid}>
          {plans.map((plan) => (
            <div
              key={plan.id}
              className={`${styles.planCard} ${
                plan.popular ? styles.popular : ""
              }`}
            >
              {plan.popular && (
                <span className={styles.badge}>Mais Popular</span>
              )}

              <h2 className={styles.planName}>{plan.name}</h2>

              <div className={styles.priceSection}>
                <span className={styles.price}>
                  {formatPrice(plan.price_cents)}
                </span>
                <span className={styles.period}>/mês</span>
              </div>

              <ul className={styles.features}>
                {plan.features.map((feature, index) => (
                  <li key={index}>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                    >
                      <path
                        d="M16.6667 5L7.50004 14.1667L3.33337 10"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                    {feature}
                  </li>
                ))}
              </ul>

              <button
                className={`${styles.selectButton} ${
                  plan.popular ? styles.primary : styles.secondary
                }`}
                onClick={() => handleSelectPlan(plan.id)}
              >
                {selectedPlan === plan.id
                  ? "Plano selecionado"
                  : "Começar agora"}
              </button>
            </div>
          ))}
        </div>

        <div className={styles.footer}>
          <p>
            Todos os planos incluem 14 dias de teste grátis. Sem cartão de
            crédito necessário.
          </p>
          <p>
            <a href="/contact">Precisa de um plano customizado?</a>
          </p>
        </div>
      </div>
    </main>
  );
}
