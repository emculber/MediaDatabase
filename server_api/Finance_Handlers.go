package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func SetupFinance(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkPermissions("admin"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Permissions Access Denied")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("execute"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Role Access Denied")
		return
	}

	financeRolePermissions := RolePermissions{}
	financeRolePermissions.Role.Role = "Finance"
	financeRolePermissions.Permission.Permission = "Financal_RW"
	financeRolePermissions.access = 3

	if err := financeRolePermissions.Role.RegisterNewRole(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}
	if err := financeRolePermissions.Permission.RegisterNewPermissions(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}
	if err := financeRolePermissions.RegisterNewRolePermissions(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Role")
		return
	}

	w.Write([]byte("OK"))
}

func NewFinancalUser(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}
	newFinancalUser := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	newFinancalUser.RolePermissions.Id, _ = strconv.Atoi(r.PostFormValue("role_permissions_id"))

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkPermissions("general"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("execute"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	newFinancalUser.User = userKeys.User
	newFinancalUser.generateKey()

	log.WithFields(log.Fields{
		"Financal User": newFinancalUser,
	}).Info("New Financal User")

	if err := newFinancalUser.RegisterNewUserKeys(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New User Key")
		return
	}

	w.Write([]byte("OK"))
	/*
		userKeys.validate()

		wallet := Wallet{
			UserKeys:         userKeys,
			Name:             "unallocated",
			RequestedPercent: -1,
			Percent:          100,
			CurrentAmount:    0,
			WalletLimit:      -1,
		}
		if err := wallet.getUnallocatedWallet(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Getting Unallocated Wallet")
			if err := wallet.RegisterNewWallet(); err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error Registering Unallocated Wallet")
			}
		}
	*/
}

func NewIncome(w http.ResponseWriter, r *http.Request) {
	transaction := Transaction{}

	r.ParseForm()
	transaction.UserKeys.Key = r.PostFormValue("key")
	transaction.Date = r.PostFormValue("date")

	float_amount, _ := strconv.ParseFloat(r.PostFormValue("amount"), 64)
	transaction.Amount = int(float_amount * 100)

	transaction.Wallet.Id, _ = strconv.Atoi(r.PostFormValue("wallet_id"))
	transaction.Note = r.PostFormValue("note")

	log.WithFields(log.Fields{
		"Transaction": transaction,
	}).Info("New Income")

	if err := transaction.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error on transaction OK")
		return
	}

	if err := transaction.UserKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Getting Wallet")
		return
	}

	if err := transaction.RegisterNewIncome(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Registering New Income")
		return
	}

	//TODO: Log Transaction
	transaction.SplitMoney()

	w.Write([]byte("OK"))
}

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("read"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	transactions := userKeys.getIncomeList()

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}

func NewExpense(w http.ResponseWriter, r *http.Request) {
	transaction := Transaction{}

	r.ParseForm()
	transaction.UserKeys.Key = r.PostFormValue("key")
	transaction.Date = r.PostFormValue("date")

	float_amount, _ := strconv.ParseFloat(r.PostFormValue("amount"), 64)
	transaction.Amount = int(float_amount * 100)

	transaction.Wallet.Id, _ = strconv.Atoi(r.PostFormValue("wallet_id"))
	transaction.Note = r.PostFormValue("note")

	log.WithFields(log.Fields{
		"Transaction": transaction,
	}).Info("New Income")

	if err := transaction.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Transaction OK")
		return
	}

	if err := transaction.UserKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	} else {
		transaction.Wallet.UserKeys = transaction.UserKeys
	}

	if err := transaction.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	if err := transaction.Wallet.getWallet(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Getting Wallet")
		return
	}

	if err := transaction.RegisterNewExpense(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Register New Expense")
		return
	}

	//TODO: Log Transaction
	transaction.TakeFromWallet()

	w.Write([]byte("OK"))
}

func NewWallet(w http.ResponseWriter, r *http.Request) {
	wallet := Wallet{}

	r.ParseForm()
	wallet.UserKeys.Key = r.PostFormValue("key")
	wallet.Name = r.PostFormValue("name")
	wallet.RequestedPercent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)
	wallet.Percent, _ = strconv.ParseFloat(r.PostFormValue("percent"), 64)

	float_limit, _ := strconv.ParseFloat(r.PostFormValue("limit"), 64)
	wallet.WalletLimit = int(float_limit * 100)

	if err := wallet.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Transaction OK")
		return
	}

	if err := wallet.UserKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := wallet.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	wallets := wallet.UserKeys.getWalletList()

	//TODO: Update to sql unallocated Wallet
	for _, unallocatedWallet := range wallets {
		if unallocatedWallet.Name == "unallocated" {
			if unallocatedWallet.Percent < wallet.Percent {
				wallet.Percent = unallocatedWallet.Percent
				unallocatedWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					log.WithFields(log.Fields{
						"Error": err,
					}).Error("Error while registering new wallet")
					return
				} else {
					if err := unallocatedWallet.updateWallet(); err != nil {
						log.WithFields(log.Fields{
							"Error": err,
						}).Error("Error updating unallocated Wallet")
						return
					}
				}
			} else {
				unallocatedWallet.Percent -= wallet.Percent
				if err := wallet.RegisterNewWallet(); err != nil {
					log.WithFields(log.Fields{
						"Error": err,
					}).Error("Error while registering new wallet")
					return
				} else {
					if err := unallocatedWallet.updateWallet(); err != nil {
						log.WithFields(log.Fields{
							"Error": err,
						}).Error("Error updating unallocated Wallet")
						return
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
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error validating user key")
		return
	}

	if err := userKeys.RolePermissions.checkAccess("read"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Access Denied")
		return
	}

	wallets := userKeys.getWalletList()

	if err := json.NewEncoder(w).Encode(wallets); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}
