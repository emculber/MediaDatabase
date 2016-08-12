package main

import log "github.com/Sirupsen/logrus"

func (transaction *Transaction) SplitMoney() {
	wallets := transaction.UserKeys.getWalletList()
	log.WithFields(log.Fields{
		"Transaction": transaction,
		"Wallets":     wallets,
	}).Info("Spliting Money")

	unallocatedWallet := Wallet{Name: "unallocated"}
	unallocatedWallet.getUnallocatedWallet()

	for _, wallet := range wallets {
		log.WithFields(log.Fields{
			"Transaction Amount": transaction.Amount,
			"Wallet Percent":     wallet.Percent,
		}).Info("Split Info")
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

			log.WithFields(log.Fields{
				"Wallet Name":            wallet.Name,
				"Wallet Decimal Percent": wallet.Percent / 100,
				"Wallet Amount":          transaction.Amount * (wallet.Percent / 100),
			}).Info("Wallet Info")

			if err := wallet.updateWallet(); err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error While Updating Wallet")
			}
			log.WithFields(log.Fields{
				"Wallet": wallet,
			}).Info("Wallet Info")
		}
	}

	if err := unallocatedWallet.updateWallet(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error While Updating Unallocated Wallet")
	}
	transaction.UserKeys.DisperseUnallocatedWalletPercent()
}

func (transaction *Transaction) TakeFromWallet() {
	log.WithFields(log.Fields{
		"Transaction Amount": transaction.Amount,
	}).Info("Take From Wallet")

	if transaction.Amount > 0 {
		transaction.Amount = -transaction.Amount
	}
	transaction.Wallet.CurrentAmount += transaction.Amount
	if err := transaction.Wallet.updateWallet(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error While Updating Unallocated Wallet")
	}
	log.WithFields(log.Fields{
		"Transaction Wallet": transaction.Wallet,
	}).Info("Wallet Info")
}

func (userKeys *UserKeys) DisperseUnallocatedWalletPercent() {
	unallocatedWallet := Wallet{Name: "unallocated"}
	unallocatedWallet.getUnallocatedWallet()

	wallets := userKeys.getWalletList()
	walletsSetToUpdate := []Wallet{}
	for _, wallet := range wallets {
		log.WithFields(log.Fields{
			"Wallet Requested Percent": wallet.RequestedPercent,
			"Wallet Percent":           wallet.Percent,
		}).Info("Wallet Percents")

		if (wallet.RequestedPercent != wallet.Percent) && wallet.RequestedPercent != -1 {
			if wallet.Percent <= 100 && wallet.Percent >= 0 && wallet.WalletLimit > wallet.CurrentAmount {
				log.WithFields(log.Fields{
					"Wallet": wallet,
				}).Info("Wallet Getting Added to Update")
				walletsSetToUpdate = append(walletsSetToUpdate, wallet)
			}
		}
	}

	if len(walletsSetToUpdate) > 0 {
		log.WithFields(log.Fields{
			"Number of Wallets": len(walletsSetToUpdate),
		}).Info("Updating Wallets")
		splitPercent := unallocatedWallet.Percent / float64(len(walletsSetToUpdate))
		var overflow float64 = 0

		log.WithFields(log.Fields{
			"Split Percent": splitPercent,
		}).Info("Spliting Wallet Percent")

		for _, wallet := range walletsSetToUpdate {
			if wallet.RequestedPercent < (wallet.Percent + splitPercent) {
				overflow += (wallet.Percent + splitPercent) - wallet.RequestedPercent
				wallet.Percent = wallet.RequestedPercent
			} else {
				wallet.Percent += splitPercent
			}
			if err := wallet.updateWallet(); err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error While Updating Wallet")
			}
		}
		unallocatedWallet.Percent = overflow
		if err := unallocatedWallet.updateWallet(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error While Updating Unallocated Wallet")
		}
	}
}
