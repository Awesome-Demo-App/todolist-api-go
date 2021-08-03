package service

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/awesome-demo-app/todolist-api/core/domain"
	ports "github.com/awesome-demo-app/todolist-api/core/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

type todoService struct {
	todoRepository ports.ToDoRepository
}

func New(repository ports.ToDoRepository) *todoService {
	return &todoService{
		todoRepository: repository,
	}
}

func (service *todoService) GetAll(ctx context.Context) ([]domain.ToDo, error) {
	tracer := otel.Tracer("github.com/awesome-demo-app/todolist-api-go")

	// we're ignoring errors here since we know these values are valid,
	// but do handle them appropriately if dealing with user-input
	foo, _ := baggage.NewMember("AwesomeMetadata", "Value")
	bag, _ := baggage.New(foo)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	return func(ctx context.Context) ([]domain.ToDo, error) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "getting-all-todos")
		defer span.End()

		//span.SetAttributes(anotherKey.String("yes"))

		fmt.Println("Domain → Getting all ToDos")
		todos, err := service.todoRepository.GetAll(ctx)

		if err != nil {
			return nil, errors.New("getting todos from repository has failed")
		}

		fmt.Printf("Domain ← Got all %d ToDos\n", len(todos))
		span.AddEvent("Got ToDos", trace.WithAttributes(attribute.Int("todos-retrieved", len(todos))))
		return todos, nil
	}(ctx)
}

func (service *todoService) Create(summary string, ctx context.Context) domain.ToDo {
	fmt.Printf("Domain → Creating a new ToDo with summary: %s\n", summary)
	return service.todoRepository.Create(domain.ToDo{
		Summary:   summary,
		Completed: false,
	}, ctx)
}

func (service *todoService) Delete(id uint, ctx context.Context) {
	fmt.Printf("Domain → Deleting ToDo %d\n", id)
	service.todoRepository.Delete(id, ctx)
}
