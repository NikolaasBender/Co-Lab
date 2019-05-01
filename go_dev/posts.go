package go_dev

import (
	"database/sql"

	_ "github.com/lib/pq"
	// "fmt"
)

/*Pushes new post to the database
If succesful, return true
Otherwise, return false
*/
func addContentPost(title, user, content, db *sql.DB) bool {
	sqlStatement := `UPDATE posts
  	SET content = $1
  	WHERE title = $2 AND user = $3;`
	_, err = db.Exec(sqlStatement, content, title, user)

	if err != nil {
		return false
	}

	return true
}

/*Removes post from the database
If succesful, return true
Otherwise, return false
*/
func deletePost(title, user, db *sql.DB) bool {
	sqlStatement := `DELETE FROM posts
  	WHERE title = $1 AND users = $2;`

	_, err = db.Exec(sqlStatement, title, user)

	if err != nil {
		return false
	}

	return true
}

/*
Gets all of the posts that a user has pinned from the databse
If succesful, return all of the user's pinned posts
If there are no pinned posts, then the userPins array returned will be empty
*/
func GetUserPins(username string, db *sql.DB) []Post {

	sqlStatement := `SELECT title, content, users FROM posts
	WHERE id = ANY(SELECT unnest(pins) FROM user_info WHERE username = $1);`

	rows, err := db.Query(sqlStatement, username)

	if err != nil {
		//Do something
	}

	var userPins = make([]Post, 5)

	defer rows.Close()

	for rows.Next() {
		var pst Post

		err = rows.Scan(&pst.Title, &pst.Content, &pst.Username)

		if err != nil {
			//Do something
		}

		userPins = append(userPins, pst)
	}

	return userPins
}

/*
Pulls all posts and tasks specific to the user
If succesful, will return everything needed to populate the user feed
If user has nothing, then the returned userFeed will be empty
*/
func GetUserFeed(username string, db *sql.DB) []Post {

	sqlStatement := `SELECT p.title, p.content, p.users, t.name
	FROM posts p INNER JOIN tasks t ON p.task = t.id
	WHERE $1 = ANY(t.users);`

	rows, err := db.Query(sqlStatement, username)

	if err != nil {
		//Do something
	}

	var userFeed = make([]Post, 5)

	defer rows.Close()

	for rows.Next() {
		var pst Post

		err = rows.Scan(&pst.Title, &pst.Content, &pst.Username, &pst.Task)
		if err != nil {
			//Do something
		}

		userFeed = append(userFeed, pst)
	}

	return userFeed
}
