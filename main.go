package main

import (
	//"os"
	// "fmt"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

const VERSION = "0.3.6-pryingeye"
const DEBUG = true
const APIVERSION = "authentication.k8s.io/v1beta1"

var db *sql.DB
var err error
var DB_DSN string
var OPT_HTTPS bool
var OPT_HTTP bool
var OPT_CERT string
var OPT_KEY string
var OPT_UNSECUREPORT int
var OPT_SECUREPORT int

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	log.Println("Starting kubernetes-authserver", VERSION)
	var DB_USER = flag.String("user", "auth", "Database user")
	var DB_PASS = flag.String("pass", "auth", "Database user")
	var DB_HOST = flag.String("host", "127.0.0.1", "Database host")
	var DB_PORT = flag.String("port", "3306", "Database port")
	var DB_NAME = flag.String("db", "auth", "Database name")
	var CERT = flag.String("cert", "/etc/ssl/tls.crt", "TLS cert path")
	var KEY = flag.String("key", "/etc/ssl/tls.key", "TLS key path")
	var NO_HTTP = flag.Bool("http", true, "Enable HTTP access")
	var NO_HTTPS = flag.Bool("https", true, "Enable HTTPS access")
	var UNSECUREPORT = flag.Int("http_port", 8087, "Unsecure HTTP port")
	var SECUREPORT = flag.Int("https_port", 8088, "Secure HTTPS port")
	flag.Parse()
	DB_DSN = *DB_USER + ":" + *DB_PASS + "@(" + *DB_HOST + ":" + *DB_PORT + ")/" + *DB_NAME + "?charset=utf8"
	OPT_HTTP = *NO_HTTP
	OPT_HTTPS = *NO_HTTPS
	OPT_CERT = *CERT
	OPT_KEY = *KEY
	OPT_SECUREPORT = *SECUREPORT
	OPT_UNSECUREPORT = *UNSECUREPORT
	log.Printf("DB DSN: %s:****@(%s:%s)/%s?charset=utf8", *DB_USER, *DB_HOST, *DB_PORT, *DB_NAME)
}

func startDBPolling() {
	for {
		time.Sleep(30 * time.Second)
		err := db.Ping()
		if err != nil {
			log.Printf("DB ERROR: %s", err)
		}
	}
}

func main() {
	errs := make(chan error)
	router := NewRouter()
	db, err = sql.Open("mysql", DB_DSN)
	checkErr(err)
	defer db.Close()
	go startDBPolling()

	if OPT_HTTP {
		// Starting HTTP server
		go func() {
			log.Printf("Staring HTTP service on %d", OPT_UNSECUREPORT)
			if err := http.ListenAndServe(":"+strconv.Itoa(OPT_UNSECUREPORT), router); err != nil {
				errs <- err
			}
		}()
	}

	if OPT_HTTPS {
		// Starting HTTPS server
		go func() {
			log.Printf("Staring HTTPS service on %d", OPT_SECUREPORT)
			if err := http.ListenAndServeTLS(":"+strconv.Itoa(OPT_SECUREPORT), OPT_CERT, OPT_KEY, router); err != nil {
				errs <- err
			}
		}()
	}

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		log.Printf("Could not start serving service due to (error: %s)", err)
	}

}
