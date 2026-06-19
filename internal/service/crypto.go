package service

type cryptoPosition struct {
	pgCryptoPositon ICryptoPositionRepo
}

func newCryptoPosition(pgCryptoPositon ICryptoPositionRepo) *cryptoPosition {
	return &cryptoPosition{pgCryptoPositon: pgCryptoPositon}
}
