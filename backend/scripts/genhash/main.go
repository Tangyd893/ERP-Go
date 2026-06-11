package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func main() {
	password := os.Getenv("GENHASH_PASSWORD")
	if password == "" {
		fmt.Fprintln(os.Stderr, "Usage: set GENHASH_PASSWORD=xxx && go run .")
		os.Exit(1)
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(h))
}
