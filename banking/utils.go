package banking

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	passwordbyte := []byte(password)
	//cost :algo used to generate hash
	hash, err := bcrypt.GenerateFromPassword(passwordbyte, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating Hash:%v", err)
	}
	return string(hash)
}

func GenerateAccountNumber() uuid.UUID {
	accountNumber := uuid.New()
	return accountNumber
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
		return nil, ErrOTP
	}
	return &otp, nil
}

var (
	ErrAccountNotFound     = errors.New("Account not found")
	ErrPINLenght           = errors.New("PIN should consist be 4 characters")
	ErrEmail               = errors.New("unable to Update email")
	ErrAddress             = errors.New("unable to update customer's address")
	ErrPhoneNumber         = errors.New("error updating customer's phone number")
	ErrOTP                 = errors.New("error generating OTP")
	ErrFetchingUser        = errors.New("error recipient fetching recipient")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func SendNotification(transaction string) string {
	switch transaction {
	case "withdrawal":
		return "Withdrawal succesful"

	case "transfer":
		return "Transfer Succesful"

	case "deposit":
		return "Deposit Succeful"
	case "deleteUserAccount":
		return "Account Deleted Succesfully"
	default:
		return "Unable to perform Operation"
	}
}
