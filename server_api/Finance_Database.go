package main

import (
	"fmt"
	"strconv"

	"github.com/emculber/database_access/postgresql"
)

var financeDatabaseSchema = []string{
	"CREATE TABLE wallet(id serial primary key, user_id integer references registered_user(id), name varchar, percent real, current_amount real)",
	"CREATE TABLE income(id serial primary key, user_id integer references registered_user(id), date varchar, amount real, wallet integer references wallet(id), note varchar)",
}

func CreateFinanceTables() {
	for _, table := range financeDatabaseSchema {
		fmt.Println(table)
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (wallet *Wallet) RegisterNewWallet() error {
	err := db.QueryRow(`insert into wallet (user_id, name, percent, current_amount) values($1, $2, $3, $4) returning id`, wallet.UserKeys.User.Id, wallet.Name, wallet.Percent, wallet.CurrentAmount).Scan(&wallet.Id)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) getWallet() error {
	err := db.QueryRow("select wallet.id, wallet.name, wallet.percent, wallet.current_amount from wallet where wallet.id = $1", wallet.Id).Scan(&wallet.Id, &wallet.Name, &wallet.Percent, &wallet.CurrentAmount)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) updateWallet() error {
	err := db.QueryRow("UPDATE wallet SET name = $1, percent = $2, current_amount = $3 WHERE id = $4 returning id", wallet.Name, wallet.Percent, wallet.CurrentAmount, wallet.Id).Scan(&wallet.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userKeys *UserKeys) getWalletList() []Wallet {
	statement := fmt.Sprintf("select wallet.id, wallet.name, wallet.percent, wallet.current_amount from wallet where user_id=%d", userKeys.User.Id)
	//TODO: Error Checking
	wallets, _, _ := postgresql_access.QueryDatabase(db, statement)
	wallet_list := []Wallet{}

	for _, wallet := range wallets {
		single_wallet := Wallet{}
		single_wallet.UserKeys = *userKeys
		single_wallet.Id, _ = strconv.Atoi(wallet[0].(string))
		single_wallet.Name = wallet[1].(string)
		single_wallet.Percent, _ = strconv.ParseFloat(wallet[2].(string), 64)
		single_wallet.CurrentAmount, _ = strconv.ParseFloat(wallet[3].(string), 64)
		wallet_list = append(wallet_list, single_wallet)
	}
	return wallet_list
}

func (income *Income) RegisterNewIncome() error {
	err := db.QueryRow(`insert into income (user_id, date, amount, wallet_id, note) values($1, $2, $3, $4, $5) returning id`, income.UserKeys.User.Id, income.Date, income.Amount, income.Wallet.Id, income.Note).Scan(&income.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userKeys *UserKeys) getIncomeList() []Income {
	statement := fmt.Sprintf("select income.id, income.date, income.amount, income.note FROM income WHERE user_id=%d", userKeys.User.Id)
	//TODO: Error Checking
	incomes, _, _ := postgresql_access.QueryDatabase(db, statement)
	income_list := []Income{}

	for _, income := range incomes {
		single_income := Income{}
		single_income.UserKeys = *userKeys
		single_income.Id, _ = strconv.Atoi(income[0].(string))
		single_income.Date = income[1].(string)
		single_income.Amount, _ = strconv.ParseFloat(income[2].(string), 64)
		single_income.Note, _ = income[3].(string)
		income_list = append(income_list, single_income)
	}
	return income_list
}
