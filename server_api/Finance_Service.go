package main

import "fmt"

func (transaction *Transaction) SplitMoney() {
	wallets := transaction.UserKeys.getWalletList()
	fmt.Println(wallets)

	var wallet_diff float64

	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Transaction Amount", transaction.Amount)
		fmt.Println("Wallet Percent", wallet.Percent)
		if wallet.Percent <= 100 && wallet.Percent > 0 {

			wallet_amount := wallet.CurrentAmount + (transaction.Amount * (wallet.Percent / 100))
			if (wallet.Limit < (wallet_amount + wallet.CurrentAmount)) && wallet.Limit != -1 {
				wallet_diff += wallet_amount - wallet.Limit
				wallet.CurrentAmount = wallet.Limit
				wallet.Percent = 0
			} else {
				wallet.CurrentAmount = wallet_amount
			}

			fmt.Println("Wallet Name", wallet.Name)
			fmt.Println("Wallet Decimal Percent", wallet.Percent/100)
			fmt.Println("Wallet Amount", transaction.Amount*(wallet.Percent/100))

			if err := wallet.updateWallet(); err != nil {
				fmt.Println("Error while Updating wallet ->", err)
			}
			fmt.Println(wallet)
		}
	}

	unallocatedWallet := Wallet{Name: "unallocated"}
	unallocatedWallet.getUnallocatedWallet()
	unallocatedWallet.CurrentAmount += wallet_diff
	if err := unallocatedWallet.updateWallet(); err != nil {
		fmt.Println("Error while Updating wallet ->", err)
	}
	fmt.Println("------")
	transaction.UserKeys.DisperseUnallocatedWalletPercent()
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

func (userKeys *UserKeys) DisperseUnallocatedWalletPercent() {
	unallocatedWallet := Wallet{Name: "unallocated"}
	unallocatedWallet.getUnallocatedWallet()

	wallets := userKeys.getWalletList()
	walletsSetToUpdate := []Wallet{}
	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Wallet Requested Percent", wallet.RequestedPercent)
		fmt.Println("Wallet Percent", wallet.Percent)

		if (wallet.RequestedPercent != wallet.Percent) && wallet.RequestedPercent != -1 {
			if wallet.Percent <= 100 && wallet.Percent >= 0 && wallet.Limit > wallet.CurrentAmount {
				walletsSetToUpdate = append(walletsSetToUpdate, wallet)
			}
		}
	}

	splitPercent := unallocatedWallet.Percent / float64(len(walletsSetToUpdate))
	var overflow float64 = 0
	fmt.Println("Split Wallet Percent ->", splitPercent)

	for _, wallet := range walletsSetToUpdate {
		if wallet.RequestedPercent < (wallet.Percent + splitPercent) {
			overflow += (wallet.Percent + splitPercent) - wallet.RequestedPercent
			wallet.Percent = wallet.RequestedPercent
		} else {
			wallet.Percent += splitPercent
			//unallocatedWallet.Percent -= splitPercent
		}
		if err := wallet.updateWallet(); err != nil {
			fmt.Println("Error while Updating wallet ->", err)
		}
	}
	unallocatedWallet.Percent = overflow
	if err := unallocatedWallet.updateWallet(); err != nil {
		fmt.Println("Error while Updating wallet ->", err)
	}
}
