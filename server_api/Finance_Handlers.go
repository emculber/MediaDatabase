package main

import (
	"encoding/json"
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

	//TODO: Remove wallet id from income. unneeded
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

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("User Validated")

	incomes := userKeys.getIncomeList()

	fmt.Println("income Grabbed")

	if err := json.NewEncoder(w).Encode(incomes); err != nil {
		fmt.Println(err)
	}

	fmt.Println("income Sent")
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
