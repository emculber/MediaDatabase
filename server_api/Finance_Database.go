package main

import (
	"fmt"

	"github.com/emculber/database_access/postgresql"
)

var financeDatabaseSchema = []string{
	"CREATE TABLE wallet(id serial primary key, user_id integer references registered_user(id), name varchar, percent real, current_amount real)",
	"CREATE TABLE income(id serial primary key, user_id integer references registered_user(id), date varchar, amount real, wallet integer references wallet(id), note varchar)",
}

func (wallet *Wallet) RegisterNewWallet() error {
	err := db.QueryRow(`insert into wallet (user_id, name, percent, current_amount) values($1, $2, $3, $4) returning id`, wallet.UserKeys.User.Id, wallet.Name, wallet.Percent, wallet.CurrentAmount).Scan(&wallet.Id)
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
		single_wallet.Id = wallet[0].(int)
		single_wallet.Name = wallet[1].(string)
		single_wallet.Percent = wallet[2].(float64)
		single_wallet.CurrentAmount = wallet[3].(float64)
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
