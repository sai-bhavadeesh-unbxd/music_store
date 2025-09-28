.PHONY: test-api up-redis run seed

up-redis:
	@docker compose up -d redis

run:
	@go run ./...

test-api: up-redis
	@bash scripts/test_api.sh

seed: up-redis
	@bash scripts/seed_data.sh
