package main

import "fmt"

func (transaction *Transaction) SplitMoney() {
	wallets := transaction.UserKeys.getWalletList()
	fmt.Println(wallets)

	unallocatedWallet := Wallet{Name: "unallocated"}
	unallocatedWallet.getUnallocatedWallet()

	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Transaction Amount", transaction.Amount)
		fmt.Println("Wallet Percent", wallet.Percent)
		if wallet.Percent <= 100 && wallet.Percent > 0 {

			wallet_amount := wallet.CurrentAmount + (transaction.Amount * (wallet.Percent / 100))
			if (wallet.WalletLimit < (wallet_amount + wallet.CurrentAmount)) && wallet.WalletLimit != -1 {

				unallocatedWallet.CurrentAmount += wallet_amount - wallet.WalletLimit
				unallocatedWallet.Percent += wallet.Percent

				wallet.CurrentAmount = wallet.WalletLimit
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
			if wallet.Percent <= 100 && wallet.Percent >= 0 && wallet.WalletLimit > wallet.CurrentAmount {
				fmt.Println("Wallet getting added to Update ->", wallet)
				walletsSetToUpdate = append(walletsSetToUpdate, wallet)
			}
		}
	}

	if len(walletsSetToUpdate) > 0 {
		fmt.Println("Updating wallets ->", len(walletsSetToUpdate))
		splitPercent := unallocatedWallet.Percent / float64(len(walletsSetToUpdate))
		var overflow float64 = 0
		fmt.Println("Split Wallet Percent ->", splitPercent)

		for _, wallet := range walletsSetToUpdate {
			if wallet.RequestedPercent < (wallet.Percent + splitPercent) {
				overflow += (wallet.Percent + splitPercent) - wallet.RequestedPercent
				wallet.Percent = wallet.RequestedPercent
			} else {
				wallet.Percent += splitPercent
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
}
