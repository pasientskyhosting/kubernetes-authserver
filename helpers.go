package main

import "golang.org/x/crypto/scrypt"

//Basic is X in slice function
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetPassword(psw string, salt []byte) string {
	dk, _ := scrypt.Key([]byte(psw), salt, 16384, 8, 1, 32)
	return string(dk)
}
