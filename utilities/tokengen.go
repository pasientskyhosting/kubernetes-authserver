package main

import (
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
	"log"
)

var (
	DB_DSN     string
	C_USERNAME string
)

func GetPassword(psw string, salt []byte) string {
	dk, _ := scrypt.Key([]byte(psw), salt, 16384, 8, 1, 32)
	return string(dk)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	flag.StringVar(&C_USERNAME, "username", "", "Username to generate token for")
	var DB_HOST = flag.String("host", "127.0.0.1", "Database host")
	var DB_PORT = flag.String("port", "3306", "Database port")
	var DB_NAME = flag.String("db", "auth", "Database name")
	var DB_USER = flag.String("user", "auth", "Database user")
	var DB_PASS = flag.String("pass", "auth", "Database user")
	flag.Parse()
	DB_DSN = *DB_USER + ":" + *DB_PASS + "@(" + *DB_HOST + ":" + *DB_PORT + ")/" + *DB_NAME + "?charset=utf8"
	//C_USERNAME = *USERNAME
	log.Printf("DB DSN: %s:****@(%s:%s)/%s?charset=utf8", *DB_USER, *DB_HOST, *DB_PORT, *DB_NAME)

}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	db, err := sql.Open("mysql", DB_DSN)
	checkErr(db.Ping())
	S_TOKEN := randToken(32)
	S_SALT := randToken(8)
	PW := GetPassword(S_TOKEN, []byte(S_SALT))
	log.Printf("Username: %s", C_USERNAME)
	log.Printf("Token: %s$%s", S_SALT, S_TOKEN)
	//log.Printf("Base16: %x", PW)
	_, err = db.Exec("UPDATE `users` SET token = ? where username = ? LIMIT 1", PW, C_USERNAME)
	db.Close()
	checkErr(err)

}
