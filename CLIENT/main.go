package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://localhost:8080/api/v1/courses"

// const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

type courseInfo struct {
	ID      string `json:"ID"`
	Title   string `json:"Title"`
	Details string `json: "Details"`
}

func getCourse(code string) {
	url := baseURL
	if code != "" {
		// url = baseURL + "/" + code + "?key=" + key
		url = baseURL + "/" + code
	}
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func addCourse(code string, jsonData courseInfo) {
	jsonValue, _ := json.Marshal(jsonData)
	fmt.Println(bytes.NewBuffer(jsonValue))
	// response, err := http.Post(baseURL+"/"+code+"?key="+key,"application/json", bytes.NewBuffer(jsonValue))
	response, err := http.Post(baseURL+"/"+code, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func updateCourse(code string, jsonData courseInfo) {
	jsonValue, _ := json.Marshal(jsonData)

	// request, err := http.NewRequest(http.MethodPut,baseURL+"/"+code+"?key="+key,bytes.NewBuffer(jsonValue))
	request, err := http.NewRequest(http.MethodPut, baseURL+"/"+code, bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}
func deleteCourse(code string) {
	// request, err := http.NewRequest(http.MethodDelete,baseURL+"/"+code+"?key="+key, nil)
	request, err := http.NewRequest(http.MethodDelete, baseURL+"/"+code, nil)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func main() {
	// getCourse("") //  get all courses

	// getCourse("IOS101") // get a specific course

	// jsonData := map[string]string{"title": "Applied Go Programming"}
	// addCourse("IOT210", jsonData)

	// jsonData := map[string]string{"title": "Go Concurrency Programming"}
	// updateCourse("IOT210", jsonData)
	// deleteCourse("IOT210")
	// getCourse("")

	var userInput string
	for {
		fmt.Println("Please enter a choice:")
		fmt.Println("1.Get all course")
		fmt.Println("2.Get specific course")
		fmt.Println("3.Add new course")
		fmt.Println("4.Update course")
		fmt.Println("5.Delete course")
		fmt.Println("6.Exit")

		fmt.Scanln(&userInput)
		switch userInput {
		case "1":
			getCourse("")
		case "2":
			fmt.Println("Please enter courseID:")
			var courseChoice string
			fmt.Scanln(&courseChoice)
			getCourse(courseChoice)
		case "3":
			jsonData := askForCourseInfo()
			// jsonData := map[string]string{"title": courseInfo}
			addCourse(jsonData.ID, jsonData)
			// fmt.Println(jsonData)
		case "4":
			jsonData := askForCourseInfo()
			updateCourse(jsonData.ID, jsonData)
		case "5":
			fmt.Println("Pleae enter courseID to delete")
			var keyTitle string
			fmt.Scanln(&keyTitle)
			deleteCourse(keyTitle)
		case "6":
			return
		}

	}
}

func askForCourseInfo() courseInfo {
	in := bufio.NewReader(os.Stdin)
	var courseID string
	var courseTitle string
	var courseDetails string
	fmt.Println("Please enter courseID:")
	courseID, _ = in.ReadString('\n')
	fmt.Println("Please enter course title:")
	courseTitle, _ = in.ReadString('\n')
	fmt.Println("Please enter courseDetails:")
	courseDetails, _ = in.ReadString('\n')

	return courseInfo{strings.TrimSuffix(courseID, "\n"), strings.TrimSuffix(courseTitle, "\n"), strings.TrimSuffix(courseDetails, "\n")}

}
