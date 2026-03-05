import { useMemo, useRef, useCallback } from "react";
import styles from "./CodeEditor.module.scss";

type CodeEditorProps = {
  value: string;
  onChange: (value: string) => void;
  language: "tsx" | "svelte" | "css" | "json" | "txt";
  flat?: boolean; // remove border/radius when embedded full-screen
};

// ─── Syntax highlighter ────────────────────────────────────────────────────

function esc(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function sp(color: string, text: string): string {
  return `<span style="color:${color}">${esc(text)}</span>`;
}

const KW = new Set([
  "import","export","from","as","default","function","const","let","var",
  "return","if","else","for","while","do","switch","case","break","continue",
  "class","extends","interface","type","enum","namespace","declare","module",
  "new","this","super","null","undefined","true","false","async","await",
  "typeof","instanceof","in","of","throw","try","catch","finally","void",
  "delete","yield","get","set","static","public","private","protected","readonly",
  "abstract","override","implements","keyof","infer","satisfies","using",
]);
const TYPES = new Set([
  "string","number","boolean","any","never","unknown","object","symbol","bigint",
  "Array","Promise","Record","Partial","Required","Readonly","Pick","Omit",
  "Exclude","Extract","ReturnType","Parameters","React","FC","ReactNode",
  "MouseEvent","KeyboardEvent","ChangeEvent","FormEvent","CSSProperties",
]);

// Color palette (VSCode Dark+)
const C = {
  comment:  "#6a9955",
  string:   "#ce9178",
  kw:       "#569cd6",
  num:      "#b5cea8",
  fn:       "#dcdcaa",
  type:     "#4ec9b0",
  jsxTag:   "#4ec9b0",
  jsxProp:  "#9cdcfe",
  ident:    "#d6deee",
  at:       "#c586c0",
  cssSel:   "#d7ba7d",
  cssProp:  "#9cdcfe",
  cssAt:    "#c586c0",
  jsonKey:  "#9cdcfe",
  punct:    "#808080",
};

// ── JS/TSX tokenizer ──────────────────────────────────────────────────────

function highlightJS(code: string): string {
  const out: string[] = [];
  let i = 0;
  const n = code.length;

  while (i < n) {
    const ch = code[i];

    // Line comment
    if (ch === "/" && code[i + 1] === "/") {
      const end = code.indexOf("\n", i);
      const s = end === -1 ? code.slice(i) : code.slice(i, end);
      out.push(sp(C.comment, s));
      i += s.length;
      continue;
    }

    // Block comment
    if (ch === "/" && code[i + 1] === "*") {
      const end = code.indexOf("*/", i + 2);
      const s = end === -1 ? code.slice(i) : code.slice(i, end + 2);
      out.push(sp(C.comment, s));
      i += s.length;
      continue;
    }

    // String: "..." or '...'
    if (ch === '"' || ch === "'") {
      const q = ch;
      let j = i + 1;
      while (j < n && code[j] !== q && code[j] !== "\n") {
        if (code[j] === "\\") j++;
        j++;
      }
      if (j < n && code[j] === q) j++;
      out.push(sp(C.string, code.slice(i, j)));
      i = j;
      continue;
    }

    // Template literal
    if (ch === "`") {
      let j = i + 1;
      let depth = 0;
      while (j < n) {
        if (code[j] === "`" && depth === 0) { j++; break; }
        if (code[j] === "$" && code[j + 1] === "{") { depth++; j += 2; continue; }
        if (code[j] === "}" && depth > 0) { depth--; j++; continue; }
        if (code[j] === "\\") j++;
        j++;
      }
      out.push(sp(C.string, code.slice(i, j)));
      i = j;
      continue;
    }

    // JSX/HTML tag  <Tag ... > or </Tag> or <br />
    if (ch === "<" && i + 1 < n && (code[i + 1].match(/[A-Za-z/!]/) !== null)) {
      let j = i + 1;
      let inStr = false;
      let strCh = "";
      while (j < n) {
        if (inStr) {
          if (code[j] === "\\" && j + 1 < n) { j += 2; continue; }
          if (code[j] === strCh) inStr = false;
        } else {
          if (code[j] === '"' || code[j] === "'") { inStr = true; strCh = code[j]; }
          if (code[j] === ">") { j++; break; }
          // stop at newline only if not inside a string
          if (code[j] === "\n") break;
        }
        j++;
      }
      out.push(colorJSXTag(code.slice(i, j)));
      i = j;
      continue;
    }

    // Decorator: @word
    if (ch === "@" && code[i + 1]?.match(/[a-zA-Z_]/)) {
      let j = i + 1;
      while (j < n && code[j].match(/[a-zA-Z0-9_]/)) j++;
      out.push(sp(C.at, code.slice(i, j)));
      i = j;
      continue;
    }

    // Numbers
    if (ch.match(/[0-9]/) || (ch === "." && code[i + 1]?.match(/[0-9]/))) {
      let j = i;
      while (j < n && code[j].match(/[0-9.xXbBoOeEn_]/)) j++;
      out.push(sp(C.num, code.slice(i, j)));
      i = j;
      continue;
    }

    // Identifiers / keywords
    if (ch.match(/[a-zA-Z_$]/)) {
      let j = i;
      while (j < n && code[j].match(/[a-zA-Z0-9_$]/)) j++;
      const word = code.slice(i, j);
      const after = code.slice(j).trimStart();
      const isFn = after.startsWith("(") && !KW.has(word);
      if (KW.has(word)) out.push(sp(C.kw, word));
      else if (TYPES.has(word)) out.push(sp(C.type, word));
      else if (isFn) out.push(sp(C.fn, word));
      else if (word.match(/^[A-Z]/)) out.push(sp(C.type, word));
      else out.push(sp(C.ident, word));
      i = j;
      continue;
    }

    out.push(esc(ch));
    i++;
  }

  return out.join("");
}

function colorJSXTag(tag: string): string {
  const parts: string[] = [];
  let i = 0;
  const n = tag.length;

  parts.push(esc("<"));
  i++;
  if (i < n && tag[i] === "/") { parts.push(esc("/")); i++; }
  if (i < n && tag[i] === "!") {
    // comment or DOCTYPE
    parts.push(sp(C.comment, tag.slice(i)));
    return "<" + parts.slice(1).join("");
  }

  let j = i;
  while (j < n && tag[j].match(/[a-zA-Z0-9._:-]/)) j++;
  const tagName = tag.slice(i, j);
  parts.push(sp(tagName.match(/^[A-Z]/) ? C.type : C.jsxTag, tagName));
  i = j;

  while (i < n) {
    const ch = tag[i];
    if (ch === ">" || (ch === "/" && tag[i + 1] === ">")) {
      parts.push(esc(tag.slice(i)));
      break;
    }
    if (ch === '"' || ch === "'") {
      let j2 = i + 1;
      while (j2 < n && tag[j2] !== ch) { if (tag[j2] === "\\") j2++; j2++; }
      if (j2 < n) j2++;
      parts.push(sp(C.string, tag.slice(i, j2)));
      i = j2;
      continue;
    }
    if (ch === "{") {
      let depth = 1; let j2 = i + 1;
      while (j2 < n && depth > 0) { if (tag[j2] === "{") depth++; if (tag[j2] === "}") depth--; j2++; }
      parts.push(esc(tag.slice(i, j2)));
      i = j2;
      continue;
    }
    if (ch.match(/[a-zA-Z_$]/)) {
      let j2 = i;
      while (j2 < n && tag[j2].match(/[a-zA-Z0-9_$:-]/)) j2++;
      parts.push(sp(C.jsxProp, tag.slice(i, j2)));
      i = j2;
      continue;
    }
    parts.push(esc(ch));
    i++;
  }

  return parts.join("");
}

// ── CSS tokenizer ─────────────────────────────────────────────────────────

function highlightCSS(code: string): string {
  const out: string[] = [];
  let i = 0;
  const n = code.length;
  let inBlock = false;
  let afterColon = false;

  while (i < n) {
    const ch = code[i];

    // Block comment
    if (ch === "/" && code[i + 1] === "*") {
      const end = code.indexOf("*/", i + 2);
      const s = end === -1 ? code.slice(i) : code.slice(i, end + 2);
      out.push(sp(C.comment, s));
      i += s.length;
      continue;
    }

    // At-rule
    if (ch === "@") {
      let j = i + 1;
      while (j < n && code[j].match(/[a-zA-Z-]/)) j++;
      out.push(sp(C.cssAt, code.slice(i, j)));
      i = j;
      continue;
    }

    // String
    if (ch === '"' || ch === "'") {
      const q = ch;
      let j = i + 1;
      while (j < n && code[j] !== q) { if (code[j] === "\\") j++; j++; }
      if (j < n) j++;
      out.push(sp(C.string, code.slice(i, j)));
      i = j;
      continue;
    }

    if (ch === "{") { inBlock = true; afterColon = false; out.push(esc(ch)); i++; continue; }
    if (ch === "}") { inBlock = false; afterColon = false; out.push(esc(ch)); i++; continue; }
    if (ch === ":") { afterColon = true; out.push(esc(ch)); i++; continue; }
    if (ch === ";") { afterColon = false; out.push(esc(ch)); i++; continue; }

    // Property name or selector
    if (ch.match(/[a-zA-Z_-]/)) {
      let j = i;
      while (j < n && code[j].match(/[a-zA-Z0-9_%-]/)) j++;
      const word = code.slice(i, j);
      if (!inBlock) {
        out.push(sp(C.cssSel, word));
      } else if (!afterColon) {
        out.push(sp(C.cssProp, word));
      } else {
        out.push(sp(C.string, word));
      }
      i = j;
      continue;
    }

    // Numbers with units
    if (ch.match(/[0-9]/) || (ch === "-" && code[i + 1]?.match(/[0-9]/))) {
      let j = i;
      if (ch === "-") j++;
      while (j < n && code[j].match(/[0-9.%a-zA-Z]/)) j++;
      out.push(sp(C.num, code.slice(i, j)));
      i = j;
      continue;
    }

    // Color hex
    if (ch === "#" && code[i + 1]?.match(/[0-9a-fA-F]/)) {
      let j = i + 1;
      while (j < n && code[j].match(/[0-9a-fA-F]/)) j++;
      out.push(sp(C.num, code.slice(i, j)));
      i = j;
      continue;
    }

    out.push(esc(ch));
    i++;
  }

  return out.join("");
}

// ── JSON tokenizer ────────────────────────────────────────────────────────

function highlightJSON(code: string): string {
  const out: string[] = [];
  let i = 0;
  const n = code.length;

  while (i < n) {
    const ch = code[i];

    if (ch === '"') {
      let j = i + 1;
      while (j < n && code[j] !== '"') { if (code[j] === "\\") j++; j++; }
      if (j < n) j++;
      const raw = code.slice(i, j);
      // Is it a key? Check if ':' follows (skipping whitespace)
      let k = j;
      while (k < n && code[k].match(/\s/)) k++;
      if (code[k] === ":") out.push(sp(C.jsonKey, raw));
      else out.push(sp(C.string, raw));
      i = j;
      continue;
    }

    if (ch.match(/[0-9\-]/)) {
      let j = i;
      while (j < n && code[j].match(/[0-9.eE+\-]/)) j++;
      out.push(sp(C.num, code.slice(i, j)));
      i = j;
      continue;
    }

    if (code.slice(i, i + 4) === "true") { out.push(sp(C.kw, "true")); i += 4; continue; }
    if (code.slice(i, i + 5) === "false") { out.push(sp(C.kw, "false")); i += 5; continue; }
    if (code.slice(i, i + 4) === "null") { out.push(sp(C.kw, "null")); i += 4; continue; }

    out.push(esc(ch));
    i++;
  }

  return out.join("");
}

// ── Svelte: split into script/style/template sections ─────────────────────

function highlightSvelte(code: string): string {
  const out: string[] = [];
  let i = 0;
  const n = code.length;

  while (i < n) {
    // <script ...>...</script>
    const scriptStart = code.indexOf("<script", i);
    const styleStart  = code.indexOf("<style", i);

    let nextTag = -1;
    let tagName = "";
    if (scriptStart !== -1 && (styleStart === -1 || scriptStart <= styleStart)) {
      nextTag = scriptStart; tagName = "script";
    } else if (styleStart !== -1) {
      nextTag = styleStart; tagName = "style";
    }

    if (nextTag === -1 || nextTag > i) {
      // Template region
      const endTemplate = nextTag === -1 ? n : nextTag;
      out.push(highlightSvelteTemplate(code.slice(i, endTemplate)));
      i = nextTag === -1 ? n : nextTag;
      continue;
    }

    // Opening tag
    const openEnd = code.indexOf(">", nextTag);
    if (openEnd === -1) { out.push(esc(code.slice(i))); break; }
    out.push(sp(C.jsxTag, esc(code.slice(nextTag, openEnd + 1))));
    i = openEnd + 1;

    // Content until </tagName>
    const closeTag = `</${tagName}>`;
    const closeIdx = code.indexOf(closeTag, i);
    const contentEnd = closeIdx === -1 ? n : closeIdx;
    const content = code.slice(i, contentEnd);
    out.push(tagName === "script" ? highlightJS(content) : highlightCSS(content));
    i = contentEnd;

    if (closeIdx !== -1) {
      out.push(sp(C.jsxTag, esc(closeTag)));
      i += closeTag.length;
    }
  }

  return out.join("");
}

function highlightSvelteTemplate(tpl: string): string {
  // Colorize {#if}, {#each}, {/if}, {/each}, {:else}, {expr}
  const out: string[] = [];
  let i = 0;
  const n = tpl.length;

  while (i < n) {
    const ch = tpl[i];

    if (ch === "<") {
      // HTML tag
      let j = i + 1;
      let inStr = false; let strCh = "";
      while (j < n) {
        if (inStr) { if (tpl[j] === "\\" && j + 1 < n) { j += 2; continue; } if (tpl[j] === strCh) inStr = false; }
        else { if (tpl[j] === '"' || tpl[j] === "'") { inStr = true; strCh = tpl[j]; } if (tpl[j] === ">") { j++; break; } }
        j++;
      }
      out.push(colorJSXTag(tpl.slice(i, j)));
      i = j;
      continue;
    }

    if (ch === "{") {
      let j = i + 1;
      let depth = 1;
      while (j < n && depth > 0) { if (tpl[j] === "{") depth++; if (tpl[j] === "}") depth--; j++; }
      const block = tpl.slice(i, j);
      const inner = block.slice(1, -1);
      if (inner.startsWith("#") || inner.startsWith("/") || inner.startsWith(":")) {
        out.push(sp(C.kw, "{") + sp(C.at, inner) + sp(C.kw, "}"));
      } else {
        out.push(sp(C.punct, "{") + highlightJS(inner) + sp(C.punct, "}"));
      }
      i = j;
      continue;
    }

    out.push(esc(ch));
    i++;
  }

  return out.join("");
}

function highlight(code: string, lang: string): string {
  try {
    switch (lang) {
      case "tsx": return highlightJS(code);
      case "svelte": return highlightSvelte(code);
      case "css": return highlightCSS(code);
      case "json": return highlightJSON(code);
      default: return esc(code);
    }
  } catch {
    return esc(code);
  }
}

// ─── Component ──────────────────────────────────────────────────────────────

export default function CodeEditor({ value, onChange, language, flat }: CodeEditorProps) {
  const gutterRef = useRef<HTMLDivElement>(null);
  const stackRef  = useRef<HTMLDivElement>(null);

  const lines = useMemo(() => {
    const count = value.split("\n").length;
    return Array.from({ length: Math.max(count, 1) }, (_, i) => i + 1);
  }, [value]);

  const highlighted = useMemo(() => highlight(value, language), [value, language]);

  // Sync gutter scroll with code scroll
  const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
    if (gutterRef.current) gutterRef.current.scrollTop = (e.currentTarget as HTMLDivElement).scrollTop;
  }, []);

  function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key !== "Tab") return;
    e.preventDefault();
    const target = e.currentTarget;
    const start = target.selectionStart;
    const end = target.selectionEnd;
    const indent = "  ";
    const next = value.slice(0, start) + indent + value.slice(end);
    onChange(next);
    requestAnimationFrame(() => {
      target.selectionStart = start + indent.length;
      target.selectionEnd   = start + indent.length;
    });
  }

  return (
    <div className={`${styles.wrap}${flat ? ` ${styles.wrapFlat}` : ""}`}>
      <div className={styles.topbar}>
        <span className={styles.lang}>{language}</span>
      </div>
      <div className={styles.body}>
        <div className={styles.gutter} ref={gutterRef}>
          {lines.map((line) => (
            <div key={line} className={styles.lineNumber}>{line}</div>
          ))}
        </div>
        <div className={styles.codeStack} ref={stackRef} onScroll={handleScroll}>
          {/* Highlighted layer */}
          <pre
            className={styles.highlightPre}
            aria-hidden="true"
            dangerouslySetInnerHTML={{ __html: highlighted + "\n" }}
          />
          {/* Editable layer */}
          <textarea
            className={styles.codeTA}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onKeyDown={handleKeyDown}
            spellCheck={false}
            autoCorrect="off"
            autoCapitalize="off"
          />
        </div>
      </div>
    </div>
  );
}
