package main

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	gormDB.First(&user, "user_id = ?", params["user_id"])
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
	gormDB.Unscoped().First(&finder, "username = ?", user.Username)

	//if username not in database
	if finder.Username != user.Username {
		if reflect.ValueOf(user).IsZero() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply user information in JSON format"))
		} else {
			gormDB.Create(&user)
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
	gormDB.Unscoped().First(&finder, "user_id = ?", params["user_id"])

	//if user_id not in database
	if reflect.ValueOf(finder).IsZero() {
		json.NewDecoder(r.Body).Decode(&user)
		//if course title is provided
		if reflect.ValueOf(user).IsZero() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply user information in JSON format"))

		} else {
			gormDB.Create(&user)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - User added: " + string(rune(user.UserID))))
		}

	} else { // update courses if existing entry found
		gormDB.First(&user, "user_id = ?", params["user_id"])
		json.NewDecoder(r.Body).Decode(&user)
		gormDB.Save(&user)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User updated successfully"))
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	gormDB.First(&user, "user_id = ?", params["user_id"])
	if reflect.ValueOf(user).IsZero() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user found"))

	} else {
		gormDB.Delete(&user, "user_id = ?", params["user_id"])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully"))
	}
}
