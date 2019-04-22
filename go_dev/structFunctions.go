package go_dev

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
)

func PopulateUserPage(username string, db *sql.DB) *UserPage {
	var page *UserPage

	page.info = *GetUserInfo(username, db)

	page.projects = GetProjects(username, db)

	page.tasks = GetUserTasks(username, db)

	page.pins = GetUserPins(username, db)

	page.feed = GetUserFeed(username, db)

	return page
}

func PopulateProjectPage(id int, db *sql.DB) *Project {
	var thing *Project

	thing.id = id

	thing.todo = GetProjectTasks(id, 0, db)

	thing.working = GetProjectTasks(id, 1, db)

	thing.done = GetProjectTasks(id, 2, db)

	thing.users = GetProjectMembers(id, db)

	return thing
}
