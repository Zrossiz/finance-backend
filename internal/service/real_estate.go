package service

type realEstate struct {
	pgRealEstate IRealEstateRepo
}

func newRealEstate(pgRealEstate IRealEstateRepo) *realEstate {
	return &realEstate{pgRealEstate: pgRealEstate}
}
