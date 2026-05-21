import { useState } from "react";
import { API_URL } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import { readFileAsDataURL } from "./blockHelpers";
import styles from "./ElementorEditor.module.scss";

const FONT_FAMILIES = [
  { value: "inherit", label: "Herdar" },
  { value: "'Inter', sans-serif", label: "Inter" },
  { value: "'Roboto', sans-serif", label: "Roboto" },
  { value: "'Open Sans', sans-serif", label: "Open Sans" },
  { value: "'Poppins', sans-serif", label: "Poppins" },
  { value: "'Lato', sans-serif", label: "Lato" },
  { value: "'Montserrat', sans-serif", label: "Montserrat" },
  { value: "'Nunito', sans-serif", label: "Nunito" },
  { value: "'Playfair Display', serif", label: "Playfair Display" },
  { value: "'Merriweather', serif", label: "Merriweather" },
  { value: "monospace", label: "Monospace" },
  { value: "serif", label: "Serif" },
];

export function Inspector({
  selectedBlock,
  editSection,
  updateBlock,
  removeBlock,
  routes,
  apis,
  doc,
  editingPopup,
  setEditingProductCard,
  setEditingProductCardSection,
  setError,
}) {
  const { websiteId } = useAuthContext();
  const [apiTesting, setApiTesting] = useState(false);
  const [apiTestResult, setApiTestResult] = useState(null);

  async function testApi(api) {
    setApiTesting(true);
    setApiTestResult(null);
    try {
      const resp = await fetch(`${API_URL}${api.path}`, {
        method: api.method,
        credentials: "include",
        headers: { "X-Website-Id": websiteId || "" },
      });
      const text = await resp.text();
      setApiTestResult(`${resp.status} ${resp.statusText}\n${text.slice(0, 400)}`);
    } catch (err) {
      setApiTestResult(err instanceof Error ? err.message : "Erro na chamada");
    } finally {
      setApiTesting(false);
    }
  }

  function update(patch) {
    updateBlock(editSection, selectedBlock.id, patch);
  }

  function updateStyle(patch) {
    update({ style: { ...selectedBlock.style, ...patch } });
  }

  const isTypographyBlock = selectedBlock && ["heading", "paragraph", "variable_text", "button"].includes(selectedBlock.type);

  return (
    <aside className={styles.right}>
      <h3 style={{ margin: "0 0 .4rem", fontSize: ".8rem", textTransform: "uppercase", letterSpacing: ".04em", color: "#4b5774" }}>
        Inspector
      </h3>

      {!selectedBlock && (
        <p className={styles.emptyInspector}>Clique em um bloco no canvas para editar.</p>
      )}

      {selectedBlock && (
        <div className={styles.form}>

          {/* CONTEUDO */}
          <div className={styles.inspectorSection}>
            <div className={styles.sectionLabelTag}>CONTEUDO · {selectedBlock.type.toUpperCase()}</div>

            {(selectedBlock.type === "heading" || selectedBlock.type === "paragraph" || selectedBlock.type === "variable_text") && (
              <>
                {selectedBlock.type === "heading" && (
                  <>
                    <label>Nível</label>
                    <select value={selectedBlock.heading_level || "h1"} onChange={(e) => update({ heading_level: e.target.value })}>
                      <option value="h1">H1 – Título Principal</option>
                      <option value="h2">H2 – Título Secundário</option>
                      <option value="h3">H3 – Subtítulo</option>
                      <option value="h4">H4 – Pequeno Título</option>
                    </select>
                  </>
                )}
                <label>Texto</label>
                <textarea value={selectedBlock.text || ""} rows={3} onChange={(e) => update({ text: e.target.value })} />
              </>
            )}

            {selectedBlock.type === "button" && (
              <>
                <label>Texto do botao</label>
                <input value={selectedBlock.label || ""} onChange={(e) => update({ label: e.target.value })} />
                <label>Acao</label>
                <select value={selectedBlock.action_type || "navigate"} onChange={(e) => update({ action_type: e.target.value })}>
                  <option value="navigate">Navegar para rota</option>
                  <option value="call_api">Chamar API</option>
                  <option value="store_login">Login na loja</option>
                  <option value="store_logout">Logout da loja</option>
                  <option value="store_register">Registrar na loja</option>
                  <option value="add_to_cart">Adicionar ao carrinho</option>
                  <option value="add_product">Adicionar produto (admin)</option>
                  <option value="open_popup">Abrir Popup</option>
                  <option value="close_popup">Fechar Popup</option>
                </select>

                {selectedBlock.action_type === "call_api" && (
                  <>
                    <label>API</label>
                    <select
                      value={selectedBlock.api_id || ""}
                      onChange={(e) => {
                        const api = apis.find((a) => a.id === e.target.value);
                        update({ api_id: e.target.value, action_target: api?.path || "", href: "#" });
                      }}
                    >
                      <option value="">— Selecione API —</option>
                      {apis.map((api) => (
                        <option key={api.id} value={api.id}>{api.method} {api.path} — {api.label}</option>
                      ))}
                    </select>
                    {selectedBlock.api_id && (() => {
                      const api = apis.find((a) => a.id === selectedBlock.api_id);
                      if (!api) return null;
                      return (
                        <div className={styles.apiTestBox}>
                          <p className={styles.apiDesc}>{api.description}</p>
                          <button type="button" className={styles.apiTestBtn} disabled={apiTesting} onClick={() => testApi(api)}>
                            {apiTesting ? "Testando..." : "▶ Testar API"}
                          </button>
                          {apiTestResult && <pre className={styles.apiResult}>{apiTestResult}</pre>}
                        </div>
                      );
                    })()}
                  </>
                )}

                {selectedBlock.action_type === "add_to_cart" && (
                  <>
                    <label>ID do produto</label>
                    <input value={selectedBlock.product_id || ""} onChange={(e) => update({ product_id: e.target.value })} placeholder="UUID do produto" />
                    <p className={styles.varHint}>Deixe vazio para usar o 1º produto da loja.</p>
                  </>
                )}

                {selectedBlock.action_type === "open_popup" && (
                  <>
                    <label>Popup</label>
                    <select value={selectedBlock.popup_id || ""} onChange={(e) => update({ popup_id: e.target.value })}>
                      <option value="">— Selecione Popup —</option>
                      {Object.entries(doc.popups).map(([id, popup]) => (
                        <option key={id} value={id}>{popup.title}</option>
                      ))}
                    </select>
                  </>
                )}

                {(selectedBlock.action_type === "store_login" || selectedBlock.action_type === "store_register") && (
                  <>
                    <label>Var email</label>
                    <input value={selectedBlock.email_var || ""} onChange={(e) => update({ email_var: e.target.value })} placeholder="variavel de email" />
                    <label>Var senha</label>
                    <input value={selectedBlock.password_var || ""} onChange={(e) => update({ password_var: e.target.value })} placeholder="variavel de senha" />
                    {selectedBlock.action_type === "store_register" && (
                      <>
                        <label>Var nome</label>
                        <input value={selectedBlock.first_name_var || ""} onChange={(e) => update({ first_name_var: e.target.value })} placeholder="variavel de nome" />
                        <label>Var sobrenome</label>
                        <input value={selectedBlock.last_name_var || ""} onChange={(e) => update({ last_name_var: e.target.value })} placeholder="variavel de sobrenome" />
                      </>
                    )}
                    <p className={styles.varHint}>Deixe vazio para abrir o modal de login/registro.</p>
                  </>
                )}

                {(!selectedBlock.action_type || selectedBlock.action_type === "navigate") && (
                  <>
                    <label>Rota destino</label>
                    <select
                      value={selectedBlock.action_target || "/"}
                      onChange={(e) => update({ action_target: e.target.value, href: e.target.value })}
                    >
                      {routes.map((r) => {
                        const path = r.path.trim().startsWith("/") ? r.path.trim() : `/${r.path.trim()}`;
                        return <option key={r.id} value={path}>{r.title} ({path})</option>;
                      })}
                    </select>
                  </>
                )}

                <div className={styles.toggleRow} style={{ marginTop: 8 }}>
                  <label>Visivel somente para admins</label>
                  <input type="checkbox" checked={!!selectedBlock.admin_only} onChange={(e) => update({ admin_only: e.target.checked })} />
                </div>
                {selectedBlock.admin_only && (
                  <p className={styles.varHint}>Oculto para visitantes; visivel apenas para o dono da loja.</p>
                )}
              </>
            )}

            {selectedBlock.type === "product_card" && (
              <>
                <label>ID do produto</label>
                <input value={selectedBlock.product_id || ""} onChange={(e) => update({ product_id: e.target.value })} placeholder="UUID do produto (opcional)" />
                <p className={styles.varHint}>Deixe vazio para exibir o 1º produto ativo da loja.</p>
                {!editingPopup && (
                  <button
                    type="button"
                    className={styles.editSectionBtn}
                    style={{ marginTop: 6 }}
                    onClick={() => { setEditingProductCard(selectedBlock.id); setEditingProductCardSection(editSection); }}
                  >
                    ✏️ Editar Layout Interno ({selectedBlock.inner_blocks?.length || 0} blocos)
                  </button>
                )}
                {(selectedBlock.inner_blocks?.length || 0) > 0 && (
                  <button
                    type="button"
                    className={styles.apiTestBtn}
                    style={{ marginTop: 4, background: "#fee2e2", color: "#b91c1c", border: "1px solid #fecaca" }}
                    onClick={() => update({ inner_blocks: [] })}
                  >
                    🗑 Limpar layout
                  </button>
                )}
                <p className={styles.varHint}>Vars: {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{product_sku}}"}</p>
              </>
            )}

            {selectedBlock.type === "product_list" && (
              <>
                <label>Produtos por pagina</label>
                <input
                  type="number" min={1} max={50}
                  value={selectedBlock.page_size || 6}
                  onChange={(e) => update({ page_size: Math.max(1, Number(e.target.value) || 6) })}
                />
                <p className={styles.varHint}>Ordenados por mais vendidos. Botoes de pagina aparecem automaticamente.</p>
              </>
            )}

            {selectedBlock.type === "image" && (
              <>
                <label>URL da imagem</label>
                <input value={selectedBlock.src || ""} onChange={(e) => update({ src: e.target.value })} placeholder="https://..." />
                <label>Var source (opcional)</label>
                <input value={selectedBlock.var_src || ""} onChange={(e) => update({ var_src: e.target.value })} placeholder="ex: {{product_image}}" />
                <label>Upload</label>
                <input
                  type="file" accept="image/*"
                  onChange={async (e) => {
                    const file = e.target.files?.[0];
                    if (!file) return;
                    try {
                      update({ src: await readFileAsDataURL(file) });
                    } catch {
                      setError("Falha ao carregar imagem");
                    }
                  }}
                />
                <label>Ajuste</label>
                <select value={selectedBlock.object_fit || "cover"} onChange={(e) => update({ object_fit: e.target.value })}>
                  <option value="cover">Cover</option>
                  <option value="contain">Contain</option>
                  <option value="fill">Fill</option>
                  <option value="none">Nenhum</option>
                </select>
              </>
            )}

            {selectedBlock.type === "user_avatar" && (
              <>
                <p className={styles.varHint}>Exibe o avatar do usuário logado. Atualiza após login/logout.</p>
                <label>Ajuste</label>
                <select value={selectedBlock.object_fit || "cover"} onChange={(e) => update({ object_fit: e.target.value })}>
                  <option value="cover">Cover</option>
                  <option value="contain">Contain</option>
                  <option value="fill">Fill</option>
                </select>
              </>
            )}

            {selectedBlock.type === "carousel" && (
              <>
                <label>Imagens (uma URL por linha)</label>
                <textarea
                  value={(selectedBlock.images || []).join("\n")} rows={4}
                  onChange={(e) => update({ images: e.target.value.split("\n").map((i) => i.trim()).filter(Boolean) })}
                />
                <div className={styles.toggleRow}>
                  <label>Autoplay</label>
                  <input type="checkbox" checked={!!selectedBlock.autoplay} onChange={(e) => update({ autoplay: e.target.checked })} />
                </div>
                {selectedBlock.autoplay && (
                  <>
                    <label>Intervalo (ms)</label>
                    <input type="number" min={500} step={100} value={selectedBlock.autoplay_interval || 3000} onChange={(e) => update({ autoplay_interval: Math.max(500, Number(e.target.value) || 3000) })} />
                  </>
                )}
                <div className={styles.toggleRow}>
                  <label>Loop infinito</label>
                  <input type="checkbox" checked={selectedBlock.loop !== false} onChange={(e) => update({ loop: e.target.checked })} />
                </div>
                <div className={styles.toggleRow}>
                  <label>Mostrar dots</label>
                  <input type="checkbox" checked={selectedBlock.show_dots !== false} onChange={(e) => update({ show_dots: e.target.checked })} />
                </div>
                <div className={styles.toggleRow}>
                  <label>Mostrar setas</label>
                  <input type="checkbox" checked={selectedBlock.show_arrows !== false} onChange={(e) => update({ show_arrows: e.target.checked })} />
                </div>
              </>
            )}

            {selectedBlock.type === "video" && (
              <>
                <label>URL do video</label>
                <input value={selectedBlock.video_url || ""} onChange={(e) => update({ video_url: e.target.value })} placeholder="https://youtube.com/watch?v=..." />
                <div className={styles.toggleRow}>
                  <label>Autoplay</label>
                  <input type="checkbox" checked={!!selectedBlock.autoplay} onChange={(e) => update({ autoplay: e.target.checked })} />
                </div>
                <div className={styles.toggleRow}>
                  <label>Loop</label>
                  <input type="checkbox" checked={!!selectedBlock.loop} onChange={(e) => update({ loop: e.target.checked })} />
                </div>
                <div className={styles.toggleRow}>
                  <label>Mudo</label>
                  <input type="checkbox" checked={!!selectedBlock.muted} onChange={(e) => update({ muted: e.target.checked })} />
                </div>
                {selectedBlock.autoplay && <p className={styles.varHint}>Autoplay requer mudo ativo na maioria dos navegadores.</p>}
              </>
            )}

            {selectedBlock.type === "input_var" && (
              <>
                <label>Nome da variavel</label>
                <input value={selectedBlock.var_name || ""} onChange={(e) => update({ var_name: e.target.value })} placeholder="ex: nome" />
                <p className={styles.varHint}>Use {"{{" + (selectedBlock.var_name || "nome") + "}}"} em blocos Texto Variavel.</p>
                <label>Label (texto acima)</label>
                <input value={selectedBlock.label_text || ""} onChange={(e) => update({ label_text: e.target.value })} placeholder="ex: Seu nome" />
                <label>Placeholder</label>
                <input value={selectedBlock.placeholder || ""} onChange={(e) => update({ placeholder: e.target.value })} />
                <label>Tipo do campo</label>
                <select value={selectedBlock.input_type || "text"} onChange={(e) => update({ input_type: e.target.value })}>
                  <option value="text">Texto</option>
                  <option value="email">Email</option>
                  <option value="password">Senha (oculta)</option>
                  <option value="number">Numero</option>
                  <option value="tel">Telefone</option>
                  <option value="date">Data</option>
                </select>
                <div className={styles.toggleRow}>
                  <label>Campo obrigatório</label>
                  <input type="checkbox" checked={!!selectedBlock.required} onChange={(e) => update({ required: e.target.checked })} />
                </div>
                <label>Máximo de caracteres</label>
                <input type="number" min={1} value={selectedBlock.max_length || ""} onChange={(e) => update({ max_length: e.target.value ? Number(e.target.value) : undefined })} placeholder="sem limite" />
              </>
            )}

            {selectedBlock.type === "profile_card" && (
              <>
                <label>Nome</label>
                <input value={selectedBlock.profile_name || ""} onChange={(e) => update({ profile_name: e.target.value })} />
                <label>Subtitulo</label>
                <input value={selectedBlock.profile_subtitle || ""} onChange={(e) => update({ profile_subtitle: e.target.value })} />
                <label>Bio (opcional)</label>
                <textarea value={selectedBlock.profile_bio || ""} rows={2} onChange={(e) => update({ profile_bio: e.target.value })} placeholder="Breve descrição..." />
                <label>URL da imagem</label>
                <input value={selectedBlock.profile_image || ""} onChange={(e) => update({ profile_image: e.target.value })} placeholder="https://..." />
                <label>Upload imagem</label>
                <input
                  type="file" accept="image/*"
                  onChange={async (e) => {
                    const file = e.target.files?.[0];
                    if (!file) return;
                    try {
                      update({ profile_image: await readFileAsDataURL(file) });
                    } catch {
                      setError("Falha ao carregar imagem");
                    }
                  }}
                />
                <label>Alinhamento</label>
                <select value={selectedBlock.align || "center"} onChange={(e) => update({ align: e.target.value })}>
                  <option value="left">Esquerda</option>
                  <option value="center">Centro</option>
                  <option value="right">Direita</option>
                </select>
                <label>Formato do avatar</label>
                <select value={selectedBlock.avatar_shape || "circle"} onChange={(e) => update({ avatar_shape: e.target.value })}>
                  <option value="circle">Circulo</option>
                  <option value="square">Quadrado arredondado</option>
                </select>
              </>
            )}
          </div>

          {/* TIPOGRAFIA */}
          {isTypographyBlock && (
            <div className={styles.inspectorSection}>
              <div className={styles.sectionLabelTag}>TIPOGRAFIA</div>
              <label>Família</label>
              <select value={selectedBlock.style["font-family"] || "inherit"} onChange={(e) => updateStyle({ "font-family": e.target.value })}>
                {FONT_FAMILIES.map((f) => (
                  <option key={f.value} value={f.value}>{f.label}</option>
                ))}
              </select>
              <div className={styles.inlineFields}>
                <span>
                  <label>Tamanho</label>
                  <input value={selectedBlock.style["font-size"] || ""} onChange={(e) => updateStyle({ "font-size": e.target.value })} placeholder="16px" />
                </span>
                <span>
                  <label>Peso</label>
                  <select value={selectedBlock.style["font-weight"] || "400"} onChange={(e) => updateStyle({ "font-weight": e.target.value })}>
                    <option value="300">Light</option>
                    <option value="400">Normal</option>
                    <option value="500">Medium</option>
                    <option value="600">SemiBold</option>
                    <option value="700">Bold</option>
                    <option value="800">ExtraBold</option>
                    <option value="900">Black</option>
                  </select>
                </span>
              </div>
              <div className={styles.inlineFields}>
                <span>
                  <label>Altura linha</label>
                  <input value={selectedBlock.style["line-height"] || ""} onChange={(e) => updateStyle({ "line-height": e.target.value })} placeholder="1.5" />
                </span>
                <span>
                  <label>Espaç. letras</label>
                  <input value={selectedBlock.style["letter-spacing"] || ""} onChange={(e) => updateStyle({ "letter-spacing": e.target.value })} placeholder="0" />
                </span>
              </div>
              <div className={styles.inlineFields}>
                <span>
                  <label>Transformar</label>
                  <select value={selectedBlock.style["text-transform"] || "none"} onChange={(e) => updateStyle({ "text-transform": e.target.value })}>
                    <option value="none">Normal</option>
                    <option value="uppercase">MAIUSCULO</option>
                    <option value="lowercase">minusculo</option>
                    <option value="capitalize">Capitalizar</option>
                  </select>
                </span>
                <span>
                  <label>Decoração</label>
                  <select value={selectedBlock.style["text-decoration"] || "none"} onChange={(e) => updateStyle({ "text-decoration": e.target.value })}>
                    <option value="none">Nenhum</option>
                    <option value="underline">Sublinhado</option>
                    <option value="line-through">Riscado</option>
                  </select>
                </span>
              </div>
            </div>
          )}

          {/* POSICAO & TAMANHO */}
          <div className={styles.inspectorSection}>
            <div className={styles.sectionLabelTag}>POSICAO &amp; TAMANHO</div>
            <div className={styles.inlineFields}>
              <span><label>X</label><input type="number" value={selectedBlock.x} onChange={(e) => update({ x: Number(e.target.value) || 0 })} /></span>
              <span><label>Y</label><input type="number" value={selectedBlock.y} onChange={(e) => update({ y: Number(e.target.value) || 0 })} /></span>
              <span><label>W</label><input type="number" value={selectedBlock.w} onChange={(e) => update({ w: Number(e.target.value) || 0 })} /></span>
              <span><label>H</label><input type="number" value={selectedBlock.h} onChange={(e) => update({ h: Number(e.target.value) || 0 })} /></span>
            </div>
            <div className={styles.inlineFields}>
              <span><label>Rot°</label><input type="number" value={selectedBlock.rotation} onChange={(e) => update({ rotation: Number(e.target.value) || 0 })} /></span>
              <span><label>Z-index</label><input type="number" value={selectedBlock.z} onChange={(e) => update({ z: Number(e.target.value) || 1 })} /></span>
            </div>
          </div>

          {/* ESTILO */}
          <div className={styles.inspectorSection}>
            <div className={styles.sectionLabelTag}>ESTILO</div>
            <div className={styles.colorRow}>
              <span>
                <label>Texto</label>
                <input
                  type="color"
                  value={selectedBlock.style.color?.startsWith("#") ? selectedBlock.style.color : "#1f2b43"}
                  onChange={(e) => updateStyle({ color: e.target.value })}
                />
              </span>
              <span>
                <label>Fundo</label>
                <input
                  type="color"
                  value={selectedBlock.style.background?.startsWith("#") ? selectedBlock.style.background : "#ffffff"}
                  onChange={(e) => updateStyle({ background: e.target.value })}
                />
              </span>
            </div>
            <label>Cor texto (hex/rgba/var)</label>
            <input value={selectedBlock.style.color || ""} onChange={(e) => updateStyle({ color: e.target.value })} placeholder="#1f2b43" />
            <label>Cor fundo</label>
            <input value={selectedBlock.style.background || ""} onChange={(e) => updateStyle({ background: e.target.value })} placeholder="transparent" />
            <label>Alinhamento texto</label>
            <select value={selectedBlock.style["text-align"] || "left"} onChange={(e) => updateStyle({ "text-align": e.target.value })}>
              <option value="left">Esquerda</option>
              <option value="center">Centro</option>
              <option value="right">Direita</option>
              <option value="justify">Justificado</option>
            </select>
            <label>Padding</label>
            <input value={selectedBlock.style.padding || ""} onChange={(e) => updateStyle({ padding: e.target.value })} placeholder="10px 12px" />
            <label>Margin</label>
            <input value={selectedBlock.style.margin || ""} onChange={(e) => updateStyle({ margin: e.target.value })} placeholder="0" />
            <label>Borda</label>
            <input value={selectedBlock.style.border || ""} onChange={(e) => updateStyle({ border: e.target.value })} placeholder="1px solid #ccc" />
            <label>Raio borda</label>
            <input value={selectedBlock.style["border-radius"] || ""} onChange={(e) => updateStyle({ "border-radius": e.target.value })} placeholder="8px" />
            <label>Sombra (box-shadow)</label>
            <input value={selectedBlock.style["box-shadow"] || ""} onChange={(e) => updateStyle({ "box-shadow": e.target.value })} placeholder="0 4px 12px rgba(0,0,0,.15)" />
            <label>Opacidade (0–1)</label>
            <input
              type="number" min={0} max={1} step={0.05}
              value={selectedBlock.style.opacity ?? "1"}
              onChange={(e) => updateStyle({ opacity: e.target.value })}
            />
          </div>

          <button type="button" className={styles.remove} onClick={() => removeBlock(selectedBlock.id)}>
            🗑 Excluir bloco
          </button>
        </div>
      )}
    </aside>
  );
}
