package main

import (
	"gloriusaiapi/config"
	"gloriusaiapi/migrations"
	"gloriusaiapi/routers"
	"log"
	"net/http"
)

func main() {
	config.InitializeDB()
	migrations.Migrate()

	router := routers.SetupRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
