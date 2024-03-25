// using Gin Framework
// Hight performance high productivity
// simple api
// Building a library api checkin a book , checkout a book , get a book by id
package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	//exported field name and we are mentioning json fields so that we can access through json file format in api when sending and receiving data
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// slice of books and collecting data we are not using database in this project
var books = []book{
	{ID: "1", Title: "Hello World Java", Author: "Havish", Quantity: 5},
	{ID: "2", Title: "Think Like A monk", Author: "Jay Shetty", Quantity: 10},
	{ID: "3", Title: "Find Smiles Inside", Author: "Romero", Quantity: 4},
}

// helper function to get book by  id
// returns pointer to book and errir  nil if no book found
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil

		}
	}
	return nil, errors.New("book not found")

}

// return all the books (json version of books in our database books)
// gin.Context is the all the essential information about the request inorder to return a response
func getBooks(c *gin.Context) {
	//Indentedjson is nicely formatted json, books is the data
	c.IndentedJSON(200, books)
}

func createBook(c *gin.Context) {
	var newBook book
	//we use c to bind the json (as part of the request payload) new struct to the book
	//error != nil means we got the error then return will return the error message
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// function to use getbookbyId and use bookByid to fetch the data
func bookByid(c *gin.Context) {
	//path parameter
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		//gin.H allows us to write a custom error json we want to return
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// function to checkout a book (this function subtracts one book from the quantity )
func checkoutBook(c *gin.Context) {
	//we get the query parameter from the id
	id, ok := c.GetQuery("id")
	//mistake in the given input
	//check if we have the id as query parameter or not
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}
	//if thats not the case we get book by id
	book, err := getBookById(id)
	// if error when trying to find the book if book is not existed
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	//reduce the quantity of the book selected
	//check the quantity of the book <= 0
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available "})
		return
	}

	//if all the things above are false then we are going to do the below task

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

// checkin book is noting but returning book to the library means adding back 1 to yhe quantity
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query id parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)

}

//delete book by id

func removeBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		//gin.H allows us to write a custom error json we want to return
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	//If a matching book is found, the function removes it from the books slice using the append function with slicing.
	//The books[:i] slice contains all the elements before the matched book, and books[i+1:] contains all the elements after the matched book.
	//By appending these two slices together, the matched book is effectively removed from the books slice.
	for i, b := range books {
		if b.ID == book.ID {
			books = append(books[:i], books[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Book Deleted Successfully"})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book Not Found"})
}

// Gin router is responsible for handling different routes and different end points of our api
func main() {
	//https://gin-gonic.com/docs/examples/http-method/
	//using router variable we can route to a specific route through a function
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/book/:id", bookByid)
	//PATCH updates only the specified fields or properties of a resource, leaving the rest of the resource unchanged.
	//UPDATE (PUT) replaces the entire resource representation with the new data provided in the request body.
	//http://localhost:8084/checkout?id=8
	router.PATCH("checkout", checkoutBook)

	router.PATCH("return", returnBook)
	//delete the book from directory
	router.DELETE("/delete/:id", removeBook)
	router.Run("localhost:8084")

}
