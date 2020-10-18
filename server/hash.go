package main

import (
	"crypto/sha256"
)

//HashPass hashes a plaintext password using sha256
func HashPass(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return string(hash.Sum(nil))
}

//VerifyPass compares a plain text passord to a hashed password using sha256
func VerifyPass(passPlain string, passHash string) bool {
	hashed := HashPass(passPlain)
	if hashed == passHash {
		return true
	}
	return false
}
