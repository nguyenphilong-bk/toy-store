up:
	docker-compose up -d --force-recreate --remove-orphans

ps:
	docker-compose ps

down:
	docker-compose down

migrate-up:
    migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/toy_store?sslmode=disable" -verbose up

migrate-down:
    migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/toy_store?sslmode=disable" -verbose down
