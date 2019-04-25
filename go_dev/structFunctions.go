package go_dev

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
)

func PopulateUserPage(username string, db *sql.DB) UserPage {
	var page UserPage

	page.Info = GetUserInfo(username, db)

	page.Projects = GetProjects(username, db)

	page.Tasks = GetUserTasks(username, db)

	page.Pins = GetUserPins(username, db)

	page.Feed = GetUserFeed(username, db)

	return page
}

func PopulateProjectPage(id int, db *sql.DB) *Project {
	var thing *Project

	thing.Id = id

	thing.Todo = GetProjectTasks(id, 0, db)

	thing.Working = GetProjectTasks(id, 1, db)

	thing.Done = GetProjectTasks(id, 2, db)

	thing.Users = GetProjectMembers(id, db)

	return thing
}
