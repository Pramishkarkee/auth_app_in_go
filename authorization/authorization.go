package authorization


import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	// "time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// "golang.org/x/crypto/bcrypt"
	"auth/test/exception"
	// "auth/test/jwt"
)

var (
	secretkey string = "secretkeyjwt"
)
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err authexception.Error
			err = authexception.SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err authexception.Error
			err = authexception.SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return

			}
		}
		var reserr authexception.Error
		reserr = authexception.SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(err)
	}
}