package banking

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName, LastName, Email, Address, PhoneNumber string
}

type Account struct {
	User
	IsRegistered  bool
	Password      string
	AccountNumber string
	CreatedAt     time.Time
}

type Bank struct{
	Customers []Account
}

func (a *Account) RegisterAccount(firstName, lastName, email, address, phoneNumber, password string, account uint64) (*Account, error) {
	return &Account{
			User: User{
				FirstName:   firstName,
				LastName:    lastName,
				Email:       email,
				Address:     address,
				PhoneNumber: phoneNumber,
			},
			IsRegistered:  true,
			Password:      HashPassword(password),
			AccountNumber: GenerateAccountNum(),
			CreatedAt:     time.Now(),
		},
		nil
}

func HashPassword(password string) string {
	passwordbyte := []byte(password)
	//cost :algo used to generate hash
	hash, err := bcrypt.GenerateFromPassword(passwordbyte, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generating Hash:%v", err)
	}
	return string(hash)
}

func GenerateAccountNum() string {
	//Use current time to create a new source everytime
	source := rand.NewSource(time.Now().UnixNano())

	//Random Number generator
	randomGenerator := rand.New(source)

	//Generate new 9-digits num
	randomNumber := randomGenerator.Intn(900000000) + 100000000

	accountNum := fmt.Sprintf("4%d", randomNumber)
	return accountNum

}

func(b *Bank) UpdateUserDetails(accounNumber string){
	for _, customer := range b.Customers{
		if customer.AccountNumber == 
	}
}
