package main

import (
	// "SoftwareDevProject/go_sql/go_dev"
	//"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"Co-Lab/go_dev"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

const appCookie = "DeleciousCoLabCookies"

var debug = true

var err error

//=====================================================================================
//FAVICON HANDLER
//=====================================================================================
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/favicon.ico")
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

	//PARSE THE FOUND FILE
	t, _ := template.ParseFiles(page)

	//GET OUR APP COOKIE FOR USE LATER
	session, _ := store.Get(r, appCookie)

	if strings.Contains(page, "userpage") == true {
		//ASK SQL TEAM FOR ALL THE USER PAGE STUFF
		p := go_dev.PopulateUserPage(session.Values["usr"].(string), db)
		t.Execute(w, p)
		return
	}

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

	//PARSE THE LOGIN PAGE
	t, _ := template.ParseFiles("auth/login.html")
	//A CHECK FOR A POST METHOD THAT MIGHT NOT BE NECESSARY ANYMORE
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	//MAKEING A NEW COOKIE FOR THE USER
	session, _ := store.Get(r, appCookie)

	//READ IN THE PASSWORD ENTERED
	pw := r.FormValue("pwd")
	//PASSWORD ENCRYPTION
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	//ALL OF THE AUTH
	if go_dev.Validate(r.FormValue("usr"), string(hash), db) == true {
		if debug == true {
			fmt.Println("user has been validated")
		}
		//SET THE USER AS LOGGED IN AND PUT THE USERNAME IN THE COOKIE
		session.Values["auth"] = true
		session.Values["usr"] = string(r.FormValue("usr"))

		if debug == true {
			fmt.Println("getUser befor save", session.Values["auth"])
		}

		//SAVE THE COOKIE TO THE USER'S BROWSER
		err := session.Save(r, w)

		if debug == true {
			fmt.Println("getUser after save", session.Values["auth"])
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//SEND THEM ALONG TO THEIR USERPAGE
		http.Redirect(w, r, "/view/userpage.html", http.StatusFound)
	} else {
		if debug == true {
			fmt.Println("user has NOT been validated")
		}

		//EXPLICITLY SET LOGIN STATUS TO FALSE
		session.Values["auth"] = false
		//SAVE THAT COOKIE TO PERSON'S BROWSER
		err := session.Save(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//JUST RELOAD THE PAGE
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
	//GET THE USER'S COOKIE
	session, err := store.Get(r, appCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//SETTING LOGIN STATUS TO FALSE AND DELETING THE COOKIE
	session.Values["auth"] = false
	session.Options.MaxAge = -1

	//SAVING THE SESSION
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//SENDING THE USER BACK TO THE INDEX PAGE
	http.Redirect(w, r, "/", http.StatusFound)
}

//=====================================================================================
//THE SIGNUP HANDLER
//=====================================================================================
func signup(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit signup")
	}

	//READ IN THE SIGNUP PAGE
	t, _ := template.ParseFiles("auth/signup.html")

	//IF ITS A GET REQUEST IT JUST SERVES THE PAGE
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	//MAKE THE USER A NEW COOKIE
	session, _ := store.Get(r, appCookie)

	//COMPARE THE TWO PASSWORDS
	pass := ""
	if r.FormValue("pwd") == r.FormValue("pwdv") {
		pass = r.FormValue("pwd")
	} else {

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	if debug == true {
		fmt.Println(r.FormValue("pwd"), r.FormValue("pwdv"), r.FormValue("name"), r.FormValue("email"))
	}

	//ADD THE USER TO THE DATABASE
	errcheck := go_dev.AddUser(r.FormValue("name"), string(hash), r.FormValue("email"), "", db)
	if errcheck != true {
		fmt.Println("Hey, the user wasn't added")
	}

	//MAKING SURE THE USER IS LOGGED OUT TO ENSURE THEY
	session.Values["auth"] = false
	err = session.Save(r, w)

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
	t, _ := template.ParseFiles("/view/404errorpage.html")

	t.Execute(w, nil)
}

//=====================================================================================
//ANY SORT OF POST WILL BE HANDLED HERE
//=====================================================================================
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

// 	p := go_dev.postGet(k)

// 	t, _ := template.ParseFiles("/view/post.html")
// 	t.Execute(w, p)
// }

//=====================================================================================
//ANY SORT OF POST WILL BE HANDLED HERE
//=====================================================================================
// func ProjectHandler(w http.ResponseWriter, r *http.Request) {
// 	if debug == true {
// 		fmt.Println("Hit ProjectHandler")
// 	}
// 	session, _ := store.Get(r, "cookie-name")

// 	if heimdall(w, r) != true {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}

// 	page := file_finder("view/", w, r)

// 	if strings.Contains(page, "create") == true {

// 		t, _ := template.ParseFiles("/view/project_create.html")
// 		//WE JUST NEED TO SERVE THE PAGE WHE ITS FIRST LOADED
// 		if r.Method != http.MethodPost {
// 			t.Execute(w, nil)
// 			return
// 		}
// 		//THIS NEED STO BE CHANGED TO THE FORM VALUE KEY
// 		pj := go_dev.CreateProject(session.Values["usr"].(string), string(r.FormValue("pjn")), db)
// 		if pj != true {
// 			//THIS IS THE BEST I COULD COME UP WITH FOR DEALING WITH A POTENTIAL ERROR
// 			varmap := map[string]interface{}{
// 				"var1": "Sorry, there was an issue creating your project",
// 			}
// 			t.Execute(w, varmap)
// 			return
// 		}
// 		if string(r.FormValue("addusrs")) != "" {
// 			//MIGHT HAVE TO MODIFY IF DB TEAM ONLY ADDS ONE MEMBER AT A TIME
// 			adu := go_dev.AddProjectMembers(session.Values["usr"].(string), string(r.FormValue("pjn")), string(r.FormValue("addusrs")), db)

// 			if adu != true {
// 				//THIS IS THE BEST I COULD COME UP WITH FOR DEALING WITH A POTENTIAL ERROR
// 				varmap := map[string]interface{}{
// 					"var1": "Sorry, there was an issue adding users to your project",
// 				}
// 				t.Execute(w, varmap)
// 				return
// 			}
// 		}
// 		//I I THINK THERE IS A BETTER WAY TO DO THIS
// 		varmap := map[string]interface{}{
// 			"var1": "Your project was successfuly created! ",
// 		}
// 		t.Execute(w, varmap)
// 		return

// 	}

// 	pathVariables := mux.Vars(r)

// 	k := pathVariables["key"]

// 	//GET THE POST FROM THE DB

// 	p := go_dev.GetProject(k)

// 	t, _ := template.ParseFiles("/view/project.html")
// 	t.Execute(w, p)

// }

// func TaskHandler(w http.ResponseWriter, r *http.Request) {
// 	if debug == true {
// 		fmt.Println("Hit TaskHandler")
// 	}
// 	//session, _ := store.Get(r, "cookie-name")

// 	if heimdall(w, r) != true {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}
// 	//
// 	t, _ := template.ParseFiles("/view/task.html")

// 	pathVariables := mux.Vars(r)

// 	k := pathVariables["key"]

// 	//GET THE POST FROM THE DB

// 	p := go_dev.GetTask(k)

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
