package service

import (
	"errors"

	domain "github.com/awesome-demo-app/todolist-api/core/domain"
	ports "github.com/awesome-demo-app/todolist-api/core/ports"
)

type todoService struct {
	todoRepository ports.ToDoRepository
}

func New(repository ports.ToDoRepository) *todoService {
	return &todoService{
		todoRepository: repository,
	}
}

func (service *todoService) GetAll() ([]domain.ToDo, error) {
	todos, err := service.todoRepository.GetAll()

	if err != nil {
		return nil, errors.New("getting todos from repository has failed")
	}

	return todos, nil
}

func (service *todoService) Create(summary string) domain.ToDo {
	return service.todoRepository.Create(domain.ToDo{
		Summary:   summary,
		Completed: false,
	})
}

func (service *todoService) Delete(id uint) {
	service.todoRepository.Delete(id)
}
