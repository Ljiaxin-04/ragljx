package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run gen_password.go <password>")
		fmt.Println("Example: go run gen_password.go 123456")
		os.Exit(1)
	}

	password := os.Args[1]

	// 生成 bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Bcrypt Hash: %s\n", string(hash))

	// 验证
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verification: OK")
}

