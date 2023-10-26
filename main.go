package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/vadym-98/simple_bank/api"
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("can't start the server:", err)
	}
}
