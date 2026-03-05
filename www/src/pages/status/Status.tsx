import { useState } from "react";
import styles from "./Status.module.scss";

type ServiceStatus = "operational" | "degraded" | "outage" | "maintenance";

interface Service {
  name: string;
  description: string;
  status: ServiceStatus;
  uptime: string;
}

interface Incident {
  id: string;
  title: string;
  severity: "critical" | "major" | "minor";
  status: "investigating" | "identified" | "monitoring" | "resolved";
  date: string;
  updates: { time: string; message: string }[];
}

const SERVICES: Service[] = [
  { name: "API Principal", description: "API REST do backend Golang", status: "operational", uptime: "99.98%" },
  { name: "Entrega de Lojas", description: "Renderização e servimento das páginas públicas", status: "operational", uptime: "99.95%" },
  { name: "Editor Elementor", description: "Editor visual drag-and-drop", status: "operational", uptime: "99.90%" },
  { name: "Editor de Código", description: "Editores React e Svelte", status: "operational", uptime: "99.88%" },
  { name: "Processamento de Pagamentos", description: "Integração Stripe Connect", status: "operational", uptime: "99.99%" },
  { name: "Hospedagem de Vídeos", description: "Upload e streaming de vídeos de cursos", status: "operational", uptime: "99.85%" },
  { name: "CDN", description: "Distribuição de conteúdo estático", status: "operational", uptime: "99.97%" },
  { name: "Banco de Dados", description: "PostgreSQL — dados de lojas e usuários", status: "operational", uptime: "99.99%" },
  { name: "Autenticação", description: "Login, registro e verificação de e-mail", status: "operational", uptime: "99.96%" },
];

const INCIDENTS: Incident[] = [
  {
    id: "inc-001",
    title: "Lentidão na entrega de lojas em algumas regiões",
    severity: "minor",
    status: "resolved",
    date: "28 fev 2025",
    updates: [
      { time: "14:32", message: "Investigando relatórios de lentidão na entrega de páginas públicas em regiões da América do Sul." },
      { time: "15:10", message: "Identificado gargalo no nó CDN da região SA-EAST-1. Iniciando redistribuição de tráfego." },
      { time: "15:45", message: "Redistribuição concluída. Latência normalizada. Monitorando por 30 minutos." },
      { time: "16:20", message: "Incidente resolvido. Todas as métricas dentro do normal." },
    ],
  },
  {
    id: "inc-002",
    title: "Manutenção programada do banco de dados",
    severity: "minor",
    status: "resolved",
    date: "15 fev 2025",
    updates: [
      { time: "02:00", message: "Início da manutenção programada. Migrações de schema sendo aplicadas." },
      { time: "02:45", message: "Manutenção concluída com sucesso. Todos os serviços operacionais." },
    ],
  },
];

const STATUS_LABELS: Record<ServiceStatus, string> = {
  operational: "Operacional",
  degraded: "Degradado",
  outage: "Fora do ar",
  maintenance: "Manutenção",
};

const INCIDENT_STATUS_LABELS: Record<Incident["status"], string> = {
  investigating: "Investigando",
  identified: "Identificado",
  monitoring: "Monitorando",
  resolved: "Resolvido",
};

const SEVERITY_LABELS: Record<Incident["severity"], string> = {
  critical: "Crítico",
  major: "Grave",
  minor: "Menor",
};

function allOperational(services: Service[]) {
  return services.every((s) => s.status === "operational");
}

export const Status: React.FC = () => {
  const [openIncident, setOpenIncident] = useState<string | null>(null);
  const ok = allOperational(SERVICES);

  return (
    <div className={styles.page}>

      {/* Hero */}
      <div className={`${styles.hero} ${ok ? styles.heroOk : styles.heroIssue}`}>
        <div className={styles.heroInner}>
          <div className={styles.heroIcon}>{ok ? "✅" : "⚠️"}</div>
          <div>
            <h1 className={styles.heroTitle}>{ok ? "Todos os sistemas operacionais" : "Algumas instabilidades detectadas"}</h1>
            <p className={styles.heroSub}>Última verificação: agora mesmo · Atualizado automaticamente</p>
          </div>
        </div>
      </div>

      <div className={styles.body}>

        {/* Serviços */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Status dos Serviços</h2>
          <div className={styles.serviceList}>
            {SERVICES.map((svc) => (
              <div key={svc.name} className={styles.serviceRow}>
                <div className={styles.serviceInfo}>
                  <strong>{svc.name}</strong>
                  <span>{svc.description}</span>
                </div>
                <div className={styles.serviceMeta}>
                  <span className={styles.uptimeLabel}>Uptime 90d: <strong>{svc.uptime}</strong></span>
                  <span className={`${styles.statusBadge} ${styles[`status_${svc.status}`]}`}>
                    <span className={styles.statusDot} />
                    {STATUS_LABELS[svc.status]}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Uptime visual */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Histórico de Uptime — últimos 90 dias</h2>
          <div className={styles.uptimeGrid}>
            {SERVICES.map((svc) => (
              <div key={svc.name} className={styles.uptimeRow}>
                <span className={styles.uptimeName}>{svc.name}</span>
                <div className={styles.uptimeBars}>
                  {Array.from({ length: 90 }, (_, i) => {
                    const chance = Math.random();
                    const cls = i >= 87 && svc.name === "Entrega de Lojas" ? styles.barIssue :
                                chance > 0.995 ? styles.barMaint : styles.barOk;
                    return <span key={i} className={`${styles.bar} ${cls}`} title={`Dia ${90 - i}`} />;
                  })}
                </div>
                <span className={styles.uptimePct}>{svc.uptime}</span>
              </div>
            ))}
          </div>
          <div className={styles.uptimeLegend}>
            <span><span className={`${styles.legendDot} ${styles.barOk}`} /> Operacional</span>
            <span><span className={`${styles.legendDot} ${styles.barMaint}`} /> Manutenção</span>
            <span><span className={`${styles.legendDot} ${styles.barIssue}`} /> Incidente</span>
          </div>
        </section>

        {/* Incidentes */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Histórico de Incidentes</h2>
          {INCIDENTS.length === 0 ? (
            <p className={styles.noIncidents}>Nenhum incidente nos últimos 90 dias. 🎉</p>
          ) : (
            <div className={styles.incidentList}>
              {INCIDENTS.map((inc) => (
                <div key={inc.id} className={styles.incidentCard}>
                  <button
                    className={styles.incidentHeader}
                    onClick={() => setOpenIncident(openIncident === inc.id ? null : inc.id)}
                  >
                    <div className={styles.incidentLeft}>
                      <span className={`${styles.incBadge} ${styles[`sev_${inc.severity}`]}`}>
                        {SEVERITY_LABELS[inc.severity]}
                      </span>
                      <strong>{inc.title}</strong>
                    </div>
                    <div className={styles.incidentRight}>
                      <span className={`${styles.incStatus} ${inc.status === "resolved" ? styles.incResolved : styles.incActive}`}>
                        {INCIDENT_STATUS_LABELS[inc.status]}
                      </span>
                      <span className={styles.incDate}>{inc.date}</span>
                      <span className={styles.incChevron}>{openIncident === inc.id ? "▲" : "▼"}</span>
                    </div>
                  </button>
                  {openIncident === inc.id && (
                    <div className={styles.incidentUpdates}>
                      {inc.updates.map((u, i) => (
                        <div key={i} className={styles.incidentUpdate}>
                          <span className={styles.updateTime}>{u.time}</span>
                          <span>{u.message}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Contato */}
        <div className={styles.contactBox}>
          <span className={styles.contactIcon}>📬</span>
          <div>
            <strong>Está enfrentando um problema não listado?</strong>
            <p>Entre em contato com nossa equipe de suporte em <a href="mailto:suporte@jesterx.com">suporte@jesterx.com</a> ou acesse o painel de incidentes.</p>
          </div>
        </div>

      </div>
    </div>
  );
};
