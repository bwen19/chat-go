DB_URL=postgresql://root:secret@localhost:5432/chat?sslmode=disable

network:
	docker network create chat-network

postgres:
	docker run --name postgres --network chat-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.4-alpine3.18

createdb:
	docker exec -it postgres createdb --username=root --owner=root chat

dropdb:
	docker exec -it postgres dropdb chat

migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

schema:
	dbml2sql --postgres -o schema/chat.sql schema/chat.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run ./src/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor

proto:
	rm -f src/pb/*.go
	rm -f swagger/*.swagger.json
	protoc --proto_path=src/protos --go_out=src/pb --go_opt=paths=source_relative \
		--go-grpc_out=src/pb --go-grpc_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
		--grpc-gateway_out=src/pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=swagger --openapiv2_opt=allow_merge=true,merge_file_name=chat \
		src/protos/*.proto
	statik -src=./swagger -dest=./src

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migration schema sqlc test server mock proto evans redis