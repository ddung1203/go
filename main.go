package main

import (
	"fmt"
	"log"

	accounts "github.com/ddung1203/go/accounts"
)

func main() {
	account := accounts.NewAccount("Joongseok")
	account.Deposit(10)
	err := account.Withdraw(5)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(account.String())
}