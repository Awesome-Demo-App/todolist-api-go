package repositories

import (
	domain "github.com/awesome-demo-app/todolist-api/core/domain"
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
	db.AutoMigrate(&domain.ToDo{})

	db.Create(&domain.ToDo{Summary: "Cloud Native Week announcement", Completed: true})
	db.Create(&domain.ToDo{Summary: "Prepare demo app", Completed: false})
	db.Create(&domain.ToDo{Summary: "Cloud Native Week Day 1", Completed: false})

	return &todoDB{db: db}
}

func (repo *todoDB) GetAll() ([]domain.ToDo, error) {
	var todos []domain.ToDo
	repo.db.Find(&todos)

	return todos, nil
}

func (repo *todoDB) Create(todo domain.ToDo) domain.ToDo {
	repo.db.Create(&todo)

	return todo
	// TODO handle failure
}

func (repo *todoDB) Delete(id uint) {
	repo.db.Delete(&domain.ToDo{}, id)
	// TODO handle failure
}
