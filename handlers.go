package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"io"
	"io/ioutil"
	"net/http"

//	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "kubernetes-authserver %s\n", VERSION)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var token Auth_token
	var r_id int
	var r_token string
	var r_username string
	var r_uid string
	var r_groups string

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		reqInvalid(w, r)
		log.Println("Got invalid TokenReview (Json decode fail)")
	} else {
		err := db.Ping()
		if err != nil {
			log.Println("Cannot reach database server!", err)
		} else {
			err := db.QueryRow("SELECT id, token, username, uid, groups FROM `users` where token = ? LIMIT 1", token.Spec.Token).Scan(&r_id, &r_token, &r_username, &r_uid, &r_groups)	
			if err == nil {	
				log.Printf("Validated token for %s", r_username)
				loginSuccess(w, r, r_id, r_username, r_uid, strings.Split(r_groups, ",")) 
			} else {
				log.Println("Invalid token received")
				invalidLogin(w, r)
			}
		}
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "418 I'm a Teapot")
}
