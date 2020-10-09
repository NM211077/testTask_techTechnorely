package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"strconv"

	"github.com/NM211077/testTask_techTechnorely/controllers"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode:   200,
			title:        "testTitle",
			author:       "testAuthor",
			price:        12.25,
			errorMessage: "",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/books", bytes.NewBufferString(v.inputJSON))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.СreateBook_router)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
	}
    //defer deleteTable()
}

func TestGetBooks(t *testing.T) {

	clearTable()

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode:   200,
			title:        "testTitle",
			author:       "testAuthor",
			price:        12.25,
			errorMessage: "",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("GET", "/books", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetAllBooks_router)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var books []models.Book

		err = json.Unmarshal([]byte(rr.Body.String()), &books)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v\n", err)
		}

		require.Equal(t, len(books), 1)
	}
}

func TestGetBookByID(t *testing.T) {

	clearTable()
	bookSample := []struct {
		id           string
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(1)),
			inputJSON:    `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode:   200,
			title:        "testTitle",
			author:       "testAuthor",
			price:        12.25,
			errorMessage: "",
		},
		{
			id:         "unknow",
			inputJSON:  `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode: 400,
		},
	}
	for _, v := range bookSample {

		req, err := http.NewRequest("GET", "/book{id}", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetBook_router)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		require.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			require.Equal(t, "testTitle", responseMap["title"])
			require.Equal(t, "testAuthor", responseMap["author"])
			require.Equal(t, 12.25, responseMap["price"])
		}
	}
}

func TestUpdateBookByID(t *testing.T) {

	clearTable()
	bookSample := []struct {
		id           string
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			id:           "1",
			inputJSON:    `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode:   200,
			title:        "testTitle",
			author:       "testAuthor",
			price:        12.25,
			errorMessage: "",
		},
		{
			id:         "unknow",
			inputJSON:  `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode: 400,
		},
	}

	updateBookSample := []struct {
		id           string
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(1)),
			inputJSON:    `{"title":"updateTestTitle", "author": "updateTestAuthor", "price": 212.25}`,
			statusCode:   200,
			title:        "updateTestTitle",
			author:       "updateTestAuthor",
			price:        212.25,
			errorMessage: "",
		},
	}
	for _, v := range updateBookSample {

		req, err := http.NewRequest("UPDATE", "/book{id}", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateBook_router)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		
		require.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			for _, prev := range bookSample {
				require.NotEqual(t, prev.title, responseMap["title"])
				require.NotEqual(t, prev.author, responseMap["author"])
				require.NotEqual(t, prev.price, responseMap["price"])
			}
		}
		if v.statusCode == 400 || v.errorMessage != "" {
			require.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteBookByID(t *testing.T) {

	clearTable()
	bookSample := []struct {
		id           string
		inputJSON    string
		statusCode   int
		title        string
		author       string
		price        float64
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(1)),
			inputJSON:    `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode:   200,
			title:        "testTitle",
			author:       "testAuthor",
			price:        12.25,
			errorMessage: "",
		},
		{
			id:         "unknow",
			inputJSON:  `{"title":"testTitle", "author": "testAuthor", "price": 12.25}`,
			statusCode: 400,
		},
	}
	for _, v := range bookSample {

		req, err := http.NewRequest("DELETE", "/book{id}", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteBook_router)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		require.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			require.Empty(t, err, "this book deleted")
		}
	}
	defer deleteTable()
}
