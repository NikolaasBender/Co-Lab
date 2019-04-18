package go_dev

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ContactDetails struct {
	Username string
	Password string
}

type UserInfo struct {
	username   string
	name       string
	bio        sql.NullString
	profileimg sql.NullString
	bannerimg  sql.NullString
}

type BigForm struct {
	nickname  string
	email     string
	password  string
	gender    string
	securityQ string
	languages string
	textbox   string
}

//WE WILL WANT TO GET MORE STUFF ABOUT EACH POST
type Post struct {
	Title   string
	Content string
	Link    string
	Type    int
}

type Feed struct {
	Posts []Post
}
