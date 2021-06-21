package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func Index(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

func Signup(res http.ResponseWriter, req *http.Request) {
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
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		if username != "" {
			// check if username exist/ taken
			if _, ok := mapUsers[username]; ok {
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
			mapSessions[myCookie.Value] = username

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = user{username, bPassword, firstname, lastname}
			mapUsers[username] = myUser
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

func Login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		myUser, ok := mapUsers[username]
		if !ok {
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func getUser(res http.ResponseWriter, req *http.Request) user {
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
	var myUser user
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username]
	return ok
}

// registerHandler serves form for registring new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// registerAuthHandler creates new user in database
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	/*
		1. check username criteria
		2. check password criteria
		3. check if username is already exists in database
		4. create bcrypt hash from password
		5. insert username and password hash in database
		(email validation will be in another video)
	*/
	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	// check username criteria
	username := r.FormValue("username")
	err := checkUsernameCriteria(username)
	if err != nil {
		tpl.ExecuteTemplate(w, "register.html", err.Error())
		return
	}
	// check password criteria
	password := r.FormValue("password")
	err = checkPasswordCriteria(password)
	if err != nil {
		tpl.ExecuteTemplate(w, "register.html", err.Error())
		return
	}
	// check email is valid
	email := r.FormValue("email")
	err = checkEmailValid(email)
	if err != nil {
		tpl.ExecuteTemplate(w, "register.html", err.Error())
		return
	}
	fmt.Println("email:", email, "is valid")
	//check if email domain exists
	err = checkEmailDomain(email)
	if err != nil {
		tpl.ExecuteTemplate(w, "register.html", err.Error())
		return
	}
	// begin transaction, every query runs successfully or else no changes are made to the database
	// func (db *DB) Begin() (*Tx, error)
	var tx *sql.Tx
	tx, err = sqlDB.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was an issue registering, please try again")
		return
	}
	// rollback will be ignored if the tx has been committed later in the function
	defer tx.Rollback()
	// check if username already exists for availability
	stmt := "SELECT id FROM users WHERE username = ?"
	row := tx.QueryRow(stmt, username)
	var uID string
	err = row.Scan(&uID)
	// we only want sql.ErrNoRows, anything else means it already exists or we encountered an error
	if err != sql.ErrNoRows {
		fmt.Println("username exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "Sorry that username is unavailable")
		// func (tx *Tx) Rollback() error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	fmt.Println("username:", username, "available")
	// create hash from password
	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "Sorry, there was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))
	// insert data into users table
	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertUserStmt *sql.Stmt
	insertUserStmt, err = tx.Prepare("INSERT INTO users (username, email, hash, created_at, is_active) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "Sorry, there was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	defer insertUserStmt.Close()
	currentTime := time.Now()
	fmt.Println("currentTime:", currentTime)
	var result sql.Result
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err = insertUserStmt.Exec(username, email, hash, currentTime, 0)
	fmt.Println("err:", err)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	// check for successfull insert
	if err != nil || rowsAff != 1 {
		fmt.Println("error inserting new user, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	// create random code for email
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(100000)
	fmt.Println("random number:", rn)
	var insertEmailVerStmt *sql.Stmt
	insertEmailVerStmt, err = tx.Prepare("INSERT INTO email_ver (username, email, ver_code) VALUES (?, ?, ?);")
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	defer insertEmailVerStmt.Close()
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err = insertEmailVerStmt.Exec(username, email, rn)
	rowsAff, _ = result.RowsAffected()
	lastIns, _ = result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	// check for successfull inserting into email_ver
	if err != nil || rowsAff != 1 {
		fmt.Println("error inserting into email_ver, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	// email the code
	err = EmailVerCode(rn, email)
	if err != nil {
		fmt.Println("error emailing verification code, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account, please try again")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	// func (tx *Tx) Commit() error
	err = tx.Commit()
	if err != nil {
		fmt.Println("error commiting changes, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	fmt.Println("successful insert queries and sent email, committing changes")
	var m Message
	m.Email = email
	tpl.ExecuteTemplate(w, "verifyemail.html", m)
}

// verifyEmailHandler serves form for registring new users
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****verifyEmailHandler running*****")
	r.ParseForm()
	// check username criteria
	email := r.FormValue("email")
	verCode := r.FormValue("vercode")
	tx, err := sqlDB.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err:", err)
		tpl.ExecuteTemplate(w, "verifyemail.html", "Sorry, there was an issue verifying email, please try again")
		return
	}
	// rollback will be ignored if the tx has been committed later in the function
	defer tx.Rollback()
	// we need to check if the verCode supplied by user in form is same as in the database
	fmt.Println("email (from form):", email)
	fmt.Println("verCode (from form):", verCode)
	stmt := "SELECT ver_code FROM email_ver WHERE email = ?"
	row := tx.QueryRow(stmt, email)
	var dbCode string
	err = row.Scan(&dbCode)
	if err != nil {
		fmt.Println("error scanning verCode err:", err)
		var m Message
		m.Email = email
		m.ErrMessage = "Sorry there was an issue verifying email, please try again"
		tpl.ExecuteTemplate(w, "verifyemail.html", m)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return
	}
	fmt.Println("dbCode:", dbCode)
	// check verification code entered into form is same as in db
	if verCode == dbCode {
		// ver_code is a match, setting account is_active to 1 (true)"
		stmt := "UPDATE users SET is_active = 1 WHERE email = ?"
		updateIsActiveStmt, err := tx.Prepare(stmt)
		if err != nil {
			fmt.Println("error preparing updateIsActiveStmt err:", err)
			var m Message
			m.Email = email
			m.ErrMessage = "Sorry, there was a problem verifying email, please try again"
			tpl.ExecuteTemplate(w, "verifyemail.html", m)
			return
		}
		defer updateIsActiveStmt.Close()
		var result sql.Result
		result, err = updateIsActiveStmt.Exec(email)
		rowsAff, _ := result.RowsAffected()
		lastIns, _ := result.LastInsertId()
		fmt.Println("rowsAff:", rowsAff)
		fmt.Println("lastIns:", lastIns)
		// check for successfull insert
		if err != nil || rowsAff != 1 {
			fmt.Println("error inserting new user, err:", err)
			var m Message
			m.Email = email
			fmt.Println("m.Email:", m.Email)
			m.ErrMessage = "There was an issue verifying email, please try again"
			tpl.ExecuteTemplate(w, "verifyemail.html", m)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
			}
			return
		}
		tpl.ExecuteTemplate(w, "login.html", "email verified, go ahead and login")
		tx.Commit()
		return
	}
	var m Message
	m.ErrMessage = "There was an issue verifying email, please try again"
	m.Email = email
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
	}
	tpl.ExecuteTemplate(w, "verifyemail.html", m)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// insert login logic here
	fmt.Fprint(w, "congrats, you are logged in")
}

func EmailVerCode(rn int, toEmail string) error {
	// sender data
	from := "wld2046@gmail.com"    //ex: "John.Doe@gmail.com"
	password := "snhovfquxxzoysse" // ex: "ieiemcjdkejspqz"
	// receiver address privided through toEmail argument
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Subject: Email Verification Code\r\n\r\n"
	verCode := strconv.Itoa(rn)
	fmt.Println("verCode:", verCode)
	body := "verification code: " + verCode
	fmt.Println("body:", body)
	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	fmt.Println("message:", string(message))
	err := smtp.SendMail(address, auth, from, to, message)
	return err
}

func checkUsernameCriteria(username string) error {
	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	if !nameAlphaNumeric {
		// func New(text string) error
		return errors.New("username must only contain letters and numbers")
	}
	// check username length
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	if !nameLength {
		return errors.New("username must be longer than 4 characters and less than 51")
	}
	return nil
}

func checkPasswordCriteria(password string) error {
	var err error
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
			err = errors.New("Pa")
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	// check password length
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	// create error for any criteria not passed
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		switch false {
		case pswdLowercase:
			err = errors.New("password must contain atleast one lower case letter")
		case pswdUppercase:
			err = errors.New("password must contain atleast one uppercase letter")
		case pswdNumber:
			err = errors.New("password must contain atleast one number")
		case pswdSpecial:
			err = errors.New("password must contain atleast one special character")
		case pswdLength:
			err = errors.New("password length must atleast 12 characters and less than 60")
		case pswdNoSpaces:
			err = errors.New("password cannot have any spaces")
		}
		return err
	}
	return nil
}

func checkEmailValid(email string) error {
	// check email syntax is valid
	//func MustCompile(str string) *Regexp
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		fmt.Println(err)
		return errors.New("sorry, something went wrong")
	}
	rg := emailRegex.MatchString(email)
	if !rg {
		return errors.New("email address is not a valid syntax, please check again")
	}
	// check email length
	if len(email) < 4 {
		return errors.New("email length is too short")
	}
	if len(email) > 253 {
		return errors.New("email length is too long")
	}
	return nil
}

func checkEmailDomain(email string) error {
	i := strings.Index(email, "@")
	host := email[i+1:]
	// func LookupMX(name string) ([]*MX, error)
	_, err := net.LookupMX(host)
	if err != nil {
		err = errors.New("could not find email's domain server, please chack and try again")
		return err
	}
	return nil
}
