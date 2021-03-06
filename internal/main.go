package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/driver"
	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/model"
	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/service"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {

	db := driver.ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/signup", service.Signup(db)).Methods("POST")
	router.HandleFunc("/login", service.Login(db)).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(service.ProtectedEndpoint)).Methods("GET")

	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// TokenVerifyMiddleWare verify the authorization token before it forwards the request to the route
func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err model.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("The signin Method cannot be verified")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				err.Message = error.Error()
				service.RespondWithError(w, http.StatusUnauthorized, err)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				err.Message = "Token is not valid"
				service.RespondWithError(w, http.StatusUnauthorized, err)
				return
			}

		} else {
			err.Message = "Token format is not correct"
			service.RespondWithError(w, http.StatusUnauthorized, err)
		}

	}
}
