import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { API_URL, apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./StoreInfo.module.scss";

// ─── Types ────────────────────────────────────────────────────────────────────

type StoreMember = { id: string; user_id: string; user_name: string; avatar_url?: string; role: string };

type Creator = {
  id: string;
  full_name: string;
  company_name?: string;
  trade_name?: string;
  avatar_url?: string;
  account_type: string;
};

type StoreInfo = {
  id: string;
  name: string;
  short_description: string;
  description: string;
  image?: string;
  type: string;
  mature_content: boolean;
  rating_avg: number;
  rating_count: number;
  editor_type: string;
  creator: Creator;
  managers?: StoreMember[];
};

type Comment = {
  id: string;
  user_id: string;
  user_name: string;
  avatar_url?: string;
  content: string;
  stars?: number;
  parent_comment_id?: string;
  replies?: Comment[];
  created_at: string;
};

type VisitDay = { date: string; count: number };

// ─── Constants ────────────────────────────────────────────────────────────────

const TYPE_LABELS: Record<string, string> = {
  ECOMMERCE: "E-commerce",
  LANDING_PAGE: "Landing Page",
  SOFTWARE_SELL: "Venda de Software",
  COURSE: "Curso",
  VIDEO: "Canal de Vídeo",
};

const EDITOR_LABELS: Record<string, { label: string; color: string }> = {
  ELEMENTOR_JSON: { label: "Elementor", color: "#e44c65" },
  REACT:          { label: "React",     color: "#61dafb" },
  SVELTE:         { label: "Svelte",    color: "var(--jx-color-primary)" },
};

// ─── Star Rating Component ────────────────────────────────────────────────────

function StarRating({
  value,
  max = 5,
  interactive = false,
  onChange,
}: {
  value: number;
  max?: number;
  interactive?: boolean;
  onChange?: (n: number) => void;
}) {
  const [hover, setHover] = useState(0);
  return (
    <span className={styles.stars}>
      {Array.from({ length: max }, (_, i) => i + 1).map((n) => (
        <span
          key={n}
          className={`${styles.star} ${(hover || value) >= n ? styles.starFilled : ""}`}
          onMouseEnter={() => interactive && setHover(n)}
          onMouseLeave={() => interactive && setHover(0)}
          onClick={() => interactive && onChange?.(n)}
          style={{ cursor: interactive ? "pointer" : "default" }}
        >
          ★
        </span>
      ))}
    </span>
  );
}

// ─── Visit Chart (SVG bar chart) ─────────────────────────────────────────────

function VisitChart({ data }: { data: VisitDay[] }) {
  if (!data.length) return null;
  const max = Math.max(...data.map((d) => d.count), 1);

  return (
    <div className={styles.chartWrap}>
      <svg viewBox={`0 0 ${data.length * 10} 60`} preserveAspectRatio="none" className={styles.chartSvg}>
        {data.map((d, i) => {
          const h = (d.count / max) * 55;
          return (
            <rect
              key={d.date}
              x={i * 10 + 0.5}
              y={60 - h}
              width={9}
              height={h}
              rx={1.5}
              className={styles.chartBar}
            />
          );
        })}
      </svg>
      <div className={styles.chartLabels}>
        {[data[0], data[Math.floor(data.length / 2)], data[data.length - 1]].filter(Boolean).map((d) => (
          <span key={d!.date}>{d!.date.slice(5)}</span>
        ))}
      </div>
    </div>
  );
}

// ─── Main Component ──────────────────────────────────────────────────────────

export function StoreInfo() {
  const { siteId } = useParams<{ siteId: string }>();
  const navigate = useNavigate();
  const { me: user, accessToken, websiteId } = useAuthContext();

  const [store, setStore] = useState<StoreInfo | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [visits, setVisits] = useState<VisitDay[]>([]);
  const [myRating, setMyRating] = useState(0);
  const [myRole, setMyRole] = useState<string>("");

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Members management (owner/manager/admin)
  const [members, setMembers] = useState<StoreMember[]>([]);
  const [addMemberId, setAddMemberId] = useState("");
  const [addMemberRole, setAddMemberRole] = useState("support");
  const [addingMember, setAddingMember] = useState(false);
  const [memberError, setMemberError] = useState("");

  // New comment form
  const [commentText, setCommentText] = useState("");
  const [commentStars, setCommentStars] = useState(0);
  const [submittingComment, setSubmittingComment] = useState(false);
  const [ratingLoading, setRatingLoading] = useState(false);

  // Reply state: which comment are we replying to?
  const [replyToId, setReplyToId] = useState<string | null>(null);
  const [replyText, setReplyText] = useState("");
  const [submittingReply, setSubmittingReply] = useState(false);

  // Edit mode (owner/manager)
  const [editMode, setEditMode] = useState(false);
  const [editName, setEditName] = useState("");
  const [editShort, setEditShort] = useState("");
  const [editDesc, setEditDesc] = useState("");
  const [editSaving, setEditSaving] = useState(false);

  const isOwner = user && store && user.id === store.creator.id;
  const canReply = myRole === "owner" || myRole === "admin" || myRole === "manager" || myRole === "support";
  const canEdit = myRole === "owner" || myRole === "admin" || myRole === "manager";
  const canManageMembers = myRole === "owner" || myRole === "admin" || myRole === "manager";

  // ── Load data ──────────────────────────────────────────────────────────────
  useEffect(() => {
    if (!siteId) return;
    Promise.all([
      fetch(`${API_URL}/api/store/${siteId}/info`).then((r) => r.json()),
      fetch(`${API_URL}/api/store/${siteId}/comments`).then((r) => r.json()),
      fetch(`${API_URL}/api/store/${siteId}/visits?days=30`).then((r) => r.json()),
    ])
      .then(([infoRes, commentsRes, visitsRes]) => {
        if (infoRes.success) setStore(infoRes.data);
        else setError("Loja não encontrada.");
        if (commentsRes.success) setComments(commentsRes.data ?? []);
        if (visitsRes.success) setVisits(visitsRes.data ?? []);
      })
      .catch(() => setError("Erro ao carregar a loja."))
      .finally(() => setLoading(false));
  }, [siteId]);

  // ── Load my rating + my role ───────────────────────────────────────────────
  useEffect(() => {
    if (!siteId || !user || !accessToken) return;
    apiFetch<{ success: boolean; data?: { stars: number } }>(
      `/api/store/${siteId}/my-rating`, { websiteId, accessToken }
    ).then((d) => { if (d.success && d.data) setMyRating(d.data.stars); }).catch(() => {});

    apiFetch<{ success: boolean; data?: { role: string } }>(
      `/api/store/${siteId}/my-role`, { websiteId, accessToken }
    ).then((d) => { if (d.success && d.data) setMyRole(d.data.role); }).catch(() => {});
  }, [siteId, user, accessToken]);

  // ── Load members (owner/manager/admin) ────────────────────────────────────
  useEffect(() => {
    if (!siteId || !accessToken || !canManageMembers) return;
    apiFetch<{ success: boolean; data?: StoreMember[] }>(
      `/api/v1/sites/${siteId}/members`, { websiteId, accessToken }
    ).then((d) => { if (d.success) setMembers(d.data ?? []); }).catch(() => {});
  }, [siteId, accessToken, canManageMembers]);

  // ── Member handlers ────────────────────────────────────────────────────────
  async function handleAddMember(e: React.FormEvent) {
    e.preventDefault();
    if (!siteId || !accessToken || !addMemberId.trim()) return;
    setAddingMember(true);
    setMemberError("");
    try {
      const d = await apiFetch<{ success: boolean; data?: StoreMember }>(
        `/api/v1/sites/${siteId}/members`,
        { websiteId, accessToken, method: "POST", body: JSON.stringify({ user_id: addMemberId.trim(), role: addMemberRole }) }
      );
      if (d.success && d.data) {
        setMembers((prev) => [...prev.filter((m) => m.user_id !== d.data!.user_id), d.data!]);
        setAddMemberId("");
      }
    } catch (err: any) {
      setMemberError(err.message ?? "Erro ao adicionar membro");
    } finally {
      setAddingMember(false);
    }
  }

  async function handleUpdateRole(memberUserId: string, newRole: string) {
    if (!siteId || !accessToken) return;
    try {
      const d = await apiFetch<{ success: boolean; data?: StoreMember }>(
        `/api/v1/sites/${siteId}/members/${memberUserId}`,
        { websiteId, accessToken, method: "PATCH", body: JSON.stringify({ role: newRole }) }
      );
      if (d.success && d.data) {
        setMembers((prev) => prev.map((m) => m.user_id === memberUserId ? d.data! : m));
      }
    } catch (err: any) {
      setMemberError(err.message ?? "Erro ao atualizar role");
    }
  }

  async function handleRemoveMember(memberUserId: string) {
    if (!siteId || !accessToken) return;
    try {
      await apiFetch(`/api/v1/sites/${siteId}/members/${memberUserId}`, {
        websiteId, accessToken, method: "DELETE",
      });
      setMembers((prev) => prev.filter((m) => m.user_id !== memberUserId));
    } catch (err: any) {
      setMemberError(err.message ?? "Erro ao remover membro");
    }
  }

  // ── Handlers ──────────────────────────────────────────────────────────────
  async function handleRate(stars: number) {
    if (!user || !accessToken || !siteId) return;
    setRatingLoading(true);
    try {
      const d = await apiFetch<{ success: boolean; data?: { stars: number } }>(
        `/api/store/${siteId}/ratings`,
        { websiteId, accessToken, method: "POST", body: JSON.stringify({ stars }) }
      );
      if (d.success) {
        setMyRating(stars);
        fetch(`${API_URL}/api/store/${siteId}/info`).then((r) => r.json()).then((d) => { if (d.success) setStore(d.data); });
      }
    } catch { /* ignore */ }
    finally { setRatingLoading(false); }
  }

  async function handlePostComment(e: React.FormEvent) {
    e.preventDefault();
    if (!commentText.trim() || commentStars === 0 || !user || !accessToken || !siteId) return;
    setSubmittingComment(true);
    try {
      const d = await apiFetch<{ success: boolean; data?: Comment }>(
        `/api/store/${siteId}/comments`,
        { websiteId, accessToken, method: "POST", body: JSON.stringify({ content: commentText.trim(), stars: commentStars }) }
      );
      if (d.success && d.data) {
        setComments((prev) => [d.data!, ...prev]);
        setCommentText("");
        setCommentStars(0);
        // Refresh store to update rating avg
        fetch(`${API_URL}/api/store/${siteId}/info`).then((r) => r.json()).then((d) => { if (d.success) setStore(d.data); });
      }
    } catch { /* ignore */ }
    finally { setSubmittingComment(false); }
  }

  async function handleReply(e: React.FormEvent) {
    e.preventDefault();
    if (!replyText.trim() || !replyToId || !user || !accessToken || !siteId) return;
    setSubmittingReply(true);
    try {
      const d = await apiFetch<{ success: boolean; data?: Comment }>(
        `/api/store/${siteId}/comments/${replyToId}/replies`,
        { websiteId, accessToken, method: "POST", body: JSON.stringify({ content: replyText.trim() }) }
      );
      if (d.success && d.data) {
        setComments((prev) =>
          prev.map((c) =>
            c.id === replyToId
              ? { ...c, replies: [...(c.replies ?? []), d.data!] }
              : c
          )
        );
        setReplyToId(null);
        setReplyText("");
      }
    } catch { /* ignore */ }
    finally { setSubmittingReply(false); }
  }

  async function handleDeleteComment(commentId: string) {
    if (!user || !accessToken || !siteId) return;
    try {
      await apiFetch(`/api/store/${siteId}/comments/${commentId}`, {
        websiteId, accessToken, method: "DELETE",
      });
      setComments((prev) => prev.filter((c) => c.id !== commentId));
    } catch { /* ignore */ }
  }

  function startEdit() {
    if (!store) return;
    setEditName(store.name);
    setEditShort(store.short_description);
    setEditDesc(store.description);
    setEditMode(true);
  }

  async function handleSaveProfile(e: React.FormEvent) {
    e.preventDefault();
    if (!siteId || !accessToken) return;
    setEditSaving(true);
    try {
      const d = await apiFetch<{ success: boolean }>(
        `/api/v1/sites/${siteId}/profile`,
        {
          websiteId, accessToken, method: "PATCH",
          body: JSON.stringify({ name: editName, short_description: editShort, description: editDesc }),
        }
      );
      if (d.success) {
        setStore((s) => s ? { ...s, name: editName || s.name, short_description: editShort || s.short_description, description: editDesc || s.description } : s);
        setEditMode(false);
      }
    } catch { /* ignore */ }
    finally { setEditSaving(false); }
  }

  // ── Render ─────────────────────────────────────────────────────────────────
  if (loading) return <main className={styles.main}><div className={styles.loader} /></main>;
  if (error || !store) return (
    <main className={styles.main}>
      <p className={styles.errorMsg}>{error ?? "Loja não encontrada."}</p>
      <button className={styles.backBtn} onClick={() => navigate(-1)}>← Voltar</button>
    </main>
  );

  const editor = EDITOR_LABELS[store.editor_type];

  return (
    <main className={styles.main}>
      <div className={styles.layout}>

        {/* ── Left column ─────────────────────────────────────────────────── */}
        <aside className={styles.sidebar}>

          {/* Logo + name */}
          <div className={styles.storeCard}>
            <div className={styles.logoWrap}>
              {store.image ? (
                <img src={`data:image/png;base64,${store.image}`} alt={store.name} className={styles.logo} />
              ) : (
                <div className={styles.logoPlaceholder}>{store.name.charAt(0).toUpperCase()}</div>
              )}
              {store.mature_content && <span className={styles.matureBadge}>+18</span>}
            </div>

            {!editMode ? (
              <div className={styles.storeNames}>
                <h1 className={styles.storeName}>{store.name}</h1>
                {store.short_description && <p className={styles.storeShort}>{store.short_description}</p>}
              </div>
            ) : (
              <form onSubmit={handleSaveProfile} className={styles.editForm}>
                <input value={editName} onChange={(e) => setEditName(e.target.value)} placeholder="Nome" className={styles.editInput} />
                <input value={editShort} onChange={(e) => setEditShort(e.target.value)} placeholder="Descrição curta" className={styles.editInput} />
                <textarea value={editDesc} onChange={(e) => setEditDesc(e.target.value)} placeholder="Descrição completa" rows={4} className={styles.editTextarea} />
                <div className={styles.editActions}>
                  <button type="button" onClick={() => setEditMode(false)} className={styles.cancelBtn}>Cancelar</button>
                  <button type="submit" disabled={editSaving} className={styles.saveBtn}>{editSaving ? "Salvando…" : "Salvar"}</button>
                </div>
              </form>
            )}

            <div className={styles.badges}>
              <span className={styles.typeBadge}>{TYPE_LABELS[store.type] ?? store.type}</span>
              {editor && (
                <span className={styles.editorBadge} style={{ background: editor.color + "22", color: editor.color }}>
                  {editor.label}
                </span>
              )}
            </div>

            {/* Rating */}
            <div className={styles.ratingRow}>
              <StarRating value={Math.round(store.rating_avg)} />
              <span className={styles.ratingNum}>{store.rating_avg.toFixed(1)}</span>
              <span className={styles.ratingCount}>({store.rating_count})</span>
            </div>

            {/* Actions */}
            <div className={styles.actions}>
              <a href={`${API_URL}/p/${store.id}`} target="_blank" rel="noreferrer" className={styles.visitBtn}>
                ↗ Visitar loja
              </a>
              {isOwner && !editMode && (
                <button type="button" onClick={startEdit} className={styles.editBtn}>✏️ Editar</button>
              )}
              {!isOwner && canEdit && !editMode && (
                <button type="button" onClick={startEdit} className={styles.editBtn}>✏️ Editar perfil</button>
              )}
              {!isOwner && (
                <Link
                  to={`/report?website_id=${store.id}&website_name=${encodeURIComponent(store.name)}`}
                  className={styles.reportBtn}
                >
                  🚩 Denunciar
                </Link>
              )}
            </div>
          </div>

          {/* Creator card */}
          <div className={styles.creatorCard}>
            <h3>Criador</h3>
            <div className={styles.creatorRow}>
              {store.creator.avatar_url ? (
                <img src={store.creator.avatar_url} alt={store.creator.full_name} className={styles.creatorAvatar} />
              ) : (
                <div className={styles.creatorAvatarPlaceholder}>{store.creator.full_name.charAt(0)}</div>
              )}
              <div className={styles.creatorInfo}>
                <strong>{store.creator.full_name}</strong>
                {(store.creator.company_name || store.creator.trade_name) && (
                  <span>{store.creator.company_name ?? store.creator.trade_name}</span>
                )}
                <span className={styles.accountType}>
                  {store.creator.account_type === "business" ? "🏢 Empresa" : "👤 Pessoal"}
                </span>
              </div>
            </div>
          </div>

          {/* Rate store (if logged in and not owner) */}
          {user && !isOwner && (
            <div className={styles.rateCard}>
              <h3>Avaliar esta loja</h3>
              <p>Sua nota:</p>
              <StarRating value={myRating} interactive onChange={ratingLoading ? undefined : handleRate} />
              {myRating > 0 && <span className={styles.ratedLabel}>Você deu {myRating}★</span>}
            </div>
          )}
        </aside>

        {/* ── Right column ────────────────────────────────────────────────── */}
        <div className={styles.content}>

          {/* Description */}
          {store.description && (
            <section className={styles.section}>
              <h2>Sobre a loja</h2>
              <p className={styles.descText}>{store.description}</p>
            </section>
          )}

          {/* Store preview */}
          <section className={styles.section}>
            <h2>Prévia da loja</h2>
            <div className={styles.previewWrap}>
              <iframe
                src={`${API_URL}/p/${store.id}`}
                className={styles.previewFrame}
                title={`Prévia de ${store.name}`}
                sandbox="allow-scripts allow-same-origin"
              />
              <div className={styles.previewOverlay}>
                <a href={`${API_URL}/p/${store.id}`} target="_blank" rel="noreferrer" className={styles.previewOpenBtn}>
                  Abrir em tela cheia ↗
                </a>
              </div>
            </div>
          </section>

          {/* Visits chart */}
          <section className={styles.section}>
            <h2>Acessos (últimos 30 dias)</h2>
            {visits.length > 0 ? (
              <>
                <p className={styles.totalVisits}>
                  Total: <strong>{visits.reduce((s, d) => s + d.count, 0)}</strong> acessos
                </p>
                <VisitChart data={visits} />
              </>
            ) : (
              <p className={styles.emptyMsg}>Nenhum acesso registrado ainda.</p>
            )}
          </section>

          {/* Members Management (owner/manager/admin only) */}
          {canManageMembers && (
            <section className={styles.section}>
              <h2>👥 Membros da Equipe</h2>

              {memberError && <p className={styles.errorMsg}>{memberError}</p>}

              {members.length === 0 ? (
                <p className={styles.emptyMsg}>Nenhum membro adicionado ainda.</p>
              ) : (
                <table className={styles.membersTable}>
                  <thead>
                    <tr>
                      <th>Usuário</th>
                      <th>Role</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    {members.map((m) => (
                      <tr key={m.user_id}>
                        <td className={styles.memberName}>
                          {m.avatar_url
                            ? <img src={m.avatar_url} alt="" className={styles.commentAvatar} />
                            : <div className={styles.commentAvatarPlaceholder}>{m.user_name.charAt(0)}</div>
                          }
                          {m.user_name}
                        </td>
                        <td>
                          <select
                            className={styles.roleSelect}
                            value={m.role}
                            onChange={(e) => handleUpdateRole(m.user_id, e.target.value)}
                          >
                            <option value="manager">Gerente</option>
                            <option value="catalog_manager">Gestor de Catálogo</option>
                            <option value="support">Suporte</option>
                            <option value="logistics">Logística</option>
                          </select>
                        </td>
                        <td>
                          <button
                            type="button"
                            className={styles.deleteCommentBtn}
                            onClick={() => handleRemoveMember(m.user_id)}
                            title="Remover membro"
                          >✕</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}

              <form onSubmit={handleAddMember} className={styles.addMemberForm}>
                <h3>Adicionar membro</h3>
                <div className={styles.addMemberRow}>
                  <input
                    value={addMemberId}
                    onChange={(e) => setAddMemberId(e.target.value)}
                    placeholder="ID do usuário"
                    className={styles.editInput}
                    required
                  />
                  <select
                    className={styles.roleSelect}
                    value={addMemberRole}
                    onChange={(e) => setAddMemberRole(e.target.value)}
                  >
                    <option value="manager">Gerente</option>
                    <option value="catalog_manager">Gestor de Catálogo</option>
                    <option value="support">Suporte</option>
                    <option value="logistics">Logística</option>
                  </select>
                  <button type="submit" className={styles.saveBtn} disabled={addingMember}>
                    {addingMember ? "Adicionando…" : "Adicionar"}
                  </button>
                </div>
              </form>
            </section>
          )}

          {/* Comments */}
          <section className={styles.section}>
            <h2>Avaliações e comentários ({comments.length})</h2>

            {user ? (
              <form onSubmit={handlePostComment} className={styles.commentForm}>
                <div className={styles.commentStarsRow}>
                  <span className={styles.commentStarsLabel}>Sua nota: <em className={styles.required}>*</em></span>
                  <StarRating value={commentStars} interactive onChange={setCommentStars} />
                  {commentStars > 0 && <span className={styles.charCount}>{commentStars}★</span>}
                </div>
                <div className={styles.commentInputRow}>
                  {user.avatar_url ? (
                    <img src={user.avatar_url} alt="" className={styles.commentAvatar} />
                  ) : (
                    <div className={styles.commentAvatarPlaceholder}>{(user.first_name ?? "?").charAt(0)}</div>
                  )}
                  <textarea
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    placeholder="Escreva sua avaliação…"
                    rows={3}
                    maxLength={1000}
                    className={styles.commentTextarea}
                  />
                </div>
                <div className={styles.commentFormFooter}>
                  <span className={styles.charCount}>{commentText.length}/1000</span>
                  <button
                    type="submit"
                    disabled={submittingComment || commentText.trim().length < 3 || commentStars === 0}
                    className={styles.commentSubmitBtn}
                  >
                    {submittingComment ? "Enviando…" : "Publicar avaliação"}
                  </button>
                </div>
              </form>
            ) : (
              <p className={styles.loginPrompt}>
                <Link to="/login">Entre</Link> para deixar uma avaliação.
              </p>
            )}

            <div className={styles.commentsList}>
              {comments.length === 0 && <p className={styles.emptyMsg}>Nenhuma avaliação ainda. Seja o primeiro!</p>}
              {comments.map((c) => (
                <div key={c.id} className={styles.commentItem}>
                  <div className={styles.commentHeader}>
                    {c.avatar_url ? (
                      <img src={c.avatar_url} alt="" className={styles.commentAvatar} />
                    ) : (
                      <div className={styles.commentAvatarPlaceholder}>{c.user_name.charAt(0)}</div>
                    )}
                    <div className={styles.commentMeta}>
                      <strong>{c.user_name}</strong>
                      <div className={styles.commentMetaRow}>
                        {c.stars && <StarRating value={c.stars} />}
                        <span>{new Date(c.created_at).toLocaleDateString("pt-BR")}</span>
                      </div>
                    </div>
                    {(user?.id === c.user_id || user?.role === "admin" || user?.role === "manager") && (
                      <button
                        type="button"
                        className={styles.deleteCommentBtn}
                        onClick={() => handleDeleteComment(c.id)}
                        title="Excluir"
                      >✕</button>
                    )}
                  </div>
                  <p className={styles.commentContent}>{c.content}</p>

                  {/* Team reply button */}
                  {canReply && replyToId !== c.id && (
                    <button
                      type="button"
                      className={styles.replyBtn}
                      onClick={() => { setReplyToId(c.id); setReplyText(""); }}
                    >
                      💬 Responder
                    </button>
                  )}

                  {/* Inline reply form */}
                  {replyToId === c.id && (
                    <form onSubmit={handleReply} className={styles.replyForm}>
                      <textarea
                        value={replyText}
                        onChange={(e) => setReplyText(e.target.value)}
                        placeholder="Escreva sua resposta…"
                        rows={2}
                        maxLength={1000}
                        className={styles.replyTextarea}
                        autoFocus
                      />
                      <div className={styles.replyActions}>
                        <button type="button" className={styles.cancelBtn} onClick={() => setReplyToId(null)}>Cancelar</button>
                        <button
                          type="submit"
                          className={styles.saveBtn}
                          disabled={submittingReply || replyText.trim().length < 3}
                        >
                          {submittingReply ? "Enviando…" : "Responder"}
                        </button>
                      </div>
                    </form>
                  )}

                  {/* Replies */}
                  {(c.replies?.length ?? 0) > 0 && (
                    <div className={styles.repliesList}>
                      {c.replies!.map((reply) => (
                        <div key={reply.id} className={styles.replyItem}>
                          <div className={styles.commentHeader}>
                            {reply.avatar_url ? (
                              <img src={reply.avatar_url} alt="" className={styles.commentAvatar} />
                            ) : (
                              <div className={styles.commentAvatarPlaceholder}>{reply.user_name.charAt(0)}</div>
                            )}
                            <div className={styles.commentMeta}>
                              <span className={styles.replyBadge}>Equipe</span>
                              <strong>{reply.user_name}</strong>
                              <span>{new Date(reply.created_at).toLocaleDateString("pt-BR")}</span>
                            </div>
                          </div>
                          <p className={styles.commentContent}>{reply.content}</p>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </section>
        </div>
      </div>
    </main>
  );
}

