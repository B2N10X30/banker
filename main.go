package main

import (
	"banker/banking"
	"fmt"
)

func main() {
	NewAccountNumber := banking.GenerateAccountNum()
	fmt.Printf("New account is : %s", NewAccountNumber)
}
