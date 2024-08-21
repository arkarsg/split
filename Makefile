.PHONY: network
network:
	docker network create split-network

.PHONY: postgres
postgres:
	docker run --name split-db --network split-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e PGUSER=root -p 5432:5432 -d postgres:12-alpine

.PHONY: createdb
createdb:
	docker exec -it split-db createdb --username=root --owner=root split_db

.PHONY: dropdb
dropdb:
	docker exec -it split-db dropdb split_db

.PHONY: migrateup
migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_db?sslmode=disable" -verbose up

.PHONY: migratedown
migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_db?sslmode=disable" -verbose down

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: server
server:
	go run ./cmd/*.go

.PHONY: mock
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/arkarsg/splitapp/db/sqlc Store
