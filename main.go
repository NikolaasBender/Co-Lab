package main

import (
	//"fmt"
	//"html/template"
	"Co-Lab/go_dev"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	//"time"

	"github.com/gorilla/mux"
)

//ROUTER CONTSRUCTER
//VERY IMPORTANT

var db *sql.DB

func newRouter() *mux.Router {
	r := mux.NewRouter()
	//http.HandleFunc("/favicon.ico", faviconHandler)
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//THIS IS 100% VOODOO - DONT FUCKING TOUCH THIS UNDER ANY CIRCUMSTANCES
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	//ALL THE REST OF THESE ARE FINE TO MESS WITH
	//VIEWS SUB ROUTER
	s := r.PathPrefix("/view").Subrouter()
	s.HandleFunc("/", ViewHandler)
	s.HandleFunc("/{page}", ViewHandler)

	//SOME SPECIFIC VIEWING HANDLERS
	r.HandleFunc("/project_view/{key}", ProjectViewHandler)
	t := r.PathPrefix("/task_view").Subrouter()
	t.HandleFunc("/view/{key}", TaskViewHandler)
	t.HandleFunc("/comment/{key}", TaskCommentHandler)
	t.HandleFunc("/status/{key}", TaskStatusHandler)

	//SESSIONS AND STUFF
	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/signup", Signup)

	//DEFAULT ROUTE WHEN SOMEONE HITS THE SITE
	r.HandleFunc("/", IndexHandler)

	//404 HANDLEING WITH CUSTOM PAGE
	r.NotFoundHandler = http.HandlerFunc(notFound)

	return r
}

func main() {
	if debug == true {
		fmt.Println("Co-Lab core starting up")
	}

	db = go_dev.Initialize()
	if db == nil {
		fmt.Println("db is bad")
	}

	//WE NEED A ROUTER
	r := newRouter()

	port := ":8080"

	fmt.Println("go to ->  http://localhost" + port)
	//RUNS THE SERVER
	log.Fatal(http.ListenAndServe(port, r))
}
