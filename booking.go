package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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
func editBooking(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	//retrieve initial data
	params := mux.Vars(req)
	BookingIDtoEdit := params["BookingID"]
	myRestaurants, _ := getRestaurants()
	myBookings, err := getBooking(myUser.Username)
	var myBooking booking
	iBooking, _ := strconv.Atoi(BookingIDtoEdit)

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodGet {
		myBooking = myBookings[iBooking]
		data := struct {
			User           user
			Booking        booking
			RestaurantList map[string]restaurant
		}{
			myUser,
			myBooking,
			myRestaurants,
		}

		tpl.ExecuteTemplate(res, "updatebooking.gohtml", data)
		return
	}

	if req.Method == http.MethodPost {
		var myBooking booking
		myBooking.BookingID = iBooking
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

			query := "SELECT TableID FROM bookings WHERE RestaurantName=? AND Date=? AND StartTime=? AND Status='Reserved' AND deletedAt IS NULL AND TableID!=?"

			results, err := db.Query(query,
				myBooking.RestaurantName,
				myBooking.Date,
				myBooking.StartTime,
				myBooking.TableID,
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
				Booking        booking
				Tables         map[int]table
				RestaurantList map[string]restaurant
			}{
				myUser,
				myBooking,
				myTables,
				myRestaurants,
			}
			tpl.ExecuteTemplate(res, "updateseatselection.gohtml", data)
			return

		} else {
			itablechoice, _ := strconv.Atoi(tablechoice)
			myBooking.TableID = itablechoice

			var checker int
			query := "SELECT BookingID FROM bookings WHERE Username=? AND RestaurantName=? AND Date=? AND StartTime=? AND TableID=? AND Status='Reserved' AND deletedAt IS NULL AND BookingID!=?"

			err := db.QueryRow(query,
				myBooking.Username,
				myBooking.RestaurantName,
				myBooking.Date,
				myBooking.StartTime,
				myBooking.TableID,
				myBooking.BookingID,
			).Scan(&checker)

			if err != nil {
				if err != sql.ErrNoRows {
					fmt.Println(err)
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				} else {
					if checker == 0 {
						err = updateBooking(myBooking)
						if err != nil {
							fmt.Println(err)
							http.Error(res, "Internal server error", http.StatusInternalServerError)
							return
						} else {
							fmt.Println("Booking Updated Successfully")
							tpl.ExecuteTemplate(res, "updatebookingsuccess.gohtml", myBooking)
							return
						}
					}
				}
			} else {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
				//http.Redirect(res, req, "/updatebooking", http.StatusSeeOther)
				//tpl.ExecuteTemplate(res, "updatebookingfail.gohtml", myUser)
			}
		}
	}
}
func deleteBooking(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	//retrieve initial data
	params := mux.Vars(req)
	BookingIDtoDel := params["BookingID"]
	iDelete, _ := strconv.Atoi(BookingIDtoDel)
	myBookings, err := getBookings(iDelete)
	var myBooking booking

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	if req.Method == http.MethodGet {
		myBooking = myBookings[iDelete]
		data := struct {
			User    user
			Booking booking
		}{
			myUser,
			myBooking,
		}

		tpl.ExecuteTemplate(res, "deletebooking.gohtml", data)
		return
	}
}

func updateBooking(myBooking booking) error {
	_, err := db.Exec("UPDATE bookings SET Username=?, RestaurantName=?, Date=?, StartTime=?, Pax=?, TableID=?, Status=?, updatedAt=? WHERE BookingID=?",
		myBooking.Username,
		myBooking.RestaurantName,
		myBooking.Date,
		myBooking.StartTime,
		myBooking.Pax,
		myBooking.TableID,
		"Reserved",
		time.Now(),
		myBooking.BookingID)
	if err != nil {
		return err
	}
	return nil
}

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

func getBookings(BookingID int) (map[int]booking, error) {
	var myBookings = map[int]booking{}
	var myBooking booking

	query := "SELECT BookingID, Username, RestaurantName, Date, StartTime, Pax, TableID FROM bookings WHERE Username=? AND deletedAt IS NULL"

	results, err := db.Query(query, BookingID)
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
