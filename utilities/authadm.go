package main

import (
	"crypto/rand"
	"database/sql"
	//	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
	"log"
	"os"
)

var (
	DB_DSN     string
	C_USERNAME string
	DB_HOST    = os.Getenv("DB_HOST")
	DB_NAME    = os.Getenv("DB_NAME")
	DB_PASS    = os.Getenv("DB_PASS")
	DB_PORT    = os.Getenv("DB_PORT")
	DB_USER    = os.Getenv("DB_USER")
	//USERNAME   = flag.String("username", "", "Username to generate token for")
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
	// flag.Parse()
	if DB_HOST == "" {
		log.Fatal("Empty ENV: DB_HOST")
	}

	if DB_PORT == "" {
		log.Fatal("Empty ENV: DB_PORT")
	}

	if DB_NAME == "" {
		log.Fatal("Empty ENV: DB_NAME")
	}

	if DB_USER == "" {
		log.Fatal("Empty ENV: DB_USER")
	}

	if DB_PASS == "" {
		log.Fatal("Empty ENV: DB_PASS")
	}

	/*if *USERNAME == "" {
		log.Fatal("You must speccify a username with --username")
	}*/

	DB_DSN = DB_USER + ":" + DB_PASS + "@(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8"
	// log.Printf("DB DSN: %s:****@(%s:%s)/%s?charset=utf8", DB_USER, DB_HOST, DB_PORT, DB_NAME)

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  add")
		fmt.Println("    user <username>")
		fmt.Println("    group <groupname>\n")
		fmt.Println("  del")
		fmt.Println("    user <username>")
		fmt.Println("    group <groupname>\n")
		fmt.Println("  group")
		fmt.Println("    add <username> <groupname>")
		fmt.Println("    del <username> <groupname>\n")
		fmt.Println("  list <users/groups>")
		fmt.Println("    Lists users or groups in db\n")
		fmt.Println("  token <username>")
		fmt.Println("    Generates a token for given username\n")
		log.Fatal("No option given")
	}
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getuserId(db *sql.DB, Username string) int {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=?", Username).Scan(&userID)
	if err != nil {
		log.Fatal(err.Error())
	}
	return userID
}

func getgroupId(db *sql.DB, Groupname string) int {
	var groupID int
	err := db.QueryRow("SELECT id FROM groups WHERE groupname=?", Groupname).Scan(&groupID)
	if err != nil {
		log.Fatal(err.Error())
	}
	return groupID
}

func main() {
	db, err := sql.Open("mysql", DB_DSN)
	checkErr(db.Ping())

	switch action := os.Args[1]; action {

	case "add": //We want to add a user or group
		if len(os.Args) < 4 {
			log.Fatal("Missing parameters for \"add\"")
		}
		switch cmdgroup := os.Args[2]; cmdgroup {
		case "user":

			res, err := db.Exec("INSERT INTO `users` (token, username) VALUES ('', ?)", os.Args[3])
			if err != nil {
				log.Fatal(err.Error())
			} else {
				id, err := res.LastInsertId()
				if err != nil {
					println("Error:", err.Error())
				} else {
					log.Printf("Created user \"%s\" with id: %d\n", os.Args[3], id)
				}
			}
		case "group":
			res, err := db.Exec("INSERT INTO `groups` (groupname) VALUES (?)", os.Args[3])
			if err != nil {
				log.Fatal(err.Error())
			} else {
				id, err := res.LastInsertId()
				if err != nil {
					println("Error:", err.Error())
				} else {
					log.Printf("Created group \"%s\" with id: %d\n", os.Args[3], id)
				}
			}
		case "default":
			log.Println("Not supported, Allowed: user, group")
		}

	case "del": //We want to delete a user or group
		if len(os.Args) < 4 {
			log.Fatal("Missing parameters for \"del\"")
		}
		switch cmdgroup := os.Args[2]; cmdgroup {
		case "user":
			_, err := db.Exec("DELETE FROM `groups_mapping` where userid IN (SELECT users.id FROM `users` WHERE username = ?)", os.Args[3])
			checkErr(err)
			_, err = db.Exec("DELETE FROM `users` where username = ? LIMIT 1", os.Args[3])
			checkErr(err)
			log.Printf("Deleted user \"%s\"\n", os.Args[3])
		case "group":
			_, err := db.Exec("DELETE FROM `groups_mapping` where groupid IN (SELECT groups.id FROM `groups` WHERE groupname = ?)", os.Args[3])
			checkErr(err)
			_, err = db.Exec("DELETE FROM `groups` where groupname = ? LIMIT 1", os.Args[3])
			checkErr(err)
			log.Printf("Deleted group \"%s\"\n", os.Args[3])
		case "default":
			log.Println("Not supported, Allowed: user, group")
		}

	case "group": //Commands to add or remove users from groups
		if len(os.Args) < 3 {
			log.Fatal("Missing action, Valid: add,del")
		}
		if len(os.Args) == 5 {
			switch grpcmd := os.Args[2]; grpcmd {
			case "add":
				userID := getuserId(db, os.Args[3])
				groupID := getgroupId(db, os.Args[4])
				_, err := db.Exec("INSERT INTO `groups_mapping` (userid, groupid) VALUES (?, ?)", userID, groupID)
				if err == nil {
					log.Printf("Added \"%s\" to \"%s\"", os.Args[3], os.Args[4])
				} else {
					log.Fatal(err.Error())
				}
			case "del":
				userID := getuserId(db, os.Args[3])
				groupID := getgroupId(db, os.Args[4])
				_, err := db.Exec("DELETE FROM `groups_mapping` WHERE userid=? AND groupid = ?", userID, groupID)
				if err == nil {
					log.Printf("Deleted \"%s\" to \"%s\"", os.Args[3], os.Args[4])
				} else {
					log.Fatal(err.Error())
				}
			default:
				log.Fatal("Unsupported")
			}
		} else {
			log.Fatal("Usage: group add <user> <group>")
		}

	case "list":
		if len(os.Args) == 2 {
			rows, err := db.Query("SELECT users.id, users.username, GROUP_CONCAT(groups.groupname SEPARATOR ',') FROM (auth.groups_mapping groups_mapping INNER JOIN auth.groups groups ON (groups_mapping.groupid = groups.id)) RIGHT JOIN auth.users users ON (groups_mapping.userid = users.id) GROUP BY users.id")
			checkErr(err)
			defer rows.Close()
			if err == nil {
				for rows.Next() {
					var uid int
					var username string
					var groupname string
					err = rows.Scan(&uid, &username, &groupname)
					if groupname == "" {
						groupname = "<none>"
					}
					fmt.Printf("%d:%s:%s\n", uid, username, groupname)
				}
			} else {
				log.Printf("DB Error: %s", err)
			}
		}

		if len(os.Args) == 3 {
			switch lscmd := os.Args[2]; lscmd {
			case "users":
				rows, err := db.Query("SELECT users.id, users.username FROM `users` ORDER BY users.username")
				checkErr(err)
				defer rows.Close()
				if err == nil {
					for rows.Next() {
						var uid int
						var username string
						err = rows.Scan(&uid, &username)
						fmt.Printf("%d:%s\n", uid, username)
					}
				} else {
					log.Printf("DB Error: %s", err)
				}
			case "groups":
				rows, err := db.Query("SELECT groups.id, groups.groupname FROM `groups` ORDER BY groups.groupname")
				checkErr(err)
				defer rows.Close()
				if err == nil {
					for rows.Next() {
						var uid int
						var groupname string
						err = rows.Scan(&uid, &groupname)
						fmt.Printf("%d:%s\n", uid, groupname)
					}
				} else {
					log.Printf("DB Error: %s", err)
				}
			default:
				log.Println("Suppported values are: users, groups")
				log.Fatalf("Unsupported option \"%s\"\n", os.Args[2])
			}
		}

	case "token": //Create a new token for a user
		if len(os.Args) < 3 {
			log.Fatal("Missing username")
		}
		S_TOKEN := randToken(32)
		S_SALT := randToken(8)
		PW := GetPassword(S_TOKEN, []byte(S_SALT))
		fmt.Printf("%s$%s\n", S_SALT, S_TOKEN)
		_, err = db.Exec("UPDATE `users` SET token = ? where username = ? LIMIT 1", PW, os.Args[2])
		checkErr(err)

	default:
		log.Println("Supported values are: add, del, list, token")
		log.Fatalf("Unsupported option \"%s\"\n", os.Args[1])
	}
	//os.Args[0]
	db.Close()
}
