import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { API_URL, apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import {
  ZOOM_LEVELS, ensurePath, newBlock, isTextBlock,
  computeSnap, makeLandingTemplate, nextId, parseDoc,
} from "./blockHelpers";
import { renderBlock } from "./BlockRenderer";
import { Inspector } from "./Inspector";
import { LeftPanel } from "./LeftPanel";
import styles from "./ElementorEditor.module.scss";

const BLANK_DOC = {
  title: "Site",
  canvas: { width: 1400, height: 900, background: "#f8f9ff" },
  header: { enabled: false, height: 80, background: "#1a2740", blocks: [] },
  footer: { enabled: false, height: 100, background: "#1a2740", blocks: [] },
  pages: { "/": { title: "Pagina Inicial", blocks: [] } },
  popups: {},
};

export const ElementorEditor = () => {
  const { siteId = "" } = useParams();
  const { websiteId } = useAuthContext();
  const navigate = useNavigate();

  const pageCanvasRef = useRef(null);
  const headerCanvasRef = useRef(null);
  const footerCanvasRef = useRef(null);
  const docRef = useRef(null);
  const activeRouteRef = useRef("/");
  const editingPopupRef = useRef(null);
  const editingProductCardRef = useRef(null);
  const editingProductCardSectionRef = useRef("page");
  const undoStack = useRef([]);
  const redoStack = useRef([]);

  const [doc, setDoc] = useState(BLANK_DOC);
  const [routes, setRoutes] = useState([{ id: "root", path: "/", title: "Inicio", requires_auth: false, position: 0 }]);
  const [activeRoute, setActiveRoute] = useState("/");
  const [selected, setSelected] = useState(null);
  const [editingPopup, setEditingPopup] = useState(null);
  const [editingProductCard, setEditingProductCard] = useState(null);
  const [editingProductCardSection, setEditingProductCardSection] = useState("page");
  const [editSection, setEditSection] = useState("page");
  const [leftTab, setLeftTab] = useState("elements");
  const [zoom, setZoom] = useState(1.0);
  const [interaction, setInteraction] = useState(null);
  const [guides, setGuides] = useState({ vertical: null, horizontal: null });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [apis, setApis] = useState([]);
  const [menu, setMenu] = useState(null);
  const [showTemplateModal, setShowTemplateModal] = useState(false);
  const [canUndo, setCanUndo] = useState(false);
  const [canRedo, setCanRedo] = useState(false);

  docRef.current = doc;
  activeRouteRef.current = activeRoute;
  editingPopupRef.current = editingPopup;
  editingProductCardRef.current = editingProductCard;
  editingProductCardSectionRef.current = editingProductCardSection;

  const currentBlocks = useMemo(() => {
    if (editingPopup) return doc.popups[editingPopup]?.blocks || [];
    if (editingProductCard) {
      const sec = editingProductCardSection;
      const src = sec === "header" ? doc.header.blocks : sec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
      return src.find((b) => b.id === editingProductCard)?.inner_blocks || [];
    }
    if (editSection === "header") return doc.header.blocks;
    if (editSection === "footer") return doc.footer.blocks;
    return doc.pages[activeRoute]?.blocks || [];
  }, [editSection, doc, activeRoute, editingPopup, editingProductCard, editingProductCardSection]);

  const pcCardBlock = useMemo(() => {
    if (!editingProductCard) return null;
    const sec = editingProductCardSection;
    const src = sec === "header" ? doc.header.blocks : sec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
    return src.find((b) => b.id === editingProductCard) ?? null;
  }, [editingProductCard, editingProductCardSection, doc, activeRoute]);

  const selectedBlock = useMemo(() => currentBlocks.find((b) => b.id === selected) || null, [currentBlocks, selected]);

  function pushHistory() {
    undoStack.current = [...undoStack.current, docRef.current].slice(-30);
    redoStack.current = [];
    setCanUndo(true);
    setCanRedo(false);
  }

  function undo() {
    if (!undoStack.current.length) return;
    const prev = undoStack.current[undoStack.current.length - 1];
    redoStack.current = [docRef.current, ...redoStack.current];
    undoStack.current = undoStack.current.slice(0, -1);
    setDoc(prev);
    setCanUndo(undoStack.current.length > 0);
    setCanRedo(true);
  }

  function redo() {
    if (!redoStack.current.length) return;
    const next = redoStack.current[0];
    undoStack.current = [...undoStack.current, docRef.current];
    redoStack.current = redoStack.current.slice(1);
    setDoc(next);
    setCanUndo(true);
    setCanRedo(redoStack.current.length > 0);
  }

  function updateBlock(section, id, patch) {
    setDoc((prev) => {
      const pid = editingPopupRef.current;
      if (pid) {
        return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: (prev.popups[pid]?.blocks || []).map((b) => b.id === id ? { ...b, ...patch } : b) } } };
      }
      const pcId = editingProductCardRef.current;
      if (pcId) {
        const pcSec = editingProductCardSectionRef.current;
        const mp = (bl) => bl.map((b) => b.id === pcId ? { ...b, inner_blocks: (b.inner_blocks || []).map((ib) => ib.id === id ? { ...ib, ...patch } : ib) } : b);
        if (pcSec === "header") return { ...prev, header: { ...prev.header, blocks: mp(prev.header.blocks) } };
        if (pcSec === "footer") return { ...prev, footer: { ...prev.footer, blocks: mp(prev.footer.blocks) } };
        const r = activeRouteRef.current;
        return { ...prev, pages: { ...prev.pages, [r]: { ...prev.pages[r], blocks: mp(prev.pages[r]?.blocks || []) } } };
      }
      if (section === "header") return { ...prev, header: { ...prev.header, blocks: prev.header.blocks.map((b) => b.id === id ? { ...b, ...patch } : b) } };
      if (section === "footer") return { ...prev, footer: { ...prev.footer, blocks: prev.footer.blocks.map((b) => b.id === id ? { ...b, ...patch } : b) } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: (prev.pages[route]?.blocks || []).map((b) => b.id === id ? { ...b, ...patch } : b) } } };
    });
  }

  function clampBlock(block, cw, ch) {
    const w = Math.min(block.w, cw);
    const h = Math.min(block.h, ch);
    return { ...block, w, h, x: Math.max(0, Math.min(block.x, cw - w)), y: Math.max(0, Math.min(block.y, ch - h)) };
  }

  function addBlock(type) {
    const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0);
    let block = newBlock(type, maxZ + 1);
    const pid = editingPopupRef.current;
    const pcId = editingProductCardRef.current;
    if (pcId) {
      const pcSec = editingProductCardSectionRef.current;
      const src = pcSec === "header" ? doc.header.blocks : pcSec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
      const pcBl = src.find((b) => b.id === pcId);
      block = clampBlock(block, pcBl?.w || 300, pcBl?.h || 280);
    } else if (pid) {
      block = clampBlock(block, doc.popups[pid]?.width || 480, doc.popups[pid]?.height || 560);
    } else if (editSection === "header") {
      block = clampBlock(block, doc.canvas.width, doc.header.height || 80);
    } else if (editSection === "footer") {
      block = clampBlock(block, doc.canvas.width, doc.footer.height || 100);
    }
    pushHistory();
    setDoc((prev) => {
      const pp = editingPopupRef.current;
      const pc = editingProductCardRef.current;
      if (pc) {
        const pcS = editingProductCardSectionRef.current;
        const mp = (bl) => bl.map((b) => b.id === pc ? { ...b, inner_blocks: [...(b.inner_blocks || []), block] } : b);
        if (pcS === "header") return { ...prev, header: { ...prev.header, blocks: mp(prev.header.blocks) } };
        if (pcS === "footer") return { ...prev, footer: { ...prev.footer, blocks: mp(prev.footer.blocks) } };
        const r = activeRouteRef.current;
        return { ...prev, pages: { ...prev.pages, [r]: { ...prev.pages[r], blocks: mp(prev.pages[r]?.blocks || []) } } };
      }
      if (pp) return { ...prev, popups: { ...prev.popups, [pp]: { ...prev.popups[pp], blocks: [...(prev.popups[pp]?.blocks || []), block] } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: [...prev.header.blocks, block] } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: [...prev.footer.blocks, block] } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: [...(prev.pages[route]?.blocks || []), block] } } };
    });
    setSelected(block.id);
  }

  function removeBlock(id) {
    pushHistory();
    setDoc((prev) => {
      const pid = editingPopupRef.current;
      const pcId = editingProductCardRef.current;
      if (pcId) {
        const pcSec = editingProductCardSectionRef.current;
        const mp = (bl) => bl.map((b) => b.id === pcId ? { ...b, inner_blocks: (b.inner_blocks || []).filter((ib) => ib.id !== id) } : b);
        if (pcSec === "header") return { ...prev, header: { ...prev.header, blocks: mp(prev.header.blocks) } };
        if (pcSec === "footer") return { ...prev, footer: { ...prev.footer, blocks: mp(prev.footer.blocks) } };
        const r = activeRouteRef.current;
        return { ...prev, pages: { ...prev.pages, [r]: { ...prev.pages[r], blocks: mp(prev.pages[r]?.blocks || []) } } };
      }
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: (prev.popups[pid]?.blocks || []).filter((b) => b.id !== id) } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: prev.header.blocks.filter((b) => b.id !== id) } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: prev.footer.blocks.filter((b) => b.id !== id) } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: (prev.pages[route]?.blocks || []).filter((b) => b.id !== id) } } };
    });
    if (selected === id) setSelected(null);
  }

  function duplicateBlock(id) {
    const source = currentBlocks.find((b) => b.id === id);
    if (!source) return;
    const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0);
    const copy = { ...source, id: nextId("blk"), x: source.x + 24, y: source.y + 24, z: maxZ + 1 };
    pushHistory();
    setDoc((prev) => {
      const pid = editingPopupRef.current;
      const pcId = editingProductCardRef.current;
      if (pcId) {
        const pcSec = editingProductCardSectionRef.current;
        const mp = (bl) => bl.map((b) => b.id === pcId ? { ...b, inner_blocks: [...(b.inner_blocks || []), copy] } : b);
        if (pcSec === "header") return { ...prev, header: { ...prev.header, blocks: mp(prev.header.blocks) } };
        if (pcSec === "footer") return { ...prev, footer: { ...prev.footer, blocks: mp(prev.footer.blocks) } };
        const r = activeRouteRef.current;
        return { ...prev, pages: { ...prev.pages, [r]: { ...prev.pages[r], blocks: mp(prev.pages[r]?.blocks || []) } } };
      }
      if (pid) return { ...prev, popups: { ...prev.popups, [pid]: { ...prev.popups[pid], blocks: [...(prev.popups[pid]?.blocks || []), copy] } } };
      if (editSection === "header") return { ...prev, header: { ...prev.header, blocks: [...prev.header.blocks, copy] } };
      if (editSection === "footer") return { ...prev, footer: { ...prev.footer, blocks: [...prev.footer.blocks, copy] } };
      const route = activeRouteRef.current;
      return { ...prev, pages: { ...prev.pages, [route]: { ...prev.pages[route], blocks: [...(prev.pages[route]?.blocks || []), copy] } } };
    });
    setSelected(copy.id);
  }

  function contextAction(action) {
    if (!menu) return;
    const block = currentBlocks.find((b) => b.id === menu.id);
    if (!block) return;
    if (action === "delete") { removeBlock(block.id); setMenu(null); return; }
    if (action === "duplicate") { duplicateBlock(block.id); setMenu(null); return; }
    if (action === "bringForward") {
      const maxZ = currentBlocks.reduce((acc, b) => Math.max(acc, b.z || 0), 0);
      updateBlock(editSection, block.id, { z: maxZ + 1 }); setMenu(null); return;
    }
    if (action === "sendBackward") {
      const minZ = currentBlocks.reduce((acc, b) => Math.min(acc, b.z || 0), Infinity);
      updateBlock(editSection, block.id, { z: Math.max(0, minZ - 1) }); setMenu(null); return;
    }
    const pid = editingPopupRef.current;
    const pcId = editingProductCardRef.current;
    const canvasW = pcId ? (pcCardBlock?.w || 300) : pid ? (doc.popups[pid]?.width || 480) : doc.canvas.width;
    const canvasH = pcId ? (pcCardBlock?.h || 280) : pid ? (doc.popups[pid]?.height || 560) : (editSection === "header" ? doc.header.height : editSection === "footer" ? doc.footer.height : doc.canvas.height);
    if (action === "width100") { updateBlock(editSection, block.id, { w: canvasW }); setMenu(null); return; }
    updateBlock(editSection, block.id, { h: canvasH }); setMenu(null);
  }

  function addPage() {
    const newPath = `/pagina-${Date.now()}`;
    const newRoute = { id: nextId("rt"), path: newPath, title: `Pagina ${routes.length + 1}`, requires_auth: false, position: routes.length };
    setRoutes((prev) => [...prev, newRoute]);
    setDoc((prev) => ({ ...prev, pages: { ...prev.pages, [newPath]: { title: newRoute.title, blocks: [] } } }));
    setActiveRoute(newPath);
  }

  function removePage(path) {
    const safePath = ensurePath(path);
    if (routes.length <= 1) return;
    setRoutes((prev) => prev.filter((r) => ensurePath(r.path) !== safePath));
    setDoc((prev) => { const p = { ...prev.pages }; delete p[safePath]; return { ...prev, pages: p }; });
    if (activeRoute === safePath) {
      const remaining = routes.filter((r) => ensurePath(r.path) !== safePath);
      setActiveRoute(ensurePath(remaining[0]?.path || "/"));
    }
  }

  function addPopup() {
    const id = nextId("popup");
    pushHistory();
    setDoc((prev) => ({ ...prev, popups: { ...prev.popups, [id]: { title: `Popup ${Object.keys(prev.popups).length + 1}`, blocks: [], width: 480, height: 560, background: "#ffffff" } } }));
    setEditingPopup(id);
    setEditSection("page");
  }

  function removePopup(id) {
    pushHistory();
    setDoc((prev) => { const p = { ...prev.popups }; delete p[id]; return { ...prev, popups: p }; });
    if (editingPopup === id) setEditingPopup(null);
  }

  function applyTemplate(type) {
    if (type === "landing") {
      const t = makeLandingTemplate();
      setDoc(t);
      setRoutes(Object.keys(t.pages).map((path, i) => ({ id: nextId("rt"), path, title: t.pages[path].title, requires_auth: false, position: i })));
      setActiveRoute("/");
    } else {
      setDoc(BLANK_DOC);
      if (!routes.length) setRoutes([{ id: nextId("rt"), path: "/", title: "Pagina Inicial", requires_auth: false, position: 0 }]);
      setActiveRoute("/");
    }
    setShowTemplateModal(false);
  }

  async function save(publish) {
    if (!siteId) return;
    setSaving(true); setError(null); setSuccess(null);
    try {
      const routeInputs = routes.map((r, i) => ({ path: ensurePath(r.path), title: r.title || ensurePath(r.path), requires_auth: r.requires_auth, position: i }));
      await apiFetch(`/api/v1/sites/${siteId}/routes`, { method: "POST", websiteId, body: JSON.stringify({ routes: routeInputs }) });
      const resp = await apiFetch(`/api/v1/sites/${siteId}/versions`, { method: "POST", websiteId, body: JSON.stringify({ source_type: "ELEMENTOR_JSON", source: JSON.stringify(doc) }) });
      if (publish) {
        if (resp.data.scan_status === "blocked") { setError("Versao bloqueada pelo scan."); return; }
        await apiFetch(`/api/v1/sites/${siteId}/publish/${resp.data.version}`, { method: "POST", websiteId });
        setSuccess("✅ Publicado com sucesso!");
      } else {
        setSuccess("💾 Salvo!");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao salvar");
    } finally {
      setSaving(false);
    }
  }

  useEffect(() => {
    async function load() {
      if (!siteId) return;
      try {
        const [routeResp, versionResp, apiResp] = await Promise.all([
          apiFetch(`/api/v1/sites/${siteId}/routes`, { method: "GET", websiteId }),
          apiFetch(`/api/v1/sites/${siteId}/versions`, { method: "GET", websiteId }),
          apiFetch("/api/v1/site-apis", { method: "GET", websiteId }),
        ]);
        const availableRoutes = routeResp.data.length > 0 ? routeResp.data : [{ id: "root", path: "/", title: "Inicio", requires_auth: false, position: 0 }];
        setRoutes(availableRoutes);
        const latestAny = versionResp.data[0];
        const latestElementor = versionResp.data.find((item) => item.source_type === "ELEMENTOR_JSON");
        if (latestAny && latestAny.source_type !== "ELEMENTOR_JSON") {
          setError(`Este site foi criado em ${latestAny.source_type}. Edite apenas nesse modo.`);
          return;
        }
        const parsed = parseDoc(latestElementor?.source, availableRoutes);
        if (parsed === null) {
          setDoc({ ...BLANK_DOC, pages: Object.fromEntries(availableRoutes.map((r) => [ensurePath(r.path), { title: r.title || ensurePath(r.path), blocks: [] }])) });
          setShowTemplateModal(true);
        } else {
          setDoc(parsed);
        }
        setActiveRoute(ensurePath(availableRoutes[0]?.path || "/"));
        setApis(apiResp.data);
      } catch (err) {
        const msg = err instanceof Error ? err.message : "Erro ao carregar editor";
        if (msg.includes("403") || msg.toLowerCase().includes("forbidden") || msg.toLowerCase().includes("unauthorized")) {
          navigate("/pages", { replace: true });
          return;
        }
        setError(msg);
      }
    }
    load();
  }, [siteId, websiteId]);

  const closeMenu = useCallback(() => setMenu(null), []);
  useEffect(() => {
    window.addEventListener("click", closeMenu);
    return () => window.removeEventListener("click", closeMenu);
  }, [closeMenu]);

  useEffect(() => {
    function onMove(e) {
      if (!interaction) return;
      const ref = interaction.section === "header" ? headerCanvasRef : interaction.section === "footer" ? footerCanvasRef : pageCanvasRef;
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      const pid = editingPopupRef.current;
      const pcId = editingProductCardRef.current;
      let blocks, canvasW, canvasH;
      if (interaction.section === "header") { blocks = doc.header.blocks; canvasH = doc.header.height; canvasW = doc.canvas.width; }
      else if (interaction.section === "footer") { blocks = doc.footer.blocks; canvasH = doc.footer.height; canvasW = doc.canvas.width; }
      else if (pcId) {
        const pcSec = editingProductCardSectionRef.current;
        const src = pcSec === "header" ? doc.header.blocks : pcSec === "footer" ? doc.footer.blocks : doc.pages[activeRoute]?.blocks || [];
        const pcBl = src.find((b) => b.id === pcId);
        blocks = pcBl?.inner_blocks || []; canvasH = pcBl?.h || 280; canvasW = pcBl?.w || 300;
      } else if (pid) {
        blocks = doc.popups[pid]?.blocks || []; canvasH = doc.popups[pid]?.height || 560; canvasW = doc.popups[pid]?.width || 480;
      } else {
        blocks = doc.pages[activeRoute]?.blocks || []; canvasH = doc.canvas.height; canvasW = doc.canvas.width;
      }

      if (interaction.mode === "drag") {
        const block = blocks.find((b) => b.id === interaction.id);
        if (!block) return;
        const nextX = (e.clientX - rect.left - interaction.dx) / zoom;
        const nextY = (e.clientY - rect.top - interaction.dy) / zoom;
        const snapped = computeSnap(nextX, nextY, block.w, block.h, block.id, blocks, canvasW, canvasH);
        setGuides(snapped.guides);
        updateBlock(interaction.section, interaction.id, { x: Math.max(0, Math.min(canvasW - block.w, snapped.x)), y: Math.max(0, Math.min(canvasH - block.h, snapped.y)) });
        return;
      }
      if (interaction.mode === "resize") {
        const dx = (e.clientX - interaction.startX) / zoom;
        const dy = (e.clientY - interaction.startY) / zoom;
        let { x, y, w, h } = interaction.start;
        if (interaction.edge === "right" || interaction.edge === "corner") w = interaction.start.w + dx;
        if (interaction.edge === "bottom" || interaction.edge === "corner") h = interaction.start.h + dy;
        if (interaction.edge === "left") { x = interaction.start.x + dx; w = interaction.start.w - dx; }
        if (interaction.edge === "top") { y = interaction.start.y + dy; h = interaction.start.h - dy; }
        w = Math.max(48, Math.min(canvasW, Math.round(w)));
        h = Math.max(24, Math.min(canvasH, Math.round(h)));
        x = Math.max(0, Math.min(canvasW - w, Math.round(x)));
        y = Math.max(0, Math.min(canvasH - h, Math.round(y)));
        const block = blocks.find((b) => b.id === interaction.id);
        if (block && isTextBlock(block.type) && interaction.baseFontSize && interaction.start.h > 0) {
          const scaledFont = Math.max(8, Math.round(interaction.baseFontSize * (h / interaction.start.h)));
          updateBlock(interaction.section, interaction.id, { x, y, w, h, style: { ...block.style, "font-size": `${scaledFont}px` } });
        } else {
          updateBlock(interaction.section, interaction.id, { x, y, w, h });
        }
        return;
      }
      const block = blocks.find((b) => b.id === interaction.id);
      if (!block) return;
      const angle = (Math.atan2(e.clientY - interaction.cy, e.clientX - interaction.cx) * 180) / Math.PI;
      updateBlock(interaction.section, interaction.id, { rotation: Math.round(interaction.startRotation + (angle - interaction.startAngle)) });
    }
    function onUp() { setInteraction(null); setGuides({ vertical: null, horizontal: null }); }
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => { window.removeEventListener("mousemove", onMove); window.removeEventListener("mouseup", onUp); };
  }, [interaction, doc, activeRoute, zoom]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    function onKeyDown(e) {
      const tag = e.target.tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      if ((e.key === "Delete" || e.key === "Backspace") && selected) removeBlock(selected);
      else if (e.key === "Escape") setSelected(null);
      else if ((e.ctrlKey || e.metaKey) && e.key === "d") { e.preventDefault(); if (selected) duplicateBlock(selected); }
      else if ((e.ctrlKey || e.metaKey) && e.key === "z") { e.preventDefault(); undo(); }
      else if ((e.ctrlKey || e.metaKey) && e.key === "y") { e.preventDefault(); redo(); }
    }
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [selected, editSection]); // eslint-disable-line react-hooks/exhaustive-deps

  function renderSectionBlocks(blocks, section) {
    const isActive = editSection === section || (section === "page" && (!!editingPopupRef.current || !!editingProductCardRef.current));
    return [...blocks].sort((a, b) => (a.z || 0) - (b.z || 0)).map((block) => (
      <article
        key={block.id}
        className={`${styles.block} ${selected === block.id && isActive ? styles.blockSelected : ""} ${!isActive ? styles.blockFrozen : ""}`}
        style={{
          left: block.x, top: block.y, width: block.w, height: block.h, zIndex: block.z,
          transform: `rotate(${block.rotation}deg)`,
          color: block.style.color, background: block.style.background,
          border: block.style.border, borderRadius: block.style["border-radius"],
          padding: block.style.padding, fontSize: block.style["font-size"],
          fontWeight: block.style["font-weight"], textAlign: block.style["text-align"],
          cursor: isActive ? "move" : "pointer",
        }}
        onMouseDown={(e) => {
          if (!isActive) { e.stopPropagation(); setEditSection(section); setSelected(null); return; }
          const blockRect = e.currentTarget.getBoundingClientRect();
          setInteraction({ mode: "drag", id: block.id, section, dx: e.clientX - blockRect.left, dy: e.clientY - blockRect.top });
          setSelected(block.id); setMenu(null); e.stopPropagation();
        }}
        onClick={(e) => e.stopPropagation()}
        onContextMenu={(e) => {
          if (!isActive) return;
          e.preventDefault(); setSelected(block.id); setMenu({ x: e.clientX, y: e.clientY, id: block.id });
        }}
      >
        {selected === block.id && isActive && <div className={styles.blockTag}>{block.type}</div>}
        {renderBlock(block)}
        {selected === block.id && isActive && (
          <>
            <button type="button" className={`${styles.handle} ${styles.rotateHandle}`} onMouseDown={(e) => { e.stopPropagation(); const r = e.currentTarget.parentElement.getBoundingClientRect(); const cx = r.left + r.width / 2; const cy = r.top + r.height / 2; setInteraction({ mode: "rotate", id: block.id, section, cx, cy, startAngle: (Math.atan2(e.clientY - cy, e.clientX - cx) * 180) / Math.PI, startRotation: block.rotation }); }} />
            {["right", "left", "top", "bottom", "corner"].map((edge) => (
              <button key={edge} type="button" className={`${styles.handle} ${styles[`${edge}Handle`]}`} onMouseDown={(e) => { e.stopPropagation(); const bfs = isTextBlock(block.type) ? parseFloat(block.style["font-size"] || "16") : undefined; setInteraction({ mode: "resize", id: block.id, section, edge, startX: e.clientX, startY: e.clientY, start: { x: block.x, y: block.y, w: block.w, h: block.h }, baseFontSize: bfs }); }} />
            ))}
          </>
        )}
      </article>
    ));
  }

  function SectionCanvas({ canvasRef, section, width, height, style, label, blocks }) {
    const isActive = editSection === section;
    return (
      <div
        ref={canvasRef}
        className={`${styles.sectionCanvas} ${isActive ? styles.sectionActive : styles.sectionInactive}`}
        style={{ width, height, ...style, ...(isActive ? { backgroundImage: "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)", backgroundSize: "20px 20px" } : {}) }}
        onClick={() => { if (!isActive) { setEditSection(section); setSelected(null); } }}
      >
        <div className={styles.sectionLabel}>{label}</div>
        {renderSectionBlocks(blocks, section)}
        {isActive && guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
        {isActive && guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
      </div>
    );
  }

  return (
    <main className={styles.main}>
      {showTemplateModal && (
        <div className={styles.modalOverlay} onClick={() => setShowTemplateModal(false)}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <button type="button" className={styles.modalClose} onClick={() => setShowTemplateModal(false)}>×</button>
            <h2 style={{ margin: "0 0 .4rem" }}>Escolha um template</h2>
            <p style={{ margin: "0 0 1rem", color: "#6a7387", fontSize: ".9rem" }}>Selecione um ponto de partida para o seu site.</p>
            <div className={styles.templateGrid}>
              <div className={styles.templateCard} onClick={() => applyTemplate("landing")}>
                <div className={styles.templateEmoji}>🚀</div>
                <strong>Landing Page</strong>
                <p>Template profissional com header, hero, features e footer</p>
              </div>
              <div className={styles.templateCard} onClick={() => applyTemplate("blank")}>
                <div className={styles.templateEmoji}>📄</div>
                <strong>Em branco</strong>
                <p>Canvas vazio para comecar do zero</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <button type="button" onClick={() => navigate("/pages")} style={{ background: "none", border: "1px solid #d2daeb", borderRadius: 8, padding: "6px 12px", cursor: "pointer", fontSize: ".82rem" }}>← Voltar</button>
          <input className={styles.siteTitleInput} value={doc.title} onChange={(e) => setDoc((prev) => ({ ...prev, title: e.target.value }))} />
          <span className={styles.siteIdLabel}>{siteId}</span>
        </div>
        <div className={styles.headerCenter}>
          <div className={styles.editSectionTabs}>
            <button type="button" className={editSection === "page" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("page"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); }}>📄 Pagina</button>
            <button type="button" className={editSection === "header" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("header"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); if (!doc.header.enabled) setDoc((prev) => ({ ...prev, header: { ...prev.header, enabled: true } })); }}>⬆ Header</button>
            <button type="button" className={editSection === "footer" ? styles.editTabActive : styles.editTab} onClick={() => { setEditSection("footer"); setEditingPopup(null); setEditingProductCard(null); setSelected(null); if (!doc.footer.enabled) setDoc((prev) => ({ ...prev, footer: { ...prev.footer, enabled: true } })); }}>⬇ Footer</button>
          </div>
          <div className={styles.zoomControl}>
            {ZOOM_LEVELS.map((z) => (
              <button key={z} type="button" className={zoom === z ? styles.zoomActive : styles.zoomBtn} onClick={() => setZoom(z)}>{Math.round(z * 100)}%</button>
            ))}
          </div>
        </div>
        <div className={styles.headerRight}>
          <button type="button" className={styles.undoBtn} disabled={!canUndo} onClick={undo} title="Desfazer (Ctrl+Z)">↩</button>
          <button type="button" className={styles.undoBtn} disabled={!canRedo} onClick={redo} title="Refazer (Ctrl+Y)">↪</button>
          <button type="button" disabled={saving} onClick={() => save(false)} style={{ background: "#fff", border: "1px solid #d2daeb", borderRadius: 8, padding: "6px 14px", cursor: "pointer", fontSize: ".82rem" }}>💾 Salvar</button>
          <button type="button" disabled={saving} onClick={() => save(true)} style={{ background: "#1a2740", color: "#fff", border: "none", borderRadius: 8, padding: "6px 14px", cursor: "pointer", fontSize: ".82rem" }}>🚀 Publicar</button>
          <a href={`${API_URL}/p/${siteId}`} target="_blank" rel="noreferrer" className={styles.openLive}>Abrir ↗</a>
        </div>
      </div>

      {error && <p className={styles.error}>{error}</p>}
      {success && <div className={styles.successBanner}>{success}</div>}

      <div className={styles.workspace}>
        <LeftPanel
          leftTab={leftTab} setLeftTab={setLeftTab}
          editingProductCard={editingProductCard}
          currentBlocks={currentBlocks}
          addBlock={addBlock}
          selected={selected} setSelected={setSelected}
          removeBlock={removeBlock}
          editSection={editSection} setEditSection={setEditSection}
          setEditingPopup={setEditingPopup} setEditingProductCard={setEditingProductCard}
          doc={doc} setDoc={setDoc}
          activeRoute={activeRoute} setActiveRoute={setActiveRoute}
          routes={routes} setRoutes={setRoutes}
          removePage={removePage} addPage={addPage}
          editingPopup={editingPopup}
          applyTemplate={applyTemplate}
          addPopup={addPopup} removePopup={removePopup}
        />

        <section className={styles.canvasWrap}>
          {editingProductCard && !editingPopup && (
            <div style={{ background: "#0ea5e9", color: "#fff", padding: "8px 16px", fontSize: ".82rem", display: "flex", alignItems: "center", gap: 10, flexWrap: "wrap" }}>
              <strong>🛍 Editando Layout do Product Card</strong>
              <span style={{ opacity: .7, fontSize: ".75rem" }}>Vars: {"{{product_name}}"} {"{{product_price}}"} {"{{product_image}}"} {"{{product_description}}"} {"{{product_short_description}}"} {"{{product_brand}}"} {"{{product_barcode}}"} {"{{cart_count}}"}</span>
              <button type="button" onClick={() => { setEditingProductCard(null); setSelected(null); }} style={{ background: "none", border: "1px solid rgba(255,255,255,.5)", color: "#fff", borderRadius: 6, padding: "2px 10px", cursor: "pointer", fontSize: ".8rem", marginLeft: "auto" }}>← Sair do Product Card</button>
            </div>
          )}
          {editingPopup && (
            <div style={{ background: "#7c3aed", color: "#fff", padding: "8px 16px", fontSize: ".82rem", display: "flex", alignItems: "center", gap: 10 }}>
              <strong>📌 Editando Popup: {doc.popups[editingPopup]?.title}</strong>
              <button type="button" onClick={() => setEditingPopup(null)} style={{ background: "none", border: "1px solid rgba(255,255,255,.5)", color: "#fff", borderRadius: 6, padding: "2px 10px", cursor: "pointer", fontSize: ".8rem" }}>← Sair do Popup</button>
            </div>
          )}

          <div className={styles.canvasScroll}>
            <div className={styles.canvasZoomHost} style={{ zoom }}>
              <div className={styles.canvasStack}>
                {editingPopup ? (
                  <div
                    ref={pageCanvasRef}
                    className={`${styles.sectionCanvas} ${styles.sectionActive}`}
                    style={{ width: doc.popups[editingPopup]?.width || 480, height: doc.popups[editingPopup]?.height || 560, backgroundColor: doc.popups[editingPopup]?.background || "#ffffff", backgroundImage: "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)", backgroundSize: "20px 20px" }}
                    onClick={() => setSelected(null)}
                  >
                    <div className={styles.sectionLabel}>POPUP · {doc.popups[editingPopup]?.title}</div>
                    {renderSectionBlocks(doc.popups[editingPopup]?.blocks || [], "page")}
                    {guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                ) : editingProductCard ? (
                  <div
                    ref={pageCanvasRef}
                    className={`${styles.sectionCanvas} ${styles.sectionActive}`}
                    style={{ width: pcCardBlock?.w || 300, height: pcCardBlock?.h || 280, backgroundColor: pcCardBlock?.style.background || "#ffffff", backgroundImage: "linear-gradient(to right,rgba(152,167,198,.14) 1px,transparent 1px),linear-gradient(to bottom,rgba(152,167,198,.14) 1px,transparent 1px)", backgroundSize: "20px 20px" }}
                    onClick={() => setSelected(null)}
                  >
                    <div className={styles.sectionLabel}>PRODUCT CARD · Layout</div>
                    {renderSectionBlocks(currentBlocks, "page")}
                    {guides.vertical !== null && <span className={styles.guideV} style={{ left: guides.vertical }} />}
                    {guides.horizontal !== null && <span className={styles.guideH} style={{ top: guides.horizontal }} />}
                  </div>
                ) : (
                  <>
                    {doc.header.enabled && (
                      <SectionCanvas canvasRef={headerCanvasRef} section="header" width={doc.canvas.width} height={doc.header.height} style={{ background: doc.header.background }} label="HEADER" blocks={doc.header.blocks} />
                    )}
                    <SectionCanvas canvasRef={pageCanvasRef} section="page" width={doc.canvas.width} height={doc.canvas.height} style={{ backgroundColor: doc.canvas.background }} label="PAGINA" blocks={doc.pages[activeRoute]?.blocks || []} />
                    {doc.footer.enabled && (
                      <SectionCanvas canvasRef={footerCanvasRef} section="footer" width={doc.canvas.width} height={doc.footer.height} style={{ background: doc.footer.background }} label="FOOTER" blocks={doc.footer.blocks} />
                    )}
                  </>
                )}
              </div>
            </div>
          </div>
        </section>

        <Inspector
          selectedBlock={selectedBlock}
          editSection={editSection}
          updateBlock={updateBlock}
          removeBlock={removeBlock}
          routes={routes}
          apis={apis}
          doc={doc}
          editingPopup={editingPopup}
          setEditingProductCard={(id) => { setEditingProductCard(id); setSelected(null); }}
          setEditingProductCardSection={setEditingProductCardSection}
          setError={setError}
        />
      </div>

      {menu && (
        <div className={styles.contextMenu} style={{ left: menu.x, top: menu.y }} onClick={(e) => e.stopPropagation()}>
          <button type="button" onClick={() => contextAction("duplicate")}>📋 Duplicar</button>
          <button type="button" onClick={() => contextAction("bringForward")}>⬆ Trazer para frente</button>
          <button type="button" onClick={() => contextAction("sendBackward")}>⬇ Enviar para tras</button>
          <button type="button" onClick={() => contextAction("width100")}>↔ Largura total</button>
          <button type="button" onClick={() => contextAction("height100")}>↕ Altura total</button>
          <button type="button" onClick={() => contextAction("delete")} style={{ color: "#d13a2f" }}>🗑 Excluir</button>
        </div>
      )}
    </main>
  );
};
