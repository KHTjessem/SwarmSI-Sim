package main

type postageContract interface {
	GetWinnings(float64) float64
}

type simpleFixedPostage struct {
	totalWithdrawn float64
}

func (sfp *simpleFixedPostage) GetWinnings(winnings float64) float64 {
	sfp.totalWithdrawn += winnings
	return winnings
}
