package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "kubernetes-authserver %s\n", VERSION)
}

func Auth(w http.ResponseWriter, r *http.Request) {

	var token Auth_token
	var r_id int
	var r_groups []string
	var r_username string
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Printf("Got invalid TokenReview (Json decode failed: %s)", err)
		reqInvalid(w, r)
	} else {
		err := db.Ping()
		if err != nil {
			log.Printf("DB Error: %s", err)
			invalidLogin(w, r)
		} else {
			s := strings.Split(token.Spec.Token, "$")
			if len(s) == 2 {
				theSalt, actuallToken := s[0], s[1]
				TokenToCheck := GetPassword(actuallToken, []byte(theSalt))
				if OPT_DEBUG {
					log.Printf("Received auth request:\n\t\t\tSalt: %s\n\t\t\tToken: %s\n\t\t\tBase16 of Scrypt hash: %x", theSalt, actuallToken, TokenToCheck)
				}
				rows, err := db.Query("SELECT users.id, users.username, groups.groupname FROM (auth.groups_mapping groups_mapping INNER JOIN auth.groups groups ON (groups_mapping.groupid = groups.id)) INNER JOIN auth.users users ON (groups_mapping.userid = users.id) WHERE BINARY (users.token = ?)", TokenToCheck)
				checkErr(err)
				defer rows.Close()
				if err == nil {
					count := 0
					for rows.Next() {	
						var uid int
						var username string
						var groupname string
						err = rows.Scan(&uid, &username, &groupname)
						checkErr(err)
						if OPT_DEBUG {
							log.Printf("Uid: %d Username: %s Groupname: %s", uid, username, groupname)
						}
						r_id = uid
						r_username = username
						r_groups = append(r_groups, groupname)
						count += 1
					}
					if OPT_DEBUG { log.Printf("Groups found for user %s: %d", r_username, count) }
					if count > 0 {
						log.Printf("Validated token for user %s", r_username)
						loginSuccess(w, r, r_id, r_username, r_groups)
					} else {
						log.Println("Invalid token received")
						invalidLogin(w, r)
					}
				} else {
					log.Printf("DB Error: %s", err)
					invalidLogin(w, r)
				}

			} else {
				log.Println("Invalid token received (cannot split)")
				invalidLogin(w, r)
			}
		}
	}
}

/*
func Authz(w http.ResponseWriter, r *http.Request) {
	var token SubjectAccessReview
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		reason := "Got invalid SubjectAccessReview (Json decode fail)"
		AuthzFailed(w, r, reason)
		log.Println(reason, err)
	} else {
		if DEBUG {
			log.Printf("POST DEBUG: %s", body)
		}

		if token.Spec.ResourceAttributes != nil {
			// var r_Namespace string
			// var r_Verb string
			// var r_Resource string
			/*
				err := db.QueryRow("SELECT namespace, verb, resource FROM `user_rules` WHERE username = ? AND (namespace = ? OR namespace = '*') AND (verb = ? OR verb = '*') LIMIT 1",
					token.Spec.User,
					token.Spec.ResourceAttributes.Namespace,
					token.Spec.ResourceAttributes.Verb).Scan(&r_Namespace, &r_Verb, &r_Resource)
*/
/*
			rows, err := db.Query("SELECT users.username, user_permissions.namespace, user_permissions.verb, user_permissions.resource, user_permissions.path FROM auth.user_permissions user_permissions INNER JOIN auth.users users ON (user_permissions.userid = users.id) WHERE (users.username = ?) AND (user_permissions.verb = ? OR user_permissions.verb = '*') AND (user_permissions.namespace = ? OR user_permissions.namespace = '*')", token.Spec.User, token.Spec.ResourceAttributes.Verb, token.Spec.ResourceAttributes.Namespace)
			defer rows.Close()
			if err == nil {
				count := 0
				for rows.Next() {
					var r_username string
					var r_namespace string
					var r_verb string
					var r_resource *string
					var r_path *string
					err = rows.Scan(&r_username, &r_namespace, &r_verb, &r_resource, &r_path)
					checkErr(err)
					fmt.Println(r_username)
					fmt.Println(r_namespace)
					fmt.Println(r_verb)
					fmt.Println(r_resource)
					fmt.Println(r_path)
					count += 1
				}
				if count > 0 {
					AuthzOK(w, r)
				} else {
					AuthzFailed(w, r, "cannot validate request")
				}
			} else {
				log.Println(err)
				AuthzFailed(w, r, "cannot validate request")
			}

		} else if token.Spec.NonResourceAttributes != nil {
			fmt.Println("it's a NonResourceAttributes")

		} else {
			AuthzFailed(w, r, "error, non valid request")
		}

		//		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		//		w.WriteHeader(http.StatusOK)
		//		fmt.Fprintf(w, "%+v", token)
		//		fmt.Printf("%s\n", token.Spec.ResourceAttributes.Namespace)
		//		fmt.Printf("%s\n", token.Spec.User)
	}

}
*/

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "418 I'm a Teapot\n")
}
