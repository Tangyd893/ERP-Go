package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func main() {
	plaintext := os.Getenv("GENHASH_PLAINTEXT")
	if plaintext == "" {
		fmt.Fprintln(os.Stderr, "Usage: GENHASH_PLAINTEXT=<secret> go run .")
		os.Exit(1)
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(h))
}
