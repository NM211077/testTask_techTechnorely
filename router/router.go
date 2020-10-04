package router

import (
	"encoding/json"
	"github.com/NM211077/testTask_techTechnorely/middleware"
	"github.com/NM211077/testTask_techTechnorely/models"
	"github.com/NM211077/testTask_techTechnorely/util"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"net/http"
	"strconv"
)

var books []models.Book

// getAllBooks will return all the books
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get all the users in the db
	books, err := middleware.GetAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all books. %v", err)
	}
	// send all the books as response
	json.NewEncoder(w).Encode(books)
}

// getBook will return a single book by its id
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get the id from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the GetBook function with book id to retrieve a single book
	book, err := middleware.GetBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(book)
}

// createBook create a book in the mysql db
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	// call insert book function and pass the book
	id, insertErr := middleware.InsertBook(book)

	if insertErr != nil {
		fmt.Println("Something wrong on our server")
		util.PanicError(insertErr)
	} else {
		msg := fmt.Sprintf("Book created successfully. New title: %v, author:%v,price: %v. ", book.Title, book.Author, book.Price)
		//format a response object
		res := middleware.Response{
			ID:      id,
			Message: msg,
		}
		// send the response
		json.NewEncoder(w).Encode(res)
	}
}

// UpdateBook update book's detail in the mysql db
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty user of type models.Book
	var book models.Book

	// decode the json request to book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update book to update the book
	updatedRows := middleware.UpdateBook(int64(id), book)

	// format the message string
	msg := fmt.Sprintf("Book updated successfully. New title: %v, author:%v,  price: %v. ", updatedRows.Title, updatedRows.Author, updatedRows.Price)

	// format the response message
	res := middleware.Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//deleteBook delete book's detail in the  db
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteBook, convert the int to int
	deletedRows := middleware.DeleteBook(int64(id))

	// format the message string
	msg := fmt.Sprintf("Book updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := middleware.Response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	return router
}
