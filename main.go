package main

import (
	"fmt"
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

	port := ":8080"

	fmt.Printf("Server is running on http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
