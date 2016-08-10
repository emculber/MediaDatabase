package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func NewIncome(w http.ResponseWriter, r *http.Request) {
	transaction := Transaction{}

	r.ParseForm()
	transaction.UserKeys.Key = r.PostFormValue("key")
	transaction.Date = r.PostFormValue("date")
	transaction.Amount, _ = strconv.ParseFloat(r.PostFormValue("amount"), 64)
	transaction.Wallet.Id, _ = strconv.Atoi(r.PostFormValue("wallet_id"))
	transaction.Note = r.PostFormValue("note")

	fmt.Println(transaction)

	if err := transaction.OK(); err != nil {
		fmt.Println(err)
	}

	if err := transaction.UserKeys.validate(); err != nil {
		fmt.Println(err)
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		fmt.Println(err)
	}

	if err := transaction.RegisterNewIncome(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(transaction)

	//TODO: Log Transaction
	transaction.SplitMoney()

	w.Write([]byte("OK"))
}

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("User Validated")

	transactions := userKeys.getIncomeList()

	fmt.Println("transaction Grabbed")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		fmt.Println(err)
	}

	fmt.Println("transaction Sent")
}

func NewExpense(w http.ResponseWriter, r *http.Request) {
	transaction := Transaction{}

	r.ParseForm()
	transaction.UserKeys.Key = r.PostFormValue("key")
	transaction.Date = r.PostFormValue("date")
	transaction.Amount, _ = strconv.ParseFloat(r.PostFormValue("amount"), 64)
	transaction.Wallet.Id, _ = strconv.Atoi(r.PostFormValue("wallet_id"))
	transaction.Note = r.PostFormValue("note")

	fmt.Println(transaction)

	if err := transaction.OK(); err != nil {
		fmt.Println(err)
	}

	if err := transaction.UserKeys.validate(); err != nil {
		fmt.Println(err)
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		fmt.Println(err)
	}

	if err := transaction.RegisterNewExpense(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(transaction)

	//TODO: Log Transaction
	transaction.TakeFromWallet()

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

	if err := wallet.UserKeys.validate(); err != nil {
		fmt.Println(err)
	}

	if err := wallet.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
	}

	wallets := wallet.UserKeys.getWalletList()

	for _, uWallet := range wallets {
		if uWallet.Name == "unallocated" {
			if uWallet.Percent < wallet.Percent {
				wallet.Percent = uWallet.Percent
				uWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					fmt.Println(err)
				} else {
					if err := uWallet.updateWallet(); err != nil {
						fmt.Println(err)
					}
				}
			} else {
				uWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					fmt.Println(err)
				} else {
					if err := uWallet.updateWallet(); err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	w.Write([]byte("OK"))
}

func GetWallets(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		fmt.Println(err)
	}

	wallets := userKeys.getWalletList()

	if err := json.NewEncoder(w).Encode(wallets); err != nil {
		fmt.Println(err)
	}
}
