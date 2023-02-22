package main

type postageContract interface {
	CollectWinnings(roundPrice float64, no *node)
}

type simpleFixedPostage struct {
	totalWithdrawn float64
}

func (sfp *simpleFixedPostage) CollectWinnings(roundPrice float64, no *node) {
	// Fixed just uses price without considering storage space in network.
	// Assumes static.
	sfp.totalWithdrawn += roundPrice
	no.earnings += roundPrice
}
