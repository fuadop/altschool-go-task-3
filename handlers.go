package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/fuadop/altschool-go-task-3/internal/atm"
	"github.com/fuadop/altschool-go-task-3/internal/helpers"
)

func opsPrompt(sig chan os.Signal) {
	if currentUser == nil {
		log.Fatal("System Failure due to some unknown reasons")
	}

	for {
		var command int
		fmt.Println("0. Exit\n1. Change PIN\n2. Check Balance\n3. Withdraw funds\n4. Deposit funds")
		fmt.Print("Make a Selection: ")
		if _, e := fmt.Scan(&command); e != nil {
			fmt.Println("Invalid Command")
			continue
		}

		switch command {
		case 0:
			sig <- syscall.SIGTERM
			<-time.After(100 * time.Second) // to prevent program exec before go-routine catches signal
		case 1:
			var newPin string
			message := "successful"
			fmt.Print("Enter New PIN: ")
			fmt.Scan(&newPin)
			ok, err := currentUser.ChangePIN(newPin)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if !ok {
				message = "failed"
			}
			fmt.Printf("PIN update %s\n", message)
		case 2:
			balance := currentUser.GetFormattedBalance()
			fmt.Printf("Your Balance is %s\n", balance)
		case 3:
			var amount float64
			fmt.Print("Enter amount to withdraw: ")
			_, err := fmt.Scan(&amount)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}
			err = currentUser.Withdraw(amount)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Printf("Your withdrawal was successful.\nYour Balance is: %s\n", currentUser.GetFormattedBalance())
		case 4:
			var amount float64
			fmt.Print("Enter amount to deposit: ")
			_, err := fmt.Scan(&amount)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}
			err = currentUser.Deposit(amount)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Printf("Your deposit was successful.\nYour Balance is: %s\n", currentUser.GetFormattedBalance())
		default:
			fmt.Println("Invalid Command")
			continue
		}

		var response string
		fmt.Print("Do you want to perform another operation? (Y/N): ")
		fmt.Scan(&response)

		if !strings.Contains(strings.ToLower(response), "y") {
			sig <- syscall.SIGTERM
			<-time.After(100 * time.Second) // to prevent program exec before go-routine catches signal
		}
	}
}

func authPrompt(sig chan os.Signal) {
	for {
		var command int
		fmt.Println("0. Exit\n1. Register\n2. Login")
		fmt.Print("Make a Selection: ")
		if _, e := fmt.Scan(&command); e != nil {
			fmt.Println("Invalid Command")
			continue
		}

		switch command {
		case 0:
			sig <- syscall.SIGTERM
			<-time.After(100 * time.Second) // to prevent program exec before go-routine catches signal
		case 1:
			registerUser()
		case 2:
			loginUser()
			return
		default:
			fmt.Println("Invalid Command")
			continue
		}
	}
}

func registerUser() {
	for {
		var username string
		var pin string

		fmt.Print("Enter username: ")
		fmt.Scan(&username)
		if len(username) < 1 {
			fmt.Println("Username should not be empty")
			continue
		}
		if _, e := users[username]; e {
			fmt.Println("User with username already exists")
			continue
		}

		fmt.Print("Enter PIN: ")
		fmt.Scan(&pin)
		if len(pin) < 4 || len(pin) > 4 {
			fmt.Println("PIN should be 4 characters")
			continue
		}

		users[username] = atm.User{
			Username: username,
			PINHash:  helpers.Hash(pin),
		}
		fmt.Println("Account registered, please login.")
		break
	}
}

func loginUser() {
	for {
		var user atm.User
		var username string
		var pin string

		fmt.Print("Enter username: ")
		fmt.Scan(&username)
		user, e := users[username]
		if !e {
			fmt.Println("User with the username not found")
			continue
		}

		fmt.Print("Enter PIN: ")
		fmt.Scan(&pin)
		if valid := user.VerifyPIN(pin); !valid {
			fmt.Println("PIN entered is incorrect")
			continue
		}

		// update the current user global variable
		currentUser = &user
		break
	}
}
