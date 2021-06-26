package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password []byte
	First    string
	Last     string
}

var tpl *template.Template
var mapUsers = map[string]user{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = user{"admin", bPassword, "admin", "admin"}
}

func main() {
	r := mux.NewRouter() //New Router Instance

	//loginlogout handlers
	r.HandleFunc("/", index)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.Handle("/favicon.ico", http.NotFoundHandler())

	//restaurant handlers
	r.HandleFunc("/restaurants", indexRestaurant)
	r.HandleFunc("/restaurants/new", createNewRestaurant)
	r.HandleFunc("/restaurants/{restaurantname}", viewRestaurant)
	r.HandleFunc("/restaurants/{restaurantname}/edit", editRestaurant)
	r.HandleFunc("/restaurants/{restaurantname}/delete", deleteRestaurant)

	log.Fatal(http.ListenAndServe(":8080", r))
}
