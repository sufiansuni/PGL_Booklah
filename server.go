package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startServer() {
	r := mux.NewRouter() //New Router Instance
	handlers(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handlers(r *mux.Router) {
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
	r.HandleFunc("/restaurants/{restaurantname}/search", searchRestaurants)
	r.HandleFunc("/restaurants/{restaurantname}/edit", editRestaurant)
	r.HandleFunc("/restaurants/{restaurantname}/delete", deleteRestaurant)

	//booking handlers
	r.HandleFunc("/booking", createBooking)
	r.HandleFunc("/viewBooking", viewBooking)
	r.HandleFunc("/viewBooking/{BookingID}", editBooking)
	r.HandleFunc("/viewBooking/{BookingID}/delete", deleteBooking)
	r.HandleFunc("/restaurants/{restaurantname}/booking", createBooking)

}
