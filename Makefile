postgres:
	docker run --name postgres12 --network bank-network -p 5436:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	# タブ文字を使う
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	# タブ文字を使う
	docker exec -it postgres12 dropdb simple_bank

migrationup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5436/simple_bank?sslmode=disable" -verbose up

migrationup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5436/simple_bank?sslmode=disable" -verbose up 1

migrationdown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5436/simple_bank?sslmode=disable" -verbose down

migrationdown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5436/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrationup migrationup1 migrationdown migrationdown1 sqlc test server mock
