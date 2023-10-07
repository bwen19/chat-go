DB_URL=postgresql://root:secret@localhost:5432/chat?sslmode=disable

network:
	docker network create chat-network

postgres:
	docker run --name postgres --network chat-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.4-alpine3.18

createdb:
	docker exec -it postgres createdb --username=root --owner=root chat

dropdb:
	docker exec -it postgres dropdb chat

schema:
	dbml2sql --postgres -o schema/chat.sql schema/chat.dbml

migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

sqlc:
	rm -f src/db/sqlc/*.sql.go
	sqlc generate

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-alpine3.18

test:
	go test -v -cover -short ./...

server:
	go run ./src/main.go

.PHONY: network postgres createdb dropdb schema migration migrateup migratedown sqlc redis test server