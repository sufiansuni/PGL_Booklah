package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type restaurants struct {
	ID        string
	Name      string
	A_table   int //available tables
	B_table   int //Booked tables
	Max_table int //Max tables
}

func GetRecords(db *sql.DB) {

	results, err := db.Query("Select * FROM my_restaurant.Restaurants")
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

func InsertRecord(db *sql.DB, ID string, Name string, A_table int, B_table int, Max_table int) {
	results, err := db.Exec("INSERT INTO my_restaurant.Restaurants VALUES (?,?,?,?,?)",
		ID, Name, A_table, B_table, Max_table)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func EditRecord(db *sql.DB, ID string, Name string, A_table int, B_table int, Max_table int) {
	results, err := db.Exec(
		"UPDATE Course SET Name=?, A_table=? , B_table=? , Max_table=? WHERE ID=?",
		Name, A_table, B_table, Max_table, ID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func DeleteRecord(db *sql.DB, Name string) {
	results, err := db.Exec("DELETE FROM Restaurants WHERE Name=?", Name)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}

func main() {
	db, err := sql.Open("mysql", "newuser:password@tcp(127.0.0.1:64810)/my_restaurant")
	defer db.Close()

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}

	//InsertRecord(db, "001", "ABC", 10, 5, 15)
	GetRecords(db)
}
