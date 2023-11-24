package banking

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
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
	ErrGeneratingReceipt   = errors.New("error Generating receipt")
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

func GeneratePDFReciept(transaction string) *os.File {
	fileName := transaction + ".pdf"
	filePath, err := os.Create(fileName)
	if err != nil {
		log.Printf("%v", err)
	}
	filePermission := 0444

	err = os.Chmod(filePath.Name(), fs.FileMode(filePermission))
	if err != nil {
		log.Println(err)
	}
	return filePath
}

func WriteReceipt(account *Account, transactionType string, amount float64) (*os.File, error) {
	file := GeneratePDFReciept(transactionType)
	Notification := fmt.Sprintf("\nTransaction Receipt\n\tTimestamp: %s\n\tTransaction Type: %s\n\tAmount:%.2f\n\tCleared Balance: %.2f\n", time.Now().Format("09-10-20 15:30:00"), transactionType, amount, account.Balance)
	_, err := file.Write([]byte(Notification))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	return file, nil
}
