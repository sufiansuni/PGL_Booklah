package main

import (
	"text/template"
)

var tpl *template.Template

var mapRestaurants = map[string]restaurant{}
var mapBookings = map[string]booking{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	connectDatabase()
	defer db.Close()
	pingDatabase()

	createUserTable()
	createSessionTable()
	createRestaurantTable()
	createTableTable()
	createBookingTable()

	createAdminAccount()

	startServer()
}
