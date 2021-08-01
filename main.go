package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&ToDo{})

	db.Create(&ToDo{Summary: "Cloud Native Week announcement", Completed: true})
	db.Create(&ToDo{Summary: "Prepare demo app", Completed: false})
	db.Create(&ToDo{Summary: "Cloud Native Week Day 1", Completed: false})

	// Read
	//var element Element
	//db.First(&element, 0)                      // find product with integer primary key
	//db.First(&element, "name = ?", "Hydrogen") // find product with code D42

	// Update - update product's price to 200
	//db.Model(&element).Update("AtomicNumber", 2)
	// Update - update multiple fields
	//db.Model(&element).Updates(Element{AtomicNumber: 2, Name: "Helium"}) // non-zero fields
	//db.Model(&element).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	//db.Delete(&element, 1)

	var todos []ToDo
	db.Find(&todos)
	// Dump all elements
	for _, todo := range todos {
		fmt.Printf("%+v\n", todo)
	}

	handleRequests()
}
