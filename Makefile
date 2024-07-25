postgres:
	docker run --name split-db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e PGUSER=root -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it split-db createdb --username=root --owner=root split_db
dropdb:
	docker exec -it split-db dropdb split_db
migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_db?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_db?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/arkarsg/splitapp/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock
