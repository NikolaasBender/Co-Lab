package go_dev

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
)


func PopulateUserPage(username string, db *sql.DB) *UserPage {
	var page *UserPage
	var usr *UserInfo

	usr = GetUserInfo(username, db)

	page.info = *usr

	

	return page
}

func PopulateProjectPage(id int,  db *sql.DB) Project {
	var thing Project

	return thing
}
