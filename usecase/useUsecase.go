package loginusecase

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	// "time"

	// "github.com/dgrijalva/jwt-go"
	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"auth/test/conn"
	"auth/test/exception"
	"auth/test/hash"
	"auth/test/jwt"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	connection := mydb.GetDatabase()
	defer mydb.CloseDatabase(connection)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err authexception.Error
		err = authexception.SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var dbuser User
	connection.Where("email = ?", user.Email).First(&dbuser)

	//check email is alredy registered or not
	if dbuser.Email != "" {
		var err authexception.Error
		err = authexception.SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = hashpassword.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("Error in password hashing.")
	}

	// insert user details in database
	connection.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	connection := mydb.GetDatabase()
	defer mydb.CloseDatabase(connection)

	var authDetails Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err authexception.Error
		err = authexception.SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authUser User
	connection.Where("email = 	?", authDetails.Email).First(&authUser)

	if authUser.Email == "" {
		var err authexception.Error
		err = authexception.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := hashpassword.CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err authexception.Error
		err = authexception.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := generatetoken.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err authexception.Error
		err = authexception.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}