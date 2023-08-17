package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(id int, firstName, lastName string) *Account {
	return &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
		Balance:   0,
	}
}
