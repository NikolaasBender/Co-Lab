package go_dev

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

/*
Adds a new user to the database
If succesful, returns true
Otherwise, returns false
*/
func AddUser(username, password, email, bio string, db *sql.DB) bool {

	pss_hash, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.MinCost)

	sqlStatement := `INSERT INTO user_login (username, password, email)
  VALUES ($1, $2, $3);`
	sqlStatement2 := `INSERT INTO user_info (username, bio)
  VALUES ($1, $2);`

	_, err = db.Exec(sqlStatement, username, string(pss_hash[:]), email)

	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = db.Exec(sqlStatement2, username, bio)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

/*
Deletes specified user from database
Returns true if succesful, false otherwise.

Be warned. Does not delete all trace of the user
*/
func DeleteUser(username, password string, db *sql.DB) bool {

	valid := Validate(username, password, db)

	if valid != true {
		return false
	}

	sqlStatement := `DELETE FROM user_info WHERE username = $1;`
	sqlStatement2 := `DELETE FROM user_login WHERE username = $1;`

	_, err := db.Exec(sqlStatement, username)
	_, er := db.Exec(sqlStatement2, username)

	if err != nil || er != nil {
		return false
	}

	return true
}

/*
Checks if a user is in the database
If the user is in the databse, return true
If user not found, return false
*/
func Exists(username string, db *sql.DB) bool {

	sqlStatement := `SELECT username FROM user_login
    WHERE username = $1;`

	var uname string

	err = db.QueryRow(sqlStatement, username).Scan(&uname)

	if err == sql.ErrNoRows {
		fmt.Println("No Rows.")
		return false
	} else if err != nil {
		fmt.Println("Other sql error.")
		return false
	}

	return true
}

/*
Validates a user login
If user entered correct password, return true
If user password combination not present in databse, return false
*/
func Validate(username, password string, db *sql.DB) bool {

	sqlStatement := `SELECT password FROM user_login
  WHERE username = $1;`

	var pss string

	err = db.QueryRow(sqlStatement, username).Scan(&pss)

	if err == sql.ErrNoRows {
		fmt.Println("No Rows.")
		return false
	} else if err != nil {
		fmt.Println("Other sql error.")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(pss),[]byte(password))

	if(err != nil) {
		fmt.Println(string(pss)+" "+password)
		return false
	} else {
		return true
	}
}

/*
Gets all necessary user info
If succesful, returns array of all user info
If user has no info, everything except for username will be null
*/
func GetUserInfo(username string, db *sql.DB) UserInfo {

	sqlStatement := `SELECT username, bio, profile_image, banner_image FROM user_info
  WHERE username = $1;`

	var info UserInfo

	err = db.QueryRow(sqlStatement, username).Scan(&info.Username, &info.Bio, &info.Profileimg, &info.Bannerimg)

	if err == sql.ErrNoRows {
		fmt.Println(err)
		return info
	} else if err != nil {
		fmt.Println(err)
		return info
	}

	return info
	//Need to figure out the best way to return some of this information.
}

/*
Updates user info in the database
If succesful, returns true
Otherwise, returns false
*/
func EditUserInfo(username, field, edit string, db *sql.DB) bool {

	sqlStatement := `UPDATE user_info
  SET $1 = $2
  WHERE username = $3;`

	var uname string

	err = db.QueryRow(sqlStatement, field, edit, username).Scan(&uname)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}
