package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/xo/dburl"
)

type Model struct {
	ID        int64      `json:"ID"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var conn *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		// TODO: Use proper logger
		fmt.Printf("Could not load environment files. %s", err.Error())
	}

	var ok bool
	conn, err, ok = tryConnectHerokuJawsDB()
	if err != nil {
		// TODO: Use proper logger
		fmt.Printf("Could not connect to JawsDB. %s", err.Error())
		return
	} else if ok {
		return
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	conn, err = createCon(username, password, dbHost, dbPort, dbName)
	if err != nil {
		fmt.Print(err)
	}
}

func GetConn() *sql.DB {
	return conn
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

func tryConnectHerokuJawsDB() (*sql.DB, error, bool){
	dbURI := os.Getenv("JAWSDB_MARIA_URL")
	if len(dbURI) == 0 {
		return nil, nil, false
	}

	db, err := dburl.Open( dbURI)
	if err != nil {
		return nil, err, false
	} else {
		fmt.Println("database is connected")
	}

	err = db.Ping()
	if err != nil {
		return nil, err, false
	}

	return db, nil, true
}