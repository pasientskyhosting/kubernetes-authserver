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

const (
	VERSION    = "0.3.9-warlock"
	APIVERSION = "authentication.k8s.io/v1beta1"
)

var (
	CERT             = flag.String("cert", "/etc/ssl/tls.crt", "TLS cert path")
	db               *sql.DB
	DB_CHARSET       = flag.String("charset", "utf8", "Database charset for DSN")
	DB_DSN           string
	DB_HOST          = flag.String("host", "127.0.0.1", "Database host")
	DB_NAME          = flag.String("db", "auth", "Database name")
	DB_PASS          = flag.String("pass", "auth", "Database user")
	DB_PORT          = flag.Int("port", 3306, "Database port")
	DB_USER          = flag.String("user", "auth", "Database user")
	DEBUG            = flag.Bool("debug", false, "Enable debugging output")
	err              error
	KEY              = flag.String("key", "/etc/ssl/tls.key", "TLS key path")
	NO_HTTP          = flag.Bool("http", true, "Enable HTTP access")
	NO_HTTPS         = flag.Bool("https", true, "Enable HTTPS access")
	OPT_CERT         string
	OPT_DEBUG        bool
	OPT_HTTP         bool
	OPT_HTTPS        bool
	OPT_KEY          string
	OPT_SECUREPORT   int
	OPT_UNSECUREPORT int
	SECUREPORT       = flag.Int("https_port", 8088, "Secure HTTPS port")
	UNSECUREPORT     = flag.Int("http_port", 8087, "Unsecure HTTP port")
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	flag.Parse()
	/*
		if *DB_HOST == "" {
			log.Fatal("Empty ENV: DB_HOST")
		}

		if *DB_PORT == "" {
			log.Fatal("Empty ENV: DB_PORT")
		}

		if *DB_NAME == "" {
			log.Fatal("Empty ENV: DB_NAME")
		}

		if *DB_USER == "" {
			log.Fatal("Empty ENV: DB_USER")
		}

		if *DB_PASS == "" {
			log.Fatal("Empty ENV: DB_PASS")
		}
	*/

	DB_DSN = *DB_USER + ":" + *DB_PASS + "@(" + *DB_HOST + ":" + strconv.Itoa(*DB_PORT) + ")/" + *DB_NAME + "?charset=" + *DB_CHARSET
	OPT_HTTP = *NO_HTTP
	OPT_HTTPS = *NO_HTTPS
	OPT_CERT = *CERT
	OPT_KEY = *KEY
	OPT_SECUREPORT = *SECUREPORT
	OPT_UNSECUREPORT = *UNSECUREPORT
	OPT_DEBUG = *DEBUG
	log.Printf("DB DSN: %s:*****@(%s:%d)/%s?charset=%s", *DB_USER, *DB_HOST, *DB_PORT, *DB_NAME, *DB_CHARSET)
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
	log.Println("Starting kubernetes-authserver", VERSION)

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
