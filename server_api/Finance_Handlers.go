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
		return
	}

	if err := transaction.UserKeys.validate(); err != nil {
		fmt.Println(err)
		return
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
		return
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		fmt.Println(err)
		return
	}

	if err := transaction.RegisterNewIncome(); err != nil {
		fmt.Println(err)
		return
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
		return
	}

	fmt.Println("User Validated")

	transactions := userKeys.getIncomeList()

	fmt.Println("transaction Grabbed")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		fmt.Println(err)
		return
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
		return
	}

	if err := transaction.UserKeys.validate(); err != nil {
		fmt.Println(err)
		return
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
		return
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		fmt.Println(err)
		return
	}

	if err := transaction.RegisterNewExpense(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(transaction)

	//TODO: Log Transaction
	transaction.TakeFromWallet()

	w.Write([]byte("OK"))
}

func NewWallet(w http.ResponseWriter, r *http.Request) {
	//TODO: Add a Requested Percent for a wallet
	wallet := Wallet{}

	r.ParseForm()
	wallet.UserKeys.Key = r.PostFormValue("key")
	wallet.Name = r.PostFormValue("name")
	wallet.RequestedPercent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)
	wallet.Percent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)

	if err := wallet.OK(); err != nil {
		fmt.Println(err)
		return
	}

	if err := wallet.UserKeys.validate(); err != nil {
		fmt.Println(err)
		return
	}

	if err := wallet.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
		return
	}

	wallets := wallet.UserKeys.getWalletList()

	for _, unallocatedWallet := range wallets {
		if unallocatedWallet.Name == "unallocated" {
			if unallocatedWallet.Percent < wallet.Percent {
				wallet.Percent = unallocatedWallet.Percent
				unallocatedWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					fmt.Println(err)
				} else {
					if err := unallocatedWallet.updateWallet(); err != nil {
						fmt.Println(err)
					}
				}
			} else {
				unallocatedWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					fmt.Println(err)
				} else {
					if err := unallocatedWallet.updateWallet(); err != nil {
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
