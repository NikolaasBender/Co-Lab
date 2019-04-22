package main

import (
	// "SoftwareDevProject/go_sql/go_dev"
	//"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"Co-Lab/go_dev"

	"github.com/gorilla/mux"
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

// store will hold all session data
var store = sessions.NewCookieStore([]byte("shhh_its_secret"))

const appCookie = "DeleciousCoLabCookies"

var debug = true

var err error

//SOME COOKIE INITIALIZATION
// func init() {
// 	authKeyOne := securecookie.GenerateRandomKey(64)
// 	encryptionKeyOne := securecookie.GenerateRandomKey(32)

// 	store = sessions.NewCookieStore(
// 		authKeyOne,
// 		encryptionKeyOne,
// 	)

// 	store.Options = &sessions.Options{
// 		Domain: "localhost",
// 		Path:   "/",
// 		//TIMOUT ON THE COOKIE OF 90min
// 		MaxAge: 60 * 90,
// 		//true so the session cannot be altered by javascript.
// 		HttpOnly: true,
// 	}

// }

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
		if debug == true {
			fmt.Println("And we're sending you back to login", heimdall(w, r))
		}
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

	session, _ := store.Get(r, appCookie)

	pw := r.FormValue("pwd")
	// hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.MinCost)
	// if err != nil {
	// 	log.Println(err)
	// }

	// Authentication goes here
	if go_dev.Validate(r.FormValue("usr"), pw, db) == true {
		if debug == true {
			fmt.Println("user has been validated")
		}

		session.Values["auth"] = true
		session.Values["usr"] = r.FormValue("usr")

		if debug == true {
			fmt.Println("getUser befor save", session.Values["auth"])
		}

		err := session.Save(r, w)

		if debug == true {
			fmt.Println("getUser after save", session.Values["auth"])
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view/userpage.html", http.StatusFound)
	} else {
		if debug == true {
			fmt.Println("user has NOT been validated")
		}

		session.Values["auth"] = false
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
	session, err := store.Get(r, appCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["auth"] = false
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

	session, _ := store.Get(r, appCookie)

	pass := ""
	if r.FormValue("pwd") == r.FormValue("pwdv") {
		pass = r.FormValue("pwd")
	}

	if debug == true {
		fmt.Println(r.FormValue("pwd"), r.FormValue("pwdv"), r.FormValue("name"), r.FormValue("email"))
	}

	errcheck := go_dev.AddUser(r.FormValue("name"), pass, r.FormValue("email"), "", db)
	if errcheck != true {
		fmt.Println("Hey, the user wasn't added")
	}

	session.Values["auth"] = true
	err := session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)

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

//=====================================================================================
//THIS DEALS WITH CHECKING FOR AUTHORIZATION
//=====================================================================================
func heimdall(w http.ResponseWriter, r *http.Request) bool {

	if debug == true {
		fmt.Println("Opening the Bifröst")
	}

	session, _ := store.Get(r, appCookie)

	if debug == true {
		fmt.Println("Bifröst: ", session, session.Values["auth"])
	}

	if session.Values["auth"] != true {
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		return false
	}
	return true
}

func foo(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, appCookie)
	session.Values["bar"] = "bar"
	session.Save(r, w)
}
