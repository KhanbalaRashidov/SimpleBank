DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres   -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:15-alpine


createdb:
	docker exec -it  postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it  postgres dropdb simple_bank

sqlc:
	docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate \


migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose down    

test:
	go test -v -cover ./...

.PHONY: network postgres createdb dropdb sqlc migrateup migratedown test
