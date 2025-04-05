package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {

	mux := http.NewServeMux()

	const port = 8084
	log.Println("user management server starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
