package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

var tpl *template.Template
var mapUsers = map[string]user{}
var mapSessions = map[string]string{}

type user struct {
	Username string
	Password []byte
	First    string
	Last     string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = user{"admin", bPassword, "admin", "admin"}
}

func InitializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index)
	r.HandleFunc("/signup", Signup)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.Handle("/favicon.ico", http.NotFoundHandler())

	r.HandleFunc("/api/v1/users/{user_id}", GetUser).Methods("GET")
	r.HandleFunc("/api/v1/users", CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users/{user_id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/api/v1/users/{user_id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	InitialMigration()
	InitializeRouter()
}
