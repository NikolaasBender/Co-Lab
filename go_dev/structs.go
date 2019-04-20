package go_dev

import (
	//"database/sql"

	_ "github.com/lib/pq"
)

type UserInfo struct {
	username string
	bio string
	profileimg string
	bannerimg string
}

type UserPage struct {
	display_name  string
	display_email string
	display_info  string
	feed          []Post
	tasks         []Task
	pins          []Task
	projects      []Project
}

type Post struct {
	Title   string
	Content string
}

type Project struct {
	project_name string
	todo         []Task
	working      []Task
	done         []Task
	users        []User
}

type User struct {
	username string
	password string
	email    string
}

//WE WILL WANT TO GET MORE STUFF ABOUT EACH POST
type Task struct {
	task_name string
	fill_task string
	key       string
	Comments  []Comment
}

type Comment struct {
	shmoo string
}

type Feed struct {
	Tasks []Task
}

// type Signup struct {
// 	email             string
// 	username          string
// 	password          string
// 	password_validate string
// }
