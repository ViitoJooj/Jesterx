# 🧪 Guia Essencial de Testes para Agentes de IA (Go + Python)

Este documento define regras obrigatórias para criação e execução de testes em um projeto com backend Go e testes externos em Python.
O objetivo é garantir testes consistentes, reprodutíveis e seguros.

## 📌 Regra Nº1: .env é obrigatório
- Sempre carregar variáveis do arquivo .env
- O arquivo .env sempre estará na raiz do projeto
- O agente deve sempre se basear no .env.example para saber quais variáveis existem
- O agente NUNCA deve alterar o .env
- O agente NUNCA deve expor valores sensíveis (tokens, secrets, senhas)

## 📍Arquivos esperados:

`/.env
/.env.example
📁 Estrutura mínima recomendada
.
├── internal/                   # Código de negócio (Go)
│   └── ...
├── tests/                      # Testes E2E/integração (Python)
│   ├── requirements.txt
│   ├── conftest.py
│   └── test_*.py
├── .env
├── .env.example
├── main.go
└── go.mod`

## ✅ Testes Unitários (Go)

- Regras obrigatórias
- Todo teste unitário deve estar em arquivo *_test.go
- Testes devem ser rápidos e sem dependências externas (HTTP real, DB real)
- Usar padrão Table-Driven Tests
- Usar t.Run para cenários
- Testar somente regras internas (domínio/service)

## 🌐 Testes de Integração / E2E (Python)
Objetivo

Simular o comportamento real do usuário consumindo a API rodando localmente.

- Regras obrigatórias
- Testes devem rodar via pytest
- Sempre validar:
- status code
- resposta JSON
- fluxo esperado do usuário
- Payloads devem ser realistas (como um frontend enviaria)
- Testes devem ser independentes (não depender da ordem)
- Testes de usuários devem usar dados aleatórios (nome/email/senha), evitando reutilização entre execuções

## 🔐 Regra Especial: Testes de Login

Para testes de login/autenticação, é permitido:

- ✅ Criar usuário diretamente na tabela (setup)
- ✅ Testar login via API
- ✅ Remover usuário criado (cleanup)

## ⚠️ Regras obrigatórias:

- Isso só pode ser feito em ambiente de teste
- Nunca reutilizar usuário real
- Sempre deletar o usuário no final do teste
- Preferir emails únicos (ex: UUID)
- Para qualquer teste envolvendo usuário, gerar dados aleatórios por execução

## ▶️ Como executar

Go
go test ./... -v
Python
cd tests
pip install -r requirements.txt
pytest -v

## 🛑 Regras de Segurança (Obrigatórias)

O agente deve sempre respeitar:

- Nunca imprimir .env em logs
- Nunca expor secrets, tokens ou senhas
- Nunca commitar .env
- Sempre usar .env.example como referência de variáveis esperadas
- Nunca alterar configurações sensíveis do ambiente
