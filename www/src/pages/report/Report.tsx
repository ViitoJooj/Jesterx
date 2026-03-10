import { useEffect, useRef, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { API_URL } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./Report.module.scss";

const REASONS = [
  { value: "FRAUD", label: "Fraude" },
  { value: "SCAM", label: "Golpe / Estelionato" },
  { value: "SPAM", label: "Spam" },
  { value: "INAPPROPRIATE", label: "Conteúdo Inapropriado" },
  { value: "COUNTERFEIT", label: "Produto Falsificado" },
  { value: "OTHER", label: "Outro" },
];

const MAX_IMAGES = 5;
const MAX_IMAGE_BYTES = 1_000_000; // 1MB raw (≈ 1.33MB base64)

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export function Report() {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { me } = useAuthContext();

  const [websiteId, setWebsiteId] = useState(params.get("website_id") ?? "");
  const [websiteName] = useState(params.get("website_name") ?? "");
  const [reporterName, setReporterName] = useState("");
  const [reporterEmail, setReporterEmail] = useState("");
  const [reason, setReason] = useState("FRAUD");
  const [description, setDescription] = useState("");
  const [evidencePreviews, setEvidencePreviews] = useState<string[]>([]);
  const [evidenceB64, setEvidenceB64] = useState<string[]>([]);
  const [imageError, setImageError] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [ticket, setTicket] = useState<{ id: string; ticketNumber: number } | null>(null);

  // Auto-fill from authenticated user
  useEffect(() => {
    if (me) {
      setReporterName(`${me.first_name ?? ""} ${me.last_name ?? ""}`.trim());
      setReporterEmail(me.email ?? "");
    }
  }, [me]);

  useEffect(() => {
    const id = params.get("website_id");
    if (id) setWebsiteId(id);
  }, [params]);

  async function handleImageChange(e: React.ChangeEvent<HTMLInputElement>) {
    setImageError(null);
    const files = Array.from(e.target.files ?? []);
    if (evidenceB64.length + files.length > MAX_IMAGES) {
      setImageError(`Máximo de ${MAX_IMAGES} imagens permitidas.`);
      return;
    }
    const newB64: string[] = [];
    const newPreviews: string[] = [];
    for (const file of files) {
      if (file.size > MAX_IMAGE_BYTES) {
        setImageError(`"${file.name}" excede o limite de 1MB.`);
        return;
      }
      if (!file.type.startsWith("image/")) {
        setImageError("Apenas imagens são permitidas (PNG, JPG, GIF, etc.).");
        return;
      }
      const b64 = await fileToBase64(file);
      newB64.push(b64);
      newPreviews.push(b64);
    }
    setEvidenceB64((prev) => [...prev, ...newB64]);
    setEvidencePreviews((prev) => [...prev, ...newPreviews]);
    // Reset input so same file can be re-added after removal
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  function removeImage(index: number) {
    setEvidenceB64((prev) => prev.filter((_, i) => i !== index));
    setEvidencePreviews((prev) => prev.filter((_, i) => i !== index));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/api/v1/reports`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          website_id: websiteId,
          reporter_name: reporterName,
          reporter_email: reporterEmail,
          reason,
          description,
          evidence_urls: evidenceB64,
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
            <label htmlFor="reporterName">
              Seu nome *{me && <span className={styles.autoFilled}> (preenchido automaticamente)</span>}
            </label>
            <input
              id="reporterName"
              value={reporterName}
              onChange={(e) => !me && setReporterName(e.target.value)}
              readOnly={!!me}
              placeholder="Nome completo"
              required
              className={`${styles.input} ${me ? styles.readOnly : ""}`}
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="reporterEmail">
              Seu email *{me && <span className={styles.autoFilled}> (preenchido automaticamente)</span>}
            </label>
            <input
              id="reporterEmail"
              type="email"
              value={reporterEmail}
              onChange={(e) => !me && setReporterEmail(e.target.value)}
              readOnly={!!me}
              placeholder="email@exemplo.com"
              required
              className={`${styles.input} ${me ? styles.readOnly : ""}`}
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

          {/* Evidence images */}
          <div className={styles.field}>
            <label>
              Evidências (opcional) — até {MAX_IMAGES} imagens, 1MB cada
            </label>
            {imageError && <p className={styles.error}>{imageError}</p>}
            {evidencePreviews.length > 0 && (
              <div className={styles.evidenceGrid}>
                {evidencePreviews.map((src, i) => (
                  <div key={i} className={styles.evidenceThumb}>
                    <img src={src} alt={`Evidência ${i + 1}`} />
                    <button
                      type="button"
                      className={styles.removeImg}
                      onClick={() => removeImage(i)}
                      title="Remover imagem"
                    >
                      ✕
                    </button>
                  </div>
                ))}
              </div>
            )}
            {evidenceB64.length < MAX_IMAGES && (
              <>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  multiple
                  onChange={handleImageChange}
                  className={styles.fileInput}
                  id="evidenceInput"
                />
                <label htmlFor="evidenceInput" className={styles.uploadBtn}>
                  📎 Adicionar imagens
                </label>
              </>
            )}
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

