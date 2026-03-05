import { useState, useRef } from "react";
import { useAuthContext } from "../../hooks/AuthContext";
import { uploadImage } from "../../lib/supabase";
import styles from "./Profile.module.scss";

export function Profile() {
  const { me, updateProfile, cancelPlan, loading } = useAuthContext();

  const [firstName, setFirstName] = useState(me?.first_name ?? "");
  const [lastName, setLastName] = useState(me?.last_name ?? "");
  const [cpfCnpj, setCpfCnpj] = useState(me?.cpf_cnpj ?? "");
  const [avatarUrl, setAvatarUrl] = useState(me?.avatar_url ?? "");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(
    me?.avatar_url ?? null
  );

  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [successMsg, setSuccessMsg] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [canceling, setCanceling] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  function handleAvatarChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  }

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setSuccessMsg("");
    setErrorMsg("");
    try {
      let finalAvatarUrl: string | null = avatarUrl || null;

      if (avatarFile) {
        setUploading(true);
        const path = `avatars/${me?.id ?? "user"}-${Date.now()}.${avatarFile.name.split(".").pop()}`;
        finalAvatarUrl = await uploadImage(avatarFile, path);
        setAvatarUrl(finalAvatarUrl);
        setUploading(false);
      }

      await updateProfile({
        first_name: firstName,
        last_name: lastName,
        cpf_cnpj: cpfCnpj || null,
        avatar_url: finalAvatarUrl,
      });
      setSuccessMsg("Perfil atualizado com sucesso!");
    } catch (err: any) {
      setErrorMsg(err.message ?? "Erro ao salvar perfil");
    } finally {
      setSaving(false);
      setUploading(false);
    }
  }

  async function handleCancelPlan() {
    setCanceling(true);
    try {
      await cancelPlan();
      setShowCancelModal(false);
      setSuccessMsg("Assinatura cancelada.");
    } catch (err: any) {
      setErrorMsg(err.message ?? "Erro ao cancelar plano");
    } finally {
      setCanceling(false);
    }
  }

  const maxSites = me?.plan_max_sites ?? 1;
  const maxRoutes = me?.plan_max_routes ?? 5;

  return (
    <div className={styles.page}>
      <div className={styles.container}>
        <h1 className={styles.title}>Meu Perfil</h1>

        {successMsg && <div className={styles.success}>{successMsg}</div>}
        {errorMsg && <div className={styles.error}>{errorMsg}</div>}

        <form className={styles.form} onSubmit={handleSave}>
          {/* Avatar */}
          <section className={styles.section}>
            <h2>Foto de Perfil</h2>
            <div className={styles.avatarRow}>
              <div
                className={styles.avatar}
                onClick={() => fileInputRef.current?.click()}
                title="Clique para alterar foto"
              >
                {avatarPreview ? (
                  <img src={avatarPreview} alt="avatar" />
                ) : (
                  <span className={styles.avatarPlaceholder}>
                    {(firstName || me?.first_name || "?")[0].toUpperCase()}
                  </span>
                )}
                <div className={styles.avatarOverlay}>✏️</div>
              </div>
              <div className={styles.avatarInfo}>
                <p>Clique na foto para alterar</p>
                <p className={styles.avatarHint}>JPG, PNG ou GIF. Máx 5MB.</p>
              </div>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                style={{ display: "none" }}
                onChange={handleAvatarChange}
              />
            </div>
          </section>

          {/* Personal Info */}
          <section className={styles.section}>
            <h2>Informações Pessoais</h2>
            <div className={styles.row}>
              <label>
                Nome
                <input
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  placeholder="Nome"
                  required
                />
              </label>
              <label>
                Sobrenome
                <input
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  placeholder="Sobrenome"
                  required
                />
              </label>
            </div>
            <label className={styles.fullWidth}>
              E-mail (não editável)
              <input value={me?.email ?? ""} disabled />
            </label>
            <label className={styles.fullWidth}>
              CPF / CNPJ
              <input
                value={cpfCnpj}
                onChange={(e) => setCpfCnpj(e.target.value)}
                placeholder="000.000.000-00 ou 00.000.000/0001-00"
                maxLength={18}
              />
            </label>
          </section>

          {/* Plan Info */}
          <section className={styles.section}>
            <h2>Plano Atual</h2>
            <div className={styles.planCard}>
              <div className={styles.planName}>
                {me?.user_plan ? (
                  <span className={styles.planBadge}>{me.user_plan}</span>
                ) : (
                  <span className={styles.noPlan}>Sem plano ativo</span>
                )}
              </div>
              <div className={styles.limits}>
                <div className={styles.limitItem}>
                  <span>Sites</span>
                  <div className={styles.bar}>
                    <div
                      className={styles.barFill}
                      style={{
                        width: `${Math.min(100, (1 / maxSites) * 100)}%`,
                      }}
                    />
                  </div>
                  <span>? / {maxSites}</span>
                </div>
                <div className={styles.limitItem}>
                  <span>Rotas por site</span>
                  <div className={styles.bar}>
                    <div className={styles.barFill} style={{ width: "0%" }} />
                  </div>
                  <span>máx. {maxRoutes}</span>
                </div>
              </div>
              <a href="/plans" className={styles.upgradeLink}>
                Ver planos →
              </a>
            </div>
          </section>

          <div className={styles.actions}>
            <button
              type="submit"
              className={styles.saveBtn}
              disabled={saving || uploading || loading}
            >
              {uploading
                ? "Enviando imagem..."
                : saving
                ? "Salvando..."
                : "Salvar alterações"}
            </button>
          </div>
        </form>

        {/* Danger Zone */}
        {me?.user_plan && (
          <section className={`${styles.section} ${styles.danger}`}>
            <h2>Zona de Perigo</h2>
            <p>Cancelar sua assinatura remove o acesso aos recursos do plano.</p>
            <button
              className={styles.cancelBtn}
              onClick={() => setShowCancelModal(true)}
            >
              Cancelar assinatura
            </button>
          </section>
        )}
      </div>

      {/* Cancel Confirmation Modal */}
      {showCancelModal && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Cancelar assinatura</h3>
            <p>
              Tem certeza que deseja cancelar sua assinatura? Você perderá
              acesso aos recursos do plano ao final do período vigente.
            </p>
            <div className={styles.modalActions}>
              <button
                className={styles.cancelBtn}
                onClick={handleCancelPlan}
                disabled={canceling}
              >
                {canceling ? "Cancelando..." : "Sim, cancelar"}
              </button>
              <button
                className={styles.secondaryBtn}
                onClick={() => setShowCancelModal(false)}
                disabled={canceling}
              >
                Voltar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
