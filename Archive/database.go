package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var sqlDB *sql.DB
var gormDB *gorm.DB
var err error

const DSN = "user:password@tcp(127.0.0.1:8081)/booklah_db?charset=utf8mb4&parseTime=True&loc=Local"

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

type Message struct {
	MessageID  uint `gorm:"primaryKey"`
	Username   string
	Email      string
	ErrMessage string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func InitialMigration() {
	sqlDB, err = sql.Open("mysql", DSN)

	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to sqlDB")
	}

	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to gormDB")
	}
	//create tables at databse based on struct
	gormDB.AutoMigrate(&User{}, &Message{})

	//create admin account
	var finder User
	gormDB.Unscoped().First(&finder, "username = ?", "admin")
	if reflect.ValueOf(finder).IsZero() {
		bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
		user := User{Username: "admin", Password: bPassword, First: "admin", Last: "admin"}

		result := gormDB.Create(&user) // pass pointer of data to Create

		fmt.Println("ID added:", user.UserID)             // returns inserted data's primary key
		fmt.Println("Error:", result.Error)               // returns error
		fmt.Println("RowsAffected:", result.RowsAffected) // returns inserted records count

	}
}
