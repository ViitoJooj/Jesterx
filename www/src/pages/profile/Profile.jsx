import { useState, useRef, useEffect, useCallback } from "react";
import { useAuthContext } from "../../hooks/AuthContext";
import { uploadImage, resolveMediaUrl } from "../../lib/storage";
import { BR_STATES } from "../../lib/constants";
import styles from "./Profile.module.scss";

export function Profile() {
  const { me, updateProfile, listAddresses, createAddress, updateAddress, deleteAddress, setDefaultAddress, cancelPlan, deleteAccount, loading, websiteId } = useAuthContext();

  const isBusiness = me?.account_type === "business";

  const [firstName, setFirstName] = useState(me?.first_name ?? "");
  const [lastName, setLastName] = useState(me?.last_name ?? "");
  const [cpfCnpj, setCpfCnpj] = useState(me?.cpf_cnpj ?? "");
  const [avatarUrl, setAvatarUrl] = useState(me?.avatar_url ?? "");
  const [avatarFile, setAvatarFile] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(resolveMediaUrl(me?.avatar_url) ?? null);
  const [displayName, setDisplayName] = useState(me?.display_name ?? "");
  const [birthDate, setBirthDate] = useState(me?.birth_date ?? "");
  const [gender, setGender] = useState(me?.gender ?? "");
  const [bio, setBio] = useState(me?.bio ?? "");
  const [instagram, setInstagram] = useState(me?.instagram ?? "");
  const [websiteUrl, setWebsiteUrl] = useState(me?.website_url ?? "");
  const [whatsapp, setWhatsapp] = useState(me?.whatsapp ?? "");

  const [companyName, setCompanyName] = useState(me?.company_name ?? "");
  const [tradeName, setTradeName] = useState(me?.trade_name ?? "");
  const [phone, setPhone] = useState(me?.phone ?? "");

  const [addresses, setAddresses] = useState([]);
  const [addrLoading, setAddrLoading] = useState(false);
  const [showAddrForm, setShowAddrForm] = useState(false);
  const [editingAddr, setEditingAddr] = useState(null);
  const [deleteAddrTarget, setDeleteAddrTarget] = useState(null);
  const emptyAddr = { label: "", zip_code: "", street: "", number: "", complement: "", district: "", city: "", state: "", country: "BR" };
  const [addrForm, setAddrForm] = useState(emptyAddr);

  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [successMsg, setSuccessMsg] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [canceling, setCanceling] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingAccount, setDeletingAccount] = useState(false);

  const fileInputRef = useRef(null);

  const loadAddresses = useCallback(async () => {
    setAddrLoading(true);
    try {
      const data = await listAddresses();
      setAddresses(data ?? []);
    } catch { /* ignore */ } finally {
      setAddrLoading(false);
    }
  }, [listAddresses]);

  useEffect(() => { loadAddresses(); }, [loadAddresses]);

  function openNewAddr() {
    setEditingAddr(null);
    setAddrForm(emptyAddr);
    setShowAddrForm(true);
  }

  function openEditAddr(addr) {
    setEditingAddr(addr.id);
    setAddrForm({
      label: addr.label ?? "",
      zip_code: addr.zip_code ?? "",
      street: addr.street ?? "",
      number: addr.number ?? "",
      complement: addr.complement ?? "",
      district: addr.district ?? "",
      city: addr.city ?? "",
      state: addr.state ?? "",
      country: addr.country ?? "BR",
    });
    setShowAddrForm(true);
  }

  async function handleAddrSubmit(e) {
    e.preventDefault();
    const payload = {
      label: addrForm.label || null,
      zip_code: addrForm.zip_code || null,
      street: addrForm.street || null,
      number: addrForm.number || null,
      complement: addrForm.complement || null,
      district: addrForm.district || null,
      city: addrForm.city || null,
      state: addrForm.state || null,
      country: addrForm.country || "BR",
    };
    try {
      if (editingAddr) {
        await updateAddress(editingAddr, payload);
      } else {
        await createAddress(payload);
      }
      setShowAddrForm(false);
      await loadAddresses();
    } catch (err) {
      setErrorMsg(err.message ?? "Erro ao salvar endereço");
    }
  }

  async function handleDeleteAddr() {
    if (!deleteAddrTarget) return;
    try {
      await deleteAddress(deleteAddrTarget);
      setDeleteAddrTarget(null);
      await loadAddresses();
    } catch (err) {
      setErrorMsg(err.message ?? "Erro ao remover endereço");
    }
  }

  async function handleSetDefault(id) {
    try {
      await setDefaultAddress(id);
      await loadAddresses();
    } catch (err) {
      setErrorMsg(err.message ?? "Erro ao definir endereço padrão");
    }
  }

  function handleAvatarChange(e) {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  }

  async function handleSave(e) {
    e.preventDefault();
    setSaving(true);
    setSuccessMsg("");
    setErrorMsg("");
    try {
      let finalAvatarUrl = avatarUrl || null;
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
      });
      setSuccessMsg("Perfil atualizado com sucesso!");
    } catch (err) {
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
    } catch (err) {
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
    } catch (err) {
      setErrorMsg(err.message ?? "Erro ao desativar conta");
    } finally {
      setDeletingAccount(false);
      setShowDeleteModal(false);
    }
  }

  const maxSites = me?.plan_max_sites ?? 1;
  const maxRoutes = me?.plan_max_routes ?? 5;

  const completionFields = [
    { label: "Nome", done: !!firstName.trim() },
    { label: "Sobrenome", done: !!lastName.trim() },
    { label: "Foto de perfil", done: !!avatarUrl },
    { label: "CPF / Telefone", done: !!(cpfCnpj.trim() || phone.trim()) },
    { label: "Data de nascimento", done: !!birthDate },
    { label: "Gênero", done: !!gender },
    { label: "Bio", done: !!bio.trim() },
    { label: "Instagram ou Website", done: !!(instagram.trim() || websiteUrl.trim()) },
    { label: "WhatsApp", done: !!whatsapp.trim() },
    { label: "Endereço", done: addresses.length > 0 },
  ];
  const completionPct = Math.round((completionFields.filter((f) => f.done).length / completionFields.length) * 100);
  const completionColor = completionPct >= 80 ? "#16a34a" : completionPct >= 50 ? "#d97706" : "#dc2626";
  const completionLabel = completionPct >= 80 ? "Perfil completo" : completionPct >= 50 ? "Em progresso" : "Incompleto";

  return (
    <div className={styles.page}>
      <div className={styles.layout}>

        {/* ── Sidebar ── */}
        <aside className={styles.sidebar}>
          <div className={styles.sidebarAvatarWrap} onClick={() => fileInputRef.current?.click()} title="Clique para alterar foto">
            {avatarPreview
              ? <img src={avatarPreview} alt="avatar" className={styles.sidebarAvatar} />
              : <div className={styles.sidebarAvatarFallback}>{(firstName || "?")[0].toUpperCase()}</div>
            }
            <div className={styles.sidebarAvatarOverlay}>✏</div>
          </div>
          <input ref={fileInputRef} type="file" accept="image/*" style={{ display: "none" }} onChange={handleAvatarChange} />

          <div className={styles.sidebarProfile}>
            <span className={styles.sidebarFullName}>
              {[firstName, lastName].filter(Boolean).join(" ") || "Seu nome"}
            </span>
            {displayName && <span className={styles.sidebarDisplayName}>{displayName}</span>}
            <span className={styles.sidebarEmail}>{me?.email}</span>
          </div>

          <span className={`${styles.sidebarBadge} ${isBusiness ? styles.sidebarBadgeBiz : ""}`}>
            {isBusiness ? "🏢 Empresa" : "👤 Pessoal"}
          </span>

          {bio && <p className={styles.sidebarBioText}>{bio}</p>}

          <div className={styles.sidebarLinks}>
            {instagram && <span className={styles.sidebarLink}>📸 {instagram}</span>}
            {websiteUrl && <span className={styles.sidebarLink}>🔗 {websiteUrl}</span>}
            {whatsapp && <span className={styles.sidebarLink}>💬 {whatsapp}</span>}
          </div>

          <div className={styles.sidebarSep} />

          <div className={styles.completionBlock}>
            <div className={styles.completionTop}>
              <span className={styles.completionLabel} style={{ color: completionColor }}>{completionLabel}</span>
              <span className={styles.completionPct} style={{ color: completionColor }}>{completionPct}%</span>
            </div>
            <div className={styles.completionBarBg}>
              <div className={styles.completionBarFill} style={{ width: `${completionPct}%`, background: completionColor }} />
            </div>
            <div className={styles.completionChecks}>
              {completionFields.map((f) => (
                <span key={f.label} className={f.done ? styles.checkDone : styles.checkMissing}>
                  {f.done ? "✓" : "○"} {f.label}
                </span>
              ))}
            </div>
          </div>

          <div className={styles.sidebarSep} />

          <div className={styles.sidebarPlanBlock}>
            <div className={styles.sidebarPlanRow}>
              {me?.user_plan
                ? <span className={styles.planBadge}>{me.user_plan}</span>
                : <span className={styles.noPlanText}>Sem plano ativo</span>
              }
              <a href="/plans" className={styles.upgradeLink}>Ver planos →</a>
            </div>
            <div className={styles.planLimits}>
              <div className={styles.planLimit}>
                <span>Sites</span>
                <strong>{maxSites}</strong>
              </div>
              <div className={styles.planLimit}>
                <span>Rotas/site</span>
                <strong>{maxRoutes}</strong>
              </div>
            </div>
          </div>
        </aside>

        {/* ── Main content ── */}
        <div className={styles.mainContent}>
          {successMsg && <div className={styles.success}>{successMsg}</div>}
          {errorMsg && <div className={styles.error}>{errorMsg}</div>}

          <form onSubmit={handleSave}>

            <section className={styles.section}>
              <div className={styles.sectionHead}>
                <h2>Informações Pessoais</h2>
                <p>Seus dados de identificação e contato</p>
              </div>
              <div className={styles.sectionBody}>
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
                  E-mail <span className={styles.fieldNote}>(não editável)</span>
                  <input value={me?.email ?? ""} disabled />
                </label>
                <div className={styles.row}>
                  <label>
                    Nome de exibição
                    <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} placeholder="Como aparecer na loja" />
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
                  <div className={styles.row}>
                    <label>
                      CPF
                      <input value={cpfCnpj} onChange={(e) => setCpfCnpj(e.target.value)} placeholder="000.000.000-00" maxLength={14} />
                    </label>
                    <label>
                      Telefone
                      <input value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="(11) 91234-5678" maxLength={16} />
                    </label>
                  </div>
                )}
              </div>
            </section>

            <section className={styles.section}>
              <div className={styles.sectionHead}>
                <h2>Perfil Público</h2>
                <p>Informações exibidas na sua loja e avaliações</p>
              </div>
              <div className={styles.sectionBody}>
                <label className={styles.fullWidth}>
                  Bio
                  <textarea value={bio} onChange={(e) => setBio(e.target.value)} placeholder="Fale um pouco sobre você..." rows={3} />
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
                <label>
                  WhatsApp
                  <input value={whatsapp} onChange={(e) => setWhatsapp(e.target.value)} placeholder="(11) 91234-5678" maxLength={20} />
                </label>
              </div>
            </section>

            {isBusiness && (
              <section className={styles.section}>
                <div className={styles.sectionHead}>
                  <h2>
                    Dados da Empresa
                    <span className={styles.businessBadge}>Empresa</span>
                  </h2>
                </div>
                <div className={styles.sectionBody}>
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
                </div>
              </section>
            )}

            <section className={styles.section}>
              <div className={styles.sectionHead}>
                <h2>Endereços</h2>
                <p>Usados automaticamente para envio e cobranças</p>
              </div>
              <div className={styles.sectionBody}>
                {addrLoading && <p className={styles.addrLoading}>Carregando...</p>}
                {addresses.map((addr) => (
                  <div key={addr.id} className={styles.addrCard}>
                    <div className={styles.addrInfo}>
                      {addr.label && <strong>{addr.label}</strong>}
                      <span>{[addr.street, addr.number].filter(Boolean).join(", ")}</span>
                      {addr.complement && <span>{addr.complement}</span>}
                      <span>{[addr.district, addr.city, addr.state].filter(Boolean).join(" · ")}</span>
                      <span>CEP {addr.zip_code ?? "—"}</span>
                    </div>
                    <div className={styles.addrActions}>
                      {addr.is_default
                        ? <span className={styles.defaultBadge}>Padrão</span>
                        : <button type="button" className={styles.ghostBtn} onClick={() => handleSetDefault(addr.id)}>Tornar padrão</button>
                      }
                      <button type="button" className={styles.ghostBtn} onClick={() => openEditAddr(addr)}>Editar</button>
                      <button type="button" className={styles.dangerGhostBtn} onClick={() => setDeleteAddrTarget(addr.id)}>Remover</button>
                    </div>
                  </div>
                ))}
                {!showAddrForm && addresses.length < 10 && (
                  <button type="button" className={styles.addAddrBtn} onClick={openNewAddr}>+ Adicionar endereço</button>
                )}
                {showAddrForm && (
                  <form className={styles.addrForm} onSubmit={handleAddrSubmit}>
                    <div className={styles.addrFormHead}>
                      <h3>{editingAddr ? "Editar endereço" : "Novo endereço"}</h3>
                    </div>
                    <div className={styles.row}>
                      <label>
                        Etiqueta
                        <input value={addrForm.label} onChange={(e) => setAddrForm((f) => ({ ...f, label: e.target.value }))} placeholder="Casa, Trabalho..." />
                      </label>
                      <label className={styles.zipField}>
                        CEP
                        <input value={addrForm.zip_code} onChange={(e) => setAddrForm((f) => ({ ...f, zip_code: e.target.value }))} placeholder="00000-000" maxLength={9} />
                      </label>
                    </div>
                    <div className={styles.row}>
                      <label>
                        Rua / Avenida
                        <input value={addrForm.street} onChange={(e) => setAddrForm((f) => ({ ...f, street: e.target.value }))} placeholder="Av. Paulista" />
                      </label>
                      <label className={styles.numberField}>
                        Número
                        <input value={addrForm.number} onChange={(e) => setAddrForm((f) => ({ ...f, number: e.target.value }))} placeholder="1000" />
                      </label>
                    </div>
                    <div className={styles.row}>
                      <label>
                        Cidade
                        <input value={addrForm.city} onChange={(e) => setAddrForm((f) => ({ ...f, city: e.target.value }))} placeholder="São Paulo" />
                      </label>
                      <label>
                        Bairro
                        <input value={addrForm.district} onChange={(e) => setAddrForm((f) => ({ ...f, district: e.target.value }))} placeholder="Centro" />
                      </label>
                    </div>
                    <div className={styles.row}>
                      <label>
                        Estado
                        <select className={styles.select} value={addrForm.state} onChange={(e) => setAddrForm((f) => ({ ...f, state: e.target.value }))}>
                          <option value="">Selecione</option>
                          {BR_STATES.map((s) => <option key={s} value={s}>{s}</option>)}
                        </select>
                      </label>
                      <label>
                        Complemento
                        <input value={addrForm.complement} onChange={(e) => setAddrForm((f) => ({ ...f, complement: e.target.value }))} placeholder="Sala 10" />
                      </label>
                    </div>
                    <div className={styles.addrFormActions}>
                      <button type="submit" className={styles.saveBtn}>{editingAddr ? "Salvar" : "Adicionar"}</button>
                      <button type="button" className={styles.ghostBtn} onClick={() => setShowAddrForm(false)}>Cancelar</button>
                    </div>
                  </form>
                )}
              </div>
            </section>

            <div className={styles.formActions}>
              <button type="submit" className={styles.saveBtn} disabled={saving || uploading || loading}>
                {uploading ? "Enviando imagem..." : saving ? "Salvando..." : "Salvar alterações"}
              </button>
            </div>
          </form>

          <section className={`${styles.section} ${styles.dangerSection}`}>
            <div className={styles.sectionHead}>
              <h2>Zona de Perigo</h2>
              <p>Ações irreversíveis — prossiga com cuidado</p>
            </div>
            <div className={styles.sectionBody}>
              {me?.user_plan && (
                <div className={styles.dangerRow}>
                  <div>
                    <strong>Cancelar assinatura</strong>
                    <p>Remove o acesso aos recursos do plano ao final do período vigente.</p>
                  </div>
                  <button type="button" className={styles.dangerOutlineBtn} onClick={() => setShowCancelModal(true)}>
                    Cancelar assinatura
                  </button>
                </div>
              )}
              <div className={`${styles.dangerRow} ${me?.user_plan ? styles.dangerRowBorder : ""}`}>
                <div>
                  <strong>Deletar conta</strong>
                  <p>Desativa o acesso imediatamente e agenda exclusão definitiva em 30 dias.</p>
                </div>
                <button type="button" className={styles.dangerSolidBtn} onClick={() => setShowDeleteModal(true)}>
                  Deletar conta
                </button>
              </div>
            </div>
          </section>
        </div>
      </div>

      {/* ── Modals ── */}

      {deleteAddrTarget && (
        <div className={styles.modalOverlay} onClick={() => setDeleteAddrTarget(null)}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <h3>Remover endereço</h3>
            <p>Tem certeza que deseja remover este endereço? Esta ação não pode ser desfeita.</p>
            <div className={styles.modalActions}>
              <button className={styles.ghostBtn} onClick={() => setDeleteAddrTarget(null)}>Cancelar</button>
              <button className={styles.dangerSolidBtn} onClick={handleDeleteAddr}>Remover</button>
            </div>
          </div>
        </div>
      )}

      {showCancelModal && (
        <div className={styles.modalOverlay} onClick={() => !canceling && setShowCancelModal(false)}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <h3>Cancelar assinatura</h3>
            <p>Tem certeza? Você perderá o acesso aos recursos do plano ao final do período vigente.</p>
            <div className={styles.modalActions}>
              <button className={styles.ghostBtn} onClick={() => setShowCancelModal(false)} disabled={canceling}>Voltar</button>
              <button className={styles.dangerOutlineBtn} onClick={handleCancelPlan} disabled={canceling}>
                {canceling ? "Cancelando..." : "Sim, cancelar"}
              </button>
            </div>
          </div>
        </div>
      )}

      {showDeleteModal && (
        <div className={styles.modalOverlay} onClick={() => !deletingAccount && setShowDeleteModal(false)}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <h3>Deletar conta</h3>
            <p>Sua conta será desativada agora e removida permanentemente em 30 dias. Durante esse período você não conseguirá fazer login.</p>
            <div className={styles.modalActions}>
              <button className={styles.ghostBtn} onClick={() => setShowDeleteModal(false)} disabled={deletingAccount}>Voltar</button>
              <button className={styles.dangerSolidBtn} onClick={handleDeleteAccount} disabled={deletingAccount}>
                {deletingAccount ? "Desativando..." : "Sim, deletar conta"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
