package go_dev

import (
	"database/sql"

	_ "github.com/lib/pq"
	// "fmt"
)

func createTask(project_name, project_owner, task_name string, db *sql.DB) bool {
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

func addTaskMembers(project_name, project_owner, task_name, newMember, db *sql.DB) bool {
	sqlStatement1 := `SELECT id FROM projects WHERE owner = $1 AND name = $2;`

	var parentID string
	err = db.QueryRow(sqlStatement1, project_owner, project_name).Scan(&parentID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	sqlStatement := `UPDATE tasks
  	SET users = users || '{$1}'
  	WHERE project = $2 AND name = $3;`

	_, err = db.Exec(sqlStatement, newMember, parentID, task_name)

	if err != nil {
		return false
	}

	return true
}

func updateStatus(project_name, project_owner, task_name, status int, db *sql.DB) bool {
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

func addDescription(project_name, project_owner, task_name, description, db *sql.DB) bool {
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

func dueDate(project_name, project_owner, task_name, dueDate, db *sql.DB) bool {
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

func deleteTask(project_name, project_owner, task_name, db *sql.DB) bool {
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

func GetUserTasks(username string, db *sql.DB) []Task {

	sqlStatement := `SELECT p.name, t.name, t.description, t.due_date, t.status
  FROM tasks t INNER JOIN  projects p ON t.project = p.id
  WHERE $1 = ANY(t.users) ORDER BY due_date ASC;`

	rows, err := db.Query(sqlStatement, username)

	if err != nil {
		//Do something
	}

	var userTasks = make([]Task, 5)

	defer rows.Close()

	for rows.Next() {
		var tsk Task

		err = rows.Scan(&tsk.project_name, &tsk.name, &tsk.description, &tsk.due_date, &tsk.status)

		if err != nil {
			//Do something
		}

		userTasks = append(userTasks, tsk)
	}

	return userTasks
}

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

	var ProjectTasks = make([]Task, 5)

	defer rows.Close()

	for rows.Next() {
		var tsk Task

		err = rows.Scan(&tsk.name, &tsk.description, &tsk.due_date)

		tsk.status = status
		tsk.project_name = projectName

		if err != nil {
			//Do something
		}

		ProjectTasks = append(ProjectTasks, tsk)
	}

	return ProjectTasks
}
