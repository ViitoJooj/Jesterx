import { useState, useRef } from "react";
import { useAuthContext } from "../../hooks/AuthContext";
import { uploadImage, resolveMediaUrl } from "../../lib/storage";
import styles from "./Profile.module.scss";

const BR_STATES = [
  "AC","AL","AP","AM","BA","CE","DF","ES","GO","MA","MT","MS",
  "MG","PA","PB","PR","PE","PI","RJ","RN","RS","RO","RR","SC",
  "SP","SE","TO",
];

export function Profile() {
  const { me, updateProfile, cancelPlan, deleteAccount, loading, websiteId } = useAuthContext();

  const isBusiness = me?.account_type === "business";

  const [firstName, setFirstName] = useState(me?.first_name ?? "");
  const [lastName, setLastName] = useState(me?.last_name ?? "");
  const [cpfCnpj, setCpfCnpj] = useState(me?.cpf_cnpj ?? "");
  const [avatarUrl, setAvatarUrl] = useState(me?.avatar_url ?? "");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(resolveMediaUrl(me?.avatar_url) ?? null);
  const [displayName, setDisplayName] = useState(me?.display_name ?? "");
  const [birthDate, setBirthDate] = useState(me?.birth_date ?? "");
  const [gender, setGender] = useState(me?.gender ?? "");
  const [bio, setBio] = useState(me?.bio ?? "");
  const [instagram, setInstagram] = useState(me?.instagram ?? "");
  const [websiteUrl, setWebsiteUrl] = useState(me?.website_url ?? "");
  const [whatsapp, setWhatsapp] = useState(me?.whatsapp ?? "");

  // business fields
  const [companyName, setCompanyName] = useState(me?.company_name ?? "");
  const [tradeName, setTradeName] = useState(me?.trade_name ?? "");
  const [phone, setPhone] = useState(me?.phone ?? "");
  const [zipCode, setZipCode] = useState(me?.zip_code ?? "");
  const [addressStreet, setAddressStreet] = useState(me?.address_street ?? "");
  const [addressNumber, setAddressNumber] = useState(me?.address_number ?? "");
  const [addressComplement, setAddressComplement] = useState(me?.address_complement ?? "");
  const [addressDistrict, setAddressDistrict] = useState(me?.address_district ?? "");
  const [addressCity, setAddressCity] = useState(me?.address_city ?? "");
  const [addressState, setAddressState] = useState(me?.address_state ?? "");
  const [addressCountry, setAddressCountry] = useState(me?.address_country ?? "");

  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [successMsg, setSuccessMsg] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [canceling, setCanceling] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingAccount, setDeletingAccount] = useState(false);

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
        finalAvatarUrl = await uploadImage(avatarFile, websiteId);
        setAvatarUrl(finalAvatarUrl);
        setUploading(false);
      }

      await updateProfile({
        first_name: firstName,
        last_name: lastName,
        cpf_cnpj: cpfCnpj || null,
        avatar_url: finalAvatarUrl,
        phone: phone || null,
        company_name: isBusiness ? companyName || null : null,
        trade_name: isBusiness ? tradeName || null : null,
        display_name: displayName || null,
        birth_date: birthDate || null,
        gender: gender || null,
        bio: bio || null,
        instagram: instagram || null,
        website_url: websiteUrl || null,
        whatsapp: whatsapp || null,
        zip_code: zipCode || null,
        address_street: addressStreet || null,
        address_number: addressNumber || null,
        address_complement: addressComplement || null,
        address_district: addressDistrict || null,
        address_city: addressCity || null,
        address_state: addressState || null,
        address_country: addressCountry || null,
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

  async function handleDeleteAccount() {
    setDeletingAccount(true);
    try {
      await deleteAccount();
      window.location.href = "/login";
    } catch (err: any) {
      setErrorMsg(err.message ?? "Erro ao desativar conta");
    } finally {
      setDeletingAccount(false);
      setShowDeleteModal(false);
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
                <input value={firstName} onChange={(e) => setFirstName(e.target.value)} placeholder="Nome" required />
              </label>
              <label>
                Sobrenome
                <input value={lastName} onChange={(e) => setLastName(e.target.value)} placeholder="Sobrenome" required />
              </label>
            </div>
            <label className={styles.fullWidth}>
              E-mail (não editável)
              <input value={me?.email ?? ""} disabled />
            </label>
            <div className={styles.row}>
              <label>
                Nome de exibição
                <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} placeholder="Como você quer aparecer na loja" />
              </label>
              <label>
                Data de nascimento
                <input type="date" value={birthDate} onChange={(e) => setBirthDate(e.target.value)} />
              </label>
            </div>
            <label className={styles.fullWidth}>
              Gênero
              <select className={styles.select} value={gender} onChange={(e) => setGender(e.target.value)}>
                <option value="">Prefiro não informar</option>
                <option value="male">Masculino</option>
                <option value="female">Feminino</option>
                <option value="other">Outro</option>
                <option value="prefer_not">Prefiro não dizer</option>
              </select>
            </label>
            {!isBusiness && (
              <label className={styles.fullWidth}>
                CPF
                <input
                  value={cpfCnpj}
                  onChange={(e) => setCpfCnpj(e.target.value)}
                  placeholder="000.000.000-00"
                  maxLength={14}
                />
              </label>
            )}
            {!isBusiness && (
              <label className={styles.fullWidth}>
                Telefone
                <input
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="(11) 91234-5678"
                  maxLength={16}
                />
              </label>
            )}
          </section>

          <section className={styles.section}>
            <h2>Perfil na Loja</h2>
            <label className={styles.fullWidth}>
              Bio
              <textarea
                value={bio}
                onChange={(e) => setBio(e.target.value)}
                placeholder="Fale um pouco sobre você"
                rows={3}
              />
            </label>
            <div className={styles.row}>
              <label>
                Instagram
                <input value={instagram} onChange={(e) => setInstagram(e.target.value)} placeholder="@seuperfil" />
              </label>
              <label>
                Website
                <input value={websiteUrl} onChange={(e) => setWebsiteUrl(e.target.value)} placeholder="https://seusite.com" />
              </label>
            </div>
            <div className={styles.row}>
              <label>
                WhatsApp
                <input value={whatsapp} onChange={(e) => setWhatsapp(e.target.value)} placeholder="(11) 91234-5678" maxLength={20} />
              </label>
            </div>
          </section>

          {/* Business Info */}
          {isBusiness && (
            <section className={styles.section}>
              <h2>
                <span className={styles.businessBadge}>Empresa</span>
                Dados da Empresa
              </h2>
              <div className={styles.row}>
                <label>
                  Razão Social
                  <input value={companyName} onChange={(e) => setCompanyName(e.target.value)} placeholder="Empresa LTDA" />
                </label>
                <label>
                  Nome Fantasia
                  <input value={tradeName} onChange={(e) => setTradeName(e.target.value)} placeholder="Nome Fantasia" />
                </label>
              </div>
              <div className={styles.row}>
                <label>
                  CNPJ
                  <input value={cpfCnpj} onChange={(e) => setCpfCnpj(e.target.value)} placeholder="00.000.000/0001-00" maxLength={18} />
                </label>
                <label>
                  Telefone
                  <input value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="(11) 91234-5678" maxLength={16} />
                </label>
              </div>
            </section>
          )}

          {/* Address — shown for all account types */}
          <section className={styles.section}>
            <h2>Endereço</h2>
            <div className={styles.row}>
              <label className={styles.zipField}>
                CEP
                <input value={zipCode} onChange={(e) => setZipCode(e.target.value)} placeholder="00000-000" maxLength={9} />
              </label>
              <label>
                Estado
                <select className={styles.select} value={addressState} onChange={(e) => setAddressState(e.target.value)}>
                  <option value="">Selecione</option>
                  {BR_STATES.map((s) => <option key={s} value={s}>{s}</option>)}
                </select>
              </label>
            </div>
            <div className={styles.row}>
              <label>
                Rua / Avenida
                <input value={addressStreet} onChange={(e) => setAddressStreet(e.target.value)} placeholder="Av. Paulista" />
              </label>
              <label className={styles.numberField}>
                Número
                <input value={addressNumber} onChange={(e) => setAddressNumber(e.target.value)} placeholder="1000" />
              </label>
            </div>
            <div className={styles.row}>
              <label>
                Cidade
                <input value={addressCity} onChange={(e) => setAddressCity(e.target.value)} placeholder="São Paulo" />
              </label>
              <label>
                Bairro
                <input value={addressDistrict} onChange={(e) => setAddressDistrict(e.target.value)} placeholder="Centro" />
              </label>
            </div>
            <div className={styles.row}>
              <label>
                Complemento
                <input value={addressComplement} onChange={(e) => setAddressComplement(e.target.value)} placeholder="Sala 10" />
              </label>
              <label>
                País
                <input value={addressCountry} onChange={(e) => setAddressCountry(e.target.value)} placeholder="Brasil" />
              </label>
            </div>
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
                    <div className={styles.barFill} style={{ width: `${Math.min(100, (1 / maxSites) * 100)}%` }} />
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
              <a href="/plans" className={styles.upgradeLink}>Ver planos →</a>
            </div>
          </section>

          <div className={styles.actions}>
            <button type="submit" className={styles.saveBtn} disabled={saving || uploading || loading}>
              {uploading ? "Enviando imagem..." : saving ? "Salvando..." : "Salvar alterações"}
            </button>
          </div>
        </form>

        {/* Danger Zone */}
        <section className={`${styles.section} ${styles.danger}`}>
          <h2>Zona de Perigo</h2>
          {me?.user_plan && (
            <>
              <p>Cancelar sua assinatura remove o acesso aos recursos do plano.</p>
              <button className={styles.cancelBtn} onClick={() => setShowCancelModal(true)}>
                Cancelar assinatura
              </button>
            </>
          )}
          <p>A ação de deletar conta desativa o acesso imediatamente e agenda exclusão definitiva em 30 dias.</p>
          <button className={styles.deleteBtn} onClick={() => setShowDeleteModal(true)}>
            Deletar conta
          </button>
        </section>
      </div>

      {showCancelModal && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Cancelar assinatura</h3>
            <p>
              Tem certeza que deseja cancelar sua assinatura? Você perderá
              acesso aos recursos do plano ao final do período vigente.
            </p>
            <div className={styles.modalActions}>
              <button className={styles.cancelBtn} onClick={handleCancelPlan} disabled={canceling}>
                {canceling ? "Cancelando..." : "Sim, cancelar"}
              </button>
              <button className={styles.secondaryBtn} onClick={() => setShowCancelModal(false)} disabled={canceling}>
                Voltar
              </button>
            </div>
          </div>
        </div>
      )}

      {showDeleteModal && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Deletar conta</h3>
            <p>
              Sua conta será desativada agora e removida permanentemente em 30 dias.
              Durante esse período você não conseguirá fazer login.
            </p>
            <div className={styles.modalActions}>
              <button className={styles.deleteBtn} onClick={handleDeleteAccount} disabled={deletingAccount}>
                {deletingAccount ? "Desativando..." : "Sim, deletar conta"}
              </button>
              <button className={styles.secondaryBtn} onClick={() => setShowDeleteModal(false)} disabled={deletingAccount}>
                Voltar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
