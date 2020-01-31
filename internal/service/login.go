package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/model"
)

// Login will login user with the email and password user provide
func Login(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var jwt model.JWT
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

		password := user.Password

		row := db.QueryRow("select * from users where email=$1", user.Email)
		err := row.Scan(&user.ID, &user.Email, &user.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist"
				respondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "The password does not match"
			respondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := GenerateToken(user)

		if err != nil {
			error.Message = "Error when generating jwt token"
			respondWithError(w, http.StatusInternalServerError, error)
			return
		}

		jwt.Token = token

		respondWithJSON(w, jwt)

		fmt.Println(user)
	}
}

// GenerateToken will generate jwt token based on user login email
func GenerateToken(user model.User) (string, error) {
	secret := "secret for tokenizing jwt"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
