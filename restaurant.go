package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type restaurant struct {
	RestaurantName string
	Tables         []table
	A_Tables       int //available tables
}

type table struct {
	Seats  int
	Status string
}

var mapRestaurants = map[string]restaurant{}

func indexRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	refreshTables(mapRestaurants)
	data := struct {
		RestaurantList map[string]restaurant
	}{
		mapRestaurants,
	}
	tpl.ExecuteTemplate(res, "restaurants.html", data)
}

func createNewRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	var myRestaurant restaurant
	var myTables []table
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		restaurantname := req.FormValue("restaurantname")
		if restaurantname != "" {
			// check if restaurant exist/ taken
			if _, ok := mapRestaurants[restaurantname]; ok {
				http.Error(res, "Restaurant name already taken", http.StatusForbidden)
				return
			}

			myRestaurant = restaurant{restaurantname, myTables, 0}
			mapRestaurants[restaurantname] = myRestaurant
			fmt.Println(mapRestaurants[restaurantname])
		}
		// redirect to main index
		http.Redirect(res, req, "/restaurants", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "restaurantnew.html", myRestaurant)
}

func viewRestaurant(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	myRestaurant := mapRestaurants[params["restaurantname"]]
	tpl.ExecuteTemplate(res, "restaurantpage.html", myRestaurant)
}

func editRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }
}

func deleteRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }
}

func refreshTables(restaurantmap map[string]restaurant) {
	for _, v := range restaurantmap {
		count := 0
		for _, tables := range v.Tables {
			if tables.Status == "Available" {
				count++
			}
		}
		v.A_Tables = count
	}
}
