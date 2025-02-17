package models

import (
	"go-bookstore/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})

}

func CreateBook(b *Book) *Book {
	db.Create(&b)
	return b

}

func GetBook() []Book {
	var books []Book
	db.Find(&books)
	return books

}
func GetBookById(Id int64) (*Book, *gorm.DB) {
	var GetBook Book
	db.Where("ID=?", Id).Find(&GetBook)
	return &GetBook, db

}

func DeleteBookById(Id int64) *Book {
	var GetBook Book
	db.Where("ID=?", Id).Delete(&GetBook)
	return &GetBook
}
