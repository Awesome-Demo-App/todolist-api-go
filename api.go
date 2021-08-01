package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")

	fmt.Fprintf(w, "Welcome to the home page of the ToDo List API!\n\n")
	fmt.Fprintf(w, "This is part of the Awesome Demo App project, and is used for the Cloud Native Week at Publicis Sapient France :)")
}

func toDoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAllToDos(w, r)
	case "POST":
		addNewToDo(w, r)
	}
}

func getAllToDos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllToDos")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	var todos []ToDo
	db.Find(&todos)
	//fmt.Fprintf(w, "%v\n", todos)
	for _, todo := range todos {
		fmt.Fprintf(w, "- %+v\n", todo.Summary)
	}
}

func addNewToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: addNewToDo")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	queryParameters := r.URL.Query()
	log.Printf("%v\n", queryParameters)

	summaryRawData, ok := queryParameters["summary"]
	log.Printf("%v\n", summaryRawData)

	if !ok || len(summaryRawData) != 1 {
		log.Println("Could not find 'summary' parameter in POST data")
		return
	}

	summary := summaryRawData[0]

	log.Printf("%v\n", summary)
	db.Create(&ToDo{Summary: summary, Completed: false})
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/todo", toDoHandler)
	// http.HandleFunc("/todo", getAllToDos)
	// http.HandleFunc("/todo", getAllToDos) // POST
	log.Fatal(http.ListenAndServe(":10000", nil))
}
