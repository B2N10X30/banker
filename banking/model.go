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
}

type Bank struct {
	Customers []Account
}

func (a *Account) RegisterAccount(firstName, lastName, email, address, phoneNumber, pin string, account uint64) (*Account, error) {
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
		Balance:       000_000,
	}
	log.Printf("Account registered: AccountNumber=%s, FirstName=%s, LastName=%s, Email=%s",
		a.AccountNumber, a.FirstName, a.LastName, a.Email)
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

func (a *Account) Deposit(amount float64) (float64, error) {
	if amount <= 0 {
		return a.Balance, fmt.Errorf("invalid deposit amount: %f", amount)
	}
	a.Balance += amount
	fmt.Println(SendNotification("deposit"))
	return a.Balance, nil
}

func (a *Account) Withdraw(amount float64) (float64, error) {
	if amount <= 0 || amount >= a.Balance {
		return a.Balance, fmt.Errorf("invalid withdrawal amount: %f", amount)
	}
	a.Balance -= amount
	notificaition := SendNotification("withdrawal")
	fmt.Println(notificaition)
	return a.Balance, nil
}

func (a *Account) Transfer(amount float64, recipientAccountNumber uuid.UUID) (string, error) {
	// success := fmt.Sprintf("Transfer of %.2f to %v was successful", amount, recipientAccountNumber)
	senderBalance := a.Balance
	if amount <= 0 || amount >= senderBalance {
		return "", fmt.Errorf("invalid Transfer amount or Insufficient Balance: %f", amount)
	}
	senderBalance -= amount

	recipient, err := a.GetUserByAccountNumber(recipientAccountNumber)
	if err != nil {
		return "", ErrFetchingUser
	}
	recipient.Balance += amount

	return SendNotification("transfer"), nil
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
