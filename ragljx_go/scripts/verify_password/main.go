package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run verify_password.go <password> <hash>")
		os.Exit(1)
	}

	password := os.Args[1]
	hash := os.Args[2]

	// 验证
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verification: OK")
	fmt.Printf("Password '%s' matches the hash\n", password)
}

