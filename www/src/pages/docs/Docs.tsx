import { useEffect, useRef, useState } from "react";
import styles from "./Docs.module.scss";

const NAV = [
  { id: "overview",      label: "Visão Geral" },
  { id: "user-resp",     label: "Responsabilidades do Usuário" },
  { id: "platform-resp", label: "Responsabilidades da Plataforma" },
  { id: "prohibited",    label: "Conteúdo Proibido" },
  { id: "moderation",    label: "Moderação de Conteúdo" },
  { id: "payments",      label: "Pagamentos e Receita" },
  { id: "liability",     label: "Limitação de Responsabilidade" },
  { id: "suspension",    label: "Suspensão e Encerramento" },
  { id: "security",      label: "Segurança e Dados" },
  { id: "reporting",     label: "Denunciar Abuso" },
];

function Section({ id, title, children }: { id: string; title: string; children: React.ReactNode }) {
  return (
    <section id={id} className={styles.section}>
      <h2 className={styles.sectionTitle}>{title}</h2>
      {children}
    </section>
  );
}

function Sub({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className={styles.sub}>
      <h3 className={styles.subTitle}>{title}</h3>
      {children}
    </div>
  );
}

function Note({ children }: { children: React.ReactNode }) {
  return <div className={styles.note}>{children}</div>;
}

function Badge({ color, children }: { color: "orange" | "red" | "blue" | "green"; children: React.ReactNode }) {
  return <span className={`${styles.badge} ${styles[`badge_${color}`]}`}>{children}</span>;
}

export const Docs: React.FC = () => {
  const [activeId, setActiveId] = useState("overview");
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) setActiveId(entry.target.id);
        });
      },
      { rootMargin: "-20% 0px -70% 0px" }
    );
    NAV.forEach(({ id }) => {
      const el = document.getElementById(id);
      if (el) observer.observe(el);
    });
    return () => observer.disconnect();
  }, []);

  function scrollTo(id: string) {
    document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  return (
    <div className={styles.page}>
      {/* Cabeçalho */}
      <div className={styles.pageHeader}>
        <div className={styles.pageHeaderInner}>
          <div className={styles.breadcrumb}>
            <span>Jesterx</span>
            <span className={styles.breadcrumbSep}>/</span>
            <span>Documentação</span>
          </div>
          <h1 className={styles.pageTitle}>Política da Plataforma &amp; Termos de Uso</h1>
          <p className={styles.pageSub}>
            Termos legais, políticas de conteúdo, regras de pagamento e responsabilidades da plataforma para todos os donos de lojas e usuários finais que operam na infraestrutura Jesterx.
          </p>
          <div className={styles.metaRow}>
            <span className={styles.metaItem}>📅 Última atualização: Março de 2025</span>
            <span className={styles.metaItem}>📖 Versão 1.0</span>
            <span className={styles.metaItem}><Badge color="green">Vigente</Badge></span>
          </div>
        </div>
      </div>

      <div className={styles.body}>
        {/* Sidebar */}
        <aside className={styles.sidebar}>
          <p className={styles.sidebarLabel}>Nesta página</p>
          <nav>
            {NAV.map((item) => (
              <button
                key={item.id}
                className={`${styles.navItem} ${activeId === item.id ? styles.navItemActive : ""}`}
                onClick={() => scrollTo(item.id)}
              >
                {item.label}
              </button>
            ))}
          </nav>
          <div className={styles.sidebarFooter}>
            <p>Dúvidas?</p>
            <a href="mailto:suporte@jesterx.com" className={styles.sidebarLink}>suporte@jesterx.com</a>
          </div>
        </aside>

        <main className={styles.content} ref={contentRef}>

          {/* 1. Visão Geral */}
          <Section id="overview" title="1. Visão Geral da Plataforma">
            <p>
              <strong>Jesterx</strong> é uma plataforma SaaS multi-tenant que fornece infraestrutura técnica para empreendedores independentes e desenvolvedores criarem, hospedarem e operarem lojas virtuais. Donos de lojas podem vender produtos físicos, downloads digitais e cursos online diretamente para seus clientes por meio de vitrines criadas no Jesterx.
            </p>
            <Sub title="Como funciona">
              <ul>
                <li>Donos de lojas criam uma conta e escolhem um plano de assinatura que determina o número de rotas, páginas e recursos disponíveis.</li>
                <li>As vitrines são construídas usando um de três modos: <code>Elementor</code> (construtor visual drag-and-drop), <code>React</code> (componentes JSX customizados) ou <code>Svelte</code> (sintaxe de componentes reativos).</li>
                <li>Cada loja é servida a partir de uma URL única na infraestrutura Jesterx (<code>jesterx.com/p/&#123;siteId&#125;</code>) ou, opcionalmente, em um domínio customizado.</li>
                <li>Os pagamentos são coletados pelo <strong>Stripe</strong> em nome dos donos de lojas, com as taxas da plataforma deduzidas automaticamente no momento da transação.</li>
                <li>Vídeos de cursos e ativos digitais são hospedados nos servidores Jesterx e entregues via CDN aos clientes finais.</li>
              </ul>
            </Sub>
            <Sub title="Stack técnica">
              <div className={styles.stackGrid}>
                <div className={styles.stackItem}><span className={styles.stackLabel}>Backend</span><span>Golang + PostgreSQL</span></div>
                <div className={styles.stackItem}><span className={styles.stackLabel}>Frontend</span><span>React + TypeScript + Vite</span></div>
                <div className={styles.stackItem}><span className={styles.stackLabel}>Pagamentos</span><span>Stripe Connect (marketplace)</span></div>
                <div className={styles.stackItem}><span className={styles.stackLabel}>Hospedagem</span><span>Infraestrutura compartilhada multi-tenant</span></div>
              </div>
            </Sub>
            <Note>
              ℹ️ A Jesterx atua exclusivamente como <strong>provedora de infraestrutura tecnológica e de marketplace</strong>. Não somos varejistas, fornecedores de cursos nem processadores de pagamento para clientes finais. Todas as relações comerciais existem entre os donos de lojas e seus clientes.
            </Note>
          </Section>

          {/* 2. Responsabilidades do Usuário */}
          <Section id="user-resp" title="2. Responsabilidades do Usuário (Donos de Loja)">
            <p>
              Ao criar uma loja no Jesterx, você concorda que é o único responsável por todos os aspectos de suas operações comerciais. O Jesterx fornece apenas a infraestrutura técnica — não a supervisão das suas práticas de negócio.
            </p>

            <Sub title="Legalidade dos Produtos">
              <p>Você é inteiramente responsável por garantir que todos os produtos, serviços e conteúdos oferecidos em sua loja estejam em conformidade com as leis de sua jurisdição e da jurisdição de seus clientes. Isso inclui:</p>
              <ul>
                <li>Verificar se os produtos são legais para venda nos mercados-alvo.</li>
                <li>Obter quaisquer licenças, autorizações ou certificações necessárias antes de listar produtos regulamentados.</li>
                <li>Cumprir restrições de importação/exportação, leis de verificação de idade e regulamentos de proteção ao consumidor aplicáveis.</li>
              </ul>
            </Sub>

            <Sub title="Propriedade Intelectual &amp; Direitos Autorais">
              <ul>
                <li>Você só pode vender, distribuir ou publicar conteúdo sobre o qual detenha os direitos de propriedade intelectual adequados ou licenças válidas.</li>
                <li>Fazer upload, revender ou distribuir conteúdo protegido por direitos autorais (incluindo músicas, vídeos, software, e-books, imagens ou materiais de cursos) sem autorização do titular dos direitos é estritamente proibido.</li>
                <li>O Jesterx responderá a notificações válidas de remoção por violação de direitos autorais (DMCA e equivalentes internacionais). Veja <button className={styles.inlineLink} onClick={() => scrollTo("reporting")}>§10 Denunciar Abuso</button>.</li>
              </ul>
            </Sub>

            <Sub title="Impostos &amp; Obrigações Fiscais">
              <ul>
                <li>Você é responsável por determinar, coletar, declarar e recolher todos os impostos aplicáveis às transações realizadas pela sua loja, incluindo ICMS, ISS, IPI, imposto de renda e quaisquer outros tributos.</li>
                <li>O Jesterx não fornece assessoria fiscal. Recomendamos consultar um contador ou profissional tributário qualificado.</li>
                <li>O Jesterx pode emitir relatórios de informação conforme exigido pela legislação aplicável, mas não gerencia sua conformidade fiscal global.</li>
              </ul>
            </Sub>

            <Sub title="Suporte ao Cliente &amp; Reembolsos">
              <ul>
                <li>Você é o único responsável por atender a todas as consultas, disputas, reclamações e solicitações de suporte relacionadas à sua loja, produtos e serviços.</li>
                <li>Você deve manter uma política de reembolso e devolução visível ao público. O não cumprimento de sua política declarada pode resultar em revisão da conta.</li>
                <li>Em caso de chargeback iniciado por um cliente, você é o único responsável por responder à disputa. Veja <button className={styles.inlineLink} onClick={() => scrollTo("payments")}>§6 Pagamentos &amp; Receita</button>.</li>
              </ul>
            </Sub>

            <Sub title="Responsabilidade pelo Conteúdo de Cursos">
              <ul>
                <li>Para lojas do tipo curso, você é responsável pela precisão, qualidade e legalidade de todo o material instrucional publicado na plataforma.</li>
                <li>É proibido publicar conteúdo de cursos que represente credenciais falsas, contenha afirmações incorretas ou promova práticas prejudiciais.</li>
                <li>Todos os vídeos e materiais de cursos enviados para os servidores Jesterx devem estar em conformidade com as políticas de conteúdo descritas em <button className={styles.inlineLink} onClick={() => scrollTo("prohibited")}>§4 Conteúdo Proibido</button>.</li>
              </ul>
            </Sub>

            <Sub title="Conformidade com Leis Locais">
              <p>Você concorda em operar sua loja em total conformidade com todas as leis locais, nacionais e internacionais aplicáveis, incluindo, sem limitação:</p>
              <ul>
                <li>Leis de proteção ao consumidor (ex.: direito de arrependimento, divulgações obrigatórias — CDC)</li>
                <li>Leis de proteção de dados e privacidade (ex.: LGPD, GDPR, CCPA)</li>
                <li>Regulamentos de comércio eletrônico</li>
                <li>Requisitos de prevenção à lavagem de dinheiro (AML) e conheça seu cliente (KYC) quando aplicável</li>
              </ul>
            </Sub>
          </Section>

          {/* 3. Responsabilidades da Plataforma */}
          <Section id="platform-resp" title="3. Responsabilidades da Plataforma">
            <p>
              O Jesterx se compromete a manter uma infraestrutura confiável, segura e escalável para todos os donos de lojas que operam na plataforma.
            </p>

            <Sub title="Hospedagem &amp; Infraestrutura">
              <ul>
                <li>O Jesterx fornece infraestrutura de hospedagem em nuvem compartilhada para todos os sites de lojas, conteúdos de cursos e ativos de produtos.</li>
                <li>Os sites das lojas são servidos via CDN globalmente distribuída para minimizar a latência para os clientes finais.</li>
                <li>Backups do banco de dados são realizados regularmente para prevenir perda de dados por falhas de infraestrutura.</li>
              </ul>
            </Sub>

            <Sub title="Disponibilidade &amp; Uptime">
              <ul>
                <li>O Jesterx tem como meta um uptime mensal de <strong>99,5%</strong> para a API principal e a infraestrutura de entrega das lojas.</li>
                <li>Janelas de manutenção programada são anunciadas com pelo menos 24 horas de antecedência sempre que possível.</li>
                <li>Manutenções emergenciais que afetem a disponibilidade serão comunicadas através da página de status da plataforma.</li>
              </ul>
            </Sub>

            <Sub title="Integração de Pagamentos">
              <ul>
                <li>O Jesterx integra o <strong>Stripe Connect</strong> para facilitar a coleta de pagamentos, gerenciamento de repasses e dedução de taxas da plataforma.</li>
                <li>O Jesterx mantém os padrões de conformidade PCI-DSS por meio da infraestrutura do Stripe.</li>
                <li>O Jesterx não armazena dados brutos de cartão em nenhum ponto do fluxo de transação.</li>
              </ul>
            </Sub>

            <Sub title="Segurança da Plataforma">
              <ul>
                <li>Todos os dados em trânsito são criptografados usando <strong>TLS 1.2+</strong>.</li>
                <li>Tokens de autenticação são gerenciados com práticas padrão da indústria, incluindo rotação de tokens e armazenamento seguro.</li>
                <li>O Jesterx realiza revisões periódicas de segurança de sua infraestrutura e código de aplicação.</li>
              </ul>
            </Sub>

            <Sub title="Ferramentas de Moderação">
              <ul>
                <li>O Jesterx fornece um mecanismo de denúncia de abuso para que usuários finais possam sinalizar lojas ou produtos que violem as políticas da plataforma.</li>
                <li>O Jesterx se reserva o direito de revisar, suspender ou remover qualquer loja ou conteúdo que viole estes termos.</li>
                <li>Ferramentas de varredura automática podem ser usadas para detectar padrões de conteúdo proibido, embora a revisão manual faça parte de todas as ações de aplicação.</li>
              </ul>
            </Sub>
          </Section>

          {/* 4. Conteúdo Proibido */}
          <Section id="prohibited" title="4. Conteúdo e Produtos Proibidos">
            <p>
              As seguintes categorias de conteúdo, produtos e serviços são <strong>estritamente proibidas</strong> na plataforma Jesterx. Lojas que violarem essas proibições estão sujeitas a suspensão imediata ou encerramento permanente.
            </p>

            <div className={styles.prohibitedGrid}>
              {[
                {
                  icon: "⚖️",
                  title: "Produtos e Serviços Ilegais",
                  items: [
                    "Substâncias controladas, entorpecentes ou precursores químicos",
                    "Mercadorias contrafeitas ou documentação falsa",
                    "Armas, armas de fogo ou munições (salvo autorização específica e conformidade legal)",
                    "Produtos ilegais de vender em qualquer jurisdição aplicável",
                  ],
                },
                {
                  icon: "🚫",
                  title: "Fraude e Golpes",
                  items: [
                    "Avaliações falsas, manipulação de contas ou bots de redes sociais",
                    "Esquemas de investimento, pirâmides ou produtos estruturados como pirâmide",
                    "Alegações de saúde enganosas ou tratamentos médicos não regulamentados",
                    "Personificação de outras marcas, empresas ou indivíduos",
                  ],
                },
                {
                  icon: "©️",
                  title: "Violação de Direitos Autorais",
                  items: [
                    "Revenda ou distribuição de software, e-books ou mídia pirateados",
                    "Uso não autorizado de logotipos, nomes de marcas ou imagens com direitos autorais",
                    "Conteúdo de cursos construído principalmente com materiais plagiados de terceiros",
                    "Ferramentas de contorno de proteção DRM",
                  ],
                },
                {
                  icon: "🦠",
                  title: "Malware e Phishing",
                  items: [
                    "Software projetado para prejudicar, infiltrar ou vigiar usuários sem consentimento",
                    "Kits de phishing, ferramentas de coleta de credenciais ou modelos de engenharia social",
                    "Produtos cujo uso principal é o acesso não autorizado a sistemas",
                    "Ferramentas de spam ou sistemas de comunicação em massa não solicitada",
                  ],
                },
                {
                  icon: "🔞",
                  title: "Conteúdo Ilegal ou Não Consensual",
                  items: [
                    "Qualquer conteúdo que retrate menores em contexto sexual — política de tolerância zero",
                    "Imagens íntimas não consensuais ou pornografia de vingança",
                    "Conteúdo que facilite ou promova tráfico humano ou exploração",
                  ],
                },
                {
                  icon: "⚠️",
                  title: "Conteúdo Prejudicial ou Abusivo",
                  items: [
                    "Conteúdo que incite violência, terrorismo ou danos em massa",
                    "Discurso de ódio contra indivíduos ou grupos com base em características protegidas",
                    "Conteúdo projetado para assediar, ameaçar ou intimidar indivíduos",
                    "Produtos que promovam automutilação ou forneçam instruções para violência",
                  ],
                },
              ].map((cat) => (
                <div key={cat.title} className={styles.prohibitedCard}>
                  <div className={styles.prohibitedCardHeader}>
                    <span className={styles.prohibitedIcon}>{cat.icon}</span>
                    <strong>{cat.title}</strong>
                  </div>
                  <ul>
                    {cat.items.map((item) => <li key={item}>{item}</li>)}
                  </ul>
                </div>
              ))}
            </div>

            <Note>
              🚨 <strong>Tolerância Zero:</strong> Qualquer loja encontrada distribuindo material de abuso sexual infantil (CSAM) será imediata e permanentemente encerrada. Todos os dados serão preservados e relatados às autoridades competentes, incluindo a Polícia Federal e equivalentes internacionais.
            </Note>
          </Section>

          {/* 5. Moderação */}
          <Section id="moderation" title="5. Política de Moderação de Conteúdo">
            <p>
              O Jesterx opera um sistema de moderação reativo e proativo. Embora não faça triagem prévia de cada produto antes da publicação, mantém um processo ativo de aplicação de políticas.
            </p>

            <Sub title="Fluxo de Denúncia">
              <div className={styles.flowSteps}>
                {[
                  { step: "01", title: "Denúncia Enviada", desc: "Um usuário, cliente ou sistema automatizado envia uma denúncia de abuso por meio da ferramenta de denúncia da plataforma ou por e-mail." },
                  { step: "02", title: "Triagem Inicial", desc: "A equipe de confiança e segurança do Jesterx analisa a denúncia em até 72 horas úteis e categoriza por gravidade (Crítica / Alta / Média / Baixa)." },
                  { step: "03", title: "Investigação", desc: "Para denúncias não críticas, o dono da loja pode ser notificado e ter a oportunidade de responder ou corrigir. Denúncias críticas seguem diretamente para ação de aplicação." },
                  { step: "04", title: "Ação de Aplicação", desc: "Com base nos resultados, o Jesterx pode emitir um aviso, remover conteúdo específico, suspender temporariamente a loja ou encerrar permanentemente a conta." },
                  { step: "05", title: "Resolução", desc: "A parte que fez a denúncia é notificada do resultado quando legalmente permitido. Recursos podem ser enviados para o e-mail de confiança e segurança em até 14 dias." },
                ].map((s) => (
                  <div key={s.step} className={styles.flowStep}>
                    <span className={styles.flowNum}>{s.step}</span>
                    <div>
                      <strong>{s.title}</strong>
                      <p>{s.desc}</p>
                    </div>
                  </div>
                ))}
              </div>
            </Sub>

            <Sub title="Níveis de Gravidade">
              <table className={styles.table}>
                <thead>
                  <tr>
                    <th>Gravidade</th>
                    <th>Exemplos</th>
                    <th>Ação Padrão</th>
                    <th>Tempo de Resposta</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><Badge color="red">Crítica</Badge></td>
                    <td>CSAM, conteúdo terrorista, distribuição de malware</td>
                    <td>Banimento permanente imediato + encaminhamento às autoridades</td>
                    <td>Em até 1 hora</td>
                  </tr>
                  <tr>
                    <td><Badge color="orange">Alta</Badge></td>
                    <td>Fraude, phishing, violação de direitos autorais</td>
                    <td>Suspensão temporária imediata + investigação</td>
                    <td>Em até 24 horas</td>
                  </tr>
                  <tr>
                    <td><Badge color="blue">Média</Badge></td>
                    <td>Violações de política, anúncios enganosos</td>
                    <td>Aviso emitido, conteúdo removido, prazo de 7 dias para correção</td>
                    <td>Em até 72 horas</td>
                  </tr>
                  <tr>
                    <td><Badge color="green">Baixa</Badge></td>
                    <td>Desvios menores de política, divulgações incompletas</td>
                    <td>Notificação enviada, ação requerida do dono da loja</td>
                    <td>Em até 7 dias</td>
                  </tr>
                </tbody>
              </table>
            </Sub>

            <Sub title="Processo de Recurso">
              <p>Donos de lojas que acreditam que uma ação de aplicação foi tomada por engano podem enviar um recurso formal em até <strong>14 dias corridos</strong> após receber o aviso. O recurso deve incluir:</p>
              <ul>
                <li>O ID da loja e o identificador do conteúdo afetado</li>
                <li>Uma explicação clara de por que a ação de aplicação foi incorreta</li>
                <li>Qualquer documentação de suporte (licenças, certificações, autorizações)</li>
              </ul>
              <p>Os recursos são analisados por um membro diferente da equipe de confiança e segurança. Os resultados do recurso são definitivos.</p>
            </Sub>
          </Section>

          {/* 6. Pagamentos */}
          <Section id="payments" title="6. Política de Pagamentos e Receita">
            <Sub title="Integração com Stripe Connect">
              <p>
                Todos os pagamentos na plataforma Jesterx são processados pelo <strong>Stripe Connect</strong>, uma solução de pagamentos para marketplace fornecida pela Stripe, Inc. Para receber repasses, donos de lojas devem:
              </p>
              <ul>
                <li>Criar e conectar uma conta Stripe à sua loja Jesterx.</li>
                <li>Concluir os requisitos de verificação de identidade (KYC) do Stripe conforme exigido pelo Stripe e pelas regulamentações financeiras aplicáveis.</li>
                <li>Manter conformidade com o <a href="https://stripe.com/ssa" target="_blank" rel="noreferrer" className={styles.link}>Contrato de Serviços do Stripe</a> e a <a href="https://stripe.com/restricted-businesses" target="_blank" rel="noreferrer" className={styles.link}>política de Negócios Restritos</a>.</li>
              </ul>
            </Sub>

            <Sub title="Taxas de Transação da Plataforma">
              <p>O Jesterx deduz uma taxa de plataforma de cada transação processada pela loja. As taxas variam de acordo com o plano de assinatura:</p>
              <table className={styles.table}>
                <thead>
                  <tr>
                    <th>Plano</th>
                    <th>Taxa da Plataforma</th>
                    <th>Observações</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Starter</td>
                    <td>5%</td>
                    <td>Por transação, deduzida antes do repasse</td>
                  </tr>
                  <tr>
                    <td>Pro / Business</td>
                    <td>2,5%</td>
                    <td>Taxa reduzida para assinantes de planos superiores</td>
                  </tr>
                  <tr>
                    <td>Enterprise / Ultra</td>
                    <td>1%</td>
                    <td>Taxa customizada disponível mediante solicitação</td>
                  </tr>
                </tbody>
              </table>
              <p>
                As taxas da plataforma são separadas e adicionais às taxas de processamento do próprio Stripe. O total deduzido de cada transação é a soma de ambas.
              </p>
            </Sub>

            <Sub title="Repasses">
              <ul>
                <li>Os repasses são gerenciados diretamente pelo Stripe de acordo com o cronograma de pagamento configurado na conta Stripe do dono da loja.</li>
                <li>O Jesterx não retém fundos nem controla o tempo de repasse além da dedução inicial da taxa de plataforma no momento da liquidação.</li>
                <li>Os donos de lojas são responsáveis por garantir que sua conta Stripe permaneça em situação regular. O Jesterx não se responsabiliza por atrasos causados pelas próprias análises de conformidade ou retenções do Stripe.</li>
              </ul>
            </Sub>

            <Sub title="Chargebacks &amp; Disputas">
              <ul>
                <li>Chargebacks iniciados por clientes finais são de <strong>responsabilidade exclusiva do dono da loja</strong>.</li>
                <li>Os donos de lojas devem responder às disputas de chargeback por meio do painel do Stripe dentro do prazo especificado pelas bandeiras de cartão (tipicamente 7 a 21 dias).</li>
                <li>O Jesterx fornecerá dados de nível de plataforma (metadados de transação, histórico de versões) para auxiliar na resolução de disputas quando formalmente solicitado.</li>
                <li>Taxas excessivas de chargeback podem resultar em revisão adicional, taxas de plataforma elevadas ou suspensão da conta.</li>
              </ul>
            </Sub>

            <Note>
              💳 O Jesterx não é processador de pagamentos, transmissor de dinheiro nem instituição financeira. Todos os fundos fluem diretamente entre o cliente e a conta Stripe do dono da loja. O Jesterx apenas deduz sua taxa de plataforma no momento da liquidação.
            </Note>
          </Section>

          {/* 7. Responsabilidade */}
          <Section id="liability" title="7. Limitação de Responsabilidade">
            <p>
              <strong>NA EXTENSÃO MÁXIMA PERMITIDA PELA LEGISLAÇÃO APLICÁVEL</strong>, o Jesterx e seus afiliados, diretores, funcionários, agentes e licenciadores não serão responsáveis por:
            </p>
            <ul>
              <li>A legalidade, qualidade, precisão, segurança ou adequação a um propósito de qualquer produto, serviço ou conteúdo vendido, distribuído ou publicado por qualquer dono de loja na plataforma.</li>
              <li>Quaisquer danos diretos, indiretos, incidentais, especiais, consequenciais ou exemplares decorrentes das atividades comerciais de um dono de loja, incluindo, sem limitação, lucros cessantes, perda de dados ou danos à reputação.</li>
              <li>Disputas entre donos de lojas e seus clientes, incluindo disputas sobre qualidade de produtos, entrega, reembolsos ou acesso digital.</li>
              <li>Perdas causadas pela falha do dono da loja em cumprir com leis aplicáveis, obrigações fiscais ou políticas da plataforma.</li>
              <li>Conduta de terceiros, incluindo fraude, chargebacks ou acesso não autorizado causado por ações fora do controle direto do Jesterx.</li>
              <li>Tempo de inatividade, perda de dados ou interrupções de serviço decorrentes de eventos de força maior, falhas de provedores upstream ou circunstâncias além do controle razoável do Jesterx.</li>
            </ul>

            <Sub title="Declaração de Provedor de Tecnologia">
              <p>
                O Jesterx opera exclusivamente como <strong>provedor de infraestrutura tecnológica e de marketplace</strong>. O Jesterx não endossa, garante ou representa qualquer loja, produto, serviço ou conteúdo disponibilizado na plataforma por donos de lojas terceiros. A presença de uma loja no Jesterx não constitui aprovação ou recomendação das ofertas dessa loja.
              </p>
            </Sub>

            <Sub title="Responsabilidade Agregada Máxima">
              <p>
                Na medida em que qualquer responsabilidade não possa ser excluída por lei, a responsabilidade agregada total do Jesterx para com qualquer dono de loja por todas as reclamações decorrentes do uso da plataforma não excederá o valor total das taxas de plataforma pagas por esse dono de loja nos <strong>doze (12) meses anteriores</strong> à reclamação.
              </p>
            </Sub>
          </Section>

          {/* 8. Suspensão */}
          <Section id="suspension" title="8. Suspensão e Encerramento de Conta">
            <Sub title="Motivos para Suspensão ou Encerramento">
              <p>O Jesterx se reserva o direito de suspender ou encerrar permanentemente qualquer conta de loja pelos seguintes motivos, sem limitação:</p>
              <ul>
                <li>Violação de qualquer seção destas Políticas &amp; Termos da Plataforma</li>
                <li>Distribuição de conteúdo proibido conforme definido em <button className={styles.inlineLink} onClick={() => scrollTo("prohibited")}>§4</button></li>
                <li>Inadimplência de taxas de assinatura ou saldos pendentes</li>
                <li>Taxas excessivas de chargeback que excedam os limites aceitáveis definidos pelo Stripe ou pelo Jesterx</li>
                <li>Declarações falsas durante o registro de conta ou verificação de identidade</li>
                <li>Recebimento de múltiplas denúncias de abuso fundamentadas em um período de 90 dias consecutivos</li>
                <li>Tentativa de contornar a segurança da plataforma, limites de taxa ou controles de acesso</li>
                <li>Qualquer ação que represente risco legal, reputacional ou financeiro para o Jesterx ou seus usuários</li>
              </ul>
            </Sub>

            <Sub title="Suspensão vs. Encerramento">
              <table className={styles.table}>
                <thead>
                  <tr><th>Ação</th><th>Efeito</th><th>Reversível?</th></tr>
                </thead>
                <tbody>
                  <tr>
                    <td><strong>Suspensão Temporária</strong></td>
                    <td>Loja fica offline; dono ainda pode acessar o painel. Nenhuma nova transação é processada.</td>
                    <td>Sim, após resolução do problema subjacente</td>
                  </tr>
                  <tr>
                    <td><strong>Encerramento Permanente</strong></td>
                    <td>Loja e todos os dados associados são agendados para exclusão. O acesso à conta é revogado.</td>
                    <td>Não, exceto em circunstâncias extraordinárias via recurso formal</td>
                  </tr>
                </tbody>
              </table>
            </Sub>

            <Sub title="Retenção de Dados Após Encerramento">
              <ul>
                <li>Após o encerramento permanente, os dados da loja (produtos, páginas, pedidos) são retidos por <strong>30 dias</strong> antes da exclusão, salvo se a lei exigir retenção por período maior.</li>
                <li>Registros de transações podem ser retidos por até <strong>7 anos</strong> para fins de conformidade financeira.</li>
                <li>Evidências relacionadas a violações de políticas ou atividades ilegais podem ser retidas indefinidamente e/ou compartilhadas com autoridades policiais.</li>
              </ul>
            </Sub>
          </Section>

          {/* 9. Segurança */}
          <Section id="security" title="9. Segurança e Proteção de Dados">
            <Sub title="Criptografia de Dados">
              <ul>
                <li>Todos os dados em trânsito entre clientes e servidores Jesterx são criptografados usando <strong>TLS 1.2 ou superior</strong>.</li>
                <li>Campos sensíveis (senhas, tokens) são armazenados usando algoritmos de hash unidirecional padrão da indústria (bcrypt com fator de custo ≥ 12).</li>
                <li>As conexões de banco de dados dentro da infraestrutura usam canais criptografados.</li>
              </ul>
            </Sub>

            <Sub title="Segurança de Autenticação">
              <ul>
                <li>A autenticação de usuários é realizada via JWT (JSON Web Tokens) seguros e de curta duração, com rotação de tokens de atualização no lado do servidor.</li>
                <li>A verificação de e-mail é obrigatória para todas as contas recém-registradas antes que o acesso completo seja concedido.</li>
                <li>Atividades suspeitas de login acionam alertas automáticos e retenções temporárias de acesso.</li>
              </ul>
            </Sub>

            <Sub title="Segurança de Infraestrutura">
              <ul>
                <li>A API de backend é construída em <strong>Golang</strong>, que elimina muitas categorias de vulnerabilidades de segurança de memória comuns em outros runtimes de servidor.</li>
                <li>O acesso ao banco de dados é restrito apenas a serviços de backend autorizados; nenhuma exposição direta do banco de dados ao público.</li>
                <li>Atualizações de dependências e patches de segurança são aplicados em um ciclo regular de revisão.</li>
              </ul>
            </Sub>

            <Sub title="Proteção de Dados &amp; Privacidade">
              <ul>
                <li>O Jesterx coleta e processa dados pessoais conforme descrito na <strong>Política de Privacidade</strong>.</li>
                <li>Os donos de lojas são controladores independentes de dados para os dados pessoais de seus próprios clientes e devem manter suas próprias políticas de privacidade em conformidade com a LGPD, GDPR e outros frameworks aplicáveis.</li>
                <li>O Jesterx atua como processador de dados para dados de clientes de lojas armazenados na infraestrutura da plataforma.</li>
                <li>O Jesterx notificará os donos de lojas afetados em até <strong>72 horas</strong> após tomar conhecimento de uma violação de dados que represente risco para seus clientes, conforme exigido pela legislação aplicável.</li>
              </ul>
            </Sub>

            <Sub title="Responsabilidades de Segurança do Dono da Loja">
              <ul>
                <li>Você é responsável por manter a confidencialidade de suas credenciais de conta.</li>
                <li>Você deve notificar imediatamente o Jesterx em <a href="mailto:seguranca@jesterx.com" className={styles.link}>seguranca@jesterx.com</a> se suspeitar de acesso não autorizado à sua conta.</li>
                <li>Você é responsável por garantir que qualquer código de terceiros ou integrações adicionadas à sua vitrine não introduzam vulnerabilidades de segurança.</li>
              </ul>
            </Sub>
          </Section>

          {/* 10. Denúncias */}
          <Section id="reporting" title="10. Denunciar Abuso">
            <p>
              O Jesterx leva as denúncias de abuso a sério e está comprometido em manter um marketplace seguro, legal e confiável. Qualquer pessoa — seja um usuário Jesterx ou membro do público — pode enviar uma denúncia de abuso.
            </p>

            <Sub title="Como Denunciar">
              <div className={styles.reportGrid}>
                <div className={styles.reportCard}>
                  <span className={styles.reportIcon}>📧</span>
                  <strong>E-mail</strong>
                  <p>Envie um relatório detalhado para <a href="mailto:abuso@jesterx.com" className={styles.link}>abuso@jesterx.com</a> incluindo a URL da loja, descrição da violação e qualquer evidência de suporte.</p>
                </div>
                <div className={styles.reportCard}>
                  <span className={styles.reportIcon}>🔘</span>
                  <strong>Botão de Denúncia na Plataforma</strong>
                  <p>Toda página pública de loja inclui um link "Denunciar esta loja" no rodapé. Clicar nele abre um formulário guiado de denúncia.</p>
                </div>
                <div className={styles.reportCard}>
                  <span className={styles.reportIcon}>©️</span>
                  <strong>Notificação de Direitos Autorais</strong>
                  <p>Para violações de direitos autorais ou propriedade intelectual, envie uma notificação formal de remoção para <a href="mailto:dmca@jesterx.com" className={styles.link}>dmca@jesterx.com</a> com as informações exigidas pela legislação aplicável.</p>
                </div>
              </div>
            </Sub>

            <Sub title="O Que Incluir em uma Denúncia">
              <ul>
                <li>A URL da loja ou produto/conteúdo específico em violação</li>
                <li>Uma descrição clara da suposta violação e por que ela infringe a política da plataforma ou a lei</li>
                <li>Capturas de tela ou outras evidências quando aplicável</li>
                <li>Suas informações de contato (mantidas confidenciais na medida do possível)</li>
                <li>Para reivindicações de PI: prova de titularidade dos direitos (certidões de registro, obra original, etc.)</li>
              </ul>
            </Sub>

            <Sub title="Denúncias Falsas">
              <p>
                Enviar uma denúncia de abuso falsa ou maliciosa com a intenção de prejudicar uma loja ou interromper atividades comerciais lícitas é em si uma violação destes termos e pode resultar em ação contra a conta denunciante.
              </p>
            </Sub>

            <Note>
              🕐 O Jesterx tem como meta acusar o recebimento de todas as denúncias não críticas em até <strong>72 horas úteis</strong>. Denúncias críticas (envolvendo conteúdo ilegal ou dano iminente) são triadas imediatamente.
            </Note>
          </Section>

          {/* Aviso de Rodapé */}
          <div className={styles.footerNotice}>
            <div className={styles.footerNoticeInner}>
              <span className={styles.footerNoticeIcon}>🏪</span>
              <div>
                <strong>Aviso de Empresas Independentes</strong>
                <p>
                  Todas as lojas hospedadas na plataforma Jesterx são operadas por <strong>empresas e indivíduos terceiros independentes</strong>. O Jesterx fornece apenas a infraestrutura técnica — incluindo hospedagem, o construtor de lojas e a integração de processamento de pagamentos — e não possui, opera, controla ou endossa qualquer loja, produto, serviço ou conteúdo disponível em qualquer vitrine desta plataforma.
                </p>
                <p>
                  O Jesterx não é parte de nenhuma transação entre um dono de loja e seu cliente. Qualquer disputa, reclamação ou responsabilidade decorrente de uma compra feita em uma loja hospedada no Jesterx é de responsabilidade exclusiva do dono daquela loja.
                </p>
                <p className={styles.footerNoticeSmall}>
                  Estes termos foram atualizados pela última vez em <strong>março de 2025</strong>. O Jesterx reserva o direito de atualizar esta política a qualquer momento. Alterações materiais serão comunicadas aos donos de lojas registrados por e-mail com pelo menos 14 dias de antecedência. O uso continuado da plataforma após a data de vigência constitui aceitação dos termos atualizados.
                </p>
              </div>
            </div>
          </div>

        </main>
      </div>
    </div>
  );
};
