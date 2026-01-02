<p align="center">
  🇧🇷 Português | <a href="./README.en.md">🇺🇸 English</a>
</p>

# Jesterx

Jesterx é um **SaaS de criação de paginas**.  
A versão inicial (v1) será aberta e pública no GitHub, servindo tanto como produto funcional quanto como experimento técnico. O projeto tem como foco **simplicidade no código**, aprendizado prático e validação de ideias.  
Ele existe para testar minhas habilidades com **Golang**, estudar arquitetura escalável e, se possível, gerar receita no futuro.

### Stack e decisões técnicas

- **Golang + Gin** no backend
- **React + TypeScript** no frontend
- **PostgreSQL** para usuários
- **MongoDB** para persistência dos sites
- **Redis** para cache e escalabilidade
- **RabbitMQ** para notificações e tarefas assíncronas
- Autenticação via **OAuth2**
- APIs públicas
- Tudo **containerizado com Docker**
- Infraestrutura na **Azure**, com **CI/CD**
- **Testes automatizados em Python**, escolhidos pela simplicidade e rapidez

---

## Visão geral

Este repositório reúne tudo o que é necessário para iniciar um e-commerce:

- API para regras de negócio
- Interface web
- Estrutura de banco de dados

A ideia é permitir que qualquer desenvolvedor consiga clonar o projeto, subir o ambiente e começar a trabalhar sem muita configuração inicial.

---

## Estrutura do projeto

O projeto segue uma arquitetura **modular**, dividida em três partes principais:

```
jesterx/
    ㄴ backend/                     # API e lógica do servidor
           ├─ config/               # Configuração dos projetos
           ├─ helpsers/             # Funções auxiliares
           ├─ middlewares/          # middlewares
           ├─ models/               # Padronização de modelos esperados nas apis
           ├─ responses/            # Padronização de respostas das apis
           ├─ services/             # serviços das apis
           ├─ go.mod                # lib
           ├─ go.sum                # lib
           ㄴ main.go               # Aplicativo principal da api
       ├─ frontend/ # Interface do usuário
       ├─ sql/ # Scripts do banco de dados
       ├─ .env.example # Exemplo de variáveis de ambiente
       ├─ docker-compose.yml
       ├─ LICENSE.md
```

---

## Backend

O backend concentra toda a lógica da aplicação, como:

- Autenticação de usuários
- Produtos
- Pedidos
- Comunicação com o banco de dados

A API segue o padrão REST, com possibilidade de adaptação futura se necessário.

---

## Frontend

O frontend é responsável pela interface da loja, incluindo:

- Listagem de produtos
- Carrinho de compras
- Login e cadastro
- Checkout
- Área administrativa (em desenvolvimento)

Ele consome diretamente a API do backend.

---

## Banco de dados

Na pasta `sql/` ficam os scripts de banco, incluindo:

- Criação das tabelas
- Relacionamentos
- Dados iniciais (quando aplicável)

---

## Configuração do ambiente

### Variáveis de ambiente

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Depois ajuste as variáveis conforme seu ambiente, como banco de dados, portas e chaves de acesso.

## Docker

O projeto possui um _docker-compose.yml_ para facilitar o setup local:

```
docker compose up -d
```

Isso irá subir o backend, frontend, redis, rabbitMQ e banco de dados.

## Funcionalidades

- Cadastro e autenticação de usuários
- (post/get) Comunicação rest nos serviços de softwares
- CRUD de produtos
- Carrinho de compras
- Sistema de pedidos
- Checkout
- Painel administrativo
- Integrações com meios de pagamento

## Desenvolvimento local

Para desenvolvimento local:

vá para o backend na pasta raiz.

```
cd backend
```
instale as dependencias
```
go mod download
```

e depois

```
cd frontend
```
e
```
npm install
```

## Contribuição

#### Quer contribuir?
É só:
1. Fazer um fork
2. Criar uma branch (autor/sua-feature)
3. Commitar suas mudanças
4. Abrir um Pull Request

- Qualquer dúvida, confere o <a href="./CONTRIBUTING.md">CONTRIBUTING.md</a>

## Licença

Este projeto está licenciado conforme o arquivo <a href="LICENSE.md">LICENSE.md</a>
