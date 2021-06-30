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
	Occupied       int
}

type booking struct {
	BookingID      int    //primary key
	Username       string //foreign key
	RestaurantName string //foreign key
	Date           string
	StartTime      string
	Pax            int
	TableID        int //foreign key
	Status         string
}

func indexRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	myRestaurants, err := getRestaurants()
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return

	}

	// separating GET and POST

	if req.Method == http.MethodGet {
		data := struct {
			User           user
			RestaurantList map[string]restaurant
		}{
			myUser,
			myRestaurants,
		}
		tpl.ExecuteTemplate(res, "restaurants.gohtml", data)
		return

	}

	if req.Method == http.MethodPost {
		var myfilteredRestaurants = map[string]restaurant{}
		var myTable table

		// get form values
		Quantity := req.FormValue("Quantity")

		if Quantity == "" {
			//look at table database
			query := "SELECT RestaurantName FROM tables WHERE Seats >=? AND deletedAt IS NULL"

			// pass in Quantity variable
			results, err := db.Query(query, Quantity)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer results.Close()
			for results.Next() {
				//store info from table database into my own variable
				err := results.Scan(&myTable.RestaurantName)
				if err != nil {
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				}
				//for all the restaurant names I received,
				//I will pick out from the entire list that I have, retrieved earlier in code
				myfilteredRestaurants[myTable.RestaurantName] = myRestaurants[myTable.RestaurantName]
			}

			data := struct {
				User           user
				RestaurantList map[string]restaurant
			}{
				myUser,
				myfilteredRestaurants,
			}
			tpl.ExecuteTemplate(res, "restaurants.gohtml", data)

		} else {
			http.Redirect(res, req, "/restaurants", http.StatusSeeOther)
		}
		return
	}
}

//retrieve your search boxes/FormValues

//use those form values to query the database

//find tables from tables database that has >= Pax
// you want restaurantname from tables database that Seats >= Pax

// make a new filtered restaurant map, myRestaurantsFiltered
// for k, v := range myRestaurants{
//myRestaurantsFiltered[restaurantname] = myrestaurants[restaurantname]
//}

func createNewRestaurant(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	if myUser.Type != "admin" {
		http.Error(res, "Access Unauthorized", http.StatusUnauthorized)
		return
	}

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
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
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
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return
						}

						err = insertTable(myTable)
						if err != nil {
							fmt.Println(err)
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return
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

	if myUser.Type != "admin" {
		http.Error(res, "Access Unauthorized", http.StatusUnauthorized)
		return
	}

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
						return
					}

					newTableSeatsConv, err := strconv.Atoi(newTableSeats)
					if err != nil {
						http.Error(res, "Internal server error", http.StatusInternalServerError)
						return
					}

					checker := 0
					query := "SELECT TableID FROM tables WHERE RestaurantName=? AND TableIndex=? AND deletedAt IS NULL"
					err = db.QueryRow(query, restaurantname, newTableIndexConv).Scan(&checker)
					if err != nil {
						if err != sql.ErrNoRows {
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return

						}
					}

					if checker != 0 {
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

					} else {
						var myTable table
						myTable.RestaurantName = restaurantname
						myTable.TableIndex = newTableIndexConv
						myTable.Seats = newTableSeatsConv

						err := insertTable(myTable)
						if err != nil {
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return
						}
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
	myUser := checkUser(res, req)

	if myUser.Type != "admin" {
		http.Error(res, "Access Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(req)

	// previously: delete(mapRestaurants, params["restaurantname"])
	statement := "UPDATE restaurants SET deletedAt=? WHERE RestaurantName=?"
	_, err := db.Exec(statement, time.Now(), params["restaurantname"])
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

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
	_, err := db.Exec("INSERT INTO tables (RestaurantName, TableIndex, Seats, Occupied, createdAt) VALUES (?,?,?,?,?)",
		myTable.RestaurantName,
		myTable.TableIndex,
		myTable.Seats,
		0,
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

	query := "SELECT TableID, TableIndex, Seats, Occupied FROM tables WHERE RestaurantName=? AND deletedAt IS NULL"
	results, err := db.Query(query, restaurantname)
	if err != nil {
		return myTables, err
	}
	defer results.Close()
	for results.Next() {
		err := results.
			Scan(&myTable.TableID,
				&myTable.TableIndex,
				&myTable.Seats,
				&myTable.Occupied,
			)
		if err != nil {
			return myTables, err
		}
		myTables[myTable.TableID] = myTable
	}
	return myTables, nil
}

func getRestaurants() (map[string]restaurant, error) {
	var myRestaurants = map[string]restaurant{}
	var myRestaurant restaurant

	query := "SELECT RestaurantName FROM restaurants WHERE deletedAt IS NULL"

	results, err := db.Query(query)
	if err != nil {
		return myRestaurants, err
	}

	defer results.Close()
	for results.Next() {
		err := results.Scan(&myRestaurant.RestaurantName)
		if err != nil {
			return myRestaurants, err
		}
		myRestaurants[myRestaurant.RestaurantName] = myRestaurant
	}
	return myRestaurants, err
}
