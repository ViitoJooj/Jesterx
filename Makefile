.PHONY: dev build test migrate docker-up docker-down docker-test cli

# Iniciar servidor de desenvolvimento
dev:
	go run cmd/api/main.go

# Build da API
build:
	go build -ldflags="-s -w" -o bin/jesterx ./cmd/api

# Build do CLI
cli:
	go build -ldflags="-s -w" -o bin/jx ./cmd/cli

# Executar testes
test:
	go test ./...

# Rodar testes com banco de teste
test-integration:
	docker compose -f docker-compose.test.yml up -d --wait
	POSTGRES_HOST=localhost POSTGRES_PORT=5433 POSTGRES_DB=jesterx_test go test ./...
	docker compose -f docker-compose.test.yml down

# Aplicar migrations
migrate:
	go run cmd/cli/main.go migrate

# Build do frontend
frontend:
	cd www && npm run build

# Iniciar Docker
docker-up:
	docker compose up -d

# Parar Docker
docker-down:
	docker compose down

# Logs da API
logs:
	docker compose logs -f api

# Formatar código
fmt:
	go fmt ./...

# Lint
lint:
	go vet ./...
