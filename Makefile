postgres:
	docker run --name split-db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it split-db createdb --username=root --owner=root split_app
dropdb:
	docker exec -it split-db dropdb split_app
migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_app?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_app?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown sqlc server
