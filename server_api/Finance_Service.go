package main

import "fmt"

func (income *Income) SplitMoney() {
	wallets := income.UserKeys.getWalletList()
	fmt.Println(wallets)

	for _, wallet := range wallets {
		fmt.Println("------")
		fmt.Println("Income Amount", income.Amount)
		fmt.Println("Wallet Percent", wallet.Percent)
		if wallet.Percent <= 100 && wallet.Percent > 0 {
			wallet.CurrentAmount += income.Amount * (wallet.Percent / 100)

			fmt.Println("Wallet Decimal Percent", wallet.Percent/100)
			fmt.Println("Wallet Amount", income.Amount*(wallet.Percent/100))

			if err := wallet.updateWallet(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(wallet)
		}
	}
	fmt.Println("------")
}
