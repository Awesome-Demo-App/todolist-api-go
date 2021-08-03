package ports

import (
	"context"

	domain "github.com/awesome-demo-app/todolist-api/core/domain"
)

type ToDoRepository interface {
	GetAll(context.Context) ([]domain.ToDo, error)
	Create(domain.ToDo, context.Context) domain.ToDo
	Delete(uint, context.Context)
}

type ToDoService interface {
	GetAll(context.Context) ([]domain.ToDo, error)
	Create(string, context.Context) domain.ToDo
	Delete(uint, context.Context)
}
