package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func ConnectDb(user string, password, host string, port int, dbname string) error {
	// Construct the MySQL connection string
	// Format: <user>:<password>@tcp(<host>:<port>)/<dbname>?parseTime=true
	connentionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbname)
	db, err := sql.Open("mysql", connentionString)
	if err != nil {
		return err
	}
	Db = db
	fmt.Println("Connected to db")
	return nil
}
