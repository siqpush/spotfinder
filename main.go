package main

import (
	"log"
	"internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080") 
	log.Fatal(srv.ListenAndServe())
}