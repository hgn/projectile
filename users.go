
package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func UsersHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter WelcomeHandler")

	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/

	p, err := loadPage("welcome")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}
	fmt.Println(p)
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}


func UserHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter WelcomeHandler")

	vars := mux.Vars(req)
	key := vars["user"]

	fmt.Println("key: %s", key)

	switch req.Method {
	case "GET":
		// Serve the resource.
	case "POST":
		// Create a new record.
	case "PUT":
		// Update an existing record.
	case "DELETE":
		// Remove the record.
	default:
		// Give an error message.
	}


	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/

	p, err := loadPage("welcome")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}
	fmt.Println(p)
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}
