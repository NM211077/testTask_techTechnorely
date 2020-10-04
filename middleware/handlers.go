package middleware

import (
	"database/sql"
	"fmt"
	"github.com/NM211077/testTask_techTechnorely/models"
	"github.com/NM211077/testTask_techTechnorely/util"
	"log"
	//"strconv"
)

// response format
type Response struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// GetAllBooks get all books from the DB
func GetAllBooks() ([]models.Book, error) {
	// create the mysql db connection
	db := util.DBConn()
	// close db connection
	defer db.Close()

	var books []models.Book
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
		var book models.Book

		// unmarshal the row object to book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the book in the books slice
		books = append(books, book)
	}
	fmt.Println(books)
	return books, err
}

// GetBook get one book from the DB by its id
func GetBook(id int64) (models.Book, error) { //дописать если нет книги с таким id

	db := util.DBConn()

	defer db.Close()

	var book models.Book

	sqlStatement := `SELECT * FROM books WHERE id=?`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return book, err
}

//InsertBook  insert new data to books table in DB.
func InsertBook(book models.Book) (int64, error) {
	db := util.DBConn()
	defer db.Close()
	
	// create the insert sql query
	sqlStatement := `INSERT INTO books (title, author, price) VALUES (?,?,? )`
	stmt, stmtErr := db.Prepare(sqlStatement)
	util.PanicError(stmtErr)

	res, queryErr := stmt.Exec(book.Title, book.Author, book.Price)
	util.PanicError(queryErr)
	// returning bookid will return the id of the inserted book
	id, getLastInsertIDErr := res.LastInsertId()
	util.PanicError(getLastInsertIDErr)

	return id, queryErr
}

// UpdateBook update book in the DB
func UpdateBook(id int64, book models.Book) (models.Book) {
	db := util.DBConn()
	defer db.Close()
	fmt.Println("book handlers", book)
	// create the update sql query
	sqlStatement := `UPDATE books SET title=?, author=?, price=? WHERE id=?`
	stmt, stmtErr := db.Prepare(sqlStatement)
	util.PanicError(stmtErr)

	_, queryErr := stmt.Exec(book.Title, book.Author, book.Price, id)
	util.PanicError(queryErr)

	return book
}

// DeleteBook delete book in the DB
func DeleteBook(id int64) int64 {
	db := util.DBConn()
	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE id=?`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
