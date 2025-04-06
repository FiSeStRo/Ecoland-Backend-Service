package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("view/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	const port = 8084
	log.Println("user management server starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
