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
	PIN           string
	AccountNumber string
	CreatedAt     time.Time
	Balance       uint64
}

type Bank struct {
	Customers []Account
}

func (a *Account) RegisterAccount(firstName, lastName, email, address, phoneNumber, pin string, account uint64) (*Account, error) {
	const PINLenght = 4
	if len(pin) < PINLenght || len(pin) > PINLenght {
		return nil, fmt.Errorf("password should consist be 4 characters")
	}
	newUserAccount := Account{
		User: User{
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			Address:     address,
			PhoneNumber: phoneNumber,
		},
		IsRegistered:  true,
		PIN:           HashPassword(pin),
		AccountNumber: GenerateAccountNum(),
		CreatedAt:     time.Now(),
		Balance:       000_000,
	}
	log.Printf("Account registered: AccountNumber=%s, FirstName=%s, LastName=%s, Email=%s",
		a.AccountNumber, a.FirstName, a.LastName, a.Email)
	return &newUserAccount, nil

}

func (a *Account) Deposit(amount uint64) (uint64, error) {
	if amount <= 0 {
		return a.Balance, fmt.Errorf("invalid deposit amount: %d", amount)
	}
	a.Balance += amount
	return a.Balance, nil
}

func (a *Account) Withdraw(amount uint64) (uint64, error) {
	if amount <= 0 || amount >= a.Balance {
		return a.Balance, fmt.Errorf("invalid withdrawal amount: %d", amount)
	}
	a.Balance -= amount
	return a.Balance, nil
}

func (a *Account) Transfer(amount uint64, accountNumber string) (uint64, error) {
	if amount <= 0 || amount >= a.Balance {
		return a.Balance, fmt.Errorf("invalid Transfer amount: %d", amount)
	}
	a.Balance -= amount
	return a.Balance, nil
}

func HashPassword(password string) string {
	passwordbyte := []byte(password)
	//cost :algo used to generate hash
	hash, err := bcrypt.GenerateFromPassword(passwordbyte, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating Hash:%v", err)
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

func GenerateOTP() (*string, error) {
	//otp lenght
	const otplenght = 6
	//Use current time to create a new source everytime
	source := rand.NewSource(time.Now().UnixNano())

	//Random Number generator
	randomGenerator := rand.New(source)
	otp := ""
	for i := 0; i < otplenght; i++ {
		otp += fmt.Sprintf("%d", randomGenerator.Intn(10))
	}

	if len(otp) != otplenght {
		return nil, fmt.Errorf("error generating otp: unexpected lenght")
	}
	return &otp, nil
}

func (b *Bank) UpdateUserDetails(accountNumber, firstName, lastName, email, address, phoneNumber, pin string) (*Account, error) {
	const PINLenght = 4
	if len(pin) < PINLenght || len(pin) > PINLenght {
		return nil, fmt.Errorf("password should consist be 4 characters")
	}
	for i, customer := range b.Customers {
		if customer.AccountNumber == accountNumber {
			//all changeable details should be changed
			b.Customers[i].FirstName = firstName
			b.Customers[i].LastName = lastName
			b.Customers[i].Email = email
			b.Customers[i].Address = address
			b.Customers[1].PhoneNumber = phoneNumber
			b.Customers[i].PIN = HashPassword(pin)
			//updated customer
			return &b.Customers[i], nil
		}
		log.Printf("User details updated:\n FirstName: %s,\nLastName: %s,\nEmail: %s,\nAddress: %s,\nPhoneNumber: %s", firstName, lastName, email, address, phoneNumber)
	}
	return nil, fmt.Errorf("customer account for account number: %s not found", accountNumber)
}

func (b *Bank) DeleteUserAccount(accountNumber string) (*[]Account, error) {
	for i, customer := range b.Customers {
		if customer.AccountNumber == accountNumber {
			b.Customers = append(b.Customers[:i], b.Customers[i+1:]...)
			log.Printf("Account deleted: %+v", customer)
			return &b.Customers, nil
		}
	}
	log.Printf("Account deleted: AccountNumber=%s", accountNumber)
	// If no account is not found
	return nil, fmt.Errorf("customer account for account number: %s not found", accountNumber)
}
