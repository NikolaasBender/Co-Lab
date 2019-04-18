package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"Co-Lab/go_dev"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// db = go_dev.Initialize()

// var templates = template.Must(template.ParseGlob("templates/*"))

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	//will eventually make random key generation
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

var debug = true

//THIS HANDLES ANYTHING IN THE TOP DIRECTORY
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit HomeHandler")
	}

	pathVariables := mux.Vars(r)
	fmt.Println("HOME HANDLER: '" + pathVariables["page"] + "'" + "'" + r.URL.Path + "'")

	t, _ := template.ParseFiles(pathVariables["page"])

	t.Execute(w, t)
}

//=====================================================================================
//SUPER BASIC INDEX HANDLER
//=====================================================================================
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit IndexHandler")
	}

	t, _ := template.ParseFiles("view/index.html")

	t.Execute(w, t)
}

//=====================================================================================
//VIEW HANDLER
//=====================================================================================
func ViewHandler(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit ViewHandler")
	}
	session, _ := store.Get(r, "cookie-name")

	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	page := file_finder("view/", w, r)

	if strings.Contains(page, "feed") == true {
		FeedHandler(w, r)
		return
	}

	t, _ := template.ParseFiles(page)
	t.Execute(w, nil)

}

//=====================================================================================
//PFFFT I HAVE NO IDEA WHAT THIS IS FOR
//=====================================================================================
func secret(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit secret")
	}

	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

//=====================================================================================
//THIS IS THE LOGIN HANDLER
//=====================================================================================
func login(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit login")
		fmt.Println(r.Method)
	}

	t, _ := template.ParseFiles("auth/login.html")
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	session, _ := store.Get(r, "cookie-name")

	pwint, _ := strconv.Atoi(r.FormValue("pwd"))
	if debug == true {
		fmt.Println(pwint)
	}
	// Authentication goes here
	if go_dev.Validate(r.FormValue("usr"), pwint, db) == true {
		if debug == true {
			fmt.Println("user has been validated")
		}
		session.Values["authenticated"] = true
		http.Redirect(w, r, "/view/index.html", http.StatusFound)
	} else {
		if debug == true {
			fmt.Println("user has NOT been validated")
		}
		session.Values["authenticated"] = false
		t.Execute(w, nil)
	}

}

//=====================================================================================
//THE LOG OUT HANDLER
//=====================================================================================
func logout(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit logout")
	}

	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	IndexHandler(w, r)
}

func signup(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit signup")
	}

}

//=====================================================================================
//THIS DISPLAYS THE CUSTOM 404 PAGE
//=====================================================================================
func notFound(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("404.html")

	t.Execute(w, nil)
}

//=====================================================================================
//ANAY SORT OF FEED WILL BE HANDLED WITH THIS
//=====================================================================================
func FeedHandler(w http.ResponseWriter, r *http.Request) {
	//POPULATE THE FEED WITH THE RIGHT POSTS
	if debug == true {
		fmt.Println("Hit FeedHandler")
	}
	session, _ := store.Get(r, "cookie-name")

	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	//page := file_finder("view/", w, r)

	//GET POSTS FOR USER
	// feedposts = go_dev.getTasks(session.Values["name"])

	// p := Feed{Title: session.Values["name"], Posts: feedposts}
	// t, _ := template.ParseFiles(page)

}

// //=====================================================================================
// //ANY SORT OF POST WILL BE HANDLED HERE
// //=====================================================================================
// func PostHnadler(w http.ResponseWriter, r *http.Request) {
// 	if debug == true {
// 		fmt.Println("Hit PostHandler")
// 	}
// 	session, _ := store.Get(r, "cookie-name")

// 	if session.Values["authenticated"] != true {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}

// 	pathVariables := mux.Vars(r)

// 	k := pathVariables["key"]

// 	//GET THE POST FROM THE DB

// 	// p := get from db in the post struct postGet(k)

// 	t, _ := template.ParseFiles("/view/post.html")
// 	t.Execute(w, p)

// }

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
