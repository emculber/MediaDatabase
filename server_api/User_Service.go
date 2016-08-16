package main

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func (userKeys *UserKeys) generateKey() {
	b := make([]rune, 64)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	//fmt.Println(string(b))
	userKeys.Key = string(b)
}
