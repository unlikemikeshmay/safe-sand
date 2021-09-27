package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// create connection with postgres db
func CreateConnection() *sql.DB {
	// load .env file
	/* 	err := godotenv.Load(".env")

	   	if err != nil {
	   		//log.Fatalf("Error loading .env file")
	   	}
	*/
	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}
