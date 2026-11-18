package api

import (
	"api/data"
	"context"
	"log"
)

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := data.SetupDB(ctx, "./data", false)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := RunHTTPServer(ctx, db); err != nil {
		log.Fatal(err)
	}
}
