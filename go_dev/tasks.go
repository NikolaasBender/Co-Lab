package go_dev

import (
	"database/sql"

	"github.com/lib/pq"
	"fmt"
)

/*
Initializes a new task for a project and adds it to the database
If succesful, returns true
Otherwise, returns false
*/
func CreateTask(project_name, project_owner, task_name string, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string

	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	sqlStatement2 := `INSERT INTO tasks(project,name,status)
  	VALUES ($1, $2, 0)`

	_, err = db.Exec(sqlStatement2, parentID, task_name)

	if err != nil {
		return false
	}

	return true
}

/*
Adds a new member to an existing task
If succesful, returns true
Otherwise, returns false
*/
func AddTaskMembers(project_name, project_owner, task_name, newMember string, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID int
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		fmt.Println("No rows")
		return false
	} else if err != nil {
		fmt.Println("Other error first statement.")
		return false
	}

	sqlStatement := `UPDATE tasks
  	SET users = users || $1
  	WHERE project = $2 AND name = $3;`

	_, err = db.Exec(sqlStatement, pq.Array([]string{newMember}), parentID, task_name)

	if err != nil {
		fmt.Println("Other error second statement.")
		fmt.Println(err)
		fmt.Printf("%T\n",newMember)
		return false
	}

	return true
}

/*
Updates a tasks status and then populates the corresponding project task lisk
If succesful, returns true
Otherwise, returns false
*/
func UpdateStatus(project_name, project_owner, task_name, status int, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}
	sqlStatement2 := `SELECT id,status FROM tasks WHERE project = $1 AND name = $2;`

	var oldStatus int
	var taskID int
	err = db.QueryRow(sqlStatement2, parentID, task_name).Scan(&taskID, &oldStatus)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}
	var oldColumn string
	if oldStatus == 0 {
		oldColumn = "inprogress_tasks"
	} else if oldStatus == 1 {
		oldColumn = "todo_tasks"
	} else {
		oldColumn = "completed_tasks"
	}

	var newColumn string
	if status == 0 {
		newColumn = "inprogress_tasks"
	} else if status == 1 {
		newColumn = "todo_tasks"
	} else {
		newColumn = "completed_tasks"
	}

	sqlStatement3 := `UPDATE tasks SET status = $1 WHERE project = $2 AND name = $3;`
	_, err = db.Exec(sqlStatement3, status, parentID, task_name)
	if err != nil {
		return false
	}

	sqlStatement4 := `UPDATE projects SET $1 = array_remove($1, $2) WHERE id = $3;`
	_, err = db.Exec(sqlStatement4, oldColumn, taskID, parentID)
	if err != nil {
		return false
	}

	sqlStatement5 := `UPDATE projects SET $1 = array_cat($1, $2) WHERE id = $3;`
	_, err = db.Exec(sqlStatement5, newColumn, taskID, parentID)
	if err != nil {
		return false
	}

	return true
}

/*
Adds a description to a task
If succesful, returns true
Otherwise, returns false
*/
func AddDescription(project_name, project_owner, task_name, description string, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	sqlStatement := `UPDATE tasks
  	SET description = $1
  	WHERE project = $2 AND name = $3;`

	_, err = db.Exec(sqlStatement, description, parentID, task_name)

	if err != nil {
		return false
	}

	return true
}

/*
Changes the due date on a task
If succesful, returns true
Otherwise, returns false
*/
func DueDate(project_name, project_owner, task_name, dueDate string, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	sqlStatement := `UPDATE tasks
  	SET due_date = $1
  	WHERE project = $2 AND name = $3;`

	_, err = db.Exec(sqlStatement, dueDate, parentID, task_name)

	if err != nil {
		return false
	}

	return true
}

/*
Removes a task from the project database
If succesful, returns true
Otherwise, returns false
*/
func DeleteTask(project_name, project_owner, task_name, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	sqlStatement := `DELETE FROM tasks
  	WHERE  project= $1 AND name = $2;`

	_, err = db.Exec(sqlStatement, parentID, task_name)

	if err != nil {
		return false
	}

	return true
}

/*
Returns all tasks specific to a user
If succesful, returns array of all tasks
If user has no tasks, then the userTasks array returned will be empty
*/
func GetUserTasks(username string, db *sql.DB) []Task {

	sqlStatement := `SELECT t.id, p.name, t.name, t.description, EXTRACT(MONTH FROM t.due_date) as month, EXTRACT(DAY FROM t.due_date) as day, t.status
  FROM tasks t INNER JOIN  projects p ON t.project = p.id
  WHERE $1 = ANY(t.users) ORDER BY t.due_date ASC;`

	sqlStatement2 := `SELECT p.title, p.users, p.content, t.name
	FROM posts p INNER JOIN tasks t ON p.task = t.id WHERE p.task = $1;`

	rows, err := db.Query(sqlStatement, username)

	if err != nil {
		//Do something
	}

	var userTasks = make([]Task, 0)
	var day, month string
	var tsk Task

	var comments = make([]Post, 0)
	var pst Post

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&tsk.Key, &tsk.Project_name, &tsk.Name, &tsk.Description, &month, &day, &tsk.Status)

		if err != nil {
			//Do something
		}

		tsk.Due_date = month + "-" + day

		rows2, er := db.Query(sqlStatement2, tsk.Key)

		if er != nil {
			//Do something
		}

		defer rows2.Close()

		for rows2.Next() {

			err = rows2.Scan(&pst.Title, &pst.Username, &pst.Content, &pst.Task)

			if err != nil {
				//Do something
			}

			comments = append(comments, pst)
		}

		tsk.Comments = comments

		userTasks = append(userTasks, tsk)
	}

	return userTasks
}

/*
Gets all tasks specific to a project
If succesful, returns array of all tasks
If user has no tasks, then the ProjectTasks array returned will be empty
*/
func GetProjectTasks(id int, status int, db *sql.DB) []Task {

	sqlStatement1 := `SELECT name FROM projects WHERE id = $1;`

	var projectName string

	err = db.QueryRow(sqlStatement1, id).Scan(&projectName)

	if err != nil {
		//Do something
	}

	sqlStatement := `SELECT name, description, due_date
  FROM tasks WHERE project = $1 AND status = $2 ORDER BY due_date ASC;`

	rows, err := db.Query(sqlStatement, id, status)

	if err != nil {
		//Do something
	}

	var ProjectTasks = make([]Task, 0)

	defer rows.Close()

	for rows.Next() {
		var tsk Task

		err = rows.Scan(&tsk.Name, &tsk.Description, &tsk.Due_date)

		tsk.Status = status
		tsk.Project_name = projectName

		if err != nil {
			//Do something
		}

		ProjectTasks = append(ProjectTasks, tsk)
	}

	return ProjectTasks
}
