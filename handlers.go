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

func reqInvalid(w http.ResponseWriter, r *http.Request) {
	response := &Json_decode_error{
		Status: "400",
		Details: "Invalid TokenReview ( Json decode failed )",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest) // unprocessable entity
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func invalidLogin(w http.ResponseWriter, r *http.Request) {
	response := Auth_unsuccessfull{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind: "TokenReview",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func loginSuccess(w http.ResponseWriter, r *http.Request, r_id int, r_username string, r_uid int, r_groups []string ) {
	response := &Auth_response_successfull{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind: "TokenReview",
		Status: Status{
			Authenticated: true,
			Userinfo: Userinfo{
				Groups: r_groups,
				UID: r_id,
				Username: r_username,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var token Auth_token
	var r_id int
	var r_token string
	var r_username string
	var r_uid int
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

