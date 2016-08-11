package main

import (
	"fmt"
	"strconv"

	"github.com/emculber/database_access/postgresql"
)

//TODO: Add in user_id to regerster rows to a single user

var financeDatabaseSchema = []string{
	"CREATE TABLE wallet(id serial primary key, user_id integer references registered_user(id), name varchar, requested_percent real, percent real, current_amount real, limit real)",
	"CREATE TABLE income(id serial primary key, user_id integer references registered_user(id), date varchar, amount real, wallet_id integer references wallet(id), note varchar)",
	"CREATE TABLE expense(id serial primary key, user_id integer references registered_user(id), date varchar, amount real, wallet_id integer references wallet(id), note varchar)",
}

func CreateFinanceTables() {
	//TODO: check if table exists
	for _, table := range financeDatabaseSchema {
		fmt.Println(table)
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (wallet *Wallet) RegisterNewWallet() error {
	err := db.QueryRow(`insert into wallet (user_id, name, requested_percent, percent, current_amount, "limit") values($1, $2, $3, $4, $5, $6) returning id`, wallet.UserKeys.User.Id, wallet.Name, wallet.RequestedPercent, wallet.Percent, wallet.CurrentAmount, wallet.Limit).Scan(&wallet.Id)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) getWallet() error {
	err := db.QueryRow("select wallet.id, wallet.name, wallet.requested_percent, wallet.percent, wallet.current_amount, wallet.limit from wallet where wallet.id = $1", wallet.Id).Scan(&wallet.Id, &wallet.Name, &wallet.RequestedPercent, &wallet.Percent, &wallet.CurrentAmount, &wallet.Limit)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) getUnallocatedWallet() error {
	err := db.QueryRow("select wallet.id, wallet.name, wallet.requested_percent, wallet.percent, wallet.current_amount, wallet.limit from wallet where wallet.name = $1", wallet.Name).Scan(&wallet.Id, &wallet.Name, &wallet.RequestedPercent, &wallet.Percent, &wallet.CurrentAmount, &wallet.Limit)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) updateWallet() error {
	err := db.QueryRow(`UPDATE wallet SET name = $1, requested_percent=$2, percent = $3, current_amount = $4, "limit"=$5 WHERE id = $6 returning id`, wallet.Name, wallet.RequestedPercent, wallet.Percent, wallet.CurrentAmount, wallet.Limit, wallet.Id).Scan(&wallet.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userKeys *UserKeys) getWalletList() []Wallet {
	statement := fmt.Sprintf("select wallet.id, wallet.name, wallet.requested_percent, wallet.percent, wallet.current_amount, wallet.limit from wallet where user_id=%d", userKeys.User.Id)
	//TODO: Error Checking
	wallets, _, err := postgresql_access.QueryDatabase(db, statement)
	if err != nil {
		fmt.Println("Error While getting wallet List ->", err)
	}
	wallet_list := []Wallet{}

	for _, wallet := range wallets {
		single_wallet := Wallet{}
		single_wallet.UserKeys = *userKeys
		single_wallet.Id, _ = strconv.Atoi(wallet[0].(string))
		single_wallet.Name = wallet[1].(string)
		single_wallet.RequestedPercent, _ = strconv.ParseFloat(wallet[2].(string), 64)
		single_wallet.Percent, _ = strconv.ParseFloat(wallet[3].(string), 64)
		single_wallet.CurrentAmount, _ = strconv.ParseFloat(wallet[4].(string), 64)
		single_wallet.Limit, _ = strconv.ParseFloat(wallet[5].(string), 64)
		wallet_list = append(wallet_list, single_wallet)
	}
	return wallet_list
}

func (transaction *Transaction) RegisterNewIncome() error {
	err := db.QueryRow(`insert into income (user_id, date, amount, wallet_id, note) values($1, $2, $3, $4, $5) returning id`, transaction.UserKeys.User.Id, transaction.Date, transaction.Amount, transaction.Wallet.Id, transaction.Note).Scan(&transaction.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userKeys *UserKeys) getIncomeList() []Transaction {
	statement := fmt.Sprintf("select income.id, income.date, income.amount, income.note FROM income WHERE user_id=%d", userKeys.User.Id)
	//TODO: Error Checking
	transactions, _, _ := postgresql_access.QueryDatabase(db, statement)
	transaction_list := []Transaction{}

	for _, transaction := range transactions {
		single_transaction := Transaction{}
		single_transaction.UserKeys = *userKeys
		single_transaction.Id, _ = strconv.Atoi(transaction[0].(string))
		single_transaction.Date = transaction[1].(string)
		single_transaction.Amount, _ = strconv.ParseFloat(transaction[2].(string), 64)
		single_transaction.Note, _ = transaction[3].(string)
		transaction_list = append(transaction_list, single_transaction)
	}
	return transaction_list
}

func (transaction *Transaction) RegisterNewExpense() error {
	err := db.QueryRow(`insert into expense (user_id, date, amount, wallet_id, note) values($1, $2, $3, $4, $5) returning id`, transaction.UserKeys.User.Id, transaction.Date, transaction.Amount, transaction.Wallet.Id, transaction.Note).Scan(&transaction.Id)
	if err != nil {
		return err
	}
	return nil
}
