//creating the docker container for mysql database
// docker run --name mycourse_db -p 64893:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:latest
// mysql -P 64893 --protocol=tcp -u root -p
// CREATE database my_db;
// CREATE TABLE Courses (ID varchar(6) NOT NULL PRIMARY KEY,Title VARCHAR(30), Details VARCHAR(30));
// add in user2 to access the database

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type courseInfo struct {
	ID      string `json:"ID"`
	Title   string `json:"Title"`
	Details string `json: "Details"`
}

// func validKey(r *http.Request) bool {
// 	v := r.URL.Query()
// 	if key, ok := v["key"]; ok {
// 		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
// 			return true
// 		} else {
// 			return false
// 		}
// 	} else {
// 		return false
// 	}
// }

// used for storing courses on the REST API

var (
	mydb *sql.DB
)

func main() {
	var err error
	defer mydb.Close()
	mydb, err = sql.Open("mysql", "user2:password@tcp(0.0.0.0:64893)/my_db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("db opened!!")
		GetRecords(mydb)
	}
	// instantiate courses
	// courses = make(map[string]courseInfo)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)

	router.HandleFunc("/api/v1/courses", allcourses)
	router.HandleFunc("/api/v1/courses/{courseid}", course).Methods("GET", "PUT", "POST", "DELETE")
	fmt.Println("Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API")
}

func allcourses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of all courses")
	allresults := GetRecords(mydb)
	json.NewEncoder(w).Encode(allresults)
}
func course(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// fmt.Fprintf(w, "Details for course "+params["courseid"])
	// fmt.Fprintf(w, "\n")
	// fmt.Fprintf(w, r.Method)
	// if !validKey(r) {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("401 - Invalid key"))
	// 	return
	// }

	if r.Method == "GET" {
		if resultcourse, ok := GetSpecificRecords(mydb, params["courseid"]); ok {
			json.NewEncoder(w).Encode(
				resultcourse)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}
	if r.Method == "DELETE" {
		if _, ok := GetSpecificRecords(mydb, params["courseid"]); ok {
			DeleteRecord(mydb, params["courseid"])
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new course
		if r.Method == "POST" {
			// read the string sent to the service
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)
			// newCourse = courseInfo{"NEW201","InternalTest","Testing Putting here"}
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " +
							"information " + "in JSON format"))
					return
				}
				// check if course exists; add only if
				// course does not exist
				if _, ok := GetSpecificRecords(mydb, params["courseid"]); !ok {
					//insert
					// courses[params["courseid"]] = newCourse
					InsertRecord(mydb, newCourse.ID, newCourse.Title, newCourse.Details)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["courseid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate course ID"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information " +
					"in JSON format"))
			}
		}
		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " +
							" information " +
							"in JSON format"))
					return
				}
				// check if course exists; add only if
				// course does not exist
				if _, ok := GetSpecificRecords(mydb, params["courseid"]); !ok {

					InsertRecord(mydb, newCourse.ID, newCourse.Title, newCourse.Details)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["courseid"]))
				} else {
					// update course
					EditRecord(mydb, newCourse.ID, newCourse.Title, newCourse.Details)
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " +
					"course information " +
					"in JSON format"))
			}
		}
	}
}

//POST - create new course
// curl -H "Content-Type:application/json" -X POST http://localhost:8080/api/v1/courses/NEW201 -d "{\"id\":\"NEW201\", \"title\":\"iOS Programming\",\"details\":\"second test course\"}"
//change existing course
// curl -H "Content-Type:application/json" -X PUT http://localhost:8080/api/v1/courses/NEW201 -d "{\"id\":\"NEW201\", \"title\":\"CHANGE2\",\"details\":\"CHANGING STUFF HERE2\"}"
//delete course
//curl -X DELETE http://localhost:8080/api/v1/courses/NEW201
//curl http://localhost:8080/api/v1/courses
