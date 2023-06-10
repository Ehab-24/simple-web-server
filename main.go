package main

import (
	"log"
	"net/http"

	"suraj.com/web_server/api"
)

func main() {
	s := api.NewServer()
	log.Fatal(http.ListenAndServe(":8080", s))
}
