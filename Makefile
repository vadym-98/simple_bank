createdb:
	docker exec -it playgroundDB createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it playgroundDB dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/vadym-98/simple_bank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock