package main

import "fmt"
import "github.com/gorilla/mux"
import "net/http"
import "io/ioutil"
import "encoding/json"

//"import encoding/json"

func userHanderGet(res http.ResponseWriter, req *http.Request) {
	//data, _ := json.Marshal("{'hello':'wercker!'}")
	content, err := ioutil.ReadFile("db/users.json")
	if err != nil {
		//Do something
		fmt.Println("Cannot open file for reading %s", err)
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write([]byte(content))
}

func RestUsersHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter RestUsersHandler")

	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/
	switch req.Method {
	case "GET":
		fmt.Println("GET request")
		userHanderGet(res, req)
		return
		// Serve the resource.
	case "POST":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Create a new record.
	case "PUT":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Update an existing record.
	case "DELETE":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Remove the record.
	default:
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Give an error message.
	}

	http.Error(res, "internal error", http.StatusInternalServerError)
	return
}

func RestUserHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter RestUserHandler")

	vars := mux.Vars(req)
	key := vars["user"]

	fmt.Println("key: ", key)

	switch req.Method {
	case "GET":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Serve the resource.
	case "POST":
		http.Error(res, "POST not valid", http.StatusInternalServerError)
		return
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

func itemsHanderGet(res http.ResponseWriter, req *http.Request) {
	//data, _ := json.Marshal("{'hello':'wercker!'}")
	content, err := ioutil.ReadFile("db/users.json")
	if err != nil {
		//Do something
		fmt.Println("Cannot open file for reading %s", err)
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write([]byte(content))
}


type item_struct struct {
	Command string;
	Data map[string]string;
}

func itemsHanderPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if (err != nil) {
		fmt.Println("Failure in reading from client %s", err)
		// FIXME: return negative status code
		return
	}

	fmt.Println(string(body))
	var t item_struct
	err = json.Unmarshal(body, &t)
	if (err != nil) {
		fmt.Println("Failure in item POST struct", err)
		return
	}

	fmt.Println(t.Command)
	fmt.Println(t.Data["Description"])


}


func RestItemsHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter RestItemsHandler")

	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/
	switch req.Method {
	case "GET":
		fmt.Println("GET request")
		//itemsHanderGet(res, req)
		return
		// Serve the resource.
	case "POST":
		fmt.Println("POST request")
		itemsHanderPost(res, req)
		return
		// Create a new record.
	case "PUT":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Update an existing record.
	case "DELETE":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Remove the record.
	default:
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Give an error message.
	}

	http.Error(res, "internal error", http.StatusInternalServerError)
	return
}
