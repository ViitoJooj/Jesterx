import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { API_URL } from "../../hooks/api";
import styles from "./Report.module.scss";

const REASONS = [
  { value: "FRAUD", label: "Fraude" },
  { value: "SCAM", label: "Golpe / Estelionato" },
  { value: "SPAM", label: "Spam" },
  { value: "INAPPROPRIATE", label: "Conteúdo Inapropriado" },
  { value: "COUNTERFEIT", label: "Produto Falsificado" },
  { value: "OTHER", label: "Outro" },
];

export function Report() {
  const [params] = useSearchParams();
  const navigate = useNavigate();

  const [websiteId, setWebsiteId] = useState(params.get("website_id") ?? "");
  const [websiteName] = useState(params.get("website_name") ?? "");
  const [reporterName, setReporterName] = useState("");
  const [reporterEmail, setReporterEmail] = useState("");
  const [reason, setReason] = useState("FRAUD");
  const [description, setDescription] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [ticket, setTicket] = useState<{ id: string; ticketNumber: number } | null>(null);

  useEffect(() => {
    const id = params.get("website_id");
    if (id) setWebsiteId(id);
  }, [params]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/api/v1/reports`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          website_id: websiteId,
          reporter_name: reporterName,
          reporter_email: reporterEmail,
          reason,
          description,
        }),
      });

      const data = await res.json();
      if (!res.ok || !data.success) {
        setError(data || "Erro ao enviar denúncia.");
      } else {
        setTicket({ id: data.data.id, ticketNumber: data.data.ticket_number });
      }
    } catch {
      setError("Erro de conexão. Tente novamente.");
    } finally {
      setLoading(false);
    }
  }

  if (ticket) {
    return (
      <main className={styles.main}>
        <div className={styles.card}>
          <div className={styles.successIcon}>✅</div>
          <h2>Denúncia enviada!</h2>
          <p>Sua denúncia foi registrada com sucesso. Guarde o número do seu ticket.</p>
          <div className={styles.ticketBox}>
            <span>Número do ticket</span>
            <strong>#{String(ticket.ticketNumber).padStart(5, "0")}</strong>
          </div>
          <p className={styles.hint}>
            Você receberá um email quando nossa equipe responder.
          </p>
          <button className={styles.backBtn} onClick={() => navigate(-1)}>
            ← Voltar
          </button>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <div className={styles.card}>
        <div className={styles.header}>
          <span className={styles.flagIcon}>🚩</span>
          <div>
            <h1>Denunciar loja</h1>
            {websiteName && <p className={styles.storeName}>{websiteName}</p>}
          </div>
        </div>

        {error && <p className={styles.error}>{String(error)}</p>}

        <form onSubmit={handleSubmit} className={styles.form} noValidate>
          {!params.get("website_id") && (
            <div className={styles.field}>
              <label htmlFor="websiteId">ID da loja *</label>
              <input
                id="websiteId"
                value={websiteId}
                onChange={(e) => setWebsiteId(e.target.value)}
                placeholder="UUID da loja"
                required
                className={styles.input}
              />
            </div>
          )}

          <div className={styles.field}>
            <label htmlFor="reporterName">Seu nome *</label>
            <input
              id="reporterName"
              value={reporterName}
              onChange={(e) => setReporterName(e.target.value)}
              placeholder="Nome completo"
              required
              className={styles.input}
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="reporterEmail">Seu email *</label>
            <input
              id="reporterEmail"
              type="email"
              value={reporterEmail}
              onChange={(e) => setReporterEmail(e.target.value)}
              placeholder="email@exemplo.com"
              required
              className={styles.input}
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="reason">Motivo da denúncia *</label>
            <select
              id="reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className={styles.select}
              required
            >
              {REASONS.map((r) => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>

          <div className={styles.field}>
            <label htmlFor="description">Descrição *</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Descreva detalhadamente o problema…"
              rows={5}
              required
              minLength={20}
              className={styles.textarea}
            />
          </div>

          <div className={styles.formActions}>
            <button
              type="button"
              className={styles.cancelBtn}
              onClick={() => navigate(-1)}
            >
              Cancelar
            </button>
            <button
              type="submit"
              className={styles.submitBtn}
              disabled={loading}
            >
              {loading ? "Enviando…" : "Enviar denúncia"}
            </button>
          </div>
        </form>
      </div>
    </main>
  );
}
