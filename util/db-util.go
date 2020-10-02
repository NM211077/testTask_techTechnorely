package util

import (
	"database/sql"
	"os"
)
// Open up our database connection.
func DBConn() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_BOOKS"))
	PanicError(err)

	return db
}
//handling with repetitive task of error-handling
func PanicError(err error) {
	if err != nil {
		panic(err.Error())
	}
}


//DB Create Operation
func CreateBook(shoppingList types.ShoppingList) (int64, error) {
	db := DBConn()
	defer db.Close()

	query := "INSERT INTO shopping_list (name, qty, unit) VALUES(?, ?, ?);"
	stmt, stmtErr := db.Prepare(query)
	util.PanicError(stmtErr)

	res, queryErr := stmt.Exec(shoppingList.Name, shoppingList.Qty, shoppingList.Unit)
	util.PanicError(queryErr)

	id, getLastInsertIdErr := res.LastInsertId()
	util.PanicError(getLastInsertIdErr)

	return id, queryErr
}



