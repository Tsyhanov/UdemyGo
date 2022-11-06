package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
	pass := "12345678"
	hashedPass, err := hashPassword(pass)
	if err != nil {
		panic(err)
	}

	err = comparePassword(pass, hashedPass)
	if err != nil {
		log.Fatalln("Not logged in")
	}

	log.Println("Logged in!")
}

func hashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while generating bcrypt hash")
	}

	return bs, nil
}

func comparePassword(password string, hashPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPass, []byte(password))
	if err != nil {
		return fmt.Errorf("Error while comparing password")
	}
	return nil
}
