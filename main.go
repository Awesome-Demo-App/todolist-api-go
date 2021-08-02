package main

import (
	services "github.com/awesome-demo-app/todolist-api/core/services"
	handlers "github.com/awesome-demo-app/todolist-api/handlers"
	repositories "github.com/awesome-demo-app/todolist-api/repositories"
)

func main() {
	todoRepository := repositories.NewSQLiteDB("file::memory:?cache=shared&mode=rwc&_journal_mode=WAL")
	todoService := services.New(todoRepository)
	todoHandler := handlers.NewHTTPHandler(todoService)

	todoHandler.HandleRequests()
}
