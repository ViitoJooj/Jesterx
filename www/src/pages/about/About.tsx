import styles from "./About.module.scss";

const TEAM = [
  {
    name: "Vítor Alves",
    role: "Fundador & CEO",
    bio: "Desenvolvedor full-stack apaixonado por produtos que empoderam empreendedores a construírem suas próprias lojas online sem depender de grandes plataformas.",
    avatar: "VA",
  },
];

const MILESTONES = [
  { year: "2024", title: "Ideia inicial", desc: "Surgiu a ideia de criar uma plataforma SaaS brasileira para criação de lojas e cursos online, com foco em flexibilidade total para o desenvolvedor." },
  { year: "2024", title: "Primeiros commits", desc: "Início do desenvolvimento com backend em Golang + PostgreSQL e frontend em React + TypeScript com Vite." },
  { year: "2025", title: "Editor visual", desc: "Lançamento do editor Elementor drag-and-drop e dos editores de código React e Svelte com syntax highlighting." },
  { year: "2025", title: "Integração de pagamentos", desc: "Integração com Stripe Connect para marketplace, permitindo que donos de lojas recebam pagamentos diretamente." },
];

const VALUES = [
  { icon: "🚀", title: "Velocidade", desc: "Infraestrutura otimizada para que suas lojas carreguem rápido em qualquer lugar do mundo." },
  { icon: "🔓", title: "Liberdade", desc: "Três modos de edição — Elementor, React e Svelte — você escolhe como construir." },
  { icon: "🛡️", title: "Segurança", desc: "TLS 1.2+, bcrypt, JWT com rotação — sua loja e seus clientes protegidos." },
  { icon: "🤝", title: "Transparência", desc: "Taxas claras, documentação aberta e sem letras miúdas escondidas." },
  { icon: "💡", title: "Inovação", desc: "Sempre adicionando novos recursos baseados no feedback real dos nossos usuários." },
  { icon: "🌱", title: "Crescimento", desc: "Da loja pequena ao negócio escalável, a plataforma cresce junto com você." },
];

export const About: React.FC = () => {
  return (
    <div className={styles.page}>

      {/* Hero */}
      <section className={styles.hero}>
        <div className={styles.heroInner}>
          <div className={styles.heroBadge}>Nossa história</div>
          <h1 className={styles.heroTitle}>
            Construindo a plataforma que<br />
            <span className={styles.heroAccent}>empreendedores merecem</span>
          </h1>
          <p className={styles.heroSub}>
            O Jesterx nasceu da frustração de não encontrar uma plataforma brasileira de e-commerce que fosse ao mesmo tempo poderosa, flexível e acessível. Criamos a que não existia.
          </p>
        </div>
        <div className={styles.heroStats}>
          <div className={styles.stat}>
            <span className={styles.statNum}>3</span>
            <span className={styles.statLabel}>modos de edição</span>
          </div>
          <div className={styles.statDivider} />
          <div className={styles.stat}>
            <span className={styles.statNum}>100%</span>
            <span className={styles.statLabel}>feito no Brasil</span>
          </div>
          <div className={styles.statDivider} />
          <div className={styles.stat}>
            <span className={styles.statNum}>Go+React</span>
            <span className={styles.statLabel}>stack moderna</span>
          </div>
        </div>
      </section>

      {/* Missão */}
      <section className={styles.mission}>
        <div className={styles.missionInner}>
          <div className={styles.missionText}>
            <span className={styles.sectionTag}>Nossa Missão</span>
            <h2 className={styles.sectionTitle}>Democratizar a criação de lojas online</h2>
            <p>
              Acreditamos que qualquer empreendedor — seja um developer experiente ou alguém que está dando os primeiros passos no digital — merece uma ferramenta que respeite sua criatividade e inteligência.
            </p>
            <p>
              No Jesterx você não fica preso em templates limitados. Você pode usar o editor visual drag-and-drop, escrever componentes React do zero, ou criar com Svelte. A escolha é sempre sua.
            </p>
            <p>
              Produtos físicos, digitais, cursos em vídeo — tudo em um lugar só, com pagamentos via Stripe e infraestrutura que escala junto com o seu negócio.
            </p>
          </div>
          <div className={styles.missionVisual}>
            <div className={styles.missionCard}>
              <div className={styles.missionCardIcon}>🏗️</div>
              <strong>Crie</strong>
              <p>Use Elementor, React ou Svelte para montar sua loja do jeito que você imaginou.</p>
            </div>
            <div className={styles.missionCard}>
              <div className={styles.missionCardIcon}>💳</div>
              <strong>Venda</strong>
              <p>Produtos físicos, digitais e cursos em vídeo com pagamentos integrados via Stripe.</p>
            </div>
            <div className={styles.missionCard}>
              <div className={styles.missionCardIcon}>📈</div>
              <strong>Cresça</strong>
              <p>Infraestrutura multi-tenant que suporta desde a primeira venda até o alto volume.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Valores */}
      <section className={styles.values}>
        <div className={styles.valuesInner}>
          <div className={styles.centerHeader}>
            <span className={styles.sectionTag}>Nossos Valores</span>
            <h2 className={styles.sectionTitle}>O que nos guia todo dia</h2>
          </div>
          <div className={styles.valuesGrid}>
            {VALUES.map((v) => (
              <div key={v.title} className={styles.valueCard}>
                <span className={styles.valueIcon}>{v.icon}</span>
                <strong>{v.title}</strong>
                <p>{v.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Timeline */}
      <section className={styles.timeline}>
        <div className={styles.timelineInner}>
          <div className={styles.centerHeader}>
            <span className={styles.sectionTag}>Linha do Tempo</span>
            <h2 className={styles.sectionTitle}>Como chegamos até aqui</h2>
          </div>
          <div className={styles.timelineList}>
            {MILESTONES.map((m, i) => (
              <div key={i} className={styles.timelineItem}>
                <div className={styles.timelineYear}>{m.year}</div>
                <div className={styles.timelineLine}>
                  <div className={styles.timelineDot} />
                  {i < MILESTONES.length - 1 && <div className={styles.timelineConnector} />}
                </div>
                <div className={styles.timelineContent}>
                  <strong>{m.title}</strong>
                  <p>{m.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Time */}
      <section className={styles.team}>
        <div className={styles.teamInner}>
          <div className={styles.centerHeader}>
            <span className={styles.sectionTag}>Time</span>
            <h2 className={styles.sectionTitle}>Quem está por trás do Jesterx</h2>
          </div>
          <div className={styles.teamGrid}>
            {TEAM.map((member) => (
              <div key={member.name} className={styles.teamCard}>
                <div className={styles.teamAvatar}>{member.avatar}</div>
                <strong className={styles.teamName}>{member.name}</strong>
                <span className={styles.teamRole}>{member.role}</span>
                <p className={styles.teamBio}>{member.bio}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className={styles.cta}>
        <div className={styles.ctaInner}>
          <h2>Pronto para criar sua loja?</h2>
          <p>Comece gratuitamente e veja como é fácil ter uma loja online profissional.</p>
          <div className={styles.ctaBtns}>
            <a href="/register" className={styles.ctaBtnPrimary}>Criar conta grátis</a>
            <a href="/docs" className={styles.ctaBtnSecondary}>Ver documentação</a>
          </div>
        </div>
      </section>

    </div>
  );
};
