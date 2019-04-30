package main

import (
	"Co-Lab/go_dev"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

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
	t, err := template.ParseFiles(page)
	if debug == true {
		fmt.Println("PASSED PAGE PARSING")
	}
	if err != nil {
		fmt.Println("PARSED THE PAGE INCORRECTLY", err)

	}
	//GET OUR APP COOKIE FOR USE LATER
	session, _ := store.Get(r, appCookie)

	//HANDLING THE USERPAGE
	if strings.Contains(page, "userpage") == true {
		//ASK SQL TEAM FOR ALL THE USER PAGE STUFF
		p := go_dev.PopulateUserPage(session.Values["usr"].(string), db)
		if debug == true {
			fmt.Println("POPULATED THE USER PAGE CORRECTLY")
			fmt.Println(p)
		}

		t.Execute(w, p)

		if debug == true {
			fmt.Println("EXECUTED THE USERPAGE")
		}

		return
	}

	//HANDLING CREATION
	if strings.Contains(page, "_create") == true {
		if r.Method != http.MethodPost {
			t.Execute(w, nil)
			return
		}
		if strings.Contains(page, "project") == true {
			worked := go_dev.CreateProject(session.Values["usr"].(string), string(r.FormValue("pjn")), db)
			if worked != true {
				fmt.Println("Error creating a project")
				t.Execute(w, nil)
			} else {
				//IF CREATION WAS SUCCESSFUL THEN REDIRECT TO USER PAGE
				http.Redirect(w, r, "/view/userpage.html", http.StatusFound)
				return
			}

		}
		// if strings.Contains(page, "task") == true {
		// 	pathVariables := mux.Vars(r)
		// 	id, _ := strconv.Atoi(string(pathVariables["key"]))
		// 	worked := go_dev.createTask(id, session.Values["usr"].(string), string(r.FormValue("pjn")), db)
		// 	if worked != true {
		// 		fmt.Println("Error creating task")
		// 	}
		// }
		//A SUCCESS MESSAGE WOULD BE BETTER
		t.Execute(w, nil)

		return
	}
	if strings.Contains(page, "allMyTasks") == true {
		
		tasks := go_dev.GetUserTasks(session.Values["usr"].(string), db)
		t.Execute(w, tasks)
		return
	}

	//THE DEFAULT FOR THE VIEW HANDLER
	//DEALS WITH REALLY BASIC STATIC PAGES
	t.Execute(w, nil)

}
