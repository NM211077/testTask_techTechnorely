package models

import (
	"fmt"
	"database/sql"
	"log"
)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
func PanicError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (b *Book) GetBook(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM books WHERE id=?",
		b.ID).Scan(&b.ID,&b.Title, &b.Author, &b.Price)
}

func (b *Book) UpdateBook(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE books SET title=?, author=?, price=? WHERE id=?",
			b.Title, b.Author, b.Price, b.ID)

	return err
}

func (b *Book) DeleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=?", b.ID)

	return err
}

func (b *Book) CreateBook(db *sql.DB) (int,error) {

    // create the insert sql query
	sqlStatement := ` INSERT INTO books (title, author, price) VALUES (?,?,? )`
	stmt, stmtErr := db.Prepare(sqlStatement)
	PanicError(stmtErr)

	res, queryErr := stmt.Exec(b.Title, b.Author, b.Price)
	PanicError(queryErr)
	// returning bookid will return the id of the inserted book
	id, getLastInsertIDErr := res.LastInsertId()
	PanicError(getLastInsertIDErr)

	var returnedID int
	returnedID = int(id)
	return returnedID, queryErr
}

func GetBooks(db *sql.DB) ([]Book, error) {
	var books []Book
	// create the select sql query
	sqlStatement := `SELECT * FROM books`
	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var book Book
		// unmarshal the row object to book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the book in the books slice
		books = append(books, book)
		fmt.Println("books",books)
	}

	return books, err
}
