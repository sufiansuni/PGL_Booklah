package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//connects to the database according to set values
func connectDatabase() {
	user := "root"
	password := "password"
	hostname := "127.0.0.1:3306"
	dbname := "booklah"
	connectionstring := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, hostname, dbname)

	connectionstring += "?parseTime=True&loc=Local" //additional parameters
	var err error
	db, err = sql.Open("mysql", connectionstring)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database Opened")
	}
}

func pingDatabase() {
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database Ping Successful")
	}
}

func createUserTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"users" +
		" (" +
		"Username VARCHAR(255) PRIMARY KEY, " +
		"Password BLOB, " +
		"Type VARCHAR(255), " +
		"createdAt DATETIME, " +
		"updatedAt DATETIME, " +
		"deletedAt DATETIME" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: users")
	}
}

func createSessionTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"sessions" +
		" (" +
		"UUID VARCHAR(255) PRIMARY KEY, " +
		"Username VARCHAR(255), " +
		"createdAt DATETIME, " +
		"updatedAt DATETIME, " +
		"deletedAt DATETIME" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: sessions")
	}
}

func createRestaurantTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"restaurants" +
		" (" +
		"RestaurantName VARCHAR(255) PRIMARY KEY, " +
		"CurrentPax INT, " +
		"createdAt DATETIME, " +
		"updatedAt DATETIME, " +
		"deletedAt DATETIME" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: restaurants")
	}
}

func createTableTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"tables" +
		" (" +
		"TableID INT PRIMARY KEY AUTO_INCREMENT, " +
		"RestaurantName VARCHAR(255), " +
		"TableIndex INT, " +
		"Seats INT, " +
		"Occupied INT, " +
		"createdAt DATETIME, " +
		"updatedAt DATETIME, " +
		"deletedAt DATETIME" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: tables")
	}
}

func createBookingTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"bookings" +
		" (" +
		"BookingID INT PRIMARY KEY AUTO_INCREMENT, " +
		"Username VARCHAR(255), " +
		"RestaurantName VARCHAR(255), " +
		"Date VARCHAR(255)," +
		"StartTime VARCHAR(255), " +
		"Pax INT, " +
		"TableID INT, " +
		"Status VARCHAR(255), " +
		"createdAt DATETIME, " +
		"updatedAt DATETIME, " +
		"deletedAt DATETIME" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: bookings")
	}
}
