# Contributing | Contribuindo

[English](#english) | [Português](#português)

---

## English

Thanks for considering contributing to Jesterx! 

### Project Setup

**What you'll need:**

- Node.js (v18+)
- Go (v1.21+)
- PostgreSQL (v14+)
- pnpm or npm

**Installation:**

```bash
# Clone the project
git clone https://github.com/ViitoJooj/Jesterx.git

# Install frontend dependencies
npm install

# Install backend dependencies
cd backend
go mod download
```

Set up your `.env` based on `.env.example` and run the project:

```bash
# Frontend
npm run dev

# Backend (in another terminal)
cd backend
go run main.go
```

### How to Contribute

1. Fork the project
2. Create a branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Commit (`git commit -m 'add: new feature'`)
5. Push to your branch (`git push origin feat/my-feature`)
6. Open a Pull Request

**Before opening a PR:**

- Test locally to make sure nothing breaks
- If it's a visual change, include screenshots in the PR
- Explain what you did clearly

### Code Standards

**TypeScript:**
- Use TypeScript, it helps a lot
- Prettier is configured, just run it
- Variable names need to make sense

**Go:**
- `gofmt` before committing
- Comment exported functions
- No hacky code

**SCSS:**
- BEM for classes
- Use existing variables
- Mobile first always

**SQL:**
- snake_case for tables and columns
- Comment complex queries

### Commits

Try to keep commits organized: 

- `feat:` for new features
- `fix:` for bugs
- `docs:` for documentation
- `refactor:` when refactoring something
- `chore:` for general tasks

Example: `feat: add stripe checkout`

### Reporting Bugs

Open an issue with:
- What happened
- What should happen
- How to reproduce
- Screenshots if helpful

### Suggesting Features

Open an issue explaining:
- What you want
- Why it would be useful
- How you imagine it working

### Questions

If you have any questions, open an issue or discussion and we'll respond.

---

Thanks for contributing! 🚀

---

## Português

Obrigado por considerar contribuir com o Jesterx!

### Setup do Projeto

**O que você vai precisar:**

- Node.js (v18+)
- Go (v1.21+)
- PostgreSQL (v14+)
- pnpm ou npm

**Instalação:**

```bash
# Clone o projeto
git clone https://github.com/ViitoJooj/Jesterx.git

# Instale as dependências do frontend
npm install

# Instale as dependências do backend
cd backend
go mod download
```

Configure o `.env` baseado no `.env.example` e rode o projeto:

```bash
# Frontend
npm run dev

# Backend (em outro terminal)
cd backend
go run main.go
```

### Como Contribuir

1. Dá um fork no projeto
2. Cria uma branch (`git checkout -b feat/minha-feature`)
3. Faz suas mudanças
4. Commita (`git commit -m 'add: nova feature'`)
5. Push pra sua branch (`git push origin feat/minha-feature`)
6. Abre um Pull Request

**Antes de abrir um PR:**

- Testa localmente pra ver se não quebrou nada
- Se for mudança visual, coloca uns prints no PR
- Explica o que você fez de forma clara

### Padrões de Código

**TypeScript:**
- Usa TypeScript mesmo, ajuda muito
- Prettier tá configurado, só rodar
- Nome de variável tem que fazer sentido

**Go:**
- `gofmt` antes de commitar
- Comenta as funções exportadas
- Nada de gambiarra

**SCSS:**
- BEM pras classes
- Usa as variáveis que já existem
- Mobile first sempre

**SQL:**
- snake_case pras tabelas e colunas
- Comenta as queries complexas

### Commits

Tenta manter os commits organizados:

- `feat:` pra features novas
- `fix:` pra bugs
- `docs:` pra documentação
- `refactor:` quando refatora algo
- `chore:` pra tarefas gerais

Exemplo: `feat: adiciona checkout com stripe`

### Reportar Bugs

Abra uma issue com: 
- O que aconteceu
- O que deveria acontecer
- Como reproduzir
- Print se ajudar

### Sugerir Features

Abre uma issue explicando:
- O que você quer
- Por que seria útil
- Como você imagina funcionando

### Dúvidas

Se tiver qualquer dúvida, abre uma issue ou discussion que a gente responde.

---

Valeu pela contribuição!
