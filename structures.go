package main

type Json_decode_error struct {
	Status string `json:"status"`
	Details string `json:"details"`
}

type Auth_token struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       Token `json:"spec"`
}

type Token struct {
	Token string `json:"token"`
}

type Authenticated struct {
	Authenticated bool `json:"authenticated"`
}

type Status struct {
	Authenticated bool `json:"authenticated"`
	Userinfo `json:"user"`
}

type Userinfo struct {
	Username string   `json:"username"`
	UID      int   `json:"uid"`
	Groups   []string `json:"groups"`
}

type Auth_response_successfull struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Status	`json:"status"`
}


type Auth_unsuccessfull struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Status     Authenticated `json:"status"`
}
