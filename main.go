package main

import (
	//"os"
	// "fmt"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

const VERSION = "0.3.1-ivanka"
const DEBUG = true
const APIVERSION = "authentication.k8s.io/v1beta1"

var db *sql.DB
var err error
var DB_DSN string

// var DEBUG_HEALTHZ int

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
	//    var DEBUG_HEALTHZ = flag.Int("debug_healthz", 0, "Enable logging of requests to /healthz")

	flag.Parse()
	DB_DSN = *DB_USER + ":" + *DB_PASS + "@(" + *DB_HOST + ":" + *DB_PORT + ")/" + *DB_NAME + "?charset=utf8"
	//    DEBUG_HEALTHZ = DEBUG_HEALTHZ
	log.Printf("DB DSN: %s:****@(%s:%s)/%s?charset=utf8", *DB_USER, *DB_HOST, *DB_PORT, *DB_NAME)
}

func Run(addr string, sslAddr string, ssl map[string]string) chan error {
	errs := make(chan error)
	router := NewRouter()

	// Starting HTTP server
	go func() {
		log.Printf("Staring HTTP service on %s", addr)

		if err := http.ListenAndServe(addr, router); err != nil {
			errs <- err
		}

	}()

	// Starting HTTPS server
	go func() {
		log.Printf("Staring HTTPS service on %s", sslAddr)
		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], router); err != nil {
			errs <- err
		}
	}()

	return errs
}

func main() {
	db, err = sql.Open("mysql", DB_DSN)
	checkErr(err)

	defer db.Close()

	checkErr(db.Ping())

	errs := Run(":8087", ":8088", map[string]string{
		"cert": "/etc/ssl/tls.crt",
		"key":  "/etc/ssl/tls.key",
	})

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		log.Printf("Could not start serving service due to (error: %s)", err)
	}

}
