package service

type bankDeposit struct {
	pgBankDeposit IBankDepositRepo
}

func newBankDeposit(pgBankDeposit IBankDepositRepo) *bankDeposit {
	return &bankDeposit{pgBankDeposit: pgBankDeposit}
}
