package main

type Wallet struct {
	Id               int
	UserKeys         UserKeys
	Name             string
	RequestedPercent float64
	Percent          float64
	CurrentAmount    int
	WalletLimit      int
}

type Transaction struct {
	Id       int
	UserKeys UserKeys
	//Date     time
	Date   string
	Amount int
	Wallet Wallet
	Note   string
}

//TODO: Complete OK functions
func (transaction *Transaction) OK() error {
	return nil
}

func (wallet *Wallet) OK() error {
	return nil
}
