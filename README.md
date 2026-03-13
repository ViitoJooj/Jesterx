<p align="center">
  🇧🇷 Português | <a href="./docs/README_en.md">🇺🇸 English</a> | <a href="./docs/README_cn.md">🇨🇳 CN</a>
</p>

# Jesterx

Jesterx é uma **plataforma SaaS multi-tenant para criação de lojas e páginas online**.  
Cada usuário pode criar e gerenciar sua própria loja, com produtos, pedidos, página pública e painel de administração — tudo hospedado na mesma infraestrutura.  
A versão v1 é aberta e pública no GitHub, servindo tanto como produto funcional quanto como experimento técnico de arquitetura escalável.

### Stack

| Camada | Tecnologia |
|--------|-----------|
| API | Go (net/http, sem frameworks externos) |
| Front-end | React 19 + TypeScript + Vite |
| Banco de dados | PostgreSQL 16 |
| E-mails | Resend |
| Pagamentos | Stripe |
| Contêineres | Docker + docker-compose |

---

## Visão geral

```
┌─────────────────────────────────────────────────────┐
│                    Jesterx Platform                 │
│                                                     │
│  ┌─────────────────┐      ┌───────────────────────┐ │
│  │   React Frontend│      │      Go API           │ │
│  │  jesterx.com.br │◄────►│  api.jesterx.com.br   │ │
│  └─────────────────┘      └──────────┬────────────┘ │
│                                      │              │
│                            ┌─────────▼────────────┐ │
│                            │     PostgreSQL        │ │
│                            └──────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

Cada **website** na plataforma é uma loja/página independente com seus próprios usuários, produtos, pedidos e configurações.

---

## Funcionalidades implementadas

### Autenticação e usuários
- ✅ Cadastro (pessoal e empresarial) com verificação de e-mail
- ✅ Login/logout com JWT (access token + refresh token via cookie)
- ✅ Atualização de perfil (dados pessoais, endereço, bio, redes sociais)
- ✅ Exclusão de conta com período de carência de 30 dias
- ✅ Limpeza automática de usuários não verificados (job em background)

### Website builder
- ✅ Criação de websites com tipos: `ECOMMERCE`, `LANDING_PAGE`, `SOFTWARE_SELL`, `COURSE`, `VIDEO`
- ✅ Gerenciamento de rotas/páginas com posicionamento
- ✅ Versionamento de conteúdo (JXML, React, Svelte, Elementor JSON)
- ✅ Publicação de versões com scan de segurança (score + findings)
- ✅ Renderização pública em `/p/{siteID}`
- ✅ Temas prontos (e-commerce, SaaS, blog, portfólio, perfil)

### E-commerce
- ✅ CRUD completo de produtos (imagens, estoque, SKU, variações, frete)
- ✅ Criação de pedidos com itens, endereço de entrega e método de envio
- ✅ Status de pedido: `pending → paid → shipped → delivered`
- ✅ Cálculo de comissão da plataforma (configurável, padrão 5%)
- ✅ Digest de vendas por e-mail a cada 2 horas para o dono da loja

### Pagamentos (Stripe)
- ✅ Listagem de planos com limites de sites e rotas
- ✅ Checkout via Stripe (link de pagamento)
- ✅ Confirmação manual de checkout por session ID
- ✅ Webhook Stripe (`checkout.session.completed`, `.expired`, etc.)
- ✅ Cancelamento de assinatura

### Social / comunidade
- ✅ Comentários públicos com respostas do dono
- ✅ Sistema de avaliação por estrelas
- ✅ Contagem de visitas à loja
- ✅ Time da loja com roles (`owner`, `manager`, `support`)
- ✅ Perfil público da loja (bio, banner, redes sociais)
- ✅ Marcação de conteúdo adulto (admin)

### Administração
- ✅ Painel admin com estatísticas globais (usuários, sites, pedidos, receita)
- ✅ Listagem de usuários, sites e pedidos (com paginação)
- ✅ Sistema de denúncias (criação pública, moderação por admin)
- ✅ E-mail de resposta ao denunciante

### Upload / armazenamento
- ✅ Upload de arquivos com validação de tipo e tamanho (limite 50 MB)
- ✅ Servidor de arquivos estático em `/files/`
- ✅ Suporte a imagens, vídeos e documentos

### Infraestrutura e segurança
- ✅ Rate limiting global (200 req/s) e por rota (pagamentos: 10, upload: 20)
- ✅ Rate limiting separado para rotas de autenticação (15 req/s)
- ✅ Banimento automático de IP por excesso de requisições
- ✅ Proteção contra path traversal
- ✅ Limite de body size (10 MB geral)
- ✅ Proteção contra paginação abusiva (máx. 100 itens)
- ✅ Cabeçalhos de segurança (HSTS, X-Frame-Options, CSP, XSS-Protection)
- ✅ Recovery de panic com log
- ✅ Graceful shutdown (SIGTERM/SIGINT)
- ✅ Health check em `GET /health`
- ✅ CLI (`jx`) para migrations, seeds, geração de handlers e listagem de rotas
- ✅ Migrations automáticas na inicialização

---

## Configuração

### Pré-requisitos
- Docker e docker-compose
- Conta Stripe (chaves públicas e privadas)
- Conta Resend (chave de API para e-mails)

### Variáveis de ambiente

Copie `.env.example` para `.env` e preencha:

```bash
cp .env.example .env
```

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `POSTGRES_*` | ✅ | Credenciais do banco |
| `JWT_ACCESS_TOKEN` | ✅ | Segredo do access token (use string longa e aleatória) |
| `JWT_REFRESH_TOKEN` | ✅ | Segredo do refresh token (diferente do access) |
| `RESEND_KEY` | ✅ | Chave da API Resend para e-mails |
| `STRIPE_PUBLIC_KEY` | ✅ | Chave pública Stripe |
| `STRIPE_SECRET_KEY` | ✅ | Chave secreta Stripe |
| `FRONTEND_URL` | ✅ em prod | URL do frontend (CORS) |
| `STRIPE_WEBHOOK_SECRET` | Recomendado | Assinar webhooks Stripe |
| `PLATFORM_COMMISSION_PCT` | Opcional | Comissão da plataforma em % (padrão 5) |

### Subir com Docker

```bash
docker compose up -d
```

### Desenvolvimento local

```bash
# Backend
make dev

# Frontend
npm run dev
```

### Migrations

```bash
make migrate
# ou
go run cmd/cli/main.go migrate
```

---

## Avaliação para lançamento como micro-SaaS

### ✅ Pronto
- Core do produto funcional (builder, e-commerce, pagamentos)
- Autenticação segura (JWT, verificação de e-mail)
- Integração com Stripe (checkout, webhook, cancelamento)
- Painel administrativo
- Infraestrutura Docker pronta para deploy
- Segurança básica (rate limit, IP ban, headers, graceful shutdown)

### ⚠️ Recomendado antes do lançamento
- [ ] **Testes automatizados** — não há nenhum teste no repositório; pelo menos testes de unidade para serviços críticos (auth, payment) são essenciais
- [ ] **Planos e preços reais no seed** — ajustar `0006_Seed.up.sql` com os valores e limites definitivos
- [ ] **Domínio e HTTPS** — configurar proxy reverso (Nginx/Caddy/Traefik) com TLS; o servidor Go não serve HTTPS diretamente
- [ ] **Backup do banco** — definir política de backup do volume PostgreSQL
- [ ] **Monitoring/alertas** — integrar Prometheus/Grafana ou Sentry para rastrear erros em produção
- [ ] **`STRIPE_WEBHOOK_SECRET`** — obrigatório em produção para validar webhooks

### 🔮 Funcionalidades futuras sugeridas
- Checkout de pedidos via Stripe (os pedidos atuais não têm cobrança automática ao comprador)
- Painel do lojista com gráficos de vendas
- Domínio personalizado por loja
- Integração com fretes (Correios, Melhor Envio)
- Plano gratuito com limites

---

## Contribuição

1. Faça um fork
2. Crie uma branch (`seu-nome/sua-feature`)
3. Commit suas mudanças
4. Abra um Pull Request

Consulte o [CONTRIBUTING.md](./CONTRIBUTING.md) para mais detalhes.

## Licença

Este projeto está licenciado conforme o arquivo [LICENSE.md](LICENSE.md)
