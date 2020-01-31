package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/service"
)

var db *sql.DB

const postgresDBURL = "postgres://xcpcmiip:6CJG7e6m1UinRwnSgBwkK-bUbDW6HNX8@rajje.db.elephantsql.com:5432/xcpcmiip"

func main() {
	pgURL, err := pq.ParseURL(postgresDBURL)

	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", pgURL)

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup", service.Signup(db)).Methods("POST")
	router.HandleFunc("/login", service.Login(db)).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("protectedEndpoint invoked")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	log.Println("TokenVerifyMiddleWare invoked")
	return next
}
