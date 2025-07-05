package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/account"
	"github.com/brunomvsouza/ynab.go/api/transaction"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const Version = "0.1.0"

func runUI(accounts []*account.Account, txns []*transaction.Transaction) {
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].Date.Time.After(txns[j].Date.Time)
	})

	txnByID := map[string][]string{}
	for _, txn := range txns {
		amStr := amtFmt(txn.Amount, 12)
		txnByID[txn.AccountID] = append(txnByID[txn.AccountID], fmt.Sprintf("%s %s %s", api.DateFormat(txn.Date), amStr, *txn.PayeeName))
	}

	rows := []string{}
	for _, account := range accounts {
		rows = append(rows, fmt.Sprintf("%-13s %s", account.Name, amtFmt(account.Balance, 10)))
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	lb := widgets.NewList()
	lb.Rows = rows
	lb.SelectedRow = 0
	lb.SelectedRowStyle = ui.NewStyle(ui.ColorBlack)
	lb.TextStyle.Fg = ui.ColorGreen
	lb.Title = fmt.Sprintf("ynab" + Version)

	detail := widgets.NewParagraph()

	pa := widgets.NewParagraph()
	pa.Text = "[Q]uit"
	pa.TextStyle.Fg = ui.ColorBlack

	w, h := ui.TerminalDimensions()
	split := w / 3
	lb.SetRect(0, 0, split, h-3)
	detail.SetRect(split+1, 0, w, h-3)
	pa.SetRect(0, h-3, w, h)
	detail.Text = strings.Join(fmtDetails(accounts[lb.SelectedRow], txnByID), "\n")
	ui.Render(lb, pa, detail)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "k", "<Up>":
				lb.ScrollUp()
			case "j", "<Down>":
				lb.ScrollDown()
			case "q":
				return
			}
			detail.Text = strings.Join(fmtDetails(accounts[lb.SelectedRow], txnByID), "\n")
			ui.Render(lb, detail)
		}
	}
}
