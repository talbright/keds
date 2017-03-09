package server

import (
	"log"
	"net/http"
)

const (
	debugEndPoint = "localhost:8081"
)

/*
	Package golang.org/x/net/trace exports two http handlers for tracing:

	/debug/requests
	/debug/events

	ex: http://localhost:8081/debug/requests
*/
func StartDebugServer() {
	log.Printf("starting http debug server on %s", debugEndPoint)
	if err := http.ListenAndServe(debugEndPoint, nil); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
