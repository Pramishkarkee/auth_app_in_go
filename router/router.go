package router


import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/gorilla/handlers"
	"auth/test/usecase"
	"auth/test/authorization"
)
var (
	router *mux.Router
	secretkey string = "secretkeyjwt"
  )
func CreateRouter() {
	router = mux.NewRouter()
}

func ServerStart() {
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "user" {
		w.Write([]byte("Not Authorized."))
		return
	}
	w.Write([]byte("Welcome, User."))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}
func InitializeRoute() {
	CreateRouter()
	router.HandleFunc("/signup", loginusecase.SignUp).Methods("POST")
	router.HandleFunc("/signin", loginusecase.SignIn).Methods("POST")
	router.HandleFunc("/admin", authorization.IsAuthorized(AdminIndex)).Methods("GET")
	router.HandleFunc("/user", authorization.IsAuthorized(UserIndex)).Methods("GET")
	router.HandleFunc("/", Index).Methods("GET")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
	ServerStart()
}