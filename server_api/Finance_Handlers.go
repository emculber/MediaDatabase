package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func NewIncome(w http.ResponseWriter, r *http.Request) {
	income := Income{}

	r.ParseForm()
	income.UserKeys.Key = r.PostFormValue("key")
	income.Date = r.PostFormValue("date")
	income.Amount, _ = strconv.ParseFloat(r.PostFormValue("amount"), 64)
	income.Wallet.Id, _ = strconv.Atoi(r.PostFormValue("wallet_id"))
	income.Note = r.PostFormValue("note")

	fmt.Println(income)

	if err := income.OK(); err != nil {
		fmt.Println(err)
	}

	if err := income.UserKeys.validate(); err != nil {
		fmt.Println(err)
	} else {
		income.Wallet.UserKeys = income.UserKeys
	}

	if err := income.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
	}

	if err := income.Wallet.getWallet(); err != nil {
		fmt.Println(err)
	}

	if err := income.RegisterNewIncome(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(income)

	//TODO: Log Transaction
	income.SplitMoney()

	w.Write([]byte("OK"))
}

func NewWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Wallet{}

	r.ParseForm()
	wallet.UserKeys.Key = r.PostFormValue("key")
	wallet.Name = r.PostFormValue("name")
	wallet.Percent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)

	if err := wallet.OK(); err != nil {
		fmt.Println(err)
	}

	//TODO: Check If Percent is avalible

	w.Write([]byte("OK"))
}
