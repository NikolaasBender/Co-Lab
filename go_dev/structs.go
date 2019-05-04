package go_dev

import (
	//"database/sql"

	_ "github.com/lib/pq"
)

var err error

type UserPage struct {
	Info     UserInfo
	Feed     []Post
	Tasks    []Task
	Pins     []Post
	Projects []Project
}

type Post struct {
	Title    string
	Content  string
	Username string
	Task     string
}

type Project struct {
	Project_name string
	Id           int
	Todo         []Task
	Working      []Task
	Done         []Task
	Users        []string
}

type Task struct {
	Project_name string
	Name         string
	Key          int
	Description  string
	Comments     []Post
	Due_date     string
	Status       int
}

type UserInfo struct {
	Username   string
	Bio        string
	Profileimg string
	Bannerimg  string
}
