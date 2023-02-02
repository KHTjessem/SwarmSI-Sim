package main

type RentOracle interface {
	GetRentPrice() float64
}

type FixedRentOracle struct {
	fixedPrice float64
}

func (ro FixedRentOracle) GetRentPrice() float64 {
	return ro.fixedPrice
}
