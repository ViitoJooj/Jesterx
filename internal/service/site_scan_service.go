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
	ProductId       string            `json:"product_id,omitempty"`
	CartItemId      string            `json:"cart_item_id,omitempty"`
	PageSize        int               `json:"page_size,omitempty"`
	PopupId         string            `json:"popup_id,omitempty"`
	EmailVar        string            `json:"email_var,omitempty"`
	PasswordVar     string            `json:"password_var,omitempty"`
	FirstNameVar    string            `json:"first_name_var,omitempty"`
	LastNameVar     string            `json:"last_name_var,omitempty"`
	RegisterMode    bool              `json:"register_mode,omitempty"`
	InputType       string            `json:"input_type,omitempty"`
	VarSrc          string            `json:"var_src,omitempty"`
	InnerBlocks     []ElementorBlock  `json:"inner_blocks,omitempty"`
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
	Title  string                    `json:"title"`
	Blocks []ElementorBlock          `json:"blocks"`
	Canvas CanvasConfig              `json:"canvas,omitempty"`
	Header GlobalSection             `json:"header,omitempty"`
	Footer GlobalSection             `json:"footer,omitempty"`
	Pages  map[string]ElementorPage  `json:"pages,omitempty"`
	Popups map[string]ElementorPopup `json:"popups,omitempty"`
}

type ElementorPage struct {
	Title  string           `json:"title"`
	Blocks []ElementorBlock `json:"blocks"`
}

type ElementorPopup struct {
	Title      string           `json:"title"`
	Blocks     []ElementorBlock `json:"blocks"`
	Width      int              `json:"width,omitempty"`
	Height     int              `json:"height,omitempty"`
	Background string           `json:"background,omitempty"`
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

	normalizedPopups := make(map[string]ElementorPopup, len(doc.Popups))
	for pid, popup := range doc.Popups {
		normalizedPopups[pid] = popup
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
		"pages":  normalizedPages,
		"popups": normalizedPopups,
	}
	payloadJSON, _ := json.Marshal(payload)

	css := `*{box-sizing:border-box}` +
		`html,body{margin:0;padding:0;overflow-x:hidden}` +
		`body{font-family:Inter,system-ui,sans-serif}` +
		`#jx-page{transform-origin:top left;position:relative;line-height:0}` +
		`.jx-header,.jx-footer{position:relative;width:100%}` +
		`.canvas{position:relative}` +
		`.blk{position:absolute;overflow:hidden;display:flex;flex-direction:column;justify-content:center;align-items:flex-start}` +
		`.blk h1,.blk h2,.blk h3,.blk h4,.blk h5,.blk h6{font-size:inherit;font-weight:inherit;margin:0;line-height:1.15;width:100%}` +
		`.blk p{font-size:inherit;font-weight:inherit;margin:0;line-height:1.45;width:100%}` +
		`.jx-btn{display:inline-block;padding:10px 18px;border-radius:10px;background:#8b1e3f;color:#fff;text-decoration:none;border:0;cursor:pointer;white-space:nowrap;font-family:inherit}` +
		`.img-blk{width:100%;height:100%;object-fit:cover;display:block}` +
		`.profile-wrap{display:flex;flex-direction:column;gap:6px;width:100%}` +
		`.profile-wrap img{width:56px;height:56px;object-fit:cover;border-radius:999px}` +
		`.input-var{width:100%;padding:8px 12px;border:1.5px solid #ccd5e8;border-radius:8px;font-size:inherit;font-family:inherit;outline:none}` +
		`.video-wrap{width:100%;height:100%;position:relative}` +
		`.video-wrap iframe{position:absolute;inset:0;width:100%;height:100%;border:0}` +
		`.admin-add{display:flex;align-items:center;justify-content:center;width:100%;height:100%;background:#e8f0fe;border:2px dashed #4a90e2;border-radius:10px;cursor:pointer;color:#1a56db;font-weight:600;font-family:inherit;font-size:inherit}` +
		// toast
		`.jx-toast{position:fixed;bottom:22px;left:50%;transform:translateX(-50%);background:#1a2740;color:#fff;padding:10px 22px;border-radius:8px;z-index:9999;font-family:Inter,system-ui,sans-serif;font-size:14px;pointer-events:none;transition:opacity .3s}` +
		// modal
		`.jx-modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.5);z-index:9990;display:flex;align-items:center;justify-content:center;padding:16px;box-sizing:border-box}` +
		`.jx-modal-box{background:#fff;border-radius:12px;padding:28px;width:360px;max-width:100%;font-family:Inter,system-ui,sans-serif;position:relative}` +
		`.jx-modal-title{margin:0 0 16px;font-size:18px;font-weight:700;color:#1a2740}` +
		`.jx-modal-input{width:100%;padding:10px 12px;border:1.5px solid #ccd5e8;border-radius:8px;font-size:14px;box-sizing:border-box;margin-bottom:10px;font-family:inherit;display:block}` +
		`.jx-modal-btn{width:100%;padding:12px;background:#8b1e3f;color:#fff;border:0;border-radius:8px;font-size:15px;font-weight:600;cursor:pointer;font-family:inherit}` +
		`.jx-modal-btn:disabled{opacity:.6;cursor:not-allowed}` +
		`.jx-modal-err{color:#e53e3e;font-size:13px;margin:8px 0 0;display:none}` +
		`.jx-modal-close{position:absolute;top:10px;right:14px;background:none;border:0;font-size:20px;cursor:pointer;color:#4b5774;line-height:1}` +
		// product card
		`.pc-wrap{display:flex;flex-direction:column;width:100%;height:100%;overflow:hidden}` +
		`.pc-img{width:100%;height:55%;object-fit:cover;border-radius:6px 6px 0 0;display:block;flex-shrink:0}` +
		`.pc-img-ph{width:100%;height:55%;background:#eef1f8;display:flex;align-items:center;justify-content:center;font-size:32px;border-radius:6px 6px 0 0;flex-shrink:0}` +
		`.pc-info{padding:8px;flex:1;overflow:hidden}` +
		`.pc-name{font-weight:700;font-size:14px;color:#1a2740;line-height:1.3;margin-bottom:4px;overflow:hidden}` +
		`.pc-price-row{display:flex;align-items:center;gap:6px;margin-bottom:4px}` +
		`.pc-price{font-weight:700;color:#8b1e3f;font-size:15px}` +
		`.pc-compare{text-decoration:line-through;color:#9aa5bc;font-size:12px}` +
		`.pc-desc{font-size:12px;color:#6a7387;line-height:1.4;overflow:hidden}` +
		`.pc-cart-btn{margin:0 8px 8px;padding:8px;background:#8b1e3f;color:#fff;border:0;border-radius:6px;cursor:pointer;font-family:inherit;font-size:13px;font-weight:600;width:calc(100% - 16px);flex-shrink:0}` +
		`.pc-state{display:flex;align-items:center;justify-content:center;width:100%;height:100%;color:#9aa5bc;font-size:13px}` +
		// product list
		`.plist-wrap{display:flex;flex-direction:column;width:100%;height:100%;overflow:hidden}` +
		`.plist-grid{flex:1;overflow-y:auto;display:grid;grid-template-columns:repeat(auto-fill,minmax(150px,1fr));gap:8px;padding:4px;align-content:start}` +
		`.plist-item{background:#fff;border:1px solid #dbe2f3;border-radius:8px;overflow:hidden;display:flex;flex-direction:column}` +
		`.plist-img{width:100%;height:96px;object-fit:cover;display:block;flex-shrink:0}` +
		`.plist-name{padding:6px 8px 2px;font-size:13px;font-weight:600;color:#1a2740;flex:1}` +
		`.plist-price{padding:0 8px 4px;font-size:12px;color:#8b1e3f;font-weight:700}` +
		`.plist-cart-btn{margin:0 8px 8px;padding:6px;background:#1a2740;color:#fff;border:0;border-radius:6px;cursor:pointer;font-family:inherit;font-size:12px;font-weight:600}` +
		`.plist-pag{display:flex;align-items:center;justify-content:center;gap:10px;padding:8px 0;flex-shrink:0;background:inherit}` +
		`.plist-pag-btn{padding:6px 14px;border:1px solid #ccd5e8;border-radius:6px;background:#fff;cursor:pointer;font-family:inherit;font-size:13px;color:#1a2740}` +
		`.plist-pag-btn:disabled{opacity:.4;cursor:not-allowed}` +
		`.plist-pag-info{font-size:13px;color:#6a7387}` +
		`.plist-state{display:flex;align-items:center;justify-content:center;width:100%;height:100%;color:#9aa5bc;font-size:13px}` +
		// cart items
		`.ci-wrap{display:flex;flex-direction:column;width:100%;height:100%;overflow-y:auto;gap:8px;padding:4px}` +
		`.ci-item{display:flex;align-items:center;gap:8px;background:#fff;border:1px solid #dbe2f3;border-radius:8px;padding:8px;font-family:inherit}` +
		`.ci-name{flex:1;font-size:13px;font-weight:600;color:#1a2740;overflow:hidden}` +
		`.ci-qty-row{display:flex;align-items:center;gap:4px;flex-shrink:0}` +
		`.ci-qty-btn{width:26px;height:26px;border:1px solid #ccd5e8;background:#f4f6fb;border-radius:6px;cursor:pointer;font-size:15px;font-weight:700;color:#1a2740;font-family:inherit;display:flex;align-items:center;justify-content:center;padding:0}` +
		`.ci-qty{min-width:22px;text-align:center;font-size:13px;font-weight:700;color:#1a2740}` +
		`.ci-total{font-size:13px;font-weight:700;color:#8b1e3f;flex-shrink:0;min-width:64px;text-align:right}` +
		`.ci-rm-btn{width:26px;height:26px;border:0;background:none;cursor:pointer;color:#9aa5bc;font-size:18px;line-height:1;flex-shrink:0;padding:0}` +
		`.ci-rm-btn:hover{color:#e53e3e}` +
		`.ci-empty{display:flex;align-items:center;justify-content:center;width:100%;height:100%;color:#9aa5bc;font-size:13px}`

	// JS renderer — uses var (ES5-compatible) and no template literals to keep Go raw-string escaping simple
	jsBody := `var state={vars:{},products:[],storeUser:{email:localStorage.getItem('jx_su_'+getSiteId())||'',id:'',role:'',first_name:'',last_name:'',display_name:'',birth_date:'',gender:'',bio:'',instagram:'',website_url:'',whatsapp:'',phone:'',zip_code:'',address_street:'',address_number:'',address_complement:'',address_district:'',address_city:'',address_state:'',address_country:'',avatar:'',is_admin:false}};
function getSiteId(){var seg=(window.location.pathname||'').split('/').filter(Boolean);return seg[0]==='p'&&seg[1]?seg[1]:'';}
function getStoreToken(){return localStorage.getItem('jx_tk_'+getSiteId())||'';}
function setStoreToken(t){if(t)localStorage.setItem('jx_tk_'+getSiteId(),t);else localStorage.removeItem('jx_tk_'+getSiteId());}
function storeFetch(url,opts){var sid=getSiteId();var headers={'Content-Type':'application/json','X-Website-Id':sid};var tk=getStoreToken();if(tk)headers['Authorization']='Bearer '+tk;var o=opts||{};var merged={credentials:'include',method:o.method||'GET',headers:headers};if(o.body!==undefined)merged.body=o.body;return fetch(url,merged);}
function showToast(msg){var t=document.createElement('div');t.className='jx-toast';t.textContent=String(msg);document.body.appendChild(t);setTimeout(function(){t.style.opacity='0';setTimeout(function(){if(t.parentNode)t.remove();},320);},2200);}
function getCart(){try{return JSON.parse(localStorage.getItem('jx_cart_'+getSiteId())||'[]');}catch(e){return[];}}
function saveCart(c){localStorage.setItem('jx_cart_'+getSiteId(),JSON.stringify(c));}
function addToCart(product){var c=getCart();var ex=null;for(var i=0;i<c.length;i++){if(c[i].id===product.id){ex=c[i];break;}}if(ex){ex.qty=(Number(ex.qty)||1)+1;}else{c.push({id:product.id,name:product.name,price:product.price,qty:1});}saveCart(c);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));showToast(String(product.name||'Produto')+' adicionado ao carrinho!');}
function showLoginModal(){var sid=getSiteId();var ov=document.createElement('div');ov.className='jx-modal-overlay';var box=document.createElement('div');box.className='jx-modal-box';box.innerHTML='<h3 class="jx-modal-title">Entrar na loja</h3><input id="jx-m-em" class="jx-modal-input" type="email" placeholder="Email"><input id="jx-m-pw" class="jx-modal-input" type="password" placeholder="Senha"><button id="jx-m-ok" class="jx-modal-btn">Entrar</button><p id="jx-m-er" class="jx-modal-err"></p><button class="jx-modal-close">&#x2715;</button>';ov.appendChild(box);document.body.appendChild(ov);ov.addEventListener('click',function(e){if(e.target===ov)ov.remove();});box.querySelector('.jx-modal-close').addEventListener('click',function(){ov.remove();});box.querySelector('#jx-m-ok').addEventListener('click',function(){var em=box.querySelector('#jx-m-em').value;var pw=box.querySelector('#jx-m-pw').value;var btn=box.querySelector('#jx-m-ok');var er=box.querySelector('#jx-m-er');btn.disabled=true;btn.textContent='Entrando...';er.style.display='none';storeFetch('/api/v1/auth/login',{method:'POST',body:JSON.stringify({email:em,password:pw})}).then(function(r){return r.text().then(function(txt){var d=null;try{d=JSON.parse(txt);}catch(e){}return{ok:r.ok,d:d,txt:txt};});}).then(function(res){if(res.ok){if(res.d&&res.d.data&&res.d.data.access_token)setStoreToken(res.d.data.access_token);localStorage.setItem('jx_su_'+sid,em);state.storeUser.email=em;ov.remove();showToast('Bem-vindo!');storeFetch('/api/v1/auth/me').then(function(r2){return r2.ok?r2.json():null;}).then(applyUserData).catch(function(){});}else{er.textContent=(res.d&&res.d.message)||res.txt||'Email ou senha incorretos';er.style.display='block';btn.disabled=false;btn.textContent='Entrar';}}).catch(function(){er.textContent='Erro de conexao';er.style.display='block';btn.disabled=false;btn.textContent='Entrar';});});}
function showAddProductModal(){var sid=getSiteId();var ov=document.createElement('div');ov.className='jx-modal-overlay';var box=document.createElement('div');box.className='jx-modal-box';box.style.width='440px';box.innerHTML='<h3 class="jx-modal-title">Novo Produto</h3><input id="jxp-nm" class="jx-modal-input" placeholder="Nome do produto"><input id="jxp-sd" class="jx-modal-input" placeholder="Descricao curta (opcional)"><textarea id="jxp-dc" class="jx-modal-input" rows="2" placeholder="Descricao completa" style="resize:vertical;line-height:1.4"></textarea><input id="jxp-ct" class="jx-modal-input" placeholder="Categoria (opcional)"><input id="jxp-br" class="jx-modal-input" placeholder="Marca (opcional)"><input id="jxp-bc" class="jx-modal-input" placeholder="Codigo de barras (opcional)"><input id="jxp-pr" class="jx-modal-input" type="number" placeholder="Preco (R$)" step="0.01" min="0"><input id="jxp-st" class="jx-modal-input" type="number" placeholder="Estoque" min="0"><input id="jxp-wg" class="jx-modal-input" type="number" placeholder="Peso (g) (opcional)" min="0"><input id="jxp-wd" class="jx-modal-input" type="number" placeholder="Largura (cm) (opcional)" min="0" step="0.01"><input id="jxp-hg" class="jx-modal-input" type="number" placeholder="Altura (cm) (opcional)" min="0" step="0.01"><input id="jxp-ln" class="jx-modal-input" type="number" placeholder="Comprimento (cm) (opcional)" min="0" step="0.01"><input id="jxp-sk" class="jx-modal-input" placeholder="SKU (opcional)"><input id="jxp-tg" class="jx-modal-input" placeholder="Tags (separadas por virgula)"><input id="jxp-im" class="jx-modal-input" placeholder="URL da imagem (opcional)"><label style="display:flex;align-items:center;gap:8px;font-size:13px;margin:6px 0 2px"><input id="jxp-rs" type="checkbox" checked> Requer envio</label><button id="jxp-sv" class="jx-modal-btn">Salvar Produto</button><p id="jxp-er" class="jx-modal-err"></p><button class="jx-modal-close">&#x2715;</button>';ov.appendChild(box);document.body.appendChild(ov);ov.addEventListener('click',function(e){if(e.target===ov)ov.remove();});box.querySelector('.jx-modal-close').addEventListener('click',function(){ov.remove();});box.querySelector('#jxp-sv').addEventListener('click',function(){var nm=box.querySelector('#jxp-nm').value.trim();var sd=box.querySelector('#jxp-sd').value.trim();var dc=box.querySelector('#jxp-dc').value.trim();var ct=box.querySelector('#jxp-ct').value.trim();var br=box.querySelector('#jxp-br').value.trim();var bc=box.querySelector('#jxp-bc').value.trim();var pr=parseFloat(box.querySelector('#jxp-pr').value)||0;var st=parseInt(box.querySelector('#jxp-st').value)||0;var wg=parseInt(box.querySelector('#jxp-wg').value);var wd=parseFloat(box.querySelector('#jxp-wd').value);var hg=parseFloat(box.querySelector('#jxp-hg').value);var ln=parseFloat(box.querySelector('#jxp-ln').value);var sk=box.querySelector('#jxp-sk').value.trim();var tg=box.querySelector('#jxp-tg').value.trim();var im=box.querySelector('#jxp-im').value.trim();var rs=!!box.querySelector('#jxp-rs').checked;var btn=box.querySelector('#jxp-sv');var er=box.querySelector('#jxp-er');if(!nm){er.textContent='Nome obrigatorio';er.style.display='block';return;}btn.disabled=true;btn.textContent='Salvando...';er.style.display='none';var body={name:nm,description:dc,price:pr,stock:st,active:true,images:im?[im]:[],requires_shipping:rs};if(sd)body.short_description=sd;if(ct)body.category=ct;if(br)body.brand=br;if(bc)body.barcode=bc;if(!isNaN(wg))body.weight_grams=wg;if(!isNaN(wd))body.width_cm=wd;if(!isNaN(hg))body.height_cm=hg;if(!isNaN(ln))body.length_cm=ln;if(sk)body.sku=sk;if(tg)body.tags=tg.split(',').map(function(v){return v.trim();}).filter(Boolean);storeFetch('/api/v1/sites/'+sid+'/products',{method:'POST',body:JSON.stringify(body)}).then(function(r){return r.text().then(function(txt){var d=null;try{d=JSON.parse(txt);}catch(e){}return{ok:r.ok,d:d,txt:txt};});}).then(function(res){if(res.ok){state.products=[];ov.remove();showToast('Produto adicionado!');}else{er.textContent=(res.d&&res.d.message)||res.txt||'Erro ao salvar';er.style.display='block';btn.disabled=false;btn.textContent='Salvar Produto';}}).catch(function(){er.textContent='Erro de conexao';er.style.display='block';btn.disabled=false;btn.textContent='Salvar Produto';});});}
function norm(v){if(!v)return '/';var p=String(v).trim();if(p.charAt(0)!='/')p='/'+p;return p||'/';}
function currentPath(){var path=window.location.pathname||'/';var seg=path.split('/').filter(Boolean);if(seg[0]==='p'&&seg.length>=2){var nested='/'+seg.slice(2).join('/');return norm(nested==='/'?'/':nested);}return norm(path);}
function currentPrefix(){var path=window.location.pathname||'/';var seg=path.split('/').filter(Boolean);if(seg[0]==='p'&&seg[1])return '/p/'+seg[1];return '';}
function fillVars(text){return String(text||'').replace(/\{\{\s*([a-zA-Z0-9_]+)\s*\}\}/g,function(_,key){return String(state.vars[key]||'');});}
function applyTextStyle(el,s){if(!s)return;['font-size','font-weight','color','text-align','letter-spacing','line-height','text-decoration'].forEach(function(k){if(s[k])el.style.setProperty(k,s[k]);});}
function applyBoxStyle(el,s){if(!s)return;['background','padding','border-radius','border','border-top','border-right','border-bottom','border-left','opacity'].forEach(function(k){if(s[k])el.style.setProperty(k,s[k]);});}
function buildLayout(b){return{x:Number(b.x||0),y:Number(b.y||0),w:Math.max(10,Number(b.w||220)),h:Math.max(10,Number(b.h||80)),r:Number(b.rotation||0)};}
function rerenderVarTexts(){document.querySelectorAll('[data-var-tmpl]').forEach(function(el){el.textContent=fillVars(el.getAttribute('data-var-tmpl')||'');});}
function rerenderAdminBlocks(){document.querySelectorAll('[data-admin-only]').forEach(function(el){el.style.display=state.storeUser.is_admin?String(el.getAttribute('data-admin-display')||'block'):'none';});}
function rerenderUserAvatars(){document.querySelectorAll('[data-user-avatar]').forEach(function(img){img.src=state.storeUser.avatar||'';});}
function syncUserVars(){var u=state.storeUser;var fn=String(u.first_name||'');var ln=String(u.last_name||'');var dn=String(u.display_name||'');state.vars['user_name']=(dn||((fn&&ln?fn+' '+ln:fn||ln)))||String(u.email||'').split('@')[0]||'';state.vars['user_first_name']=fn;state.vars['user_last_name']=ln;state.vars['user_display_name']=dn;state.vars['user_email']=String(u.email||'');state.vars['user_avatar']=String(u.avatar||'');state.vars['user_birth_date']=String(u.birth_date||'');state.vars['user_gender']=String(u.gender||'');state.vars['user_bio']=String(u.bio||'');state.vars['user_instagram']=String(u.instagram||'');state.vars['user_website_url']=String(u.website_url||'');state.vars['user_whatsapp']=String(u.whatsapp||'');state.vars['user_phone']=String(u.phone||'');state.vars['user_zip_code']=String(u.zip_code||'');state.vars['user_address_street']=String(u.address_street||'');state.vars['user_address_number']=String(u.address_number||'');state.vars['user_address_complement']=String(u.address_complement||'');state.vars['user_address_district']=String(u.address_district||'');state.vars['user_address_city']=String(u.address_city||'');state.vars['user_address_state']=String(u.address_state||'');state.vars['user_address_country']=String(u.address_country||'');var c=getCart();var cnt=0;for(var i=0;i<c.length;i++){cnt+=Number(c[i].qty)||1;}state.vars['cart_count']=String(cnt);rerenderVarTexts();rerenderUserAvatars();}
function applyUserData(me){if(!me)return;state.storeUser.email=String(me.email||state.storeUser.email);state.storeUser.id=String(me.id||'');state.storeUser.role=String(me.role||'');state.storeUser.first_name=String(me.first_name||'');state.storeUser.last_name=String(me.last_name||'');state.storeUser.display_name=String(me.display_name||'');state.storeUser.birth_date=String(me.birth_date||'');state.storeUser.gender=String(me.gender||'');state.storeUser.bio=String(me.bio||'');state.storeUser.instagram=String(me.instagram||'');state.storeUser.website_url=String(me.website_url||'');state.storeUser.whatsapp=String(me.whatsapp||'');state.storeUser.phone=String(me.phone||'');state.storeUser.zip_code=String(me.zip_code||'');state.storeUser.address_street=String(me.address_street||'');state.storeUser.address_number=String(me.address_number||'');state.storeUser.address_complement=String(me.address_complement||'');state.storeUser.address_district=String(me.address_district||'');state.storeUser.address_city=String(me.address_city||'');state.storeUser.address_state=String(me.address_state||'');state.storeUser.address_country=String(me.address_country||'');state.storeUser.avatar=String(me.avatar_url||'');syncUserVars();if(state.storeUser.role==='admin'){state.storeUser.is_admin=true;rerenderAdminBlocks();}else{storeFetch('/api/v1/sites/'+getSiteId()+'/products?limit=1').then(function(r2){state.storeUser.is_admin=(r2.status===200);rerenderAdminBlocks();}).catch(function(){rerenderAdminBlocks();});}}
function fetchProducts(){if(state.products.length>0)return Promise.resolve(state.products);var sid=getSiteId();return storeFetch('/api/store/'+sid+'/products').then(function(r){if(!r.ok)return[];return r.json().then(function(d){state.products=(d&&d.data)||[];return state.products;});}).catch(function(){return[];});}
function fetchProduct(pid){var sid=getSiteId();return storeFetch('/api/store/'+sid+'/products/'+String(pid)).then(function(r){if(!r.ok)return null;return r.json().then(function(d){return d&&d.data?d.data:null;});}).catch(function(){return null;});}
var _popupOverlay=null;
function showPopup(pid){if(_popupOverlay)_popupOverlay.remove();var popup=DATA.popups&&DATA.popups[pid];if(!popup)return;var ov=document.createElement('div');ov.className='jx-modal-overlay';ov.style.zIndex='9980';var box=document.createElement('div');var pw=Math.max(200,Number(popup.width||480));var ph=Math.max(100,Number(popup.height||560));var bg=String(popup.background||'#ffffff');box.style.cssText='position:relative;width:'+pw+'px;max-width:95vw;height:'+ph+'px;background:'+bg+';border-radius:12px;overflow:hidden;';var closeBtn=document.createElement('button');closeBtn.style.cssText='position:absolute;top:8px;right:12px;background:none;border:0;font-size:18px;cursor:pointer;z-index:2;color:#4b5774;line-height:1';closeBtn.innerHTML='&#x2715;';closeBtn.addEventListener('click',function(){if(_popupOverlay)_popupOverlay.remove();_popupOverlay=null;});box.appendChild(closeBtn);var inner=document.createElement('div');inner.style.cssText='position:relative;width:'+pw+'px;height:'+ph+'px;overflow:hidden;';var blocks=Array.isArray(popup.blocks)?popup.blocks:[];blocks.sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(blk){var n=renderBlock(blk);if(n)inner.appendChild(n);});box.appendChild(inner);ov.appendChild(box);ov.addEventListener('click',function(e){if(e.target===ov){ov.remove();_popupOverlay=null;}});document.body.appendChild(ov);_popupOverlay=ov;}
function closePopup(){if(_popupOverlay){_popupOverlay.remove();_popupOverlay=null;}}
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
    var a=document.createElement('a');a.className='jx-btn';a.textContent=String(block.label||'Botao');a.href='#';
    var aType=String(block.action_type||'navigate');
    var aTarget=String(block.action_target||block.href||'/');
    var aAPI=String(block.api_id||'');
    var aPid=String(block.product_id||'');
    applyTextStyle(a,s);
    if(s['background'])a.style.background=s['background'];if(s['color'])a.style.color=s['color'];if(s['border-radius'])a.style.borderRadius=s['border-radius'];
    a.addEventListener('click',function(e){
      e.preventDefault();
      if(aType==='store_login'){var eVar=String(block.email_var||'');var pVar=String(block.password_var||'');if(eVar&&pVar){var emV=String(state.vars[eVar]||'');var pwV=String(state.vars[pVar]||'');var sid2=getSiteId();storeFetch('/api/v1/auth/login',{method:'POST',body:JSON.stringify({email:emV,password:pwV})}).then(function(r){return r.text().then(function(txt){var d=null;try{d=JSON.parse(txt);}catch(e){}return{ok:r.ok,d:d,txt:txt};});}).then(function(res){if(res.ok){if(res.d&&res.d.data&&res.d.data.access_token)setStoreToken(res.d.data.access_token);localStorage.setItem('jx_su_'+sid2,emV);state.storeUser.email=emV;showToast('Bem-vindo!');storeFetch('/api/v1/auth/me').then(function(r2){return r2.ok?r2.json():null;}).then(applyUserData).catch(function(){syncUserVars();});}else{showToast((res.d&&res.d.message)||res.txt||'Email ou senha incorretos');}}).catch(function(){showToast('Erro de conexao');});}else{showLoginModal();}return;}
      if(aType==='store_logout'){setStoreToken(null);localStorage.removeItem('jx_su_'+getSiteId());state.storeUser.email='';state.storeUser.id='';state.storeUser.first_name='';state.storeUser.last_name='';state.storeUser.avatar='';state.storeUser.is_admin=false;syncUserVars();rerenderAdminBlocks();showToast('Saiu da conta!');return;}
      if(aType==='open_popup'){showPopup(String(block.popup_id||''));return;}
      if(aType==='close_popup'){closePopup();return;}
      if(aType==='store_register'){var reVar=String(block.email_var||'');var rpVar=String(block.password_var||'');var fnVar=String(block.first_name_var||'');var lnVar=String(block.last_name_var||'');var emReg=String(state.vars[reVar]||'');var pwReg=String(state.vars[rpVar]||'');var fnReg=String(state.vars[fnVar]||'');var lnReg=String(state.vars[lnVar]||'');var sidR=getSiteId();storeFetch('/api/v1/auth/register',{method:'POST',body:JSON.stringify({email:emReg,password:pwReg,first_name:fnReg,last_name:lnReg})}).then(function(r){return r.text().then(function(txt){var d=null;try{d=JSON.parse(txt);}catch(e){}return{ok:r.ok,d:d,txt:txt};});}).then(function(res){if(res.ok){if(res.d&&res.d.data&&res.d.data.access_token)setStoreToken(res.d.data.access_token);if(emReg)localStorage.setItem('jx_su_'+sidR,emReg);state.storeUser.email=emReg;showToast('Conta criada!');storeFetch('/api/v1/auth/me').then(function(r2){return r2.ok?r2.json():null;}).then(applyUserData).catch(function(){syncUserVars();});}else{showToast((res.d&&res.d.message)||res.txt||'Erro ao registrar');}}).catch(function(){showToast('Erro de conexao');});return;}
      if(aType==='add_to_cart'){var effPid=aPid||String(state.vars['product_id']||'');if(effPid){fetchProduct(effPid).then(function(prod){if(prod)addToCart(prod);else showToast('Produto nao encontrado');});}else{showToast('Configure o ID do produto no botao');}return;}
      if(aType==='add_product'){if(!state.storeUser.is_admin){showToast('Sem permissao. Faca login como admin.');return;}showAddProductModal();return;}
      if(aType==='cart_remove_item'){var crid=String(block.cart_item_id||'');if(crid){var crc=getCart();var crnc=[];for(var cri=0;cri<crc.length;cri++){if(crc[cri].id!==crid)crnc.push(crc[cri]);}saveCart(crnc);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));}return;}
      if(aType==='cart_increment_qty'){var ciid=String(block.cart_item_id||'');if(ciid){var cic=getCart();for(var cii=0;cii<cic.length;cii++){if(cic[cii].id===ciid){cic[cii].qty=(Number(cic[cii].qty)||1)+1;break;}}saveCart(cic);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));}return;}
      if(aType==='cart_decrement_qty'){var cdid=String(block.cart_item_id||'');if(cdid){var cdc=getCart();var cdnc=[];for(var cdi=0;cdi<cdc.length;cdi++){var cdit=cdc[cdi];if(cdit.id===cdid){var cdnq=(Number(cdit.qty)||1)-1;if(cdnq>0)cdnc.push({id:cdit.id,name:cdit.name,price:cdit.price,qty:cdnq});}else{cdnc.push(cdit);}}saveCart(cdnc);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));}return;}
      if(aType==='navigate'||aType==='link'){if(aTarget&&aTarget.charAt(0)==='/'){window.location.href=currentPrefix()+aTarget;}return;}
      if(!aTarget)return;
      storeFetch(aTarget,{method:'POST',body:JSON.stringify({})}).then(function(res){showToast(res.ok?'Executado com sucesso':'Erro ao executar');}).catch(function(){showToast('Erro de rede');});
    });
    node.appendChild(a);
    if(block.admin_only){node.setAttribute('data-admin-only','1');node.setAttribute('data-admin-display','block');node.style.display='none';}
  }else if(t==='image'){
    var img=document.createElement('img');img.className='img-blk';var imgSrc=block.var_src?fillVars(String(block.var_src)):String(block.src||'');img.src=imgSrc;img.alt='imagem';
    if(block.object_fit)img.style.objectFit=String(block.object_fit);
    node.appendChild(img);
  }else if(t==='carousel'){
    var imgs=Array.isArray(block.images)?block.images.filter(Boolean):[];
    var cimg=document.createElement('img');cimg.className='img-blk';var cidx=0;cimg.src=String(imgs[0]||'');
    if(imgs.length>1)setInterval(function(){cidx=(cidx+1)%imgs.length;cimg.src=String(imgs[cidx]||'');},2600);
    node.appendChild(cimg);
  }else if(t==='input_var'){
    var input=document.createElement('input');input.className='input-var';input.placeholder=String(block.placeholder||'Digite aqui');input.type=String(block.input_type||'text');
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
    var pcw=document.createElement('div');pcw.className='pc-wrap';
    var pid=String(block.product_id||'');
    function renderPCard(prod){
      pcw.innerHTML='';
      if(!prod){pcw.innerHTML='<div class="pc-state">Produto nao encontrado</div>';return;}
      if(Array.isArray(block.inner_blocks)&&block.inner_blocks.length>0){var sVars={};['product_name','product_price','product_compare_price','product_image','product_description','product_short_description','product_sku','product_id','product_brand','product_model','product_barcode','product_weight','product_color','product_size','product_material','product_category','product_slug','product_condition','product_origin_country','product_tags','product_attributes'].forEach(function(k){sVars[k]=state.vars[k];});state.vars['product_name']=String(prod.name||'');state.vars['product_price']=prod.price!=null?'R$ '+Number(prod.price).toFixed(2):'';state.vars['product_compare_price']=prod.compare_price?'R$ '+Number(prod.compare_price).toFixed(2):'';state.vars['product_image']=String((prod.images&&prod.images[0])||'');state.vars['product_description']=String(prod.description||'');state.vars['product_short_description']=String(prod.short_description||'');state.vars['product_sku']=String(prod.sku||'');state.vars['product_id']=String(prod.id||'');state.vars['product_brand']=String(prod.brand||'');state.vars['product_model']=String(prod.model||'');state.vars['product_barcode']=String(prod.barcode||'');state.vars['product_weight']=prod.weight_grams!=null?String(prod.weight_grams):'';state.vars['product_color']=String(prod.color||'');state.vars['product_size']=String(prod.size||'');state.vars['product_material']=String(prod.material||'');state.vars['product_category']=String(prod.category||'');state.vars['product_slug']=String(prod.slug||'');state.vars['product_condition']=String(prod.condition||'');state.vars['product_origin_country']=String(prod.origin_country||'');state.vars['product_tags']=Array.isArray(prod.tags)?prod.tags.join(', '):'';state.vars['product_attributes']=prod.attributes?JSON.stringify(prod.attributes):'';var pcIn=document.createElement('div');pcIn.style.cssText='position:relative;width:100%;height:100%;overflow:hidden;';block.inner_blocks.slice().sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(ib){var renderedIb=ib;if(String((ib.type||'')).toLowerCase()==='button'&&String(ib.action_type||'')==='add_to_cart'&&!ib.product_id){renderedIb=Object.assign({},ib,{product_id:prod.id});}var n=renderBlock(renderedIb);if(n)pcIn.appendChild(n);});Object.keys(sVars).forEach(function(k){state.vars[k]=sVars[k];});pcw.style.cssText='position:relative;width:100%;height:100%;overflow:hidden;';pcw.appendChild(pcIn);return;}
      var pimgs=Array.isArray(prod.images)?prod.images:[];
      if(pimgs.length){var pci=document.createElement('img');pci.className='pc-img';pci.src=String(pimgs[0]);pci.alt=String(prod.name||'');pcw.appendChild(pci);}
      else{var pph=document.createElement('div');pph.className='pc-img-ph';pph.textContent='📦';pcw.appendChild(pph);}
      var pinf=document.createElement('div');pinf.className='pc-info';
      var pnm=document.createElement('div');pnm.className='pc-name';pnm.textContent=String(prod.name||'Produto');
      var ppr=document.createElement('div');ppr.className='pc-price-row';
      if(prod.compare_price){var pcmp=document.createElement('span');pcmp.className='pc-compare';pcmp.textContent='R$ '+Number(prod.compare_price).toFixed(2);ppr.appendChild(pcmp);}
      var ppc=document.createElement('span');ppc.className='pc-price';ppc.textContent='R$ '+Number(prod.price||0).toFixed(2);ppr.appendChild(ppc);
      var pdc=document.createElement('div');pdc.className='pc-desc';pdc.textContent=String(prod.description||'');
      pinf.append(pnm,ppr,pdc);pcw.appendChild(pinf);
      var pcb=document.createElement('button');pcb.className='pc-cart-btn';pcb.textContent='Adicionar ao carrinho';
      pcb.addEventListener('click',function(){addToCart(prod);});
      pcw.appendChild(pcb);
    }
    pcw.innerHTML='<div class="pc-state">Carregando...</div>';
    if(pid){fetchProduct(pid).then(renderPCard);}else{fetchProducts().then(function(prods){renderPCard(prods&&prods.length?prods[0]:null);});}
    node.appendChild(pcw);
  }else if(t==='product_list'){
    var pgSize=Math.max(1,Number(block.page_size||6));
    var curPage=[0];
    var plw=document.createElement('div');plw.className='plist-wrap';
    plw.innerHTML='<div class="plist-state">Carregando produtos...</div>';
    var plHasInner=Array.isArray(block.inner_blocks)&&block.inner_blocks.length>0;
    fetchProducts().then(function(prods){
      if(!prods||!prods.length){plw.innerHTML='<div class="plist-state">Nenhum produto disponivel</div>';return;}
      var total=Math.ceil(prods.length/pgSize);
      function renderPage(pg){
        var items=prods.slice(pg*pgSize,(pg+1)*pgSize);
        plw.innerHTML='';
        var grid=document.createElement('div');grid.className='plist-grid';
        items.forEach(function(prod){
          var item=document.createElement('div');item.className='plist-item';
          if(plHasInner){
            var plSv={};['product_name','product_price','product_compare_price','product_image','product_description','product_short_description','product_sku','product_id','product_brand','product_model','product_barcode','product_weight','product_color','product_size','product_material','product_category','product_slug','product_condition','product_origin_country','product_tags','product_attributes'].forEach(function(k){plSv[k]=state.vars[k];});
            state.vars['product_name']=String(prod.name||'');
            state.vars['product_price']=prod.price!=null?'R$ '+Number(prod.price).toFixed(2):'';
            state.vars['product_compare_price']=prod.compare_price?'R$ '+Number(prod.compare_price).toFixed(2):'';
            state.vars['product_image']=String((prod.images&&prod.images[0])||'');
            state.vars['product_description']=String(prod.description||'');
            state.vars['product_short_description']=String(prod.short_description||'');
            state.vars['product_sku']=String(prod.sku||'');
            state.vars['product_id']=String(prod.id||'');
            state.vars['product_brand']=String(prod.brand||'');
            state.vars['product_model']=String(prod.model||'');
            state.vars['product_barcode']=String(prod.barcode||'');
            state.vars['product_weight']=prod.weight_grams!=null?String(prod.weight_grams):'';
            state.vars['product_color']=String(prod.color||'');
            state.vars['product_size']=String(prod.size||'');
            state.vars['product_material']=String(prod.material||'');
            state.vars['product_category']=String(prod.category||'');
            state.vars['product_slug']=String(prod.slug||'');
            state.vars['product_condition']=String(prod.condition||'');
            state.vars['product_origin_country']=String(prod.origin_country||'');
            state.vars['product_tags']=Array.isArray(prod.tags)?prod.tags.join(', '):'';
            state.vars['product_attributes']=prod.attributes?JSON.stringify(prod.attributes):'';
            var plIn=document.createElement('div');plIn.style.cssText='position:relative;width:100%;height:100%;overflow:hidden;';
            block.inner_blocks.slice().sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(ib){
              var rib=ib;
              if(String((ib.type||'')).toLowerCase()==='button'&&String(ib.action_type||'')==='add_to_cart'&&!ib.product_id){rib=Object.assign({},ib,{product_id:prod.id});}
              var n=renderBlock(rib);if(n)plIn.appendChild(n);
            });
            Object.keys(plSv).forEach(function(k){state.vars[k]=plSv[k];});
            item.style.cssText='position:relative;overflow:hidden;';
            item.appendChild(plIn);
          }else{
            var pimgs=Array.isArray(prod.images)?prod.images:[];
            if(pimgs.length){var limg=document.createElement('img');limg.className='plist-img';limg.src=String(pimgs[0]);limg.alt=String(prod.name||'');item.appendChild(limg);}
            var lnm=document.createElement('div');lnm.className='plist-name';lnm.textContent=String(prod.name||'Produto');
            var lpr=document.createElement('div');lpr.className='plist-price';lpr.textContent='R$ '+Number(prod.price||0).toFixed(2);
            var lcb=document.createElement('button');lcb.className='plist-cart-btn';lcb.textContent='+ Carrinho';
            (function(pr){lcb.addEventListener('click',function(){addToCart(pr);});})(prod);
            item.append(lnm,lpr,lcb);
          }
          grid.appendChild(item);
        });
        plw.appendChild(grid);
        if(total>1){
          var pag=document.createElement('div');pag.className='plist-pag';
          var prev=document.createElement('button');prev.className='plist-pag-btn';prev.textContent='← Anterior';prev.disabled=pg===0;
          var info=document.createElement('span');info.className='plist-pag-info';info.textContent=(pg+1)+' / '+total;
          var next=document.createElement('button');next.className='plist-pag-btn';next.textContent='Proximo →';next.disabled=pg>=total-1;
          prev.addEventListener('click',function(){if(curPage[0]>0){curPage[0]--;renderPage(curPage[0]);}});
          next.addEventListener('click',function(){if(curPage[0]<total-1){curPage[0]++;renderPage(curPage[0]);}});
          pag.append(prev,info,next);plw.appendChild(pag);
        }
      }
      renderPage(0);
    });
    node.appendChild(plw);
  }else if(t==='cart_items'){
    var ciw=document.createElement('div');ciw.className='ci-wrap';
    var ciHasInner=Array.isArray(block.inner_blocks)&&block.inner_blocks.length>0;
    function renderCartList(){
      ciw.innerHTML='';
      var cartItems=getCart();
      if(!cartItems.length){ciw.innerHTML='<div class="ci-empty">Carrinho vazio</div>';return;}
      cartItems.forEach(function(item){
        if(ciHasInner){
          var ciSv={};['cart_item_name','cart_item_price','cart_item_qty','cart_item_total','cart_item_id'].forEach(function(k){ciSv[k]=state.vars[k];});
          state.vars['cart_item_name']=String(item.name||'');
          state.vars['cart_item_price']='R$ '+Number(item.price||0).toFixed(2);
          state.vars['cart_item_qty']=String(Number(item.qty)||1);
          state.vars['cart_item_total']='R$ '+(Number(item.price||0)*(Number(item.qty)||1)).toFixed(2);
          state.vars['cart_item_id']=String(item.id||'');
          var ciIn=document.createElement('div');ciIn.style.cssText='position:relative;width:100%;overflow:hidden;';
          block.inner_blocks.slice().sort(function(a,b){return Number(a.z||0)-Number(b.z||0);}).forEach(function(ib){
            var rib=ib;
            if(String((ib.type||'')).toLowerCase()==='button'){
              var ibat=String(ib.action_type||'');
              if((ibat==='cart_remove_item'||ibat==='cart_increment_qty'||ibat==='cart_decrement_qty')&&!ib.cart_item_id){rib=Object.assign({},ib,{cart_item_id:item.id});}
            }
            var n=renderBlock(rib);if(n)ciIn.appendChild(n);
          });
          Object.keys(ciSv).forEach(function(k){state.vars[k]=ciSv[k];});
          ciw.appendChild(ciIn);
        }else{
          var ci=document.createElement('div');ci.className='ci-item';
          var ciNm=document.createElement('div');ciNm.className='ci-name';ciNm.textContent=String(item.name||'Produto');
          var ciQr=document.createElement('div');ciQr.className='ci-qty-row';
          var ciDec=document.createElement('button');ciDec.className='ci-qty-btn';ciDec.textContent='−';
          var ciQt=document.createElement('span');ciQt.className='ci-qty';ciQt.textContent=String(Number(item.qty)||1);
          var ciInc=document.createElement('button');ciInc.className='ci-qty-btn';ciInc.textContent='+';
          var ciTot=document.createElement('div');ciTot.className='ci-total';ciTot.textContent='R$ '+(Number(item.price||0)*(Number(item.qty)||1)).toFixed(2);
          var ciRm=document.createElement('button');ciRm.className='ci-rm-btn';ciRm.textContent='×';
          (function(it){
            ciDec.addEventListener('click',function(){var c=getCart();var nc=[];for(var i=0;i<c.length;i++){var x=c[i];if(x.id===it.id){var nq=(Number(x.qty)||1)-1;if(nq>0)nc.push({id:x.id,name:x.name,price:x.price,qty:nq});}else{nc.push(x);}}saveCart(nc);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));});
            ciInc.addEventListener('click',function(){var c=getCart();for(var i=0;i<c.length;i++){if(c[i].id===it.id){c[i].qty=(Number(c[i].qty)||1)+1;break;}}saveCart(c);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));});
            ciRm.addEventListener('click',function(){var c=getCart();var nc=[];for(var i=0;i<c.length;i++){if(c[i].id!==it.id)nc.push(c[i]);}saveCart(nc);syncUserVars();document.dispatchEvent(new CustomEvent('jx:cartupdate'));});
          })(item);
          ciQr.append(ciDec,ciQt,ciInc);ci.append(ciNm,ciQr,ciTot,ciRm);ciw.appendChild(ci);
        }
      });
    }
    renderCartList();
    document.addEventListener('jx:cartupdate',renderCartList);
    node.appendChild(ciw);
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
  }else if(t==='user_avatar'){
    var uaImg=document.createElement('img');uaImg.className='img-blk';uaImg.setAttribute('data-user-avatar','1');uaImg.src=state.storeUser.avatar||'';uaImg.alt='avatar';if(block.object_fit)uaImg.style.objectFit=String(block.object_fit);node.appendChild(uaImg);
  }else if(t==='admin_add_btn'){
    var ab=document.createElement('button');
    ab.className='admin-add';
    ab.textContent=String(block.label||'+ Adicionar');if(s['font-size'])ab.style.fontSize=s['font-size'];
    var bact=String(block.btn_action_type||'add_product');
    ab.addEventListener('click',function(){if(bact==='add_product')showAddProductModal();else showToast('Acao nao configurada');});
    node.appendChild(ab);
    node.setAttribute('data-admin-only','1');node.setAttribute('data-admin-display','block');node.style.display='none';
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
  rerenderVarTexts();syncUserVars();
}
render();
(function(){storeFetch('/api/v1/auth/me').then(function(r){if(!r.ok){rerenderAdminBlocks();return null;}return r.json();}).then(function(me){if(!me||!me.id){rerenderAdminBlocks();return;}applyUserData(me);}).catch(function(){rerenderAdminBlocks();});})();
function scaleToViewport(){var dw=DATA.canvas.width;if(!dw||dw<=0)return;var s=window.innerWidth/dw;var pg=document.getElementById('jx-page');if(!pg)return;pg.style.transform='scale('+s+')';document.body.style.height=pg.offsetHeight*s+'px';}
scaleToViewport();
window.addEventListener('resize',scaleToViewport);`

	js := "const DATA=" + string(payloadJSON) + ";" + jsBody

	return "<!doctype html><html><head>" +
		"<meta charset=\"utf-8\" />" +
		"<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" />" +
		"<title>" + html.EscapeString(title) + "</title>" +
		"<style>" + css + "</style>" +
		"</head><body>" +
		"<div id=\"jx-page\">" +
		"<div id=\"jx-header\" class=\"jx-header\"></div>" +
		"<div id=\"canvas\" class=\"canvas\"></div>" +
		"<div id=\"jx-footer\" class=\"jx-footer\"></div>" +
		"</div>" +
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

	return "<!doctype html><html><head><meta charset=\"utf-8\" /><title>" + html.EscapeString(title) + "</title><style>body{margin:0;background:#f4f5f8;font-family:Inter,system-ui,sans-serif}main{max-width:980px;margin:28px auto;background:#fff;border:1px solid #dadde5;border-radius:16px;padding:22px}.btn{display:inline-block;padding:12px 16px;border-radius:10px;background:#8b1e3f;color:#fff;text-decoration:none}</style></head><body><main>" + strings.Join(body, "\n") + "</main></body></html>"
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

