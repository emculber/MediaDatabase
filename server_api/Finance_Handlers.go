package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func NewIncome(w http.ResponseWriter, r *http.Request) {
	income := Income{}

	r.ParseForm()
	income.UserRole.Key = r.PostFormValue("key")
	income.Date = r.PostFormValue("date")
	income.Amount, _ = strconv.ParseFloat(r.PostFormValue("amount"), 64)
	income.Wallet.Name = r.PostFormValue("wallet")
	income.Note = r.PostFormValue("note")

	if err := income.OK(); err != nil {
		fmt.Println(err)
	}

	//TODO: Log Transaction
	//TODO: Split Amount

	w.Write([]byte("OK"))
}

func NewWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Wallet{}

	r.ParseForm()
	wallet.UserRole.Key = r.PostFormValue("key")
	wallet.Name = r.PostFormValue("name")
	wallet.Percent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)

	if err := wallet.OK(); err != nil {
		fmt.Println(err)
	}

	//TODO: Check If Percent is avalible

	w.Write([]byte("OK"))
}
