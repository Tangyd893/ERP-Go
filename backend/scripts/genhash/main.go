package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func main() {
	h, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcryptCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(h))
}
