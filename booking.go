package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func createBooking(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	params := mux.Vars(req)
	FromRestaurant := params["restaurantname"]

	myRestaurants, err := getRestaurants()
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodGet {
		data := struct {
			User           user
			RestaurantList map[string]restaurant
			FromRestaurant string
		}{
			myUser,
			myRestaurants,
			FromRestaurant,
		}
		tpl.ExecuteTemplate(res, "booking.gohtml", data)
	}

	if req.Method == http.MethodPost {
		var myBooking booking
		myBooking.Username = req.FormValue("username")
		myBooking.RestaurantName = req.FormValue("restaurantname")
		myBooking.Date = req.FormValue("date")
		myBooking.StartTime = req.FormValue("time")

		pax := req.FormValue("pax")
		ipax, _ := strconv.Atoi(pax)
		myBooking.Pax = ipax

		myTables, err := getTables(myBooking.RestaurantName)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		tablechoice := req.FormValue("tablechoice")

		if tablechoice == "" {

			var toremove int

			query := "SELECT TableID FROM bookings WHERE RestaurantName=? AND Date=? AND StartTime=? AND Status='Reserved' AND deletedAt IS NULL"

			results, err := db.Query(query,
				myBooking.RestaurantName,
				myBooking.Date,
				myBooking.StartTime,
			)
			if err != nil {
				if err != sql.ErrNoRows {
					fmt.Println(err)
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
			defer results.Close()
			for results.Next() {
				err := results.Scan(&toremove)
				if err != nil {
					if err != sql.ErrNoRows {
						fmt.Println(err)
						http.Error(res, "Internal server error", http.StatusInternalServerError)
						return
					}
				}
				//delete tables that are reserved based on user criterias
				delete(myTables, toremove)
			}

			//delete tables that are < user intended pax
			for k, v := range myTables {
				if v.Seats < ipax {
					delete(myTables, k)
				}
			}

			data := struct {
				User           user
				RestaurantList map[string]restaurant
				FromRestaurant string
				Booking        booking
				Tables         map[int]table
			}{
				myUser,
				myRestaurants,
				FromRestaurant,
				myBooking,
				myTables,
			}
			tpl.ExecuteTemplate(res, "seatselection.gohtml", data)

		} else {
			itablechoice, _ := strconv.Atoi(tablechoice)
			myBooking.TableID = itablechoice

			var checker int
			query := "SELECT BookingID FROM bookings WHERE Username=? AND RestaurantName=? AND Date=? AND StartTime=? AND TableID=? AND Status='Reserved' AND deletedAt IS NULL"

			err := db.QueryRow(query,
				myBooking.Username,
				myBooking.RestaurantName,
				myBooking.Date,
				myBooking.StartTime,
				myBooking.TableID,
			).Scan(&checker)

			if err != nil {
				if err != sql.ErrNoRows {
					fmt.Println(err)
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				} else {
					if checker == 0 {
						err = insertBooking(myBooking)
						if err != nil {
							fmt.Println(err)
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return
						} else {
							fmt.Println("New Booking Created")
							tpl.ExecuteTemplate(res, "bookingsuccess.gohtml", myBooking)
						}
					}
				}
			} else {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
				//http.Redirect(res, req, "/booking", http.StatusSeeOther)
				//tpl.ExecuteTemplate(res, "bookingfail.gohtml", myUser)
			}
		}
	}
}

/*func updateBooking(myBooking booking) error {
	_, err := db.Exec("UPDATE INTO bookings (Username, RestaurantName, Date, StartTime, Pax, TableID, Status, createdAt) VALUES (?,?,?,?,?,?,?,?)",
		myBooking.Username,
		myBooking.RestaurantName,
		myBooking.Date,
		myBooking.StartTime,
		myBooking.Pax,
		myBooking.TableID,
		"Reserved",
		time.Now())
	if err != nil {
		return err
	}
	return nil
}*/

func insertBooking(myBooking booking) error {
	_, err := db.Exec("INSERT INTO bookings (Username, RestaurantName, Date, StartTime, Pax, TableID, Status, createdAt) VALUES (?,?,?,?,?,?,?,?)",
		myBooking.Username,
		myBooking.RestaurantName,
		myBooking.Date,
		myBooking.StartTime,
		myBooking.Pax,
		myBooking.TableID,
		"Reserved",
		time.Now())
	if err != nil {
		return err
	}
	return nil
}

/*func getBooking(restaurantname string) (map[string]booking, error) {
	myBookings := make(map[string]booking)
	var myBooking booking

	query := "SELECT BookingID, Username, Date, StartTime, Pax, TableID FROM bookings WHERE RestaurantName=? AND deletedAt IS NULL"
	results, err := db.Query(query, restaurantname)
	if err != nil {
		return myBookings, err
	}
	defer results.Close()
	for results.Next() {
		err := results.
			Scan(&myBooking.BookingID,
				&myBooking.Username,
				&myBooking.Date,
				&myBooking.StartTime,
				&myBooking.Pax,
				&myBooking.TableID,
			)
		if err != nil {
			return myBookings, err
		}
		myBookings[myBooking.Date] = myBooking
	}
	return myBookings, nil
}*/

/*func viewBooking(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	params := mux.Vars(req)

	var myRestaurant restaurant
	var myBookings = map[string]booking{}

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

	myBookings, err = getBooking(params["restaurantname"])
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			http.Error(res, "Booking Doesnt Exist", http.StatusForbidden)
			return
		}
	}

	data := struct {
		User       user
		Restaurant restaurant
		Bookings   map[string]booking
	}{
		myUser,
		myRestaurant,
		myBookings,
	}

	tpl.ExecuteTemplate(res, "viewBooking.gohtml", data)
}*/

func viewBooking(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)

	myBooking, err := getBooking(myUser.Username)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data := struct {
		User    user
		Booking map[int]booking
	}{
		myUser,
		myBooking,
	}
	tpl.ExecuteTemplate(res, "viewBooking.gohtml", data)
}

func getBooking(Username string) (map[int]booking, error) {
	var myBookings = map[int]booking{}
	var myBooking booking

	query := "SELECT BookingID, Username, RestaurantName, Date, StartTime, Pax, TableID FROM bookings WHERE Username=? AND deletedAt IS NULL"

	results, err := db.Query(query, Username)
	if err != nil {
		return myBookings, err
	}

	defer results.Close()
	for results.Next() {
		err := results.Scan(
			&myBooking.BookingID,
			&myBooking.Username,
			&myBooking.RestaurantName,
			&myBooking.Date,
			&myBooking.StartTime,
			&myBooking.Pax,
			&myBooking.TableID)
		if err != nil {
			return myBookings, err
		}
		myBookings[myBooking.BookingID] = myBooking
	}
	return myBookings, err
}
