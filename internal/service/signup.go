package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func respondWithError(w http.ResponseWriter, statusCode int, error model.Error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(error)
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

// Signup lets user sign up with email and password
func Signup(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var error model.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			respondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			respondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		statement := "insert into users (email, password) values($1, $2) RETURNING id;"

		err = db.QueryRow(statement, user.Email, user.Password).Scan(&user.ID)

		if err != nil {
			error.Message = "There was an issue inserting the user to database"
			respondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		respondWithJSON(w, user)
	}
}
