package go_dev

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var err error

func AddUser(username, password, email, bio string, db *sql.DB) bool {

	sqlStatement := `INSERT INTO user_login (username, password, email)
  VALUES ($1, $2, $3);`
	sqlStatement2 := `INSERT INTO user_info (username, bio)
  VALUES ($1, $2);`

	_, err = db.Exec(sqlStatement, username, password, email)

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

func Validate(info User, db *sql.DB) bool {

	sqlStatement := `SELECT username FROM user_login
  WHERE username = $1 AND password = $2;`

	var uname string

	err = db.QueryRow(sqlStatement, info.username, info.password).Scan(&uname)

	if err == sql.ErrNoRows {
		fmt.Println("No Rows.")
		return false
	} else if err != nil {
		fmt.Println("Other sql error.")
		return false
	}

	return true
}

func GetUserInfo(username string, db *sql.DB) (UserInfo) {

	sqlStatement := `SELECT * FROM user_info
  WHERE username = $1;`

	var info UserInfo

	err = db.QueryRow(sqlStatement, username).Scan(&info.username, &info.bio, &info.profileimg, &info.bannerimg)

	if err == sql.ErrNoRows {
		fmt.Println(err)
		return nil
	} else if err != nil {
		fmt.Println(err)
		return nil
	}

	return info
	//Need to figure out the best way to return some of this information.
}

func EditUserInfo(username, field, edit string, db *sql.DB) bool {

	sqlStatement := `UPDATE user_info
  SET $1 = $2
  WHERE username = $3;`

	var uname string
	var err error

	err = db.QueryRow(sqlStatement, field, edit, username).Scan(&uname)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}
