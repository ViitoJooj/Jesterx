import { BLOCK_CATEGORIES, ensurePath, getBlockIcon } from "./blockHelpers";
import styles from "./ElementorEditor.module.scss";

function BlockPreview({ type }) {
  const base = { width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center", background: "#f0f4fc", borderRadius: 5, overflow: "hidden" };
  switch (type) {
    case "heading":
      return (
        <div style={base}>
          <div style={{ display: "flex", flexDirection: "column", gap: 4, padding: "0 8px", width: "100%" }}>
            <div style={{ height: 7, background: "#2d3f62", borderRadius: 3, width: "88%" }} />
            <div style={{ height: 3.5, background: "#b8c6de", borderRadius: 3, width: "65%" }} />
          </div>
        </div>
      );
    case "paragraph":
      return (
        <div style={base}>
          <div style={{ display: "flex", flexDirection: "column", gap: 3, padding: "0 8px", width: "100%" }}>
            <div style={{ height: 3, background: "#8a96b0", borderRadius: 2, width: "95%" }} />
            <div style={{ height: 3, background: "#8a96b0", borderRadius: 2, width: "88%" }} />
            <div style={{ height: 3, background: "#8a96b0", borderRadius: 2, width: "72%" }} />
          </div>
        </div>
      );
    case "variable_text":
      return (
        <div style={base}>
          <span style={{ fontSize: ".62rem", color: "#3b6fd4", fontWeight: 700, fontFamily: "monospace", letterSpacing: "-.02em" }}>{"{ {var} }"}</span>
        </div>
      );
    case "button":
      return (
        <div style={base}>
          <div style={{ background: "#3b6fd4", color: "#fff", borderRadius: 5, padding: "4px 10px", fontSize: ".62rem", fontWeight: 700, letterSpacing: ".02em" }}>BOTÃO</div>
        </div>
      );
    case "image":
      return (
        <div style={{ ...base, background: "#e0e8f5" }}>
          <svg width="32" height="26" viewBox="0 0 24 18" fill="none">
            <rect x="1" y="1" width="22" height="16" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
            <circle cx="7" cy="6" r="2" fill="#8a96b0" />
            <path d="M1 13l5-4 4 3 4-4 5 5" stroke="#8a96b0" strokeWidth="1.2" strokeLinecap="round" />
          </svg>
        </div>
      );
    case "carousel":
      return (
        <div style={{ ...base, background: "#e0e8f5", position: "relative" }}>
          <svg width="36" height="26" viewBox="0 0 32 22" fill="none">
            <rect x="5" y="1" width="22" height="16" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
            <rect x="1" y="3" width="4" height="12" rx="1" fill="#c0cde0" />
            <rect x="27" y="3" width="4" height="12" rx="1" fill="#c0cde0" />
            <circle cx="12" cy="21" r="1.5" fill="#8a96b0" />
            <circle cx="16" cy="21" r="2" fill="#3b6fd4" />
            <circle cx="20" cy="21" r="1.5" fill="#8a96b0" />
          </svg>
        </div>
      );
    case "video":
      return (
        <div style={{ ...base, background: "#1a2740" }}>
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="9" stroke="rgba(255,255,255,.35)" strokeWidth="1.5" />
            <path d="M10 8.5l6 3.5-6 3.5V8.5z" fill="rgba(255,255,255,.75)" />
          </svg>
        </div>
      );
    case "input_var":
      return (
        <div style={base}>
          <div style={{ width: "82%", height: 16, border: "1.5px solid #b0bdd4", borderRadius: 4, background: "#fff", display: "flex", alignItems: "center", padding: "0 5px" }}>
            <div style={{ height: 2.5, width: "55%", background: "#c8d4e4", borderRadius: 2 }} />
            <div style={{ height: 8, width: 1, background: "#3b6fd4", marginLeft: "auto", animation: "none" }} />
          </div>
        </div>
      );
    case "profile_card":
      return (
        <div style={base}>
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 4 }}>
            <div style={{ width: 22, height: 22, borderRadius: "50%", background: "#8a96b0", display: "flex", alignItems: "center", justifyContent: "center" }}>
              <svg viewBox="0 0 20 20" width="14" height="14">
                <circle cx="10" cy="7" r="4" fill="rgba(255,255,255,.8)" />
                <ellipse cx="10" cy="18" rx="7" ry="5" fill="rgba(255,255,255,.8)" />
              </svg>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: 2.5, alignItems: "center" }}>
              <div style={{ height: 3.5, width: 28, background: "#2d3f62", borderRadius: 2 }} />
              <div style={{ height: 2.5, width: 20, background: "#b8c6de", borderRadius: 2 }} />
            </div>
          </div>
        </div>
      );
    case "product_card":
      return (
        <div style={{ ...base, flexDirection: "column", background: "#fff", border: "1px solid #e2e8f3", gap: 0, padding: 0 }}>
          <div style={{ width: "100%", flex: "0 0 54%", background: "#e0e8f5", display: "flex", alignItems: "center", justifyContent: "center" }}>
            <div style={{ width: 18, height: 14, background: "#c0cde0", borderRadius: 2 }} />
          </div>
          <div style={{ flex: 1, padding: "3px 5px", display: "flex", flexDirection: "column", gap: 2.5 }}>
            <div style={{ height: 3.5, background: "#b8c6de", borderRadius: 2, width: "80%" }} />
            <div style={{ height: 3, background: "#dde6f3", borderRadius: 2, width: "55%" }} />
            <div style={{ height: 5, background: "#a8bcd8", borderRadius: 2, width: "40%", marginTop: 1 }} />
          </div>
        </div>
      );
    case "product_list":
      return (
        <div style={{ ...base, gap: 2.5, padding: 5 }}>
          {[...Array(4)].map((_, i) => (
            <div key={i} style={{ flex: 1, height: "100%", background: "#fff", border: "1px solid #e2e8f3", borderRadius: 3, display: "flex", flexDirection: "column", overflow: "hidden" }}>
              <div style={{ flex: "0 0 52%", background: "#e0e8f5" }} />
              <div style={{ flex: 1, padding: "2px 3px", display: "flex", flexDirection: "column", gap: 1.5 }}>
                <div style={{ height: 2.5, background: "#b8c6de", borderRadius: 1 }} />
                <div style={{ height: 2, background: "#dde6f3", borderRadius: 1, width: "70%" }} />
              </div>
            </div>
          ))}
        </div>
      );
    case "user_avatar":
      return (
        <div style={{ ...base, background: "#e0e8f5", borderRadius: "50%", width: 46, height: 46, margin: "0 auto" }}>
          <svg viewBox="0 0 40 40" width="32" height="32">
            <circle cx="20" cy="14" r="9" fill="#8a96b0" />
            <ellipse cx="20" cy="38" rx="14" ry="10" fill="#8a96b0" />
          </svg>
        </div>
      );
    case "divider":
      return (
        <div style={base}>
          <div style={{ width: "80%", display: "flex", flexDirection: "column", gap: 4, alignItems: "center" }}>
            <div style={{ width: "100%", height: 2, background: "#8a96b0", borderRadius: 1 }} />
          </div>
        </div>
      );
    default:
      return <div style={base} />;
  }
}

export function LeftPanel({
  leftTab,
  setLeftTab,
  editingProductCard,
  currentBlocks,
  addBlock,
  selected,
  setSelected,
  removeBlock,
  editSection,
  setEditSection,
  setEditingPopup,
  setEditingProductCard,
  doc,
  setDoc,
  activeRoute,
  setActiveRoute,
  routes,
  setRoutes,
  removePage,
  addPage,
  editingPopup,
  applyTemplate,
  addPopup,
  removePopup,
}) {
  return (
    <aside className={styles.left}>
      <div className={styles.leftTabs}>
        <button type="button" className={leftTab === "elements" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("elements")} title="Elementos">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
        </button>
        <button type="button" className={leftTab === "layers" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("layers")} title="Camadas">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><polygon points="12 2 22 8.5 12 15 2 8.5"/><polyline points="2 15 12 21 22 15"/><polyline points="2 11.5 12 18 22 11.5"/></svg>
        </button>
        <button type="button" className={leftTab === "settings" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("settings")} title="Configuracoes">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
        </button>
        <button type="button" className={leftTab === "templates" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("templates")} title="Templates">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/></svg>
        </button>
        <button type="button" className={leftTab === "popups" ? styles.leftTabActive : styles.leftTabBtn} onClick={() => setLeftTab("popups")} title="Popups">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M8 12h8M12 8v8"/></svg>
        </button>
      </div>

      {leftTab === "elements" && (
        <div className={styles.blockList}>
          {editingProductCard ? (
            <div className={styles.varHint} style={{ padding: "6px 10px", margin: "0 0 6px", background: "#e0f2fe", borderRadius: 6, fontSize: ".75rem", color: "#0369a1" }}>
              🛍 <strong>Vars produto:</strong> {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{product_short_description}}"} {"{{product_sku}}"} {"{{product_brand}}"} {"{{product_model}}"} {"{{product_barcode}}"} {"{{product_weight}}"} {"{{product_color}}"} {"{{product_size}}"}
            </div>
          ) : (
            <div className={styles.varHint} style={{ padding: "6px 10px", margin: "0 0 6px", background: "#f1f5f9", borderRadius: 6, fontSize: ".75rem" }}>
              👤 <strong>Vars usuário:</strong> {"{{user_name}}"} {"{{user_display_name}}"} {"{{user_email}}"} {"{{user_phone}}"} {"{{user_address_city}}"} {"{{user_address_state}}"} {"{{cart_count}}"}
            </div>
          )}
          {BLOCK_CATEGORIES.map((cat) => (
            <div key={cat.label}>
              <div className={styles.catLabel}>{cat.label}</div>
              <div className={styles.catItems}>
                {cat.items.map((item) => (
                  <button key={item.type} type="button" className={styles.blockItem} onClick={() => addBlock(item.type)}>
                    <div className={styles.blockItemPreview}>
                      <BlockPreview type={item.type} />
                    </div>
                    <span>{item.label}</span>
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}

      {leftTab === "layers" && (
        <div className={styles.layersList}>
          <div className={styles.layersHeader}>Camadas · {editSection === "page" ? activeRoute : editSection}</div>
          {currentBlocks.length === 0 && <div className={styles.emptyLayers}>Sem blocos nesta secao</div>}
          {[...currentBlocks].sort((a, b) => (b.z || 0) - (a.z || 0)).map((block) => (
            <div
              key={block.id}
              className={`${styles.layerItem} ${selected === block.id ? styles.layerSelected : ""}`}
              onClick={() => setSelected(block.id)}
            >
              <span className={styles.layerIcon}>{getBlockIcon(block.type)}</span>
              <span className={styles.layerName}>{block.text || block.label || block.profile_name || block.type}</span>
              <button type="button" className={styles.layerDelete} onClick={(e) => { e.stopPropagation(); removeBlock(block.id); }}>🗑</button>
            </div>
          ))}
        </div>
      )}

      {leftTab === "settings" && (
        <div className={styles.settingsPanel}>
          <div className={styles.settingsGroup}>
            <div className={styles.settingsGroupTitle}>Pagina</div>
            <select value={activeRoute} onChange={(e) => { setActiveRoute(e.target.value); }}>
              {routes.map((r) => (
                <option key={r.id} value={ensurePath(r.path)}>{r.title} ({ensurePath(r.path)})</option>
              ))}
            </select>
            <label>Titulo</label>
            <input
              value={doc.pages[activeRoute]?.title || ""}
              onChange={(e) => {
                const newTitle = e.target.value;
                setDoc((prev) => ({ ...prev, pages: { ...prev.pages, [activeRoute]: { ...prev.pages[activeRoute], title: newTitle } } }));
                setRoutes((prev) => prev.map((r) => ensurePath(r.path) === activeRoute ? { ...r, title: newTitle } : r));
              }}
            />
          </div>

          <div className={styles.settingsGroup}>
            <div className={styles.settingsGroupTitle}>Canvas</div>
            <div className={styles.colorPickerRow}>
              <label>Fundo</label>
              <input type="color" value={doc.canvas.background || "#f8f9ff"} onChange={(e) => setDoc((prev) => ({ ...prev, canvas: { ...prev.canvas, background: e.target.value } }))} />
            </div>
            <label>Largura</label>
            <input type="number" value={doc.canvas.width} min={900} max={2800} onChange={(e) => setDoc((prev) => ({ ...prev, canvas: { ...prev.canvas, width: Math.max(900, Math.min(2800, Number(e.target.value) || 900)) } }))} />
            <label>Altura</label>
            <input type="number" value={doc.canvas.height} min={700} max={2800} onChange={(e) => setDoc((prev) => ({ ...prev, canvas: { ...prev.canvas, height: Math.max(700, Math.min(2800, Number(e.target.value) || 700)) } }))} />
          </div>

          <div className={styles.settingsGroup}>
            <div className={styles.settingsGroupTitle}>Header Global</div>
            <div className={styles.toggleRow}>
              <label>Ativo</label>
              <input type="checkbox" checked={doc.header.enabled} onChange={(e) => setDoc((prev) => ({ ...prev, header: { ...prev.header, enabled: e.target.checked } }))} />
            </div>
            <label>Altura</label>
            <input type="number" value={doc.header.height} min={40} max={400} onChange={(e) => setDoc((prev) => ({ ...prev, header: { ...prev.header, height: Math.max(40, Number(e.target.value) || 80) } }))} />
            <div className={styles.colorPickerRow}>
              <label>Fundo</label>
              <input type="color" value={doc.header.background || "#1a2740"} onChange={(e) => setDoc((prev) => ({ ...prev, header: { ...prev.header, background: e.target.value } }))} />
            </div>
            <button type="button" className={editSection === "header" ? styles.editSectionBtnActive : styles.editSectionBtn} onClick={() => { setEditSection("header"); }}>
              ✏️ Editar Header
            </button>
          </div>

          <div className={styles.settingsGroup}>
            <div className={styles.settingsGroupTitle}>Footer Global</div>
            <div className={styles.toggleRow}>
              <label>Ativo</label>
              <input type="checkbox" checked={doc.footer.enabled} onChange={(e) => setDoc((prev) => ({ ...prev, footer: { ...prev.footer, enabled: e.target.checked } }))} />
            </div>
            <label>Altura</label>
            <input type="number" value={doc.footer.height} min={40} max={400} onChange={(e) => setDoc((prev) => ({ ...prev, footer: { ...prev.footer, height: Math.max(40, Number(e.target.value) || 100) } }))} />
            <div className={styles.colorPickerRow}>
              <label>Fundo</label>
              <input type="color" value={doc.footer.background || "#1a2740"} onChange={(e) => setDoc((prev) => ({ ...prev, footer: { ...prev.footer, background: e.target.value } }))} />
            </div>
            <button type="button" className={editSection === "footer" ? styles.editSectionBtnActive : styles.editSectionBtn} onClick={() => { setEditSection("footer"); }}>
              ✏️ Editar Footer
            </button>
          </div>

          <div className={styles.settingsGroup}>
            <div className={styles.settingsGroupTitle}>Paginas</div>
            {routes.map((route) => (
              <div key={route.id} className={styles.pageRow}>
                <span style={{ fontSize: ".78rem", flex: 1 }}>{ensurePath(route.path)}</span>
                <button type="button" className={styles.removePageBtn} disabled={routes.length <= 1} onClick={() => removePage(route.path)}>✕</button>
              </div>
            ))}
            <button type="button" className={styles.addPageBtn} onClick={addPage}>+ Nova Pagina</button>
          </div>
        </div>
      )}

      {leftTab === "templates" && (
        <div className={styles.templatesPanel}>
          <div className={styles.templateCardSmall} onClick={() => applyTemplate("landing")}>
            <span className={styles.templateEmoji}>🚀</span>
            <div>
              <div style={{ fontWeight: 600, fontSize: ".82rem" }}>Landing Page</div>
              <div style={{ fontSize: ".72rem", color: "#8a96b0", marginTop: 1 }}>Header, hero, cards, footer</div>
            </div>
          </div>
          <div className={styles.templateCardSmall} onClick={() => applyTemplate("blank")}>
            <span className={styles.templateEmoji}>📄</span>
            <div>
              <div style={{ fontWeight: 600, fontSize: ".82rem" }}>Em branco</div>
              <div style={{ fontSize: ".72rem", color: "#8a96b0", marginTop: 1 }}>Canvas vazio</div>
            </div>
          </div>
        </div>
      )}

      {leftTab === "popups" && (
        <div className={styles.layersList}>
          <div className={styles.layersHeader}>Popups</div>
          {Object.entries(doc.popups).map(([id, popup]) => (
            <div
              key={id}
              className={`${styles.layerItem} ${editingPopup === id ? styles.layerSelected : ""}`}
              onClick={() => { setEditingPopup(id); setEditingProductCard(null); setEditSection("page"); }}
            >
              <span className={styles.layerIcon}>📌</span>
              <span className={styles.layerName}>{popup.title}</span>
              <button type="button" className={styles.layerDelete} onClick={(e) => { e.stopPropagation(); removePopup(id); }}>🗑</button>
            </div>
          ))}
          {Object.keys(doc.popups).length === 0 && <div className={styles.emptyLayers}>Nenhum popup criado</div>}
          <button type="button" className={styles.addPageBtn} onClick={addPopup}>+ Novo Popup</button>
        </div>
      )}
    </aside>
  );
}
