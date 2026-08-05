force = false
kubernetes = false

# --- Global(all services) ---

# Build all application
verkoupe-build:

# Run all tests
verkoupe-test:

# Make run all services
verkoupe-run:

# --- Core(Daemon) ---

# Run core tests
daemon-test:

# Run verkoupe core
daemon-run:

# Build core
daemon-build:

# --- Front end(View) ---

# Run verkoupe front end
view-run:

# Run front end tests
view-tests:

# Build front end
view-build:

# --- database ---

# Run postgres, redis and mongo
database-run:

# Go to next migration
migration-next:

# Return on migration
migration-return:

# --- gerals ---

# Terraform apply
terraform-apply:

# Drop all databases and destroy clusters and containers
nuke:
up-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml up -d --build

down-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml down

up-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml up -d

down-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml down
