package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func reqInvalid(w http.ResponseWriter, r *http.Request) {
	response := &Json_decode_error{
		Status:  "400",
		Details: "Invalid TokenReview ( Json decode failed )",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest) // unprocessable entity
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func invalidLogin(w http.ResponseWriter, r *http.Request) {
	response := &Auth_unsuccessfull{
		APIVersion: APIVERSION,
		Kind:       "TokenReview",
		Status: &Authenticated{
			Authenticated: false,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func loginSuccess(w http.ResponseWriter, r *http.Request, r_id int, r_username string, r_groups []string) {
	response := &Auth_response_successfull{
		APIVersion: APIVERSION,
		Kind:       "TokenReview",
		Status: &AStatus{
			Authenticated: true,
			Userinfo: &Userinfo{
				Groups:   r_groups,
				UID:      strconv.Itoa(r_id),
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
