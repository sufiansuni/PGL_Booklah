package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type restaurant struct {
	RestaurantName string //primary key
	//address
	//hours
	//contact
	//summary
}

// Potential New Staff Management Feature

// type staff struct {
// 	ID             int    //primarykey
// 	Username       string //foreign key
// 	RestaurantName string //foreign key
// 	Position       string
// 	createdAt      time.Time
// 	updatedAt      time.Time
// 	deletedAt      time.Time
// }

type table struct {
	TableID        int    //primary key
	RestaurantName string //foreign key
	TableIndex     int
	Seats          int
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
}

var mapRestaurants = map[string]restaurant{}
var mapBookings = map[string]booking{}

func indexRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	var myRestaurants = map[string]restaurant{}
	var myRestaurant restaurant

	query := "SELECT RestaurantName FROM restaurants WHERE deletedAt IS NULL"

	results, err := db.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	defer results.Close()
	for results.Next() {
		err := results.Scan(&myRestaurant.RestaurantName)
		if err != nil {
			panic("error getting results from sql select")
		}
		myRestaurants[myRestaurant.RestaurantName] = myRestaurant
	}

	data := struct {
		User           user
		RestaurantList map[string]restaurant
	}{
		myUser,
		myRestaurants,
	}
	tpl.ExecuteTemplate(res, "restaurants.gohtml", data)
}

func createNewRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	var myRestaurant restaurant

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		restaurantname := req.FormValue("restaurantname")

		if restaurantname != "" {
			// check if restaurant exist/ taken
			var checker string

			query := "SELECT RestaurantName FROM restaurants WHERE RestaurantName=? AND deletedAt IS NULL"
			err := db.QueryRow(query, restaurantname).Scan(&checker)

			if err != nil {
				if err != sql.ErrNoRows {
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(res, "Restaurant already taken", http.StatusForbidden)
				return
			}

			myRestaurant.RestaurantName = restaurantname
			err = insertRestaurant(myRestaurant) //previously: mapRestaurants[restaurantname] = myRestaurant
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("New Restaurant Created:", myRestaurant.RestaurantName)
			}

			for i := 1; i < 21; i++ {
				var myTable table
				iString := strconv.Itoa(i)
				myTable.RestaurantName = restaurantname

				tableindexform := req.FormValue("table" + iString)
				if tableindexform != "" {
					myTable.TableIndex, err = strconv.Atoi(tableindexform)
					if err != nil {
						fmt.Println(err)
					}

					tableseatsform := req.FormValue("table" + iString + "seats")
					if tableseatsform != "" {
						myTable.Seats, err = strconv.Atoi(tableseatsform)
						if err != nil {
							fmt.Println(err)
						}

						err = insertTable(myTable)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("Table Created for", restaurantname)
						}
					}
				}
			}
		}
		// redirect to main index
		http.Redirect(res, req, "/restaurants", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "restaurantnew.gohtml", myUser)
}

func viewRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	params := mux.Vars(req)

	var myRestaurant restaurant
	var myTables = map[int]table{}

	myRestaurant, err := getRestaurant(params["restaurantname"])
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			http.Error(res, "Restaurant Doesnt Exist", http.StatusForbidden)
			return
		}
	}

	myTables, err = getTables(params["restaurantname"])
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			http.Error(res, "Tables Doesnt Exist", http.StatusForbidden)
			return
		}
	}

	data := struct {
		User       user
		Restaurant restaurant
		Tables     map[int]table
	}{
		myUser,
		myRestaurant,
		myTables,
	}
	tpl.ExecuteTemplate(res, "restaurantpage.gohtml", data)
}

func editRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	//retrieve initial data
	params := mux.Vars(req)

	var myRestaurant restaurant
	var myTables = map[int]table{}

	myRestaurant, err := getRestaurant(params["restaurantname"])
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			http.Error(res, "Restaurant Doesnt Exist", http.StatusForbidden)
			return
		}
	}

	myTables, err = getTables(params["restaurantname"])
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			http.Error(res, "Tables Doesnt Exist", http.StatusForbidden)
			return
		}
	}

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		restaurantname := req.FormValue("restaurantname")
		if restaurantname != "" {
			// check if restaurant exist/ taken

			restaurantCheck, err := getRestaurant(restaurantname)
			if err != nil {
				if err != sql.ErrNoRows {
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				}
			} else {
				if restaurantCheck.RestaurantName != params["restaurantname"] {
					http.Error(res, "Restaurant already taken", http.StatusForbidden)
					return
				}
			}

			statement := "UPDATE restaurants SET RestaurantName=?, updatedAt=? WHERE RestaurantName=?"
			_, err = db.Exec(statement, restaurantname, time.Now(), params["restaurantname"])
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			for i := 1; i < 21; i++ {
				iString := strconv.Itoa(i)
				newTableIndex := req.FormValue("table" + iString)
				newTableSeats := req.FormValue("table" + iString + "seats")

				if newTableIndex != "" && newTableSeats != "" {
					newTableIndexConv, err := strconv.Atoi(newTableIndex)
					if err != nil {
						http.Error(res, "Internal server error", http.StatusInternalServerError)
					}

					newTableSeatsConv, err := strconv.Atoi(newTableSeats)
					if err != nil {
						http.Error(res, "Internal server error", http.StatusInternalServerError)
					}

					statement := "UPDATE tables SET Seats=?, updatedAt=? WHERE RestaurantName=? AND TableIndex=?"
					_, err = db.Exec(statement,
						newTableSeatsConv,
						time.Now(),
						restaurantname,
						newTableIndexConv)
					if err != nil {
						http.Error(res, "Internal server error", http.StatusInternalServerError)
						return
					}

				}
			}
		}
		// redirect to main index
		http.Redirect(res, req, "/restaurants/"+restaurantname, http.StatusSeeOther)
		return
	}
	data := struct {
		User       user
		Restaurant restaurant
		Tables     map[int]table
	}{
		myUser,
		myRestaurant,
		myTables,
	}
	tpl.ExecuteTemplate(res, "restaurantedit.gohtml", data)
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

func insertRestaurant(myRestaurant restaurant) error {
	_, err := db.Exec("INSERT INTO restaurants (RestaurantName, createdAt) VALUES (?,?)",
		myRestaurant.RestaurantName,
		time.Now())
	if err != nil {
		return err
	}
	return nil
}

func insertTable(myTable table) error {
	_, err := db.Exec("INSERT INTO tables (RestaurantName, TableIndex, Seats, createdAt) VALUES (?,?,?,?)",
		myTable.RestaurantName,
		myTable.TableIndex,
		myTable.Seats,
		time.Now())
	if err != nil {
		return err
	}
	return nil
}

func getRestaurant(restaurantname string) (restaurant, error) {
	var myRestaurant restaurant

	query := "SELECT RestaurantName FROM restaurants WHERE RestaurantName=? AND deletedAt IS NULL"
	err := db.QueryRow(query, restaurantname).Scan(&myRestaurant.RestaurantName)

	return myRestaurant, err
}

func getTables(restaurantname string) (map[int]table, error) {
	myTables := make(map[int]table)
	var myTable table

	query := "SELECT TableIndex, Seats FROM tables WHERE RestaurantName=? AND deletedAt IS NULL"
	results, err := db.Query(query, restaurantname)
	if err != nil {
		return myTables, err
	}
	defer results.Close()
	for results.Next() {
		err := results.Scan(&myTable.TableIndex, &myTable.Seats)
		if err != nil {
			return myTables, err
		}
		myTables[myTable.TableIndex] = myTable
	}
	return myTables, nil
}
