package main

type RentOracle interface {
	GetRentPrice() int
}

type FixedRentOracle struct {
	fixedPrice int
}

func (ro FixedRentOracle) GetRentPrice() int {
	return ro.fixedPrice
}
