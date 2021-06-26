package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type rest_table struct {
	Name     string
	TableID  string
	Booking  string
	Capacity int
}

type booking struct {
	BookingID string
	Username  string
	Name      string
	TableID   string
	Pax       int
	Start     time.Time
	End       time.Time
}

type restaurants struct {
	ID        string
	Name      string
	A_table   int //available tables
	B_table   int //Booked tables
	Max_table int //Max tables
}

func GetRecords_table(db *sql.DB) {

	results, err := db.Query("Select * FROM my_restaurant.Table")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		// map this type to the record in the table
		var Restaurants rest_table
		err = results.Scan(&Restaurants.Name, &Restaurants.TableID, &Restaurants.Booking, &Restaurants.Capacity)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(Restaurants.Name, Restaurants.TableID, Restaurants.Booking, Restaurants.Capacity)
	}
}

func InsertRecord_table(db *sql.DB, Name string, TableID string, Booking string, Capacity int) {
	results, err := db.Exec("INSERT INTO my_restaurant.Table VALUES (?,?,?,?)",
		Name, TableID, Booking, Capacity)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func EditRecord_table(db *sql.DB, Name string, TableID string, Booking string, Capacity int) {
	results, err := db.Exec(
		"UPDATE Table SET TableID=?, Booking=? , Capacity=?  WHERE Name=?",
		TableID, Booking, Capacity, Name)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func DeleteRecord_table(db *sql.DB, Name string) {
	results, err := db.Exec("DELETE FROM Table WHERE Name=?", Name)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func GetRecords_rest(db *sql.DB) {

	results, err := db.Query("Select * FROM my_restaurant.Restaurant")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		// map this type to the record in the table
		var Restaurants restaurants
		err = results.Scan(&Restaurants.ID, &Restaurants.Name, &Restaurants.A_table, &Restaurants.B_table, &Restaurants.Max_table)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(Restaurants.ID, Restaurants.Name, Restaurants.A_table, Restaurants.B_table, Restaurants.Max_table)
	}
}

func InsertRecord_rest(db *sql.DB, ID string, Name string, A_table int, B_table int, Max_table int) {
	results, err := db.Exec("INSERT INTO my_restaurant.Restaurant VALUES (?,?,?,?,?)",
		ID, Name, A_table, B_table, Max_table)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func EditRecord_rest(db *sql.DB, ID string, Name string, A_table int, B_table int, Max_table int) {
	results, err := db.Exec(
		"UPDATE Restaurant SET Name=?, A_table=? , B_table=? , Max_table=? WHERE ID=?",
		Name, A_table, B_table, Max_table, ID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func DeleteRecord_rest(db *sql.DB, Name string) {
	results, err := db.Exec("DELETE FROM Restaurant WHERE Name=?", Name)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func GetRecords(db *sql.DB) {

	results, err := db.Query("Select * FROM my_restaurant.Booking")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		// map this type to the record in the table
		var Bookings booking
		err = results.Scan(&Bookings.BookingID, &Bookings.Username, &Bookings.Name, &Bookings.TableID, &Bookings.Pax, &Bookings.Start, &Bookings.End)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(Bookings.BookingID, Bookings.Username, Bookings.Name, Bookings.TableID, Bookings.Pax, Bookings.Start, Bookings.End)
	}
}

func InsertRecord(db *sql.DB, BookingID string, Username string, Name string, TableID string, Pax int, Start time.Time, End time.Time) {
	results, err := db.Exec("INSERT INTO my_restaurant.Booking VALUES (?,?,?,?,?,?,?)",
		BookingID, Username, Name, TableID, Pax, Start, End)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func EditRecord(db *sql.DB, BookingID string, Username string, Name string, TableID string, Pax int, Start time.Time, End time.Time) {
	results, err := db.Exec(
		"UPDATE Booking SET Username=?, Name=?, TableID =? , Pax=? , Start=?, End=? WHERE BookingID=?",
		Username, Name, TableID, Pax, Start, End, BookingID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func DeleteRecord(db *sql.DB, BookingID string) {
	results, err := db.Exec("DELETE FROM Booking WHERE BookingID=?", BookingID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}
