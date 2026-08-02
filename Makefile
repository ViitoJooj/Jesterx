up-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml up -d --build

down-dev:
	docker compose -f ./infra/compose/docker-compose-dev.yaml down

up-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml up -d

down-prod:
	docker compose -f ./infra/compose/docker-compose-prod.yaml down