package main

type postageContract interface {
	CollectWinnings(roundPrice int, no *node)
	GetTotalPayout() int
}

type simpleFixedPostage struct {
	totalWithdrawn int
}

func (sfp *simpleFixedPostage) CollectWinnings(roundPrice int, no *node) {
	// Fixed just uses price without considering storage space in network.
	// Assumes static.
	sfp.totalWithdrawn += roundPrice
	no.Earnings += roundPrice
}
func (sfp *simpleFixedPostage) GetTotalPayout() int {
	return sfp.totalWithdrawn
}
