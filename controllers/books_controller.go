package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NM211077/testTask_techTechnorely/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite database driver
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = sql.Open("mysql", DBURL)
	fmt.Println("Successfully connected!")
	if err != nil {
		log.Fatal(err)
	}

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8002", server.Router))
}

func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		respondWithJSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	respondWithJSON(w, http.StatusBadRequest, nil)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func (server *Server) createBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&b); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	id, err := b.CreateBook(server.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, id)
}

func (server *Server) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetBooks(server.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, books)
}

func (server *Server) getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	b := models.Book{ID: id}
	if err := b.GetBook(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, err)
		default:
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, b)
}

func (server *Server) updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var b models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&b); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	b.ID = id

	if err := b.UpdateBook(server.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, b)
}

func (server *Server) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	b := models.Book{ID: id}
	if err := b.DeleteBook(server.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/book/{id}", s.getBook).Methods("GET")
	s.Router.HandleFunc("/books", s.getAllBooks).Methods("GET")
	s.Router.HandleFunc("/book", s.createBook).Methods("POST", "PUT")
	s.Router.HandleFunc("/book/{id}", s.updateBook).Methods("PUT", "POST")
	s.Router.HandleFunc("/book/{id}", s.deleteBook).Methods("DELETE")
}

