package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kovarike/protocol/src/go/db"
	"github.com/kovarike/protocol/src/go/tls"
	"github.com/kovarike/protocol/src/go/worker"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	worker.Init(db.Pg, db.Rdb)

	ctx := context.Background()
	modemPort := "/dev/ttyUSB0"
	go worker.WorkerLoop(ctx, modemPort)

	err = tls.ListenTLS(":2525", "server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server listening on :2525")
}
