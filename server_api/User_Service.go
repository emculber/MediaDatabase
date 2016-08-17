package main

import (
	"math/rand"
	"time"
)

func (userKeys *UserKeys) generateKey() {
	strlen := 64
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	userKeys.Key = string(result)
}
