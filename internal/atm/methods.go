package atm

import (
	"errors"
	"fmt"

	"github.com/fuadop/altschool-go-task-3/internal/helpers"
)

func (u *User) VerifyPIN(pin string) bool {
	pinHash := helpers.Hash(pin)
	return pinHash == u.PINHash
}

func (u *User) ChangePIN(newPin string) (bool, error) {
	if len(newPin) != 4 {
		return false, errors.New("PIN must be four characters")
	}

	pinHash := helpers.Hash(newPin)
	u.PINHash = pinHash
	return true, nil
}

func (u *User) GetFormattedBalance() string {
	return fmt.Sprintf("₦%f", u.Balance)
}

func (u *User) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("invalid amount")
	}
	u.Balance += amount
	return nil
}

func (u *User) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("invalid amount")
	}

	if amount > u.Balance {
		return errors.New("insufficient balance")
	}
	u.Balance -= amount
	return nil
}
