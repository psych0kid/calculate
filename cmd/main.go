package main

import (
	"log"
	"net/http"

	"github.com/exxxception/rpn/internal/handler"
)

const addr = ":8080"

// main runs the server.
//
// It registers the handler for /api/v1/calculate and starts the server
// on the address defined in the addr constant.
func main() {
	http.HandleFunc("/api/v1/calculate", handler.CalculateHandler)

	log.Printf("**Server starting on: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
