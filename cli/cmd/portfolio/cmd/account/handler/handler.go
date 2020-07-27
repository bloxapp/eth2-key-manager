package handler

import "github.com/bloxapp/KeyVault/cli/util/printer"

// Account contains handler functions of the CLI commands related to portfolio account.
type Account struct {
	printer printer.Printer
}

// New is the constructor of Account handler.
func New(printer printer.Printer) *Account {
	return &Account{
		printer: printer,
	}
}
