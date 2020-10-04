package util

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// DBConn open up our database connection.
func DBConn() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("MYSQL_BOOKS"))
	PanicError(err)
	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

// PanicError handling with repetitive task of error-handling
func PanicError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

