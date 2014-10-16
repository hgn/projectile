package main

import "encoding/base64"
import "crypto/rand"
import "fmt"
import "os"
import "io/ioutil"
import "errors"
import "code.google.com/p/go.crypto/bcrypt"

const saltDbFilePath string = "db/salt.db"

func initializeSaltDb() error {
	size := 128

	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	f, err := os.Create(saltDbFilePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	n3, err := f.WriteString(rs)
	fmt.Printf("wrote %d bytes\n", n3)
	f.Sync()

	return nil
}

func ReadSaltDB() ([]byte, error) {
	if _, err := os.Stat(saltDbFilePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Password salt DB not available - BUG")
			return []byte(" "), errors.New("only add supported")
		}
	}

	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	return b, nil
}

func CheckPassword(plain, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plain)
	if err == nil {
		return true
	}
	return false
}

func CryptPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), err
	}

	return hashedPassword, nil
}

func InitialzeCryptSystem() error {

	if _, err := os.Stat(saltDbFilePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Password salt DB not available - generate now")
			err = initializeSaltDb()
			if err != nil {
				fmt.Println("Failed to generate passwd salt file!")
				return err
			}
		}
	}

	CryptPassword([]byte("foooba"))

	return nil
}
