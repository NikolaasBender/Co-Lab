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
//THIS STARTED REALLY LEAN THEN IT GOT TO FINALS AND STRESS ATE LOGIC
//First checks if user is logged in
//Then finds the page
//Then parses the file
//After, it gets our cookie for the user
//Then it decides if the route needs any special treatment and decides how to fill the page or ingest a form
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

	//SPECIFYING THE USERNAME UP HERE BECAUSE IT'S USED SO MUCH
	user := session.Values["usr"].(string)

	//HANDLING CREATION
	if strings.Contains(page, "_create") == true {

		if strings.Contains(page, "project") == true {
			if r.Method != http.MethodPost {
				t.Execute(w, nil)
				return
			}
			worked := go_dev.CreateProject(user, string(r.FormValue("pjn")), db)
			adusr := string(r.FormValue("addusrs"))
			if adusr != "" {
				usrs := strings.Split(adusr, ",")
				for _, usr := range usrs {
					err := go_dev.AddProjectMembers(user, string(r.FormValue("pjn")), usr, db)
					if err != true {
						fmt.Println("Error adding ", usr, " to project")
					}
				}
			}

			if worked != true {
				fmt.Println("Error creating a project")
				t.Execute(w, nil)
			} else {
				//IF CREATION WAS SUCCESSFUL THEN REDIRECT TO USER PAGE
				http.Redirect(w, r, "/view/userpage.html", http.StatusFound)
				return
			}

		}
		if strings.Contains(page, "task") == true {
			pjcts := go_dev.GetProjects(user, db)
			if r.Method != http.MethodPost {
				t.Execute(w, pjcts)
				return
			}
			pjs := string(r.FormValue("pjc_sel"))
			nam := string(r.FormValue("name"))
			due := string(r.FormValue("dd"))
			des := string(r.FormValue("des"))
			ok := go_dev.CreateTask(pjs, user, nam, db)
			if ok != true {
				fmt.Println("error creating task")
			}
			ok = go_dev.AddDescription(pjs, user, nam, des, db)
			if ok != true {
				fmt.Println("error adding description to task")
			}
			ok = go_dev.DueDate(pjs, user, nam, due, db)
			if ok != true {
				fmt.Println("error adding due date to task")
			}
			t.Execute(w, pjcts)
			return
		}
		//A SUCCESS MESSAGE WOULD BE BETTER
		t.Execute(w, nil)

		return
	}
	if strings.Contains(page, "allMyTasks") == true {

		tasks := go_dev.GetUserTasks(session.Values["usr"].(string), db)
		if debug == true {
			fmt.Println("POPULATING THE USER TASK PAGE RETURNED: ", tasks)
		}
		t.Execute(w, tasks)
		return
	}

	//THE DEFAULT FOR THE VIEW HANDLER
	//DEALS WITH REALLY BASIC STATIC PAGES
	t.Execute(w, nil)

}
