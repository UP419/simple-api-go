package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "1984", Author: "Goerge Orwell", Quantity: 4},
	{ID: "2", Title: "Animal Farm", Author: "George Orwell", Quantity: 3},
	{ID: "3", Title: "Harry Potter: Goblet of fire", Author: "J.K. Rowling", Quantity: 17},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByIdHelper(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func getBookByIdHelper(id string) (*book, error) {
	for i := range books {
		if books[i].ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id query parameter missing"})
		return
	}

	book, err := getBookByIdHelper(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity < 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is not available at the moment"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id query parameter missing"})
		return
	}

	book, err := getBookByIdHelper(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("books/:id", getBookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/returnBook", returnBook)
	router.POST("/books/create", createBook)
	router.Run("localhost:8080")
}
