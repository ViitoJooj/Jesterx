# Checklist de Lançamento — Jesterx micro-SaaS

Documento com **cada alteração necessária** antes de lançar o Jesterx como produto pago.  
Organizado por prioridade: 🔴 Crítico → 🟠 Alta → 🟡 Média → 🟢 Baixa.

---

## 🔴 Crítico — Bloqueia o lançamento

### 1. Rate limiter vulnerável a IP spoofing

**Arquivo:** `pkg/ratelimit/ratelimit.go`  
**Problema:** A função `realIP()` confia cegamente no header `X-Forwarded-For` enviado pelo cliente. Qualquer usuário pode enviar `X-Forwarded-For: 1.2.3.4` e contornar completamente o rate limiting.

**O que mudar:**

O problema está na função `realIP()` no final do arquivo. Ela lê o header `X-Forwarded-For` diretamente do cliente sem verificar se veio de um proxy confiável. A solução mais segura é fazer o proxy **sobrescrever** o header antes de chegar no Go — assim o código não precisa ser alterado.

```nginx
# nginx.conf — bloco location da API:
location / {
    proxy_pass       http://localhost:8080;
    # Sobrescreve X-Forwarded-For com o IP real do cliente (não o que ele enviou)
    proxy_set_header X-Forwarded-For $remote_addr;
    proxy_set_header X-Real-IP       $remote_addr;
    proxy_set_header Host            $host;
}
```

Se preferir corrigir no código sem proxy, substitua a função `realIP` completa em `pkg/ratelimit/ratelimit.go`:

```go
// pkg/ratelimit/ratelimit.go — substitua a função realIP inteira por:
var trustedProxyCIDRs = []string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12"}

func realIP(r *http.Request) string {
    remoteHost, _, _ := net.SplitHostPort(r.RemoteAddr)
    for _, cidr := range trustedProxyCIDRs {
        _, network, err := net.ParseCIDR(cidr)
        if err != nil {
            ip := net.ParseIP(cidr)
            if ip != nil && ip.Equal(net.ParseIP(remoteHost)) {
                goto trusted
            }
            continue
        }
        if network.Contains(net.ParseIP(remoteHost)) {
            goto trusted
        }
    }
    return remoteHost
trusted:
    if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
        return strings.TrimSpace(strings.SplitN(fwd, ",", 2)[0])
    }
    return remoteHost
}
```

---

### 2. STRIPE_WEBHOOK_SECRET obrigatório em produção

**Arquivo:** `internal/service/payment_service.go`  
**Problema:** Se `STRIPE_WEBHOOK_SECRET` estiver vazio, o webhook aceita qualquer payload sem validar a assinatura. Qualquer pessoa pode simular um `checkout.session.completed` falso e ativar planos sem pagar.

**O que mudar:**

```go
// internal/service/payment_service.go — ProcessStripeWebhook
func (s *PaymentService) ProcessStripeWebhook(rawBody []byte, signature string) error {
    // MUDAR: tornar obrigatório em produção
    if config.StripeWebhookSecret == "" {
        return errors.New("webhook secret not configured")  // bloquear em prod
    }
    if !validateStripeSignature(rawBody, signature, config.StripeWebhookSecret) {
        return errors.New("invalid stripe signature")
    }
    ...
}
```

**E no `.env` de produção:**
```bash
STRIPE_WEBHOOK_SECRET=whsec_XXXXXXXXXXXXXXXXXX  # obrigatório
```

Para obter o secret: Stripe Dashboard → Developers → Webhooks → seu endpoint → Signing secret.

---

### 3. Remover Redis do docker-compose (não está sendo usado)

**Arquivo:** `docker-compose.yml`  
**Problema:** O serviço `redis` está declarado no docker-compose e no docker-compose.test.yml, e o container roda em produção consumindo memória, mas **nenhuma linha de código Go conecta ao Redis**.

**O que fazer:**  
- **Opção A (simples):** Remova o serviço `redis` e o volume `redis_data` do `docker-compose.yml`. Remova também `REDIS_URL` do env da api e de `internal/config/dotenv.go`.  
- **Opção B (se você planeja usar Redis):** Implemente o cache/sessão que justifica o Redis antes do lançamento.

```yaml
# docker-compose.yml — REMOVER o bloco abaixo:
  redis:
    image: redis:7-alpine
    ...
    volumes:
      - redis_data:/data
    ...

# E remover redis_data de volumes:
volumes:
  postgres_data:
  redis_data:  # ← remover

# E remover da seção api:
      - REDIS_URL=redis://redis:6379
      ...
    depends_on:
      redis:           # ← remover
        condition: service_healthy
```

```go
// internal/config/dotenv.go — REMOVER:
var RedisURL string       // linha 19
RedisURL = getEnvOrDefault("REDIS_URL", "redis://localhost:6379")  // linha 50
```

---

### 4. Criar usuário administrador

**Problema:** Não existe nenhum mecanismo para criar um usuário com `role = 'admin'` no banco. Sem isso, as rotas `/api/v1/admin/*` são inacessíveis.

**O que fazer:** Após o primeiro `docker compose up`, registre sua conta pela tela normal de cadastro no frontend. Em seguida, execute no PostgreSQL:

```sql
-- Substitua o e-mail pelo que você usou no cadastro
UPDATE users
SET role = 'admin'
WHERE email = 'seu@email.com'
  AND website_id = '00000000-0000-0000-0000-000000000001';

-- Verifique:
SELECT id, email, role FROM users WHERE email = 'seu@email.com';
```

**Para executar no Docker:**
```bash
docker exec -it jesterx_postgres psql -U postgres -d jesterx \
  -c "UPDATE users SET role='admin' WHERE email='seu@email.com' AND website_id='00000000-0000-0000-0000-000000000001';"
```

> ⚠️ Não crie migrations com e-mail ou senha hardcoded — use sempre o fluxo acima de forma manual no servidor de produção.

---

## 🟠 Alta — Deve ser feito antes do lançamento

### 5. Inconsistência na validação de senha

**Arquivos:** `internal/http/handlers/auth_handler.go` e `internal/service/auth_service.go`  
**Problema:** O handler valida mínimo de 8 caracteres, mas o serviço valida mínimo de 6. A regra é contraditória.

```go
// auth_handler.go linha ~120 — define mínimo 8:
MinLen("password", req.Password, 8)

// auth_service.go linha ~75 — define mínimo 6:
if input.Password == "" || len(input.Password) < 6 || len(input.Password) > 50 {
```

**O que mudar em `internal/service/auth_service.go`:**
```go
// ANTES:
if input.Password == "" || len(input.Password) < 6 || len(input.Password) > 50 {
// DEPOIS:
if input.Password == "" || len(input.Password) < 8 || len(input.Password) > 50 {
```

---

### 6. Configurar proxy reverso com HTTPS

**Problema:** O servidor Go escuta na porta 8080 em HTTP puro. Em produção, usuários devem acessar via HTTPS. Sem isso, tokens JWT e cookies transitam em texto plano.

**O que fazer:** Instale o Caddy (mais simples) ou Nginx na sua VPS.

**Exemplo com Caddy (`/etc/caddy/Caddyfile`):**
```
api.seudominio.com.br {
    reverse_proxy localhost:8080
}

seudominio.com.br {
    reverse_proxy localhost:5173   # ou sirva o build estático
}
```
O Caddy gerencia o certificado TLS automaticamente (Let's Encrypt).

**Exemplo com Nginx (`/etc/nginx/sites-available/jesterx`):**
```nginx
server {
    listen 443 ssl;
    server_name api.seudominio.com.br;

    ssl_certificate     /etc/letsencrypt/live/api.seudominio.com.br/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.seudominio.com.br/privkey.pem;

    location / {
        proxy_pass         http://localhost:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Forwarded-For $remote_addr;  # IP real do cliente
        proxy_set_header   X-Real-IP $remote_addr;
    }
}
```

---

### 7. Configurar IDs e URLs de produção no `.env`

**Arquivo:** `.env` (cópia do `.env.example`)

```bash
# Valores que DEVEM ser alterados antes de ir ao ar:

ENVIRONMENT=prod                         # não "dev"

FRONTEND_URL=https://seudominio.com.br   # URL real do frontend
BACKEND_URL=https://api.seudominio.com.br

JWT_ACCESS_TOKEN=<string-aleatória-longa-e-segura>   # nunca "teste"
JWT_REFRESH_TOKEN=<string-aleatória-diferente>        # nunca "teste1"

STRIPE_PUBLIC_KEY=pk_live_XXXXX          # chave live, não test
STRIPE_SECRET_KEY=sk_live_XXXXX
STRIPE_WEBHOOK_SECRET=whsec_XXXXX        # obrigatório (item 2)

RESEND_KEY=re_XXXXX
```

**Para gerar secrets seguros:**
```bash
openssl rand -base64 64
```

---

### 8. Configurar o frontend com o ID da plataforma

**Arquivo:** `www/.env` (crie baseado no `.env.example`)

```bash
VITE_WEBSITE_ID=00000000-0000-0000-0000-000000000001
VITE_API_URL=https://api.seudominio.com.br
```

O `VITE_WEBSITE_ID` é o ID do website da plataforma Jesterx criado pelo seed (`0006_Seed.up.sql`).  
Sem isso, o header `X-Website-Id` não é enviado e nenhum cadastro/login funciona.

---

### 9. Atualizar os planos no seed com preços reais

**Arquivo:** `migrations/0006_Seed.up.sql`

Os planos atuais são placeholders. Atualize com seus preços definitivos e crie os produtos correspondentes no Stripe Dashboard.

```sql
-- 0006_Seed.up.sql — atualize os valores:
INSERT INTO plans (id, name, description, description_md, price, billing_cycle, active, max_sites, max_routes, ...)
VALUES
  ('plan-starter', 'Starter', '...', '...', 29.90, 'monthly', TRUE, 1, 7, ...),
  ('plan-pro',     'Pro',     '...', '...', 79.90, 'monthly', TRUE, 5, 35, ...),
  ('plan-business','Business','...', '...', 199.90,'monthly', TRUE, 15, 105, ...)
```

**Importante:** Os IDs dos planos (`plan-starter`, `plan-pro`, `plan-business`) devem bater com os Price IDs ou Product IDs que você usar no Stripe. Atualmente o checkout usa `price_data` dinâmico (cria o preço na hora), então os IDs do banco são independentes do Stripe.

---

### 10. Definir política de backup do PostgreSQL

**Problema:** Os dados ficam no volume Docker `postgres_data`. Se o servidor cair sem backup, tudo é perdido.

**Opção A — Backup automático com cron na VPS:**
```bash
# /etc/cron.d/jesterx-backup
0 3 * * * root docker exec jesterx_postgres pg_dump -U postgres jesterx | gzip > /backups/jesterx_$(date +\%Y\%m\%d).sql.gz
# manter os últimos 7 dias:
0 4 * * * root find /backups -name "jesterx_*.sql.gz" -mtime +7 -delete
```

**Opção B — Use um banco gerenciado** (Supabase, Railway, Neon, PlanetScale-pg) que já inclui backup.

---

## 🟡 Média — Importante mas não bloqueia o lançamento imediato

### 11. Adicionar testes automatizados para os serviços críticos

**Problema:** Nenhum `*_test.go` existe no repositório. Qualquer refatoração pode quebrar silenciosamente autenticação ou pagamentos.

**O que criar** (em ordem de impacto):

```
internal/service/auth_service_test.go     — Register, Login, VerifyEmail
internal/service/payment_service_test.go  — validateStripeSignature, ConfirmCheckoutSession
internal/service/order_service_test.go    — CreateOrder (validação de estoque, endereço)
```

**Exemplo mínimo** (`internal/service/auth_service_test.go`):
```go
package service_test

import (
    "testing"
    "github.com/ViitoJooj/Jesterx/internal/service"
)

func TestRegister_InvalidEmail(t *testing.T) {
    // Use um repositório mock (implemente a interface UserRepository)
    svc := service.NewAuthService(mockUserRepo{}, mockWebsiteRepo{}, mockPaymentRepo{})
    _, err := svc.Register(service.RegisterInput{Email: "nao-e-email", Password: "senha123", ...})
    if err == nil {
        t.Fatal("esperava erro para email inválido")
    }
}
```

O projeto já tem `docker-compose.test.yml` e `make test-integration` prontos — só falta escrever os testes.

---

### 12. Adicionar monitoramento de erros

**Problema:** Erros em produção aparecem apenas no `docker compose logs` — sem alerta, sem stack trace, sem contexto.

**Opção A — Sentry (gratuito até 5k eventos/mês):**
```bash
go get github.com/getsentry/sentry-go
```
```go
// cmd/api/main.go — após config.LoadEnv():
import "github.com/getsentry/sentry-go"
sentry.Init(sentry.ClientOptions{Dsn: os.Getenv("SENTRY_DSN")})
defer sentry.Flush(2 * time.Second)
```
```bash
# .env
SENTRY_DSN=https://xxxxx@sentry.io/xxxxx
```

**Opção B — UptimeRobot (gratuito):**  
Configure um monitor HTTP para `https://api.seudominio.com.br/health` e receba e-mail quando o servidor cair.

---

### 13. Substituir logs por logging estruturado

**Problema:** `log.Println(err)` espalhado pelo código mistura informações sem nível (debug, info, error) e sem campos estruturados, dificultando análise em produção.

**O que fazer:** O `pkg/logger/logger.go` já existe com middleware de request logging. Considere usar `log/slog` (Go 1.21+, já disponível na versão do projeto) para logs estruturados nos services:

```go
// ANTES (em qualquer service/handler):
log.Println("cleanup error:", err)

// DEPOIS:
slog.Error("cleanup failed", "error", err)
```

Sem adicionar dependências, apenas substitua `log.Println`/`log.Printf` por `slog.Info`/`slog.Error` nos arquivos de serviço.

---

### 14. Documentar API (README ou Swagger)

**Problema:** Não existe documentação das rotas da API. Dificulta integração com outros frontends ou parceiros.

**Todas as rotas disponíveis** (para referência):

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| POST | `/api/v1/auth/register` | Público | Cadastro |
| POST | `/api/v1/auth/login` | Público | Login |
| GET | `/api/v1/auth/refresh` | Cookie | Renovar access token |
| GET | `/api/v1/auth/me` | JWT | Dados do usuário |
| PATCH | `/api/v1/auth/me` | JWT | Atualizar perfil |
| DELETE | `/api/v1/auth/me` | JWT | Deletar conta |
| GET | `/api/v1/auth/logout` | Público | Logout |
| GET | `/api/v1/plans` | Público | Listar planos |
| POST | `/api/v1/payments/checkout` | JWT | Criar checkout Stripe |
| GET | `/api/v1/payments/confirm` | JWT | Confirmar checkout |
| POST | `/api/v1/payments/cancel` | JWT | Cancelar assinatura |
| POST | `/api/v1/payments/webhook` | Stripe | Webhook eventos |
| POST | `/api/v1/websites` | JWT | Criar website |
| GET | `/api/v1/websites` | JWT | Listar websites do usuário |
| DELETE | `/api/v1/sites/{id}` | JWT | Deletar website |
| POST | `/api/v1/sites/{id}/routes` | JWT | Atualizar rotas |
| POST | `/api/v1/sites/{id}/versions` | JWT | Criar versão |
| POST | `/api/v1/sites/{id}/publish/{v}` | JWT | Publicar versão |
| POST | `/api/v1/sites/{id}/products` | JWT | Criar produto |
| GET | `/api/v1/sites/{id}/products` | JWT | Listar produtos (owner) |
| PATCH | `/api/v1/sites/{id}/products/{pid}` | JWT | Atualizar produto |
| DELETE | `/api/v1/sites/{id}/products/{pid}` | JWT | Deletar produto |
| GET | `/api/v1/sites/{id}/orders` | JWT | Listar pedidos (owner) |
| POST | `/api/v1/upload` | JWT | Upload de arquivo |
| GET | `/files/{path}` | Público | Servir arquivo |
| GET | `/api/v1/themes` | Público | Listar temas |
| GET | `/api/store/{id}/products` | Público | Produtos da loja |
| GET | `/api/store/{id}/products/{pid}` | Público | Detalhe do produto |
| POST | `/api/store/{id}/orders` | JWT | Criar pedido |
| GET | `/api/store/{id}/info` | Público | Info da loja |
| POST | `/api/store/{id}/comments` | JWT | Comentar |
| POST | `/api/store/{id}/ratings` | JWT | Avaliar loja |
| GET | `/api/v1/admin/stats` | Admin | Estatísticas globais |
| GET | `/api/v1/admin/users` | Admin | Listar usuários |
| GET | `/api/v1/admin/sites` | Admin | Listar sites |
| GET | `/api/v1/admin/orders` | Admin | Listar pedidos |
| POST | `/api/v1/reports` | Público | Criar denúncia |
| PATCH | `/api/v1/admin/reports/{id}` | Admin | Moderar denúncia |
| GET | `/health` | Público | Health check |

**Forma mais simples de gerar docs:** execute `go run cmd/cli/main.go routes` para listar todas as rotas registradas.

---

## 🟢 Baixa — Melhorias para após o lançamento

### 15. Cobrar o comprador via Stripe (não só o plano do lojista)

**Problema:** Quando um usuário faz um pedido na loja (`POST /api/store/{id}/orders`), o pedido é criado com `status: pending` mas **não há cobrança no Stripe**. O fluxo de pagamento do comprador não está implementado.

**O que fazer:**
1. Após `s.orderRepo.Create(order)` em `internal/service/order_service.go`, criar um Stripe Payment Intent ou Checkout Session para o comprador.
2. Adicionar rota `GET /api/store/{id}/orders/{orderID}/pay` que retorna o link de checkout.
3. Webhook para atualizar `order.status` para `paid` após o pagamento.

---

### 16. Painel do lojista com gráficos de vendas

**Problema:** Existe `GET /api/v1/sites/{id}/orders` que retorna pedidos, mas o frontend não tem página de dashboard com gráficos.

**O que fazer no backend:**
- Adicionar endpoint `GET /api/v1/sites/{id}/stats` retornando: total de pedidos, receita total, pedidos por status, receita por dia (últimos 30 dias).

**O que fazer no frontend (`www/src/pages/`):**
- Criar `dashboard/Dashboard.tsx` com gráfico de vendas usando uma lib leve (recharts ou chart.js).

---

### 17. Domínio personalizado por loja

**Problema:** Atualmente as lojas são acessíveis apenas em `/p/{siteID}`. Para um SaaS sério, cada lojista deveria poder usar `minha-loja.com`.

**O que fazer:**
1. Adicionar coluna `custom_domain TEXT UNIQUE` na tabela `websites`.
2. No middleware, resolver o domínio da requisição para o `siteID` correspondente.
3. Instruir o lojista a apontar o CNAME para o servidor Jesterx.

---

### 18. Integração com cálculo de frete

**Problema:** O campo `shipping_cost` existe no domain `Order` mas é sempre `0`. Não há integração com Correios/Melhor Envio.

**O que fazer:**
- Integrar a API da [Melhor Envio](https://docs.melhorenvio.com.br) para calcular frete antes de finalizar o pedido.
- Adicionar endpoint `POST /api/store/{id}/shipping/quote` com CEP de origem (da loja) e destino (do comprador).

---

## Resumo por arquivo

| Arquivo | Mudança |
|---------|---------|
| `pkg/ratelimit/ratelimit.go` | Fix IP spoofing — não confiar em X-Forwarded-For do cliente (item 1) |
| `internal/service/payment_service.go` | STRIPE_WEBHOOK_SECRET obrigatório em prod (item 2) |
| `docker-compose.yml` | Remover Redis não utilizado (item 3) |
| `internal/config/dotenv.go` | Remover RedisURL (item 3) |
| `internal/service/auth_service.go` | Senha mínimo 8 chars (item 5) |
| `.env` | Valores de produção reais (item 7) |
| `www/.env` | VITE_WEBSITE_ID e VITE_API_URL (item 8) |
| `migrations/0006_Seed.up.sql` | Preços e limites reais dos planos (item 9) |
| `nginx.conf` / `Caddyfile` | TLS + proxy reverso (item 6) |
| SQL manual | Criar usuário admin no banco (item 4) |
| `internal/service/*_test.go` | Testes unitários (item 11) |
| `cmd/api/main.go` | Sentry DSN opcional (item 12) |
