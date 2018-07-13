package connection

import (
	"database/sql"
	"fmt"
)

// Connect to the Database
func Connect() *sql.DB {
	//Connect to database
	db, err := sql.Open("mysql", "root:root@/movies")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully connect to the database!")
	}

	db.Ping()

	return db
}
