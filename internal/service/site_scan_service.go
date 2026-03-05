package service

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"unicode/utf8"
)

type ScanReport struct {
	Status     string
	Score      int
	Findings   []string
	Errors     []string
	Summary    string
	SourceType string
}

type ScanResult struct {
	Report       ScanReport
	CompiledHTML string
}

type ElementorBlock struct {
	Type            string            `json:"type"`
	Text            string            `json:"text,omitempty"`
	Label           string            `json:"label,omitempty"`
	Href            string            `json:"href,omitempty"`
	ActionType      string            `json:"action_type,omitempty"`
	ActionTarget    string            `json:"action_target,omitempty"`
	APIID           string            `json:"api_id,omitempty"`
	Style           map[string]string `json:"style,omitempty"`
	API             string            `json:"api,omitempty"`
	Src             string            `json:"src,omitempty"`
	ObjectFit       string            `json:"object_fit,omitempty"`
	Images          []string          `json:"images,omitempty"`
	VarName         string            `json:"var_name,omitempty"`
	Placeholder     string            `json:"placeholder,omitempty"`
	ProfileName     string            `json:"profile_name,omitempty"`
	ProfileSubtitle string            `json:"profile_subtitle,omitempty"`
	ProfileImage    string            `json:"profile_image,omitempty"`
	VideoURL        string            `json:"video_url,omitempty"`
	AdminOnly       bool              `json:"admin_only,omitempty"`
	BtnActionType   string            `json:"btn_action_type,omitempty"`
	X               int               `json:"x,omitempty"`
	Y               int               `json:"y,omitempty"`
	W               int               `json:"w,omitempty"`
	H               int               `json:"h,omitempty"`
	Rot             int               `json:"rotation,omitempty"`
	Z               int               `json:"z,omitempty"`
}

type CanvasConfig struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Background string `json:"background,omitempty"`
}

type GlobalSection struct {
	Enabled    bool             `json:"enabled"`
	Height     int              `json:"height,omitempty"`
	Background string           `json:"background,omitempty"`
	Blocks     []ElementorBlock `json:"blocks,omitempty"`
}

type ElementorDoc struct {
	Title  string                   `json:"title"`
	Blocks []ElementorBlock         `json:"blocks"`
	Canvas CanvasConfig             `json:"canvas,omitempty"`
	Header GlobalSection            `json:"header,omitempty"`
	Footer GlobalSection            `json:"footer,omitempty"`
	Pages  map[string]ElementorPage `json:"pages,omitempty"`
}

type ElementorPage struct {
	Title  string           `json:"title"`
	Blocks []ElementorBlock `json:"blocks"`
}

type CodeBundle struct {
	Component string `json:"component"`
	CSS       string `json:"css"`
}

func ScanWebsiteSource(sourceType string, source string) ScanResult {
	normalizedType := strings.ToUpper(strings.TrimSpace(sourceType))
	report := ScanReport{
		Status:     "clean",
		Score:      100,
		Findings:   make([]string, 0),
		Errors:     make([]string, 0),
		SourceType: normalizedType,
	}

	trimmedSource := strings.TrimSpace(source)
	if trimmedSource == "" {
		report.Status = "blocked"
		report.Score = 0
		report.Errors = append(report.Errors, "source vazio")
		report.Summary = "scan bloqueado: sem conteudo"
		return ScanResult{Report: report}
	}

	if !utf8.ValidString(trimmedSource) {
		report.Status = "blocked"
		report.Score = 0
		report.Errors = append(report.Errors, "source invalido (utf-8)")
		report.Summary = "scan bloqueado: encoding invalido"
		return ScanResult{Report: report}
	}

	if len(trimmedSource) > 100*1024 {
		report.Status = "blocked"
		report.Score = 0
		report.Errors = append(report.Errors, "source excede 100kb")
		report.Summary = "scan bloqueado: tamanho acima do limite"
		return ScanResult{Report: report}
	}

	blockedPatterns := []string{
		"<iframe", "javascript:", "onerror=", "onload=",
		"eval(", "new function(", "document.cookie", "child_process",
	}
	warnPatterns := []string{
		"fetch(\"http://", "fetch('http://", "xmlhttprequest", "innerhtml",
		"window.location", "localstorage",
	}

	lowered := strings.ToLower(trimmedSource)
	for _, pattern := range blockedPatterns {
		if strings.Contains(lowered, pattern) {
			report.Status = "blocked"
			report.Score -= 70
			report.Errors = append(report.Errors, fmt.Sprintf("padrao proibido detectado: %s", pattern))
		}
	}
	for _, pattern := range warnPatterns {
		if strings.Contains(lowered, pattern) {
			if report.Status != "blocked" {
				report.Status = "warning"
			}
			report.Score -= 15
			report.Findings = append(report.Findings, fmt.Sprintf("padrao sensivel encontrado: %s", pattern))
		}
	}

	if report.Score < 0 {
		report.Score = 0
	}

	compiled := compileSourceToSafeHTML(normalizedType, trimmedSource)
	if report.Status == "clean" {
		report.Summary = "scan aprovado"
	} else if report.Status == "warning" {
		report.Summary = "scan aprovado com alertas"
	} else {
		report.Summary = "scan bloqueado por risco alto"
	}

	return ScanResult{
		Report:       report,
		CompiledHTML: compiled,
	}
}

func compileSourceToSafeHTML(sourceType string, source string) string {
	switch sourceType {
	case "JXML":
		return compileJXMLToHTML(source)
	case "REACT":
		return compileReactToHTML(source)
	case "SVELTE":
		return compileSvelteToHTML(source)
	case "ELEMENTOR_JSON":
		return compileElementorJSONToHTML(source)
	default:
		escaped := html.EscapeString(source)
		return "<!doctype html><html><head><meta charset=\"utf-8\" /></head><body><main><pre>" + escaped + "</pre></main></body></html>"
	}
}

func compileElementorJSONToHTML(source string) string {
	var doc ElementorDoc
	if err := json.Unmarshal([]byte(source), &doc); err != nil {
		return "<!doctype html><html><head><meta charset=\"utf-8\" /></head><body><main><h1>Template invalido</h1></main></body></html>"
	}

	title := strings.TrimSpace(doc.Title)
	if title == "" {
		title = "Pagina"
	}

	if len(doc.Pages) == 0 {
		doc.Pages = map[string]ElementorPage{
			"/": {
				Title:  title,
				Blocks: doc.Blocks,
			},
		}
	}

	canvasW := 1400
	canvasH := 980
	if doc.Canvas.Width >= 900 && doc.Canvas.Width <= 2800 {
		canvasW = doc.Canvas.Width
	}
	if doc.Canvas.Height >= 700 && doc.Canvas.Height <= 2800 {
		canvasH = doc.Canvas.Height
	}

	canvasBG := doc.Canvas.Background
	if canvasBG == "" {
		canvasBG = "#ffffff"
	}

	headerH := doc.Header.Height
	if headerH < 40 {
		headerH = 80
	}
	headerBG := doc.Header.Background
	if headerBG == "" {
		headerBG = "#1a2740"
	}

	footerH := doc.Footer.Height
	if footerH < 40 {
		footerH = 80
	}
	footerBG := doc.Footer.Background
	if footerBG == "" {
		footerBG = "#1a2740"
	}

	normalizedPages := make(map[string]ElementorPage, len(doc.Pages))
	for rawPath, page := range doc.Pages {
		path := normalizeRoutePath(rawPath)
		if path == "" {
			path = "/"
		}
		normalizedPages[path] = page
	}
	if _, ok := normalizedPages["/"]; !ok {
		normalizedPages["/"] = ElementorPage{Title: title, Blocks: []ElementorBlock{}}
	}

	payload := map[string]interface{}{
		"canvas": map[string]interface{}{
			"width":      canvasW,
			"height":     canvasH,
			"background": canvasBG,
		},
		"header": map[string]interface{}{
			"enabled":    doc.Header.Enabled,
			"height":     headerH,
			"background": headerBG,
			"blocks":     doc.Header.Blocks,
		},
		"footer": map[string]interface{}{
			"enabled":    doc.Footer.Enabled,
			"height":     footerH,
			"background": footerBG,
			"blocks":     doc.Footer.Blocks,
		},
		"pages": normalizedPages,
	}
	payloadJSON, _ := json.Marshal(payload)

	// CSS: reset h1/h4 font-size so applyTextStyle has full control (browser default is 2em which doubles size)
	css := `*{box-sizing:border-box}` +
		`body{font-family:Inter,system-ui,sans-serif;margin:0;padding:0;overflow-x:auto}` +
		`.jx-header,.jx-footer{position:relative;width:100%}` +
		`.canvas{position:relative}` +
		`.blk{position:absolute;overflow:hidden;display:flex;flex-direction:column;justify-content:center;align-items:flex-start}` +
		`.blk h1,.blk h2,.blk h3,.blk h4,.blk h5,.blk h6{font-size:inherit;font-weight:inherit;margin:0;line-height:1.15;width:100%}` +
		`.blk p{font-size:inherit;font-weight:inherit;margin:0;line-height:1.45;width:100%}` +
		`.jx-btn{display:inline-block;padding:10px 18px;border-radius:10px;background:#ff5d1f;color:#fff;text-decoration:none;border:0;cursor:pointer;white-space:nowrap;font-family:inherit}` +
		`.img-blk{width:100%;height:100%;object-fit:cover;display:block}` +
		`.profile-wrap{display:flex;flex-direction:column;gap:6px;width:100%}` +
		`.profile-wrap img{width:56px;height:56px;object-fit:cover;border-radius:999px}` +
		`.product-list{display:flex;flex-direction:column;gap:8px;width:100%;overflow-y:auto}` +
		`.product-item{padding:8px;border:1px solid #dbe2f3;border-radius:8px;font-size:14px}` +
		`.input-var{width:100%;padding:8px 12px;border:1.5px solid #ccd5e8;border-radius:8px;font-size:inherit;font-family:inherit;outline:none}` +
		`.video-wrap{width:100%;height:100%;position:relative}` +
		`.video-wrap iframe{position:absolute;inset:0;width:100%;height:100%;border:0}` +
		`.admin-add{display:none;align-items:center;justify-content:center;width:100%;height:100%;background:#e8f0fe;border:2px dashed #4a90e2;border-radius:10px;cursor:pointer;color:#1a56db;font-weight:600;font-family:inherit;font-size:inherit}` +
		`.admin-add.show{display:flex}`

	// JS renderer — uses var (ES5-compatible) and no template literals to keep Go raw-string escaping simple
	jsBody := `var state={vars:{},products:[]};
function norm(v){if(!v)return '/';var p=String(v).trim();if(p.charAt(0)!='/')p='/'+p;return p||'/';}
function currentPath(){var path=window.location.pathname||'/';var seg=path.split('/').filter(Boolean);if(seg[0]==='p'&&seg.length>=2){var nested='/'+seg.slice(2).join('/');return norm(nested==='/'?'/':nested);}return norm(path);}
function currentPrefix(){var path=window.location.pathname||'/';var seg=path.split('/').filter(Boolean);if(seg[0]==='p'&&seg[1])return '/p/'+seg[1];return '';}
function fillVars(text){return String(text||'').replace(/\{\{\s*([a-zA-Z0-9_]+)\s*\}\}/g,function(_,key){return String(state.vars[key]||'');});}
function applyTextStyle(el,s){if(!s)return;['font-size','font-weight','color','text-align','letter-spacing','line-height','text-decoration'].forEach(function(k){if(s[k])el.style.setProperty(k,s[k]);});}
function applyBoxStyle(el,s){if(!s)return;['background','padding','border-radius','border','border-top','border-right','border-bottom','border-left','opacity'].forEach(function(k){if(s[k])el.style.setProperty(k,s[k]);});}
function buildLayout(b){return{x:Number(b.x||0),y:Number(b.y||0),w:Math.max(10,Number(b.w||220)),h:Math.max(10,Number(b.h||80)),r:Number(b.rotation||0)};}
function rerenderVarTexts(){document.querySelectorAll('[data-var-tmpl]').forEach(function(el){el.textContent=fillVars(el.getAttribute('data-var-tmpl')||'');});}
function fetchProducts(){if(state.products.length>0)return Promise.resolve(state.products);return fetch('/api/store/products').then(function(r){return r.json();}).then(function(d){state.products=(d&&d.data)||[];return state.products;}).catch(function(){return[];});}
function renderBlock(block){
  var t=String((block&&block.type)||'').toLowerCase();
  var l=buildLayout(block);
  var s=block.style||{};
  var node=document.createElement('div');
  node.className='blk';
  node.style.left=l.x+'px';node.style.top=l.y+'px';node.style.width=l.w+'px';node.style.height=l.h+'px';
  node.style.zIndex=String(Number(block.z||0));
  if(l.r)node.style.transform='rotate('+l.r+'deg)';
  applyBoxStyle(node,s);
  if(t==='heading'){
    var h=document.createElement('h1');h.textContent=String(block.text||'');applyTextStyle(h,s);node.appendChild(h);
  }else if(t==='paragraph'){
    var p=document.createElement('p');p.textContent=String(block.text||'');applyTextStyle(p,s);node.appendChild(p);
  }else if(t==='button'){
    var a=document.createElement('a');a.className='jx-btn';a.href=String(block.href||'#');a.textContent=String(block.label||'Botao');
    var aType=String(block.action_type||block.btn_action_type||'navigate');
    var aTarget=String(block.action_target||block.href||'');
    var aAPI=String(block.api_id||'');
    applyTextStyle(a,s);
    if(s['background'])a.style.background=s['background'];if(s['color'])a.style.color=s['color'];if(s['border-radius'])a.style.borderRadius=s['border-radius'];
    a.addEventListener('click',function(e){
      if(aType==='navigate'||aType==='link'){if(aTarget&&aTarget.charAt(0)==='/'){e.preventDefault();window.location.href=currentPrefix()+aTarget;}return;}
      e.preventDefault();if(!aTarget)return;
      var method=(aAPI&&(aAPI.indexOf('login')>=0||aAPI.indexOf('shipping')>=0||aAPI.indexOf('download')>=0))?'POST':'GET';
      fetch(aTarget,{method:method,headers:{'Content-Type':'application/json'},body:method==='POST'?JSON.stringify({}):undefined}).then(function(res){alert(res.ok?'Executado com sucesso':'Erro ao executar');}).catch(function(){alert('Erro de rede');});
    });
    node.appendChild(a);
  }else if(t==='image'){
    var img=document.createElement('img');img.className='img-blk';img.src=String(block.src||'');img.alt='imagem';
    if(block.object_fit)img.style.objectFit=String(block.object_fit);
    node.appendChild(img);
  }else if(t==='carousel'){
    var imgs=Array.isArray(block.images)?block.images.filter(Boolean):[];
    var cimg=document.createElement('img');cimg.className='img-blk';var cidx=0;cimg.src=String(imgs[0]||'');
    if(imgs.length>1)setInterval(function(){cidx=(cidx+1)%imgs.length;cimg.src=String(imgs[cidx]||'');},2600);
    node.appendChild(cimg);
  }else if(t==='input_var'){
    var input=document.createElement('input');input.className='input-var';input.placeholder=String(block.placeholder||'Digite aqui');
    var vkey=String(block.var_name||'var');
    if(s['font-size'])input.style.fontSize=s['font-size'];
    input.addEventListener('input',function(){state.vars[vkey]=input.value;rerenderVarTexts();});
    node.appendChild(input);
  }else if(t==='variable_text'){
    var vp=document.createElement('p');var tmpl=String(block.text||'');vp.setAttribute('data-var-tmpl',tmpl);vp.textContent=fillVars(tmpl);applyTextStyle(vp,s);node.appendChild(vp);
  }else if(t==='profile_card'){
    var pw=document.createElement('div');pw.className='profile-wrap';
    var pimg=document.createElement('img');pimg.src=String(block.profile_image||'');pimg.alt='perfil';
    var ph=document.createElement('h4');ph.textContent=String(block.profile_name||'Usuario');applyTextStyle(ph,s);
    var pp=document.createElement('p');pp.textContent=String(block.profile_subtitle||'');applyTextStyle(pp,s);
    pw.append(pimg,ph,pp);node.appendChild(pw);
  }else if(t==='product_card'){
    var list=document.createElement('div');list.className='product-list';list.textContent='Carregando...';
    fetchProducts().then(function(products){
      if(!products||!products.length){list.textContent='Sem produtos';return;}
      list.innerHTML='';
      products.slice(0,3).forEach(function(item){
        var card=document.createElement('div');card.className='product-item';
        var ttl=document.createElement('strong');ttl.textContent=String(item.name||item.title||'Produto');
        var prc=document.createElement('span');prc.textContent=' - R$ '+String(item.price||'--');
        card.append(ttl,prc);list.appendChild(card);
      });
    });
    node.appendChild(list);
  }else if(t==='divider'){
    node.style.background=s['background']||'#cfd7ea';
  }else if(t==='video'){
    var vurl=String(block.video_url||'');
    if(vurl){
      var vwrap=document.createElement('div');vwrap.className='video-wrap';
      var iframe=document.createElement('iframe');
      var eurl=vurl;
      var ytM=vurl.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})/);if(ytM)eurl='https://www.youtube.com/embed/'+ytM[1]+'?rel=0';
      var vmM=vurl.match(/vimeo\.com\/(\d+)/);if(vmM)eurl='https://player.vimeo.com/video/'+vmM[1];
      iframe.src=eurl;iframe.setAttribute('allow','accelerometer;autoplay;clipboard-write;encrypted-media;gyroscope;picture-in-picture');iframe.setAttribute('allowfullscreen','1');
      vwrap.appendChild(iframe);node.appendChild(vwrap);
    }
  }else if(t==='admin_add_btn'){
    var ab=document.createElement('button');
    ab.className='admin-add'+(localStorage.getItem('jx_admin')==='1'?' show':'');
    ab.textContent=String(block.label||'+ Adicionar');if(s['font-size'])ab.style.fontSize=s['font-size'];
    node.appendChild(ab);
  }
  return node;
}
function render(){
  var path=currentPath();
  var page=DATA.pages[path]||DATA.pages['/']||{blocks:[]};
  var blocks=Array.isArray(page.blocks)?page.blocks:[];
  var hdr=document.getElementById('jx-header');
  var cvs=document.getElementById('canvas');
  var ftr=document.getElementById('jx-footer');
  if(hdr){hdr.innerHTML='';if(DATA.header&&DATA.header.enabled){hdr.style.width=String(DATA.canvas.width)+'px';hdr.style.height=String(DATA.header.height||80)+'px';hdr.style.background=String(DATA.header.background||'#1a2740');var hblocks=Array.isArray(DATA.header.blocks)?DATA.header.blocks:[];hblocks.sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(block){var n=renderBlock(block);if(n)hdr.appendChild(n);});}}
  if(cvs){cvs.innerHTML='';cvs.style.width=String(DATA.canvas.width)+'px';cvs.style.height=String(DATA.canvas.height)+'px';cvs.style.background=String(DATA.canvas.background||'#fff');blocks.sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(block){var n=renderBlock(block);if(n)cvs.appendChild(n);});}
  if(ftr){ftr.innerHTML='';if(DATA.footer&&DATA.footer.enabled){ftr.style.width=String(DATA.canvas.width)+'px';ftr.style.height=String(DATA.footer.height||80)+'px';ftr.style.background=String(DATA.footer.background||'#1a2740');var fblocks=Array.isArray(DATA.footer.blocks)?DATA.footer.blocks:[];fblocks.sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(block){var n=renderBlock(block);if(n)ftr.appendChild(n);});}}
  rerenderVarTexts();
}
render();`

	js := "const DATA=" + string(payloadJSON) + ";" + jsBody

	return "<!doctype html><html><head>" +
		"<meta charset=\"utf-8\" />" +
		"<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" />" +
		"<title>" + html.EscapeString(title) + "</title>" +
		"<style>" + css + "</style>" +
		"</head><body>" +
		"<div id=\"jx-header\" class=\"jx-header\"></div>" +
		"<div id=\"canvas\" class=\"canvas\"></div>" +
		"<div id=\"jx-footer\" class=\"jx-footer\"></div>" +
		"<script>" + js + "</script>" +
		"</body></html>"
}

func compileReactToHTML(source string) string {
	bundle := parseCodeBundle(source)
	componentCode := strings.ReplaceAll(bundle.Component, "export default ", "")
	escapedCode := escapeForScriptTag(componentCode)
	return "<!doctype html><html><head><meta charset=\"utf-8\" /><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" /><title>React Page</title><script crossorigin src=\"https://unpkg.com/react@18/umd/react.development.js\"></script><script crossorigin src=\"https://unpkg.com/react-dom@18/umd/react-dom.development.js\"></script><script src=\"https://unpkg.com/@babel/standalone/babel.min.js\"></script><style>body{margin:0;background:#f4f5f8;font-family:Inter,system-ui,sans-serif}#app{max-width:980px;margin:28px auto;background:#fff;border:1px solid #dadde5;border-radius:16px;padding:22px}" + bundle.CSS + "</style></head><body><div id=\"app\"></div><script type=\"text/babel\">const React = window.React; const ReactDOM = window.ReactDOM;" + escapedCode + "; const __App = (typeof App !== 'undefined') ? App : (() => <main><h1>React sem App</h1></main>); ReactDOM.createRoot(document.getElementById('app')).render(<__App />);</script></body></html>"
}

func compileSvelteToHTML(source string) string {
	bundle := parseCodeBundle(source)
	sfc := strings.TrimSpace(bundle.Component)
	if !strings.Contains(sfc, "<style") && strings.TrimSpace(bundle.CSS) != "" {
		sfc = "<style>\n" + bundle.CSS + "\n</style>\n" + sfc
	}
	escapedSource := jsTemplateLiteralEscape(sfc)
	return "<!doctype html><html><head><meta charset=\"utf-8\" /><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" /><title>Svelte Page</title><style>body{margin:0;background:#f4f5f8;font-family:Inter,system-ui,sans-serif}#app{max-width:980px;margin:28px auto;background:#fff;border:1px solid #dadde5;border-radius:16px;padding:22px}</style></head><body><div id=\"app\"></div><script type=\"module\">import { compile } from 'https://unpkg.com/svelte@4.2.19/compiler.mjs'; try { const source = `" + escapedSource + "`; const result = compile(source, { generate: 'dom', format: 'esm' }); let code = result.js.code; code = code.replaceAll('from \"svelte/internal\"', 'from \"https://unpkg.com/svelte@4.2.19/src/runtime/internal/index.mjs\"'); code = code.replaceAll(\"from 'svelte/internal'\", \"from 'https://unpkg.com/svelte@4.2.19/src/runtime/internal/index.mjs'\"); const blob = new Blob([code], { type: 'text/javascript' }); const url = URL.createObjectURL(blob); const mod = await import(url); const App = mod.default; new App({ target: document.getElementById('app') }); } catch (e) { document.getElementById('app').innerHTML = '<pre style=\"white-space:pre-wrap\">' + String(e).replace(/</g,'&lt;') + '</pre>'; }</script></body></html>"
}

func compileJXMLToHTML(source string) string {
	lines := strings.Split(source, "\n")
	title := "Pagina"
	body := make([]string, 0, len(lines))

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "title "):
			if quoted, ok := extractQuoted(strings.TrimPrefix(line, "title ")); ok {
				title = quoted
			}
		case strings.HasPrefix(line, "h1 "):
			if quoted, ok := extractQuoted(strings.TrimPrefix(line, "h1 ")); ok {
				body = append(body, "<h1>"+html.EscapeString(quoted)+"</h1>")
			}
		case strings.HasPrefix(line, "p "):
			if quoted, ok := extractQuoted(strings.TrimPrefix(line, "p ")); ok {
				body = append(body, "<p>"+html.EscapeString(quoted)+"</p>")
			}
		case strings.HasPrefix(line, "button "):
			rest := strings.TrimSpace(strings.TrimPrefix(line, "button "))
			parts := strings.Split(rest, "->")
			if len(parts) != 2 {
				continue
			}
			label, ok := extractQuoted(strings.TrimSpace(parts[0]))
			if !ok {
				continue
			}
			href, ok := extractQuoted(strings.TrimSpace(parts[1]))
			if !ok {
				continue
			}
			safeHref := html.EscapeString(href)
			body = append(body, "<a class=\"btn\" href=\""+safeHref+"\">"+html.EscapeString(label)+"</a>")
		}
	}

	return "<!doctype html><html><head><meta charset=\"utf-8\" /><title>" + html.EscapeString(title) + "</title><style>body{margin:0;background:#f4f5f8;font-family:Inter,system-ui,sans-serif}main{max-width:980px;margin:28px auto;background:#fff;border:1px solid #dadde5;border-radius:16px;padding:22px}.btn{display:inline-block;padding:12px 16px;border-radius:10px;background:#ff5d1f;color:#fff;text-decoration:none}</style></head><body><main>" + strings.Join(body, "\n") + "</main></body></html>"
}

func extractQuoted(v string) (string, bool) {
	trimmed := strings.TrimSpace(v)
	if len(trimmed) < 2 {
		return "", false
	}
	if !strings.HasPrefix(trimmed, "\"") || !strings.HasSuffix(trimmed, "\"") {
		return "", false
	}
	return strings.Trim(trimmed, "\""), true
}

func escapeForScriptTag(v string) string {
	return strings.ReplaceAll(v, "</script>", "<\\/script>")
}

func jsTemplateLiteralEscape(v string) string {
	out := strings.ReplaceAll(v, "\\", "\\\\")
	out = strings.ReplaceAll(out, "`", "\\`")
	out = strings.ReplaceAll(out, "${", "\\${")
	out = strings.ReplaceAll(out, "</script>", "<\\/script>")
	return out
}

func parseCodeBundle(source string) CodeBundle {
	bundle := CodeBundle{
		Component: source,
		CSS:       "",
	}
	var parsed CodeBundle
	if err := json.Unmarshal([]byte(source), &parsed); err == nil && strings.TrimSpace(parsed.Component) != "" {
		return parsed
	}
	return bundle
}

func inlineStyle(style map[string]string) string {
	if len(style) == 0 {
		return ""
	}
	allowed := []string{
		"color", "background", "padding", "margin", "font-size", "font-weight", "border-radius", "border",
		"border-top", "border-right", "border-bottom", "border-left",
		"text-align",
	}
	parts := make([]string, 0, len(allowed))
	for _, key := range allowed {
		if value, ok := style[key]; ok {
			parts = append(parts, key+":"+html.EscapeString(value))
		}
	}
	return strings.Join(parts, ";")
}

func blockLayoutStyle(block ElementorBlock) string {
	x := block.X
	y := block.Y
	w := block.W
	h := block.H
	if w <= 0 {
		w = 260
	}
	if h <= 0 {
		h = 80
	}
	return fmt.Sprintf("position:absolute;left:%dpx;top:%dpx;width:%dpx;min-height:%dpx;transform:rotate(%ddeg);box-sizing:border-box;", x, y, w, h, block.Rot)
}
