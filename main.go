package main

import (
	"banker/banking"
	"fmt"
	"log"
)

//16:33

func main() {

	user1, err := banking.NewAccount("Jon", "Doe", "jondoe@boolbank.com", "12 Main Str", "555-335-4055", "2030")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New Account:\n%+v\n", *user1)

	deposit, err := user1.Deposit(1000)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%s\n", deposit)

	withdraw, err := user1.Withdraw(350)
	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Printf("%s\n", withdraw)

	user2, err := banking.NewAccount("Mary", "Anne", "maryanne@boolbank.com", "12 Pkg Str", "555-533-1250", "1745")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New Account:\n%+v\n", *user2)

	transfer, err := user1.Transfer(300, user2.AccountNumber)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(transfer)
	fmt.Printf("New Account:\n%+v\n", *user2)
}
