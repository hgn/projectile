package main

import "os"
import "fmt"
import "github.com/gorilla/mux"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "errors"
import "time"

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

// JSON encoded content in form of
// { desc: "foo", bar: "baz" }
// no comas are added, go logic will
// add these as a prepartion before 
// transport
type item_file_line struct {
	Id string
	Description string
}

type item_struct struct {
	Command string
	Data    map[string]string
}

func checkIfDataisValid(data item_struct) error {
	// read over all items and get the highest
	// item id, we will get the new item id+1
	// item recycling is not implemented now
	desc, ok := data.Data["Description"]
	if ok != true {
		return errors.New("Item Description mission from struct")
	}

	fmt.Println(desc)

	return nil
}

// Open file in append mode and add
// JSON encoded line
func appendItemData(data item_file_line) error {

	_, err := os.Stat("db/items.json")
	if err != nil {
		// no such file or dir
	}

	f, err := os.OpenFile("db/items.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s, _ := json.Marshal(data)

	if _, err = f.WriteString(string(s) + "\n"); err != nil {
		panic(err)
	}

	return nil
}

// return UNIX time in milliseconds,
// this should be unique enough for one user
// for now
func getNewItemId() string {
	ct := time.Now()
	now := ct.Nanosecond()
	miliSeconds := (now % 1e9) / 1e6
	sec := ct.UTC().Format("20060102150405")
	return fmt.Sprintf("%s%03d", sec, miliSeconds)
}

func addItem(data item_struct) error {

	err := checkIfDataisValid(data)
	if err != nil {
		return err
	}

	var new_data item_file_line
	new_data.Id = getNewItemId()
	new_data.Description = data.Data["Description"]

	err = appendItemData(new_data)
	if err != nil {
		return err
	}

	return nil
}

type client_items_response_msg struct {
	Status string
}


func itemsHanderPost(w http.ResponseWriter, r *http.Request) {
	var t item_struct

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failure in reading from client %s", err)
		// FIXME: return negative status code
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("Failure in item POST struct", err)
		return
	}

	fmt.Println(t.Command)
	err = errors.New("Not implemented")
	switch t.Command {
	case "add":
		err = addItem(t)
		// Serve the resource.
	case "del":
	default:
		// Give an error message.
	}

	if err != nil {
		fmt.Println("item error occured", err)
		return
	}

	var res_msg client_items_response_msg
	res_msg.Status = "success"

	msg, _ := json.Marshal(res_msg)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(msg))
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
