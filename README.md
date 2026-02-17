<p align="center">
  🇧🇷 Português | <a href="./docs/README_en.md">🇺🇸 English</a> | <a href="./docs/README_cn.md">🇨🇳 CN</a>
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
- Área administrativa (dashboard para admins)

Ele consome diretamente a API do backend.

---

## Área Administrativa

- Apenas usuários com e-mail listado em `ADMIN_EMAILS` recebem a role `platform_admin`.
- Dashboard com métricas de usuários criados, ticket médio, receita e planos mais usados.
- Gestão de usuários (visualizar, editar dados de perfil e plano sem expor senha, banir/desbanir, deletar).
- Gestão de planos (alterar preços, descrições, limites e benefícios). O checkout usa sempre os valores atualizados.
- Exportação de usuários em XLSX direto da interface e via API.

---

## Banco de dados

Na pasta `migrations/` ficam os scripts de banco, incluindo:

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
Inclua a lista de e-mails administrativos em `ADMIN_EMAILS` (separados por vírgula) para liberar o dashboard `/admin` e as rotas `/v1/admin`.

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

instale as dependencias
```
go mod tidy
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
