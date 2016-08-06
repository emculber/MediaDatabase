package main

import "fmt"

func (income *Income) SplitMoney() {
	wallets := income.UserKeys.getWalletList()
	fmt.Println(wallets)

	for _, wallet := range wallets {
		if wallet.Percent <= 100 && wallet.Percent > 0 {
			wallet.CurrentAmount += income.Amount * (wallet.Percent / 100)
			income.Amount -= income.Amount * (wallet.Percent / 100)
			if err := wallet.updateWallet(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(wallet)
		}
	}
}
