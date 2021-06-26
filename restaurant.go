package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type restaurant struct {
	RestaurantName string //primary key
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      time.Time
	//address
	//hours
	//contact
	//summary
}

type table struct {
	TableID        int //primary key
	RestaurantName string //foreign key
	TableIndex     int
	Seats          int
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      time.Time
}

type booking struct {
	BookingID      int    //primary key
	Username       string //foreign key
	RestaurantName string //foreign key
	Pax            int
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	TableID        int //foreign key
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      time.Time
}

func indexRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	data := struct {
		RestaurantList map[string]restaurant
	}{
		mapRestaurants,
	}
	tpl.ExecuteTemplate(res, "restaurants.gohtml", data)
}

func createNewRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	var myRestaurant restaurant
	// var myTables []table
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

			// for i := 1; i < 21; i++ {
			// 	var myTable table
			// 	iString := strconv.Itoa(i)
			// 	mySeats, _ := strconv.Atoi(req.FormValue("Table" + iString + "Seats"))
			// 	myOccupied, _ := strconv.ParseBool(req.FormValue("Table" + iString + "Occupied"))
			// 	if mySeats != 0 {
			// 		myTable = table{i, mySeats, myOccupied}
			// 		myTables = append(myTables, myTable)
			// 	}
			// }

			myRestaurant = restaurant{RestaurantName: restaurantname}
			mapRestaurants[restaurantname] = myRestaurant
			fmt.Println("New Restaurant Created:", myRestaurant.RestaurantName)
			fmt.Println(mapRestaurants[restaurantname])
		}
		// redirect to main index
		http.Redirect(res, req, "/restaurants", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "restaurantnew.gohtml", myRestaurant)
}

func viewRestaurant(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	myRestaurant := mapRestaurants[params["restaurantname"]]
	tpl.ExecuteTemplate(res, "restaurantpage.gohtml", myRestaurant)
}

func editRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	//retrieve initial data
	params := mux.Vars(req)
	myRestaurant := mapRestaurants[params["restaurantname"]]
	// var myTables []table

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		restaurantname := req.FormValue("restaurantname")
		if restaurantname != "" {
			// check if restaurant exist/ taken
			if _, ok := mapRestaurants[restaurantname]; ok {
				if params["restaurantname"] != restaurantname {
					http.Error(res, "Restaurant name already taken", http.StatusForbidden)
					return
				}
			}

			// for i := 1; i < 21; i++ {
			// 	var myTable table
			// 	iString := strconv.Itoa(i)
			// 	mySeats, _ := strconv.Atoi(req.FormValue("Table" + iString + "Seats"))
			// 	myOccupied, _ := strconv.ParseBool(req.FormValue("Table" + iString + "Occupied"))
			// 	if mySeats != 0 {
			// 		myTable = table{i, mySeats, myOccupied}
			// 		myTables = append(myTables, myTable)
			// 	}
			// }

			myRestaurant.RestaurantName = restaurantname
			// myRestaurant.Tables = myTables
			mapRestaurants[restaurantname] = myRestaurant

			if params["restaurantname"] != restaurantname {
				delete(mapRestaurants, params["restaurantname"])
				fmt.Println(params["restaurantname"], "updated to", myRestaurant.RestaurantName)
			} else {
				fmt.Println(params["restaurantname"], "updated")
			}
			fmt.Println(mapRestaurants[restaurantname])
		}
		// redirect to main index
		http.Redirect(res, req, "/restaurants/"+restaurantname, http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "restaurantedit.gohtml", myRestaurant)
}

func deleteRestaurant(res http.ResponseWriter, req *http.Request) {
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }
	params := mux.Vars(req)
	delete(mapRestaurants, params["restaurantname"])
	fmt.Println(params["restaurantname"], "deleted")

	http.Redirect(res, req, "/restaurants", http.StatusSeeOther)
}
