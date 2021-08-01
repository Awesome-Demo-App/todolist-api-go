package main

import "gorm.io/gorm"

type ToDo struct {
	gorm.Model
	Summary   string
	Completed bool
}
