package main

import (
	"time"

	"github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/transaction"
	"github.com/mschilli/go-murmur"
)

func main() {
	apiToken, err := murmur.NewMurmur().Lookup("ynab-test")
	if err != nil {
		panic(err)
	}

	c := ynab.NewClient(apiToken)
	budgets, err := c.Budget().GetBudgets()
	if err != nil {
		panic(err)
	}

	budgetID := budgets[0].ID
	snp, err := c.Account().GetAccounts(budgetID, nil)
	if err != nil {
		panic(err)
	}

	since := time.Now().AddDate(0, -2, 0)
	txns, err := c.Transaction().GetTransactions(budgetID, &transaction.Filter{Since: &api.Date{Time: since}})
	if err != nil {
		panic(err)
	}

	runUI(snp.Accounts, txns)
}
