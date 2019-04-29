package main

import (
	"Co-Lab/go_dev"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

// type Person struct {
// 	usr  string
// 	auth bool
// }

// store will hold all session data
// securecookie.GenerateRandomKey(64)
var store = sessions.NewCookieStore([]byte("super-secret-key"))

const appCookie = "DeleciousCoLabCookies"

//=====================================================================================
//THIS INITIALIZES COOKIE STUFF
//=====================================================================================
func init() {
	//gob.Register(&Person{})

	store.Options = &sessions.Options{
		// Domain: "localhost",
		MaxAge:   3600 * 72, // 3days
		HttpOnly: true,
	}
}

//=====================================================================================
//THIS IS THE LOGIN HANDLER (he thicc)
//=====================================================================================
func login(w http.ResponseWriter, r *http.Request) {

	if debug == true {
		fmt.Println("Hit login")
		fmt.Println(r.Method)
	}

	//MAKEING A NEW COOKIE FOR THE USER
	session, err := store.Get(r, appCookie)
	if err != nil {
		fmt.Println("ERROR WITH store.Get", err)
	}
	//PARSE THE LOGIN PAGE
	t, err := template.ParseFiles("auth/login.html")
	if err != nil {
		fmt.Println("Login Handler parsing error", err)
	}
	//A CHECK FOR A POST METHOD THAT MIGHT NOT BE NECESSARY ANYMORE
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	//READ IN THE PASSWORD ENTERED
	pw := r.FormValue("pwd")
	//PASSWORD ENCRYPTION
	// hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	// if err != nil {
	// 	log.Println(err)
	// }

	// if debug == true {
	// 	fmt.Println(r.FormValue("usr"), string(hash))
	// }

	//ALL OF THE AUTH
	if go_dev.Validate(r.FormValue("usr"), string(pw), db) == true {
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

		//SANITY CHECK TO MAKE SURE THE COOKIE WAS ACTUALLY SAVED
		session, err := store.Get(r, appCookie)
		if err != nil {
			fmt.Println("ERROR WITH store.Get", err)
		}

		if debug == true {
			fmt.Println("getUser after save", session.Values["auth"])
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//SEND THEM ALONG TO THEIR USERPAGE
		http.Redirect(w, r, "/view/userpage.html", 302)
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
	t, err := template.ParseFiles("auth/signup.html")
	if err != nil {
		fmt.Println("task view Handler parsing error", err)
	}

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

	// hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	// if err != nil {
	// 	log.Println(err)
	// }

	if debug == true {
		fmt.Println(r.FormValue("pwd"), r.FormValue("pwdv"), r.FormValue("name"), r.FormValue("email"))
	}

	//ADD THE USER TO THE DATABASE
	errcheck := go_dev.AddUser(r.FormValue("name"), string(pass), r.FormValue("email"), "", db)
	if errcheck != true {
		fmt.Println("Hey, the user wasn't added")
		//fmt.Println("This is the password", string(hash), len(string(hash)))
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
		session.Values["auth"] = false
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error saving cookie")
			return false
		}
		return false
	}
	//session.Values["auth"] = true
	return true
}
