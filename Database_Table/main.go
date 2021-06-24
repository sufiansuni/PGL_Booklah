package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type restaurants struct {
	Name     string
	TableID  string
	Booking  string
	Capacity int
}

func GetRecords(db *sql.DB) {

	results, err := db.Query("Select * FROM my_restaurant.Table")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		// map this type to the record in the table
		var Restaurants restaurants
		err = results.Scan(&Restaurants.Name, &Restaurants.TableID, &Restaurants.Booking, &Restaurants.Capacity)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(Restaurants.Name, Restaurants.TableID, Restaurants.Booking, Restaurants.Capacity)
	}
}

func InsertRecord(db *sql.DB, Name string, TableID string, Booking string, Capacity int) {
	results, err := db.Exec("INSERT INTO my_restaurant.Table VALUES (?,?,?,?)",
		Name, TableID, Booking, Capacity)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func EditRecord(db *sql.DB, Name string, TableID string, Booking string, Capacity int) {
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

func DeleteRecord(db *sql.DB, Name string) {
	results, err := db.Exec("DELETE FROM Table WHERE Name=?", Name)
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

	InsertRecord(db, "ABC", "001", "Available", 10)
	GetRecords(db)
}
