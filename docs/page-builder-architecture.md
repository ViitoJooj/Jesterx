# Arquitetura Proposta: Builder de Paginas Multi-Rota

## Objetivo
Permitir que cada usuario com plano ativo crie paginas com varias rotas, usando:
- editor visual (Elementor-like),
- codigo React,
- codigo Svelte.

Com pipeline de seguranca antes da publicacao e templates prontos para:
- `ECOMMERCE`
- `LANDING_PAGE`
- `SOFTWARE_SELL`
- `COURSE`

## Regras de negocio principais
1. Usuario sem plano ativo nao cria pagina.
2. Limite de rotas por plano (exemplo):
- Starter: 8 rotas
- Pro: 30 rotas
- Enterprise: 100 rotas
3. Toda versao publicada passa por scans obrigatorios.
4. Publicacao bloqueada se algum scan critico falhar.

## DSL propria (simples) para render
Criar uma linguagem declarativa minima chamada `JXML`, compilada para HTML seguro.

Exemplo:
```jxml
page "/"
title "Minha landing"

section hero {
  h1 "Venda mais em 7 dias"
  p "Checkout rapido com frete automatico"
  button "Comecar" -> "/checkout"
}
```

Fluxo:
1. Parser -> AST.
2. Validador semantico (tags, atributos, binds, limites).
3. Sanitizador (remove scripts inline perigosos).
4. Renderer AST -> HTML SSR.
5. Cache/CDN.

Observacao: React/Svelte podem virar "modo avancado" e gerar AST intermediaria igual, para reaproveitar o mesmo pipeline de seguranca e deploy.

## Pipeline de seguranca (scans)
1. Scan de assinatura: regex + hash para payloads conhecidos.
2. Scan AST: bloquear `eval`, `Function`, `document.cookie`, redirects suspeitos.
3. Scan de dependencias (quando houver React/Svelte): SCA (CVE/licencas).
4. Scan de comportamento: executar build/render em sandbox sem rede e monitorar syscalls.
5. Politica de saida:
- sem risco: publica;
- risco medio: quarentena para revisao;
- risco alto: bloqueia e notifica.

## Modelo MongoDB sugerido
Colecoes:
1. `sites`
- `_id`, `owner_id`, `plan_snapshot`, `template_type`, `editor_mode`, `status`, `created_at`
2. `site_routes`
- `_id`, `site_id`, `path`, `layout_id`, `visibility`, `order`
3. `site_versions`
- `_id`, `site_id`, `version`, `source_type` (`JXML|REACT|SVELTE|ELEMENTOR_JSON`), `source`, `compiled_html`, `scan_report`, `published_at`
4. `templates`
- `_id`, `type`, `name`, `schema_version`, `default_routes`, `blocks`
5. `api_tokens`
- `_id`, `site_id`, `scopes`, `token_hash`, `expires_at`

Indices:
- `sites.owner_id`
- `site_routes.site_id + path` (unique)
- `site_versions.site_id + version` (unique)

## APIs de produto (MVP)
1. `POST /api/v1/sites` cria site (valida plano ativo).
2. `POST /api/v1/sites/:id/routes` cria/edita rotas.
3. `POST /api/v1/sites/:id/versions` envia codigo/DSL para scan.
4. `POST /api/v1/sites/:id/publish/:version` publica versao aprovada.
5. `GET /api/v1/sites/:id/scan-reports/:version` detalhes de seguranca.

APIs de negocio para as lojas:
1. `GET /api/store/products`
2. `POST /api/store/login`
3. `POST /api/store/shipping/quote`
4. `POST /api/software/download-token`
5. `GET /api/course/modules`

## Roadmap de entrega
1. Fase 1: templates + limite por plano + rotas + scan basico + publish manual.
2. Fase 2: DSL `JXML` + compilador + cache por versao.
3. Fase 3: modo React/Svelte com build sandboxado.
4. Fase 4: marketplace de templates e plugins.
