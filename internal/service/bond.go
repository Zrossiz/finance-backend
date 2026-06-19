package service

type bond struct {
	pgBond IBondRepo
}

func newBond(pgBond IBondRepo) *bond {
	return &bond{pgBond: pgBond}
}
