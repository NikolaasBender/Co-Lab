package go_dev

import (
	"database/sql"
	//"fmt"
	_ "github.com/lib/pq"
)

/*
Initializes a new project and adds it to the database
If succesful, return true
Otherwise, return false
*/
func CreateProject(owner, name string, db *sql.DB) bool {

	sqlStatement := `INSERT INTO projects(owner,name)
  VALUES ($1, $2)`

	_, err = db.Exec(sqlStatement, owner, name)

	if err != nil {
		return false
	}

	return true
}

/*
Adds a new user to an existing project
If succesful, return true
Otherwise, return false
*/
func AddProjectMembers(owner, name, newuser string, db *sql.DB) bool {

	sqlStatement := `UPDATE projects
  SET users = users || '{$1}'
  WHERE owner = $2 AND name = $3;`

	_, err = db.Exec(sqlStatement, newuser, owner, name)

	if err != nil {
		return false
	}

	return true
}

/*
Removes an existing project from the datbase
If succesful, return true
Otherwise, return false
*/
func DeleteProject(owner, name string, db *sql.DB) bool {

	sqlStatement := `DELETE FROM projects
  WHERE owner = $1 AND name = $2;`

	_, err = db.Exec(sqlStatement, owner, name)

	if err != nil {
		return false
	}

	return true
}

/*
Gets all projects specific to a user
If succesful, returns all projects a user is a part of
If a user has no projects, the usrProjects array returned will be empty
*/
func GetProjects(owner string, db *sql.DB) []Project {

	sqlStatement := `SELECT id, name FROM projects
  WHERE owner = $1 OR $1 = ANY(users);`

	rows, err := db.Query(sqlStatement, owner)

	if err != nil {
		//Do something
	}

	var usrProjects = make([]Project, 5)

	defer rows.Close()

	for rows.Next() {
		var p Project

		err = rows.Scan(&p.Id, &p.Project_name)

		if err != nil {
			//Do something
		}

		usrProjects = append(usrProjects, p)
	}

	return usrProjects
}

/*
Gets all users specific to a project
If succesful, returns all project members
If no project members, members will be returned empty
*/
func GetProjectMembers(id int, db *sql.DB) []string {
	sqlStatement1 := `SELECT users FROM projects WHERE id = $1;`

	var members []string

	err = db.QueryRow(sqlStatement1, id).Scan(&members)

	if err != nil {
		//Do something
	}

	return members
}

func GetProjectName(id int, db *sql.DB) string {
	sqlStatement1 := `SELECT name FROM projects WHERE id = $1;`

	var name string

	err = db.QueryRow(sqlStatement1, id).Scan(&name)

	if err != nil {
		//Do something
	}

	return name
}
