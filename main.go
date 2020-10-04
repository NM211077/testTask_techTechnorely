package main

import (
	//"database/sql"
	//"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"

	"log"
	//"math/rand"
	"net/http"
	//"os"
	//"github.com/NM211077/testTask_techTechnorely/middleware"
	//"github.com/NM211077/testTask_techTechnorely/models"
	"github.com/NM211077/testTask_techTechnorely/router"
	//"strconv"
)

func main() {
	r := router.Router()
	fmt.Println("Go MySQL Tutorial")

	log.Fatal(http.ListenAndServe(":8001", r))
}
