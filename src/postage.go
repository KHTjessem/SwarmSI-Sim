package main

type postageContract interface {
	CollectWinnings(roundPrice float64, no *node)
	GetTotalPayout() float64
}

type simpleFixedPostage struct {
	totalWithdrawn float64
}

func (sfp *simpleFixedPostage) CollectWinnings(roundPrice float64, no *node) {
	// Fixed just uses price without considering storage space in network.
	// Assumes static.
	sfp.totalWithdrawn += roundPrice
	no.Earnings += roundPrice
}
func (sfp *simpleFixedPostage) GetTotalPayout() float64 {
	return sfp.totalWithdrawn
}
