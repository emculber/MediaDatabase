package main

import "fmt"

func (transaction *Transaction) SplitMoney() {
	wallets := transaction.UserKeys.getWalletList()
	fmt.Println(wallets)

	var wallet_diff float64
	var unallocatedWallet Wallet
	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Transaction Amount", transaction.Amount)
		fmt.Println("Wallet Percent", wallet.Percent)
		if wallet.Percent <= 100 && wallet.Percent > 0 {

			wallet_amount := wallet.CurrentAmount + (transaction.Amount * (wallet.Percent / 100))
			if wallet.Limit < (wallet_amount + wallet.CurrentAmount) {
				wallet_diff += wallet_amount - wallet.Limit
				wallet.CurrentAmount = wallet.Limit
				wallet.Percent = 0
			} else {
				wallet.CurrentAmount = wallet_amount
			}

			fmt.Println("Wallet Decimal Percent", wallet.Percent/100)
			fmt.Println("Wallet Amount", transaction.Amount*(wallet.Percent/100))

			if err := wallet.updateWallet(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(wallet)
			//TODO: Move this to its own sql statment
			if wallet.Name == "unallocated" {
				unallocatedWallet = wallet
			}
		}
	}
	unallocatedWallet.CurrentAmount += wallet_diff
	if err := unallocatedWallet.updateWallet(); err != nil {
		fmt.Println(err)
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
