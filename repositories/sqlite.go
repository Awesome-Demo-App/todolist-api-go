package repositories

import (
	"context"
	"fmt"

	domain "github.com/awesome-demo-app/todolist-api/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type todoDB struct {
	db *gorm.DB
}

func NewSQLiteDB(dbPath string) *todoDB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	fmt.Println("Repositories → Migrating SQLite DB schema")
	db.AutoMigrate(&domain.ToDo{})

	fmt.Println("Repositories → Bootstraping sample data")
	db.Create(&domain.ToDo{Summary: "Cloud Native Week announcement", Completed: true})
	db.Create(&domain.ToDo{Summary: "Prepare demo app", Completed: false})
	db.Create(&domain.ToDo{Summary: "Cloud Native Week Day 1", Completed: false})

	return &todoDB{db: db}
}

func (repo *todoDB) GetAll(ctx context.Context) ([]domain.ToDo, error) {
	tracer := otel.Tracer("github.com/awesome-demo-app/todolist-api-go")
	var todos []domain.ToDo
	fmt.Println("Repositories → Fetching all ToDo from DB")
	func(ctx context.Context) {
		var span trace.Span
		_, span = tracer.Start(ctx, "Querying DB...")
		defer span.End()
		repo.db.Find(&todos)
	}(ctx)

	return todos, nil
}

func (repo *todoDB) Create(todo domain.ToDo) domain.ToDo {
	fmt.Println("Repositories → Creating new ToDo in DB")
	repo.db.Create(&todo)

	return todo
	// TODO handle failure
}

func (repo *todoDB) Delete(id uint) {
	fmt.Println("Repositories → Soft deleting existing ToDo in DB")
	repo.db.Delete(&domain.ToDo{}, id)
	// TODO handle failure
}
