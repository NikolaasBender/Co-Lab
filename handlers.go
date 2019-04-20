package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"Co-Lab/go_dev"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	// "golang.org/x/crypto/bcrypt"
)

// db = go_dev.Initialize()

// var templates = template.Must(template.ParseGlob("templates/*"))

// var (
// 	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
// 	//will eventually make random key generation
// 	key   = []byte("super-secret-key")
// 	store = sessions.NewCookieStore(key)
// )

// User holds a users account information
type User struct {
	Username      string
	Authenticated bool
}

// store will hold all session data
var store *sessions.CookieStore

var debug = true

var err error

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		//TIMOUT ON THE COOKIE OF 90min
		MaxAge: 60 * 90,
		//true so the session cannot be altered by javascript.
		HttpOnly: true,
	}

	gob.Register(User{})
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

	if heimdall(w, r) != true {
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
//THIS IS THE LOGIN HANDLER (he thicc)
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

	pw := r.FormValue("pwd")
	// hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.MinCost)
	// if err != nil {
	// 	log.Println(err)
	// }

	user := &User{
		Username:      r.FormValue("usr"),
		Authenticated: false,
	}

	// Authentication goes here
	if go_dev.Validate(r.FormValue("usr"), pw, db) == true {
		if debug == true {
			fmt.Println("user has been validated")
		}

		user.Authenticated = true

		session.Values["user"] = user
		err := session.Save(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view/userpage.html", http.StatusFound)
	} else {
		if debug == true {
			fmt.Println("user has NOT been validated")
		}

		session.Values["user"] = user
		err := session.Save(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//=====================================================================================
//THE SIGNUP HANDLER
//=====================================================================================
func signup(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit signup")
	}

	t, _ := template.ParseFiles("auth/signup.html")
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	//session, _ := store.Get(r, "cookie-name")

}

//=====================================================================================
//THIS DISPLAYS THE CUSTOM 404 PAGE
//=====================================================================================
func notFound(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("404.html")

	t.Execute(w, nil)
}

//=====================================================================================
//ANY SORT OF FEED WILL BE HANDLED WITH THIS
//=====================================================================================
func FeedHandler(w http.ResponseWriter, r *http.Request) {
	//POPULATE THE FEED WITH THE RIGHT POSTS
	if debug == true {
		fmt.Println("Hit FeedHandler")
	}

	if heimdall(w, r) != true {
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

// getUser returns a user from session s
// on error returns an empty user
func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

//=====================================================================================
//THIS DEALS WITH CHECKING FOR AUTHORIZATION
//=====================================================================================
func heimdall(w http.ResponseWriter, r *http.Request) bool {

	if debug == true {
		fmt.Println("Opening the Bifröst")
	}

	session, _ := store.Get(r, "cookie-name")

	user := getUser(session)

	if auth := user.Authenticated; !auth {
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		return false
	}
	return true
}
