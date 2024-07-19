package main

import (
	"database/sql"
	"log"

	"github.com/arkarsg/splitapp/api"
	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/split_app?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	dbEnvs := u.GetDevDbEnvs()
	serverEnvs := u.GetServerEnvs()

	conn, err := sql.Open(dbEnvs.DbDriver, u.GetDevDbSource())

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverEnvs.Address)
	if err != nil {
		log.Fatal("🛑 Cannot start server: ", err)
	}
}
