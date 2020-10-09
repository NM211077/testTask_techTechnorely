package models_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NM211077/testTask_techTechnorely/controllers"
	"github.com/NM211077/testTask_techTechnorely/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var server = controllers.Server{}
var bookInstance = models.Book{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	Database()

	ensureTableExists()
	clearTable()
	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = sql.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	} else {
		fmt.Printf("!!Cannot connect to %s database\n", TestDbDriver)
	}
}

func ensureTableExists() {
	if _, err := server.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	server.DB.Exec("DELETE FROM book")
	server.DB.Exec("ALTER SEQUENCE book_id_seq RESTART WITH 1")
}

func deleteTable() {
	if _, err := server.DB.Exec(dropTable); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS books
(
	id int  NOT NULL AUTO_INCREMENT ,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT books_bkey PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;`

const dropTable = `DROP TABLE books`

func TestCreateBook(t *testing.T) {

	newBook := models.Book{
		Title:  "testBook",
		Author: "testAuthor",
		Price:  125.20,
	}

	savedId, err := newBook.CreateBook(server.DB)

	if err != nil {
		t.Errorf("Error while saving a book: %v\n", err)
		return
	}

	require.NotEmptyf(t, savedId, "The new book creat, by id: %v.")
	//defer deleteTable()
}

func TestGetAllBooks(t *testing.T) {
	clearTable()
	newBook := models.Book{
		Title:  "testBook",
		Author: "testAuthor",
		Price:  125.00,
	}

	savedId, err := newBook.CreateBook(server.DB)

	if err != nil {
		t.Errorf("Error while saving a book: %v\n", err)
		return
	}

	require.NotEmptyf(t, savedId, "The new book creat, by id: %v.")

	books, err := models.GetBooks(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the books: %v\n", err)
		return
	}
	require.NotEmpty(t, books, "All books")
	//defer deleteTable()
}

func TestGetBook(t *testing.T) {
	clearTable()
	newBook := models.Book{
		Title:  "testBook",
		Author: "testAuthor",
		Price:  126.00,
	}

	savedId, err := newBook.CreateBook(server.DB)

	if err != nil {
		t.Errorf("Error while saving a book: %v\n", err)
		return
	}

	b := models.Book{ID: savedId}

	errID := b.GetBook(server.DB)

	if err != nil {
		t.Errorf("this is the error getting the books: %v\n", errID)
		return
	}

	require.Equal(t, b.ID, savedId)
	require.Equal(t, b.Author, "testAuthor")
	require.Equal(t, b.Title, "testBook")
	require.Equal(t, b.Price, 126.00)
	//defer deleteTable()
}

func TestUpdateBook(t *testing.T) {
	clearTable()
	newBook := models.Book{
		Title:  "testBook",
		Author: "testAuthor",
		Price:  126.00,
	}

	savedId, err := newBook.CreateBook(server.DB)

	if err != nil {
		t.Errorf("Error while saving a book: %v\n", err)
		return
	}

	updateBook := models.Book{
		Title:  "updateTestBook",
		Author: "updateTestAuthor",
		Price:  526.00,
	}
	updateBook.ID = savedId
	errID := updateBook.UpdateBook(server.DB)

	if err != nil {
		t.Errorf("this is the error updating the book: %v\n", errID)
		return
	}

	require.Equal(t, updateBook.ID, savedId)
	require.NotEqual(t, updateBook.Author, newBook.Author)
	require.NotEqual(t, updateBook.Title, newBook.Title)
	require.NotEqual(t, updateBook.Price, newBook.Price)
	//defer deleteTable()
}
func TestDeleteBook(t *testing.T) {
	clearTable()
	newBook := models.Book{
		Title:  "testBook",
		Author: "testAuthor",
		Price:  125.00,
	}

	savedId, err := newBook.CreateBook(server.DB)

	if err != nil {
		t.Errorf("Error while saving a book: %v\n", err)
		return
	}

	b := models.Book{ID: savedId}

	errID := b.DeleteBook(server.DB)

	if err != nil {
		t.Errorf("this is the error deleting the book: %v\n", errID)
		return
	}

	require.Empty(t, errID, "this book deleted")

	//defer deleteTable()
}
