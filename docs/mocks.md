# Mocks

O que está mockado no frontend (`www/`) e onde ligar o real depois.

## Autenticação

- **Arquivo:** `www/src/auth/index.tsx`
- Usuário lido de `localStorage["auth-user"]` (JSON `{ name, email, avatar? }`).
- `logout()` apenas remove a chave. Não existe `login()` — botões Entrar/Registrar só navegam para `/login` e `/register` (rotas inexistentes, caem em NotFound).
- **Ligar real:** trocar provider por chamadas ao backend Go (`internal/`) e sessão/token.

## Header — Notificações

- **Arquivo:** `www/src/components/Header/Header.tsx`
- Dropdown do sino mostra sempre estado vazio (`header.noNotifications`). Sem fetch, sem badge de não-lidas.
- **Ligar real:** endpoint de notificações + estado de leitura.

## Footer — Links e redes sociais

- **Arquivo:** `www/src/components/Footer/Footer.tsx`
- Colunas `mainPages` e `utilityPages`: somente `/` e `/design` existem. Demais (`/about`, `/blog`, `/pricing`, `/contact`, `/careers`, `/integrations`, `/login`, `/register`, `/licenses`, `/changelog`, `/404`) caem em NotFound.
- Array `socials`: URLs genéricas (`facebook.com`, `twitter.com`, `instagram.com`, `linkedin.com`), não perfis reais.
- **Ligar real:** criar rotas/páginas e trocar URLs sociais pelos perfis oficiais.

## Tema

- **Arquivo:** `www/src/theme/index.ts`
- Funcional (persiste em `localStorage["theme"]`), não é mock — registrado aqui apenas como referência de wiring.
