package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	ports "github.com/awesome-demo-app/todolist-api/core/ports"
)

type HTTPHandler struct {
	todoService ports.ToDoService
}

func NewHTTPHandler(service ports.ToDoService) *HTTPHandler {
	return &HTTPHandler{
		todoService: service,
	}
}

func (handler *HTTPHandler) HandleRequests() {
	http.HandleFunc("/todo", handler.ToDoDispatchHandler)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

// Dispatch depending on method
func (handler *HTTPHandler) ToDoDispatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dispatching ToDo endpoint based on HTTP verb")
	switch r.Method {
	case "GET":
		fmt.Println("GET")
		handler.GetAllToDos(w, r)
	case "POST":
		fmt.Println("POST")
		handler.CreateNewToDo(w, r)
	case "DELETE":
		fmt.Println("DELETE")
		handler.DeleteToDoByID(w, r)
	}
}

func (handler *HTTPHandler) GetAllToDos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetAllToDos")

	todos, err := handler.todoService.GetAll()
	if err != nil {
		fmt.Fprintf(w, "Failed lol")
	}

	for _, todo := range todos {
		fmt.Fprintf(w, "- %+v\n", todo.Summary)
	}
}

func (handler *HTTPHandler) CreateNewToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: CreateNewToDo")

	queryParameters := r.URL.Query()
	log.Printf("%v\n", queryParameters)

	summaryRawData, ok := queryParameters["summary"]
	log.Printf("%v\n", summaryRawData)

	if !ok || len(summaryRawData) != 1 {
		log.Println("Could not find 'summary' parameter in POST data")
		return
	}

	summary := summaryRawData[0]
	todo := handler.todoService.Create(summary)
	serializedToDo, err := json.Marshal(&todo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%v\n", string(serializedToDo))
}

func (handler *HTTPHandler) DeleteToDoByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: DeleteToDoByID")

	queryParameters := r.URL.Query()
	log.Printf("%v\n", queryParameters)

	idRawData, ok := queryParameters["id"]
	log.Printf("%v\n", idRawData)

	if !ok || len(idRawData) != 1 {
		log.Println("Could not find 'id' parameter in POST data")
		return
	}

	id, err := strconv.ParseUint(idRawData[0], 0, 64)
	if err != nil {
		fmt.Println(err)
	}
	handler.todoService.Delete(uint(id))
}
