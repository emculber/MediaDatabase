package main

type Wallet struct {
	Id            int
	UserKeys      UserKeys
	Name          string
	Percent       float64
	CurrentAmount float64
}

type Transaction struct {
	Id       int
	UserKeys UserKeys
	//Date     time
	Date   string
	Amount float64
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
