package go_dev

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
)

/*
Gets everything necesary to populate the user page
Returns a page object that will be used with templating
*/
func PopulateUserPage(username string, db *sql.DB) UserPage {
	var page UserPage

	page.Info = GetUserInfo(username, db)

	page.Projects = GetProjects(username, db)

	page.Tasks = GetUserTasks(username, db)

	page.Pins = GetUserPins(username, db)

	page.Feed = GetUserFeed(username, db)

	return page
}

/*
Gets everything necesary to populate a project page
Returns a "thing" object to be used with templating
*/
func PopulateProjectPage(id int, db *sql.DB) *Project {
	var thing *Project

	thing.Id = id

	thing.Project_name = GetProjectName(id, db)

	thing.Todo = GetProjectTasks(id, 0, db)

	thing.Working = GetProjectTasks(id, 1, db)

	thing.Done = GetProjectTasks(id, 2, db)

	thing.Users = GetProjectMembers(id, db)

	return thing
}
