package main

import (
	"database/sql"
	"fmt"
)

func GetRecords(db *sql.DB) map[string]courseInfo {
	results, err := db.Query("Select * FROM my_db.Courses")
	result := make(map[string]courseInfo)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var course courseInfo
		err = results.Scan(&course.ID, &course.Title,
			&course.Details)
		if err != nil {
			panic(err.Error())
		}
		result[course.ID] = course
		fmt.Println(course)
	}
	// fmt.Println("I'm here at GetRecords")
	// fmt.Println(result)
	return result
}
func GetSpecificRecords(db *sql.DB, courseID string) (courseInfo, bool) {
	results, err := db.Query("Select * FROM my_db.Courses WHERE ID=?", courseID)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		// map this type to the record in the table
		var course courseInfo
		err = results.Scan(&course.ID, &course.Title,
			&course.Details)
		if err != nil {
			panic(err.Error())
		}
		return course, true
	} else {
		return courseInfo{}, false
	}
}

func DeleteRecord(db *sql.DB, courseID string) {
	results, err := db.Exec("DELETE FROM my_db.Courses WHERE ID=?", courseID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}
func InsertRecord(db *sql.DB, courseID string, title string, details string) {
	results, err := db.Exec("INSERT INTO my_db.Courses VALUES (?,?,?)", courseID, title, details)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}
func EditRecord(db *sql.DB, courseID string, title string, details string) {
	results, err := db.Exec("UPDATE my_db.Courses SET Title=?, Details=? WHERE ID=?", title, details, courseID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		fmt.Println(rows)
	}
}
