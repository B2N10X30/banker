package banking

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	FirstName, LastName, Email, Address, PhoneNumber string
}

type Account struct {
	User
	IsRegistered  bool
	PIN           string
	AccountNumber uuid.UUID
	CreatedAt     time.Time
	Balance       float64
	Notificaition string
}

type Bank struct {
	Customers []Account
}

func NewAccount(firstName, lastName, email, address, phoneNumber, pin string) (*Account, error) {
	const PINLenght = 4
	if len(pin) < PINLenght || len(pin) > PINLenght {
		return nil, fmt.Errorf("PIN should consist be 4 characters")
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
		AccountNumber: GenerateAccountNumber(),
		CreatedAt:     time.Now(),
		Balance:       0.0,
	}
	return &newUserAccount, nil

}

func (a *Account) GetUserByAccountNumber(accountNumber uuid.UUID) (*Account, error) {
	return &Account{}, nil
}

func (a *Account) ChangePIN(pin string) error {
	const PINLenght = 4
	if len(pin) < PINLenght || len(pin) > PINLenght {
		return ErrPINLenght
	}
	a.PIN = pin
	return nil
}

func (a *Account) ChangeEmail(email string) error {
	a.Email = email
	return ErrEmail
}

func (a *Account) ChangeAddress(address string) error {
	a.Address = address
	return ErrAddress
}

func (a *Account) Deposit(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid deposit amount: %f", amount)
	}
	a.Balance += amount
	a.Notificaition = SendNotification("deposit")
	return a.Notificaition, nil
}

func (a *Account) Withdraw(amount float64) (string, error) {
	if amount >= a.Balance {
		return "", ErrInsufficientBalance
	}
	if amount <= 0 {
		return "", fmt.Errorf("invalid withdrawal Amount: %f", amount)
	}
	a.Balance -= amount
	a.Notificaition = SendNotification("withdrawal")
	return a.Notificaition, nil
}

func (a *Account) Transfer(amount float64, recipientAccountNumber uuid.UUID) (string, error) {
	// success := fmt.Sprintf("Transfer of %.2f to %v was successful", amount, recipientAccountNumber)
	if amount <= 0 {
		return "", fmt.Errorf("invalid Transfer amount: %f", amount)
	}
	if amount >= a.Balance {
		return "", ErrInsufficientBalance
	}

	recipient, err := a.GetUserByAccountNumber(recipientAccountNumber)
	if err != nil {
		return "", ErrFetchingUser
	}
	a.Balance -= amount
	recipient.Balance += amount
	a.Notificaition = SendNotification("transfer")
	return a.Notificaition, nil
}

func (b *Bank) RegisterNewUser(firstName, lastName, email, address, phoneNumber, pin string) (*Bank, error) {
	newUser, err := NewAccount(firstName, lastName, email, address, phoneNumber, pin)
	if err != nil {
		log.Fatal(err)
	}
	b.Customers = append(b.Customers, *newUser)
	return b, nil
}

// return only error here
func (b *Bank) DeleteUserAccount(accountNumber uuid.UUID) (string, error) {
	for i, customer := range b.Customers {
		if customer.AccountNumber == accountNumber {
			b.Customers = append(b.Customers[:i], b.Customers[i+1:]...)
			log.Printf("Account deleted: %+v", customer)
			return SendNotification("deleteUserAccount"), nil
		}
	}
	log.Printf("Account deleted: AccountNumber=%s", accountNumber)
	// If no account is not found
	return "", ErrAccountNotFound
}
