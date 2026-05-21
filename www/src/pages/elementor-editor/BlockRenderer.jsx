import styles from "./ElementorEditor.module.scss";

function textStyle(s = {}) {
  return {
    fontSize: s["font-size"],
    fontWeight: s["font-weight"],
    fontFamily: s["font-family"],
    color: s.color,
    lineHeight: s["line-height"],
    letterSpacing: s["letter-spacing"],
    textAlign: s["text-align"],
    textTransform: s["text-transform"],
    textDecoration: s["text-decoration"],
    margin: 0,
    overflow: "hidden",
  };
}

export function renderBlock(block) {
  const s = block.style || {};

  if (block.type === "heading") {
    const Tag = block.heading_level || "h1";
    return (
      <Tag style={{ ...textStyle(s), lineHeight: s["line-height"] || 1.15 }}>
        {block.text}
      </Tag>
    );
  }

  if (block.type === "paragraph") {
    return (
      <p style={{ ...textStyle(s), lineHeight: s["line-height"] || 1.55 }}>
        {block.text}
      </p>
    );
  }

  if (block.type === "variable_text") {
    return (
      <p style={{ ...textStyle(s), lineHeight: s["line-height"] || 1.4 }}>
        {block.text || "Texto dinamico {{var}}"}
      </p>
    );
  }

  if (block.type === "button") {
    const justify = s["text-align"] === "left" ? "flex-start" : s["text-align"] === "right" ? "flex-end" : "center";
    return (
      <div style={{ width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: justify }}>
        <span style={{
          display: "inline-block",
          background: s.background || "var(--jx-color-primary)",
          color: s.color || "#fff",
          border: s.border !== "0" ? s.border : "none",
          borderRadius: s["border-radius"] || "8px",
          padding: s.padding || "12px 24px",
          fontSize: s["font-size"] || "16px",
          fontWeight: s["font-weight"] || "600",
          fontFamily: s["font-family"] || "inherit",
          letterSpacing: s["letter-spacing"] || "0",
          textTransform: s["text-transform"] || "none",
          lineHeight: s["line-height"] || "1.4",
          boxShadow: s["box-shadow"] && s["box-shadow"] !== "none" ? s["box-shadow"] : undefined,
          cursor: "default",
          whiteSpace: "nowrap",
        }}>
          {block.label || "Botao"}
        </span>
      </div>
    );
  }

  if (block.type === "user_avatar") {
    return (
      <div style={{ width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center", background: s.background || "#e6ebf5", borderRadius: s["border-radius"] || "50%", overflow: "hidden" }}>
        <svg viewBox="0 0 100 100" style={{ width: "65%", height: "65%", opacity: 0.45 }}>
          <circle cx="50" cy="36" r="21" fill="#7a8fac" />
          <ellipse cx="50" cy="86" rx="32" ry="22" fill="#7a8fac" />
        </svg>
      </div>
    );
  }

  if (block.type === "image") {
    const fit = block.object_fit || "cover";
    if (block.var_src) {
      return (
        <div className={styles.imagePlaceholder} style={{ flexDirection: "column", gap: 5 }}>
          <svg width="30" height="30" viewBox="0 0 24 24" fill="none" style={{ opacity: 0.4 }}>
            <rect x="3" y="3" width="18" height="18" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
            <circle cx="8.5" cy="8.5" r="1.5" fill="#8a96b0" />
            <path d="M3 16l5-5 4 4 3-3 6 6" stroke="#8a96b0" strokeWidth="1.5" strokeLinecap="round" />
          </svg>
          <span style={{ fontSize: ".72rem", color: "#8a96b0" }}>{block.var_src}</span>
        </div>
      );
    }
    return block.src
      ? <img src={block.src} alt="imagem" style={{ width: "100%", height: "100%", objectFit: fit, display: "block" }} />
      : (
        <div className={styles.imagePlaceholder}>
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" style={{ opacity: 0.3 }}>
            <rect x="3" y="3" width="18" height="18" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
            <circle cx="8.5" cy="8.5" r="1.5" fill="#8a96b0" />
            <path d="M3 16l5-5 4 4 3-3 6 6" stroke="#8a96b0" strokeWidth="1.5" strokeLinecap="round" />
          </svg>
        </div>
      );
  }

  if (block.type === "carousel") {
    const first = block.images?.[0];
    const count = block.images?.length || 0;
    return (
      <div style={{ position: "relative", width: "100%", height: "100%", overflow: "hidden" }}>
        {first
          ? <img src={first} alt="carousel" style={{ width: "100%", height: "100%", objectFit: "cover", display: "block" }} />
          : (
            <div className={styles.imagePlaceholder} style={{ flexDirection: "column", gap: 6 }}>
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none" style={{ opacity: 0.3 }}>
                <rect x="3" y="3" width="18" height="18" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
                <path d="M9 3v18M15 3v18" stroke="#8a96b0" strokeWidth="1" strokeDasharray="2 2" />
              </svg>
              <span style={{ fontSize: ".72rem", color: "#8a96b0" }}>Carousel vazio</span>
            </div>
          )
        }
        {count > 1 && (
          <>
            <div style={{ position: "absolute", left: 8, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,.38)", borderRadius: "50%", width: 24, height: 24, display: "flex", alignItems: "center", justifyContent: "center", color: "#fff", fontSize: 14, lineHeight: 1 }}>‹</div>
            <div style={{ position: "absolute", right: 8, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,.38)", borderRadius: "50%", width: 24, height: 24, display: "flex", alignItems: "center", justifyContent: "center", color: "#fff", fontSize: 14, lineHeight: 1 }}>›</div>
            <div style={{ position: "absolute", bottom: 8, left: "50%", transform: "translateX(-50%)", display: "flex", gap: 4 }}>
              {Array.from({ length: Math.min(count, 5) }).map((_, i) => (
                <div key={i} style={{ width: i === 0 ? 14 : 6, height: 6, borderRadius: 3, background: i === 0 ? "#fff" : "rgba(255,255,255,.5)" }} />
              ))}
            </div>
          </>
        )}
      </div>
    );
  }

  if (block.type === "input_var") {
    return (
      <div className={styles.inputVarWrap}>
        {block.label_text && (
          <div style={{ fontSize: ".72rem", color: s.color || "#596582", marginBottom: 4, fontWeight: 500, fontFamily: s["font-family"] || "inherit" }}>
            {block.label_text}{block.required && <span style={{ color: "#e53e3e", marginLeft: 2 }}>*</span>}
          </div>
        )}
        <input
          className={styles.previewInput}
          type={block.input_type === "password" ? "password" : "text"}
          placeholder={block.placeholder || "Digite"}
          readOnly
          style={{
            background: s.background || "#fff",
            border: s.border || "1px solid #cad3e7",
            borderRadius: s["border-radius"] || "7px",
            color: s.color || "#1f2b43",
            fontSize: s["font-size"] || "14px",
            fontFamily: s["font-family"] || "inherit",
          }}
        />
        {block.var_name && <span className={styles.varBadge}>📌 {block.var_name}</span>}
      </div>
    );
  }

  if (block.type === "profile_card") {
    const align = block.align || "center";
    const avatarRadius = block.avatar_shape === "square" ? "10px" : "50%";
    return (
      <div style={{ display: "flex", flexDirection: "column", alignItems: align === "left" ? "flex-start" : align === "right" ? "flex-end" : "center", justifyContent: "center", gap: 8, width: "100%", height: "100%", padding: 8, boxSizing: "border-box", textAlign: align }}>
        {block.profile_image
          ? <img src={block.profile_image} alt="perfil" style={{ width: 64, height: 64, objectFit: "cover", borderRadius: avatarRadius, flexShrink: 0 }} />
          : (
            <div style={{ width: 64, height: 64, borderRadius: avatarRadius, background: "#e0e8f5", display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0 }}>
              <svg viewBox="0 0 40 40" width="34" height="34">
                <circle cx="20" cy="14" r="9" fill="#8a96b0" />
                <ellipse cx="20" cy="38" rx="14" ry="10" fill="#8a96b0" />
              </svg>
            </div>
          )
        }
        <div>
          <div style={{ fontWeight: 600, fontSize: "0.92em", color: s.color || "#1f2b43", fontFamily: s["font-family"] || "inherit" }}>{block.profile_name || "Nome do usuário"}</div>
          {block.profile_subtitle && <div style={{ fontSize: "0.78em", color: "#596582", marginTop: 2 }}>{block.profile_subtitle}</div>}
          {block.profile_bio && <div style={{ fontSize: "0.72em", color: "#6a7387", marginTop: 5, lineHeight: 1.45 }}>{block.profile_bio}</div>}
        </div>
      </div>
    );
  }

  if (block.type === "product_card") {
    return (
      <div style={{ display: "flex", flexDirection: "column", width: "100%", height: "100%", overflow: "hidden" }}>
        <div style={{ flex: "0 0 55%", background: "#e8edf7", display: "flex", alignItems: "center", justifyContent: "center", position: "relative" }}>
          <svg width="52" height="40" viewBox="0 0 24 18" fill="none" style={{ opacity: 0.25 }}>
            <rect x="1" y="1" width="22" height="16" rx="2" stroke="#8a96b0" strokeWidth="1.5" />
            <circle cx="7" cy="6" r="2" fill="#8a96b0" />
            <path d="M1 13l5-4 4 3 4-4 5 5" stroke="#8a96b0" strokeWidth="1.2" strokeLinecap="round" />
          </svg>
          {block.product_id && (
            <div style={{ position: "absolute", top: 5, right: 5, background: "rgba(0,0,0,.35)", color: "#fff", fontSize: ".55rem", padding: "1px 5px", borderRadius: 4 }}>ID</div>
          )}
        </div>
        <div style={{ flex: 1, padding: "8px 10px", display: "flex", flexDirection: "column", justifyContent: "space-between" }}>
          <div>
            <div style={{ height: 8, background: "#c8d4e8", borderRadius: 4, width: "80%", marginBottom: 5 }} />
            <div style={{ height: 6, background: "#dde6f3", borderRadius: 4, width: "60%" }} />
          </div>
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
            <div style={{ height: 11, background: "#b8c8e0", borderRadius: 4, width: "35%" }} />
            <div style={{ height: 22, width: 30, background: "#d2daeb", borderRadius: 5 }} />
          </div>
        </div>
      </div>
    );
  }

  if (block.type === "product_list") {
    const cardCount = Math.min(block.page_size || 6, 6);
    return (
      <div style={{ width: "100%", height: "100%", overflow: "hidden", display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 6, padding: 6, boxSizing: "border-box", alignContent: "start" }}>
        {Array.from({ length: cardCount }).map((_, i) => (
          <div key={i} style={{ background: "#fff", border: "1px solid #e2e8f3", borderRadius: 6, overflow: "hidden" }}>
            <div style={{ background: "#e8edf7", height: 52, display: "flex", alignItems: "center", justifyContent: "center" }}>
              <div style={{ width: 20, height: 16, background: "#d2daeb", borderRadius: 3 }} />
            </div>
            <div style={{ padding: "5px 6px", display: "flex", flexDirection: "column", gap: 3 }}>
              <div style={{ height: 5, background: "#d2daeb", borderRadius: 3, width: "80%" }} />
              <div style={{ height: 4, background: "#e8edf7", borderRadius: 3, width: "55%" }} />
              <div style={{ height: 7, background: "#c0cde0", borderRadius: 3, width: "40%", marginTop: 1 }} />
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (block.type === "video") {
    const url = block.video_url || "";
    const ytMatch = url.match(/(?:youtu\.be\/|youtube\.com\/(?:watch\?v=|embed\/))([\w-]+)/);
    const params = new URLSearchParams();
    if (block.autoplay) { params.set("autoplay", "1"); params.set("mute", "1"); }
    if (block.loop && ytMatch) params.set("loop", "1");
    const embedSrc = ytMatch
      ? `https://www.youtube.com/embed/${ytMatch[1]}${params.size ? `?${params}` : ""}`
      : url;
    return url
      ? <iframe src={embedSrc} allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowFullScreen style={{ width: "100%", height: "100%", border: "none", display: "block" }} title="video" />
      : (
        <div className={styles.videoPlaceholder} style={{ flexDirection: "column", gap: 10 }}>
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" style={{ opacity: 0.45 }}>
            <circle cx="12" cy="12" r="10" stroke="#aaa" strokeWidth="1.5" />
            <path d="M10 8.5l6 3.5-6 3.5V8.5z" fill="#aaa" />
          </svg>
          <span style={{ fontSize: ".78rem", color: "#777" }}>Adicione uma URL de video</span>
        </div>
      );
  }

  return <div style={{ width: "100%", height: "100%", background: s.background || "#cad3e7" }} />;
}
