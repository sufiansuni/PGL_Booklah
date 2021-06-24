package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type booking struct {
	BookingID string
	Username  string
	Name      string
	TableID   string
	Pax       int
	Start     time.Time
	End       time.Time
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

func main() {
	db, err := sql.Open("mysql", "newuser:password@tcp(127.0.0.1:55033)/my_restaurant")
	defer db.Close()

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}

	//InsertRecord(db, "ABC", "001", "Available", 10)
	GetRecords(db)
}
