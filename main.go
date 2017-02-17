package main

import (
	"log"
    //"os"
    // "fmt"
    "flag"
	"net/http"
	"database/sql"
  _ "github.com/go-sql-driver/mysql"
)

const VERSION = "0.1-earlybird"

var db *sql.DB 
var err error
var DB_DSN string


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func init() {
    log.Println("Starting kubernetes-authserver", VERSION)
    var DB_HOST = flag.String("host", "127.0.0.1", "Database host")
    var DB_PORT = flag.String("port", "3306", "Database port")
    var DB_NAME = flag.String("db", "auth", "Database name")
    var DB_USER = flag.String("user", "auth", "Database user")
    var DB_PASS = flag.String("pass", "auth", "Database user")
    flag.Parse()
    DB_DSN = *DB_USER + ":" + *DB_PASS + "@(" + *DB_HOST + ":" + *DB_PORT + ")/" + *DB_NAME + "?charset=utf8"
    log.Printf("DB DSN: %s:****@(%s:%s)/%s?charset=utf8", *DB_USER, *DB_HOST, *DB_PORT, *DB_NAME)
}
 

func main() {
    db, err = sql.Open("mysql", DB_DSN)
    checkErr(err)
    defer db.Close()
    checkErr(db.Ping())
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8087", router))
}
