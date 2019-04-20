package go_dev

import (
	//"database/sql"

	_ "github.com/lib/pq"
)

type UserPage struct {
	username  string
	email string
	bio  string
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
	users        []string
}

type Task struct {
	name string
	description string
	comments  []Post
}
