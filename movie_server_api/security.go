package main

var permissions = make(map[string][]string)

func InitSecurity() {
	permissions["movie"] = []string{"admin"}
}
