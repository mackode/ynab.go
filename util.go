package main

import (
	"fmt"

	"github.com/brunomvsouza/ynab.go/api/account"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func amtFmt(amount int64, wide int) string {
	amt := float64(amount) / 1000.0
	color := "(fg:green)"
	sign := ""
	if amt < 0 {
		sign = "-"
		amt = -amt
		color = "(fg:red)"
	}

	p := message.NewPrinter(language.English)
	amStr := p.Sprintf("%s$%.2f", sign, amt)
	return fmt.Sprintf("[%*s]%s", wide, amStr, color)
}

func fmtDetails(account *account.Account, txnByID map[string][]string) []string {
	details := []string{
		fmt.Sprintf("%-11s%s", "Balance", amtFmt(account.Balance, 12)),
	}
	for _, e := range txnByID[account.ID] {
		details = append(details, e)
	}
	return details
}
