postgres:
	@docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=changeme -d postgres:16-alpine

createdb:
	@docker exec -it postgres16 createdb --username=root --owner=root simple_bank

enterdb:
	@docker exec -it postgres16 psql simple_bank

dropdb:
	@docker exec -it postgres16 dropdb simple_bank

migrate-up:
	@migrate -path infra/db/migration -database "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-up-1:
	@migrate -path infra/db/migration -database "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate-down:
	@migrate -path infra/db/migration -database "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate-down-1:
	@migrate -path infra/db/migration -database "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	@sqlc generate

mock:
	@mockgen -package mockdb  -destination infra/db/mock/store.go github.com/pe-Gomes/simple-bank-go/infra/db/repository Store

test:
	@go test -v -cover ./...

server:
	@go run main.go

.PHONY: postgres createdb dropdb migrate-up migrate-down migrate-up-1 migrate-down-1 sqlc test server mock