package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Model struct {
	ID        int64      `json:"ID"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var con *sql.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	con, e = createCon(username, password, dbHost, dbPort, dbName)
	if e != nil {
		fmt.Print(e)
	}
}

func GetConn() *sql.DB {
	return con
}

/*Create mysql connection*/
func createCon(username string, password string, dbHost string, dbPort string, dbName string) (db *sql.DB, err error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("database is connected")
	}
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Printf("MySQL db is not connected %s", err.Error())
	}
	return db, err
}
