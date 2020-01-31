package service

import (
	"log"
	"net/http"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("protectedEndpoint invoked")
}
