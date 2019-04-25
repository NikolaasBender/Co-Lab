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
//VIEW HANDLER
//=====================================================================================
func ViewHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit ViewHandler")
	}

	//CHECK IF THE USER IS LOGGED IN
	if heimdall(w, r) != true {
		if debug == true {
			fmt.Println("And we're sending you back to login", heimdall(w, r))
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return

	}

	//FIND THE RIGHT PAGE, JUST TO DEAL WITH ".html" MAYBE NOT BEING THERE
	page := file_finder("view/", w, r)
	if debug == true {
		fmt.Println("FOUND THE RIGHT PAGE", page)
	}

	//PARSE THE FOUND FILE
	t, _ := template.ParseFiles(page)
	if debug == true {
		fmt.Println("PARSED THE PAGE CORRECTLY")
	}

	//GET OUR APP COOKIE FOR USE LATER
	session, _ := store.Get(r, appCookie)

	if strings.Contains(page, "userpage") == true {
		//ASK SQL TEAM FOR ALL THE USER PAGE STUFF
		p := go_dev.PopulateUserPage(session.Values["usr"].(string), db)
		if debug == true {
			fmt.Println("POPULATED THE USER PAGE CORRECTLY")
		}

		t.Execute(w, p)

		if debug == true {
			fmt.Println("EXECUTED THE USERPAGE")
		}

		return
	}

	t.Execute(w, nil)

}

//=====================================================================================
//THIS DISPLAYS THE CUSTOM 404 PAGE
//=====================================================================================
func notFound(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/view/404errorpage.html")

	t.Execute(w, nil)
}

//=====================================================================================
//THIS DEALS WITH VIEWING POSTS
//=====================================================================================
// func PostViewHandler(w http.ResponseWriter, r *http.Request) {
// 	if debug == true {
// 		fmt.Println("Hit PostViewHandler")
// 	}
// 	//session, _ := store.Get(r, "cookie-name")

// 	if heimdall(w, r) != true {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}

// 	pathVariables := mux.Vars(r)

// 	id, _ := strconv.Atoi(string(pathVariables["key"]))

// 	p := go_dev.PopulateProjectPage(id, db)

// 	t, _ := template.ParseFiles("/view/task_view.html")

// 	t.Execute(w, p)

// }

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

	t, _ := template.ParseFiles("/view/project.html")

	t.Execute(w, p)

}

//=====================================================================================
//THIS DEALS WITH CREATING PROJECTS
//=====================================================================================
func ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	if debug == true {
		fmt.Println("Hit ProjectCreateHandler")
	}
	session, _ := store.Get(r, "cookie-name")

	if heimdall(w, r) != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	
	t, _ := template.ParseFiles("/view/project_create.html")
	if debug == true {
		fmt.Println("PARSED PROJECT CREATE CORRECTLY")
	}


	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	worked := go_dev.CreateProject(session.Values["usr"].(string), string(r.FormValue("pjn")), db)
	if worked != true {
		fmt.Println("Error creating a project")
	}

	t.Execute(w, nil)

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
