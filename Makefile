force = false
kubernetes = false

# --- Global(all services) ---

# Build all application
verkoupe-build: daemon-build view-build

# Run all tests
verkoupe-test: daemon-test view-tests

# Make run all services
verkoupe-run: daemon-run view-run

# --- Core(Daemon) ---

# Run core tests
daemon-test:
	go test ./...

# Run verkoupe core
daemon-run:
	go run ./cmd/api

# Build core
daemon-build:
	CGO_ENABLED=0 go build -o ./bin/daemon ./cmd/api

# --- Front end(View) ---

# Run verkoupe front end
view-run:
	cd www && npm run dev

# Run front end tests
view-tests:
	cd www && npm test

# Build front end
view-build:
	cd www && npm run build

# --- database ---

# Run postgres, redis and mongo
database-run:
	docker compose -f ./infra/compose/docker-compose-dev.yaml up -d postgres

# Go to next migration
migration-next:
	@echo "Run migrations manually: psql -f migrations/001_init.sql"

# Return on migration
migration-return:
	@echo "Rollback not implemented yet"

# --- gerals ---

# Terraform apply
terraform-apply:
	@echo "Terraform not configured yet"

# Drop all databases and destroy clusters and containers
nuke:
	docker compose -f ./infra/compose/docker-compose-dev.yaml down -v

up-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml up -d --build

down-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml down

up-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml up -d

down-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml down
