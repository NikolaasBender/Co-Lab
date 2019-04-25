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
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// store will hold all session data
var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

const appCookie = "DeleciousCoLabCookies"

var debug = true

var err error

//=====================================================================================
//FAVICON HANDLER
//=====================================================================================
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

//=====================================================================================
//SUPER BASIC INDEX HANDLER
//=====================================================================================
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit IndexHandler")
	}

	//PARSE THE INDEX FILE
	t, _ := template.ParseFiles("view/index.html")

	//SERVE INDEX
	t.Execute(w, t)
}

//=====================================================================================
//THIS DISPLAYS THE CUSTOM 404 PAGE
//=====================================================================================
func notFound(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/view/404errorpage.html")

	t.Execute(w, nil)
}

//=====================================================================================
//THIS DEALS WITH VIEWING PROJECTS
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

	pathVariables := mux.Vars(r)

	id, _ := strconv.Atoi(string(pathVariables["key"]))

	p := go_dev.PopulateProjectPage(id, db)

	t, _ := template.ParseFiles("/view/project_view.html")

	t.Execute(w, p)

}

//=====================================================================================
//THIS DEALS WITH VIEWING TASKS
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

	t, _ := template.ParseFiles("/view/task_view.html")

	t.Execute(w, p)

}

//=====================================================================================
//THIS DEALS WITH FINDING THE RIGHT FILES
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
		return ""
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		fmt.Println(err)
		return ""
	}
}
