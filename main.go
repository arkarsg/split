package main

import (
	"database/sql"
	"log"

	"github.com/arkarsg/splitapp/api"
	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
	_ "github.com/lib/pq"
)

func main() {
	config := u.GetConfig()
	dbEnvs := u.GetDevDbEnvs()
	serverEnvs := u.GetServerEnvs()

	conn, err := sql.Open(dbEnvs.DbDriver, u.GetDevDbSource())

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("🛑 Cannot create server: ", err)
	}

	err = server.Start(serverEnvs.Address)
	if err != nil {
		log.Fatal("🛑 Cannot start server: ", err)
	}
}
