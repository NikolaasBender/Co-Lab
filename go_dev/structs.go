package go_dev

import (
	//"database/sql"

	_ "github.com/lib/pq"
)

var err error

type UserPage struct {
	info     UserInfo
	feed     []Post
	tasks    []Task
	pins     []Post
	projects []Project
}

type Post struct {
	title    string
	content  string
	username string
	task string
}

type Project struct {
	project_name string
	id           int
	todo         []Task
	working      []Task
	done         []Task
	users        []string
}

type Task struct {
	project_name string
	name string
	description string
	comments    []Post
	due_date    string
	status      int
}

type UserInfo struct {
	username   string
	bio        string
	profileimg string
	bannerimg  string
}
