package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DNS = "user:password@tcp(127.0.0.1:8081)/booklah_db?charset=utf8mb4&parseTime=True&loc=Local"

type User struct {
	UserID    uint `gorm:"primaryKey"`
	Username  string
	Password  []byte
	First     string
	Last      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	//create tables at databse based on struct
	DB.AutoMigrate(&User{})

	//create admin account
	var finder User
	DB.Unscoped().First(&finder, "username = ?", "admin")
	if reflect.ValueOf(finder).IsZero() {
		bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
		user := User{Username: "admin", Password: bPassword, First: "admin", Last: "admin"}

		result := DB.Create(&user) // pass pointer of data to Create

		fmt.Println("ID added:", user.UserID)             // returns inserted data's primary key
		fmt.Println("Error:", result.Error)               // returns error
		fmt.Println("RowsAffected:", result.RowsAffected) // returns inserted records count

	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	DB.First(&user, "user_id = ?", params["user_id"])
	if reflect.ValueOf(user).IsZero() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user found"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	var finder User
	DB.Unscoped().First(&finder, "username = ?", user.Username)

	//if username not in database
	if finder.Username != user.Username {
		if reflect.ValueOf(user).IsZero() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply user information in JSON format"))
		} else {
			DB.Create(&user)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - User added: " + string(rune(user.UserID))))
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 - Duplicate User"))
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	var finder User
	DB.Unscoped().First(&finder, "user_id = ?", params["user_id"])

	//if user_id not in database
	if reflect.ValueOf(finder).IsZero() {
		json.NewDecoder(r.Body).Decode(&user)
		//if course title is provided
		if reflect.ValueOf(user).IsZero() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply user information in JSON format"))

		} else {
			DB.Create(&user)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - User added: " + string(rune(user.UserID))))
		}

	} else { // update courses if existing entry found
		DB.First(&user, "user_id = ?", params["user_id"])
		json.NewDecoder(r.Body).Decode(&user)
		DB.Save(&user)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User updated successfully"))
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	DB.First(&user, "user_id = ?", params["user_id"])
	if reflect.ValueOf(user).IsZero() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user found"))

	} else {
		DB.Delete(&user, "user_id = ?", params["user_id"])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully"))
	}
}
