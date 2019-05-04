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
	id, err := strconv.Atoi(string(pathVariables["key"]))
	if debug == true {
		fmt.Println("PROJECT ID IS: ", id)
	}
	if err != nil {
		fmt.Println("THERE WAS AN ISSUE TURNING THE KEY INTO AN INT ", err)
	}
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

	//HAVE TO DO THIS SEQUENCE FOR GET AND POST
	pathVariables := mux.Vars(r)

	id, _ := strconv.Atoi(string(pathVariables["key"]))

	p := go_dev.GetTask(id, db)

	t, err := template.ParseFiles("/view/task_view.html")

	if err != nil {
		fmt.Println("task view Handler parsing error", err)
	}

	//IF ITS A GET REQUEST IT JUST SHOWS THE TASK AND ITS COMMENTS
	if r.Method != http.MethodPost {
		t.Execute(w, p)
		return
	}

	t.Execute(w, p)
	return
}

//DEALS WITH ADDING COMMENTS TO A TASK
func TaskCommentHandler(w http.ResponseWriter, r *http.Request) {
	if debug == true {
		fmt.Println("Hit TaskViewHandler")
	}
	//session, _ := store.Get(r, "cookie-name")

	if heimdall(w, r) != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	//GET OUR APP COOKIE FOR USE LATER
	session, _ := store.Get(r, appCookie)
	//SPECIFYING THE USERNAME UP HERE BECAUSE IT'S USED SO MUCH
	user := session.Values["usr"].(string)

	//HAVE TO DO THIS SEQUENCE FOR GET AND POST
	pathVariables := mux.Vars(r)

	id, _ := strconv.Atoi(string(pathVariables["key"]))

	p := go_dev.GetTask(id, db)

	t, err := template.ParseFiles("/view/ac2t.html")

	if err != nil {
		fmt.Println("task view Handler parsing error", err)
	}

	//IF ITS A GET REQUEST IT JUST SHOWS THE TASK AND ITS COMMENTS
	if r.Method != http.MethodPost {
		t.Execute(w, p)
		return
	}

	title := string(r.FormValue("title"))
	cont := string(r.FormValue("content"))
	if debug == true {
		fmt.Println("got the form values", title, cont)
	}
	ok := go_dev.CreatePost(id, user, title, cont, db)
	if ok != true {
		fmt.Println("error adding comment to task")
	}

	t.Execute(w, p)
	return
}

func TaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	if debug == true {
		fmt.Println("Hit TaskStatusHandler")
	}
	//session, _ := store.Get(r, "cookie-name")

	if heimdall(w, r) != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	//HAVE TO DO THIS SEQUENCE FOR GET AND POST
	pathVariables := mux.Vars(r)

	id, _ := strconv.Atoi(string(pathVariables["key"]))

	p := go_dev.GetTask(id, db)

	t, err := template.ParseFiles("/view/task_stat.html")

	if err != nil {
		fmt.Println("task view Handler parsing error", err)
	}

	//IF ITS A GET REQUEST IT JUST SHOWS THE TASK AND ITS COMMENTS
	if r.Method != http.MethodPost {
		t.Execute(w, p)
		return
	}

	//change, err := strconv.ParseInt(r.FormValue("change")[0:], 10, 64)
	// if err != nil {
	// 	// handle the error in some way
	// }
	del := string(r.FormValue("del"))

	if del == "true" {
		ok := go_dev.DeleteTask(id, db)
		if ok != true {
			fmt.Println("ERROR DELETING TASK")
		}
	}
	// if change != p["status"] {
	// 	ok := go_dev.UpdateStatus(id, change, db)
	// 	if ok != true {
	// 		fmt.Println("ERROR UPDATING TASK STATUS")
	// 	}
	// }

	t.Execute(w, p)
	return
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
