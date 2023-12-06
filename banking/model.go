package banking

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type User struct {
	FirstName, LastName, Email, Address, PhoneNumber string
	DateOfBirth                                      string
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

func NewAccount(firstName, lastName, email, address, phoneNumber, pin, dob string) (*Account, error) {
	const PINLenght = 4
	if len(pin) != PINLenght {
		return nil, fmt.Errorf("PIN should consist be 4 characters")
	}
	Date, err := ParseDate(dob)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date")
	}
	DOB := Date.Format("2006 January 02") //date format converts date to specified format

	newUserAccount := Account{
		User: User{
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			Address:     address,
			PhoneNumber: phoneNumber,
			DateOfBirth: DOB,
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
	if isEmpty(email) {
		return ErrEmail
	}
	a.Email = email
	return nil
}

func (a *Account) ChangeAddress(address string) error {
	if isEmpty(address) {
		return ErrAddress
	}
	a.Address = address
	return nil
}

func (a *Account) ChangeDOB(dob string) (string, error) {
	Date, err := ParseDate(dob)
	if err != nil {
		return "", fmt.Errorf("%v", ErrDOB)
	}
	DOB := Date.Format("2006 January 02")
	a.DateOfBirth = DOB
	return "", nil
}

func (a *Account) Deposit(amount float64) (string, *os.File, error) {
	transactionType := "Deposit"
	if amount <= 0 {
		return "", nil, fmt.Errorf("invalid deposit amount: %f", amount)
	}
	a.Balance += amount
	a.Notificaition = SendNotification("deposit")
	receipt, err := WriteReceipt(a, transactionType, amount)
	if err != nil {
		log.Printf("Error generating reciept: %v", err)
	}
	return a.Notificaition, receipt, nil
}

func (a *Account) Withdraw(amount float64) (string, *os.File, error) {
	transactionType := "Withdrawal"
	if amount >= a.Balance {
		return "", nil, ErrInsufficientBalance
	}
	if amount <= 0 {
		return "", nil, fmt.Errorf("invalid withdrawal Amount: %f", amount)
	}
	a.Balance -= amount
	a.Notificaition = SendNotification("withdrawal")
	receipt, err := WriteReceipt(a, transactionType, amount)
	if err != nil {
		log.Printf("Error generating reciept: %v", err)
	}
	return a.Notificaition, receipt, nil
}

func (a *Account) Transfer(amount float64, recipientAccountNumber uuid.UUID) (string, *os.File, error) {
	// success := fmt.Sprintf("Transfer of %.2f to %v was successful", amount, recipientAccountNumber)
	transactionType := "Transfer"
	if amount <= 0 {
		return "", nil, fmt.Errorf("invalid Transfer amount: %f", amount)
	}
	if amount >= a.Balance {
		return "", nil, ErrInsufficientBalance
	}

	recipient, err := a.GetUserByAccountNumber(recipientAccountNumber)
	if err != nil {
		return "", nil, ErrFetchingUser
	}
	a.Balance -= amount
	recipient.Balance += amount
	a.Notificaition = SendNotification("transfer")
	file, err := WriteReceipt(recipient, transactionType, amount)
	if err != nil {
		log.Printf("Error generatingn receipt: %v", err)
	}

	return a.Notificaition, file, nil
}

func (a *Account) SendEmailNotification() {
	Dailer := gomail.NewDialer("contact@bool.com", 25, "Admin", "^R0b0+@b00l")
	Mailer := gomail.NewMessage()
	//Messages
	Dailer.DialAndSend(Mailer)

}
func (b *Bank) RegisterNewUser(firstName, lastName, email, address, phoneNumber, pin, dob string) (*Bank, error) {
	newUser, err := NewAccount(firstName, lastName, email, address, phoneNumber, pin, dob)
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
