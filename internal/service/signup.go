package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Signup lets user sign up with email and password
func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var error model.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			RespondWithError(w, http.StatusBadRequest, error)
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
			RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		RespondWithJSON(w, user)
	}
}
