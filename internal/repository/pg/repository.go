package pgrepo

import (
	"database/sql"
	"fmt"
)

type Postgres struct {
	User           *user
	Stock          *stock
	Bond           *bond
	BankDeposit    *bankDeposit
	CryptoPosition *cryptoPosition
	RealEstate     *realEstate
}

func Connect(uri string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("postgres conn err: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping postgres db err: %w", err)
	}

	return conn, nil
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{
		User:           newUser(conn),
		Stock:          newStock(conn),
		RealEstate:     newRealEstate(conn),
		CryptoPosition: newCryptoPosition(conn),
		Bond:           newBond(conn),
		BankDeposit:    newBankDeposit(conn),
	}
}
