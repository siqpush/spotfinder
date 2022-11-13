package main

import (
	"log"
	"spotfinder/internal/server"
)

func main() {
	srv := server.NewHTTPServer("0.0.0.0:3000") 
	log.Fatal(srv.ListenAndServe())
}