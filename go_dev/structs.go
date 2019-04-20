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
	pins          []Post
	projects      []Project
}

type Post struct {
	Title   string
	Content string
	username string
}

type Project struct {
	project_name string
	id 			 int
	todo         []Task
	working      []Task
	done         []Task
	users        []string
}

type Task struct {
	name string
	description string
	comments  []Post
	due_date string
	status int
}

type UserInfo struct {
	username string
	bio string
	profileimg string
	bannerimg string
}
