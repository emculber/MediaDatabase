package main

import "fmt"

func (transaction *Transaction) SplitMoney() {
	wallets := transaction.UserKeys.getWalletList()
	fmt.Println(wallets)

	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Transaction Amount", transaction.Amount)
		fmt.Println("Wallet Percent", wallet.Percent)
		if wallet.Percent <= 100 && wallet.Percent > 0 {
			wallet.CurrentAmount += transaction.Amount * (wallet.Percent / 100)

			fmt.Println("Wallet Decimal Percent", wallet.Percent/100)
			fmt.Println("Wallet Amount", transaction.Amount*(wallet.Percent/100))

			if err := wallet.updateWallet(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(wallet)
		}
	}
	fmt.Println("------")
}

func (transaction *Transaction) TakeFromWallet() {
	fmt.Println("Transaction Amount", transaction.Amount)

	if transaction.Amount > 0 {
		transaction.Amount = -transaction.Amount
	}
	transaction.Wallet.CurrentAmount += transaction.Amount
	if err := transaction.Wallet.updateWallet(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(transaction.Wallet)
}
