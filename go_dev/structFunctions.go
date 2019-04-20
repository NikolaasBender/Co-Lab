package go_dev

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
)


func PopulateUserPage(username string, db *sql.DB) UserPage {
	var thing UserPage

	return thing
}

func PopulateProjectPage(id int,  db *sql.DB) Project {
	var thing Project 

	return thing
}
