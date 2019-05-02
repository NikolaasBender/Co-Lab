package main

import (
	// "SoftwareDevProject/go_sql/go_dev"
	//"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"Co-Lab/go_dev"

	"github.com/gorilla/mux"
)

var debug = true

var err error

//=====================================================================================
//FAVICON HANDLER
//=====================================================================================
// func faviconHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "favicon.ico")
// }

//=====================================================================================
//This handler is called when the user enters the home page
//It simply just displays the homepage using a template
//=====================================================================================
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit IndexHandler")
	}

	//PARSE THE INDEX FILE
	t, err := template.ParseFiles("view/index.html")

	if err != nil {
		fmt.Println("IndexHandler parsing error", err)
	}

	//SERVE INDEX
	t.Execute(w, t)
}

//=====================================================================================
//This is called when we get a 404 error
//When a user enters a restricted or unfound url on our site, the 404 page is served
//=====================================================================================
func notFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("404.html")

	if err != nil {
		fmt.Println("404 Handler parsing error", err)
	}

	t.Execute(w, nil)
}

//=====================================================================================
//This is called when a project page is entered
//First it checks if the user is logged in
//If the user is logged in and is part of the project then the project data is pulled
//The project data is passed into a templated which then populates the page
//=====================================================================================
func ProjectViewHandler(w http.ResponseWriter, r *http.Request) {
	if debug == true {
		fmt.Println("Hit ProjectViewHandler")
	}
	//session, _ := store.Get(r, "cookie-name")

	if heimdall(w, r) != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	t, err := template.ParseFiles("view/project_view.html")
	if err != nil {
		fmt.Println("THERE WAS AN ERROR PARSING project_view.html ", err)
	} else {
		fmt.Println("THERE WAS *NO* ERROR IN  PARSING project_view.html ")
	}

	//HANDLING VIEWING
	pathVariables := mux.Vars(r)
	id, _ := strconv.Atoi(string(pathVariables["key"]))
	p := go_dev.PopulateProjectPage(id, db)
	if debug == true {
		fmt.Println("POPULATING THE PROJECT PAGE RETURNED: ", p)
	}
	t.Execute(w, p)

}

//=====================================================================================
//This is called when a task page needs to be displayed
//First it verifies that the user is logged in
//Then it pulls the necessary task info from the database
//Then it passes the task info into a template and populates the page
//=====================================================================================
func TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	if debug == true {
		fmt.Println("Hit TaskViewHandler")
	}
	//session, _ := store.Get(r, "cookie-name")

	if heimdall(w, r) != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	pathVariables := mux.Vars(r)

	id, _ := strconv.Atoi(string(pathVariables["key"]))

	p := go_dev.PopulateProjectPage(id, db)

	t, err := template.ParseFiles("/view/task_view.html")

	if err != nil {
		fmt.Println("task view Handler parsing error", err)
	}

	t.Execute(w, p)

}

//=====================================================================================
//This is called for finding a file
//It checks if path/to/whatever exists
//It does not check if whatever actually exists
//=====================================================================================
func file_finder(folder string, w http.ResponseWriter, r *http.Request) string {

	pathVariables := mux.Vars(r)
	if debug == true {
		fmt.Println("File Finder: '" + pathVariables["page"] + "'")
	}

	page := ""

	if strings.Contains(pathVariables["page"], ".html") == true {
		page = folder + pathVariables["page"]
	} else {
		page = folder + pathVariables["page"] + ".html"
	}

	if debug == true {
		fmt.Println("corrected path: '" + page + "'")
	}

	if _, err := os.Stat(page); err == nil {
		// path/to/whatever exists
		return page

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println("Can't find file")
		notFound(w, r)
		return ""
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		fmt.Println(err)
		return ""
	}
	return ""
}
