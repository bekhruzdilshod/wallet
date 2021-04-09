package main

import (
	"fmt"

	"github.com/bekhruzdilshod/wallet/pkg/wallet")



func main() {

	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992900801441")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}

	}

	fmt.Println(account.Balance)

}