package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username  string //primary key
	Password  []byte
	Type      string
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

type session struct {
	UUID      string //primary key
	Username  string //foreign key
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

// The following maps are no longer used
// var mapUsers = map[string]user{}
// var mapSessions = map[string]string{}

func createAdminAccount() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	myUser := user{
		Username: "admin",
		Password: bPassword,
		Type:     "admin",
	}
	err := insertUser(myUser) //previously mapUsers["admin"] = myUser
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Admin Account Created")
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	var myUser user
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")

		if username != "" {
			// check if username exist/ taken
			var query string

			err := db.QueryRow("SELECT Username FROM users WHERE Username=? AND deletedAt IS NULL", username).Scan(&query)
			if err == sql.ErrNoRows {
				//
			} else if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}

			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			err = insertSession(myCookie.Value, username) // previously: mapSessions[myCookie.Value] = username
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Session Created")
			}

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = user{
				Username: username,
				Password: bPassword,
				Type:     "regular",
			}

			err = insertUser(myUser) // previouslymapUsers[username] = myUser
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("User Created:", username)
			}
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

func login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		var query user
		err := db.QueryRow("SELECT Username, Password FROM users WHERE Username=? AND deletedAt IS NULL", username).Scan(
			&query.Username,
			&query.Password,
		)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		// Matching of password entered
		err = bcrypt.CompareHashAndPassword(query.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		err = insertSession(myCookie.Value, username) // previously: mapSessions[myCookie.Value] = username
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Session Created")
			http.Redirect(res, req, "/", http.StatusSeeOther)
		}
		return
	}
	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	_, err := db.Exec("UPDATE sessions SET deletedAt=? WHERE UUID=? AND deletedAt IS NULL", time.Now(), myCookie.Value)
	if err != nil {
		fmt.Println(err)
	}
	// previously: delete(mapSessions, myCookie.Value)

	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func checkUser(res http.ResponseWriter, req *http.Request) user {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}

	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var query string
	var myUser user

	err = db.QueryRow("SELECT Username FROM sessions WHERE UUID=? AND deletedAt IS NULL", myCookie.Value).Scan(&query)
	if err == sql.ErrNoRows {
		//
	} else if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
	} else {
		err = db.QueryRow("SELECT Username, Type FROM users WHERE Username=? AND deletedAt IS NULL", query).Scan(
			&myUser.Username,
			&myUser.Type,
		)
		if err == sql.ErrNoRows {
			//
		} else if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		}
	}
	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	var query string

	err = db.QueryRow("SELECT Username FROM sessions WHERE UUID=? AND deletedAt IS NULL", myCookie.Value).Scan(&query)
	if err == sql.ErrNoRows {
		//
	} else if err != nil {
		fmt.Print(err)
	} else {
		err = db.QueryRow("SELECT Username FROM users WHERE Username=? AND deletedAt IS NULL", query).Scan(&query)
		if err == sql.ErrNoRows {
			//
		} else if err != nil {
			fmt.Print(err)
		} else {
			return true
		}
	}
	return false
}

func insertUser(myUser user) error {
	_, err := db.Exec("INSERT INTO users (Username, Password, Type, createdAt) VALUES (?,?,?,?)",
		myUser.Username, myUser.Password, myUser.Type, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func insertSession(cookievalue string, username string) error {
	_, err := db.Exec("INSERT INTO sessions(UUID, Username, createdAt) VALUES(?, ?, ?)",
		cookievalue, username, time.Now())
	if err != nil {
		return err
	}
	return nil
}
