package main

import (
	"math"
	"math/rand"
)

// For creating different stakes.

type StakeCreator interface {
	GetStake(nodeID int) int
}

type EqualStake struct {
	amount int
}

func (st EqualStake) GetStake(nodeID int) int {
	return st.amount
}

// The stake will be distributed by power series
type PowerDistStake struct {
	alpha      float64
	minStake   int
	rounding   bool
	roundBy    int
	limitStake bool
	maxStake   int
}

func (st PowerDistStake) GetStake(nodeID int) int {
	r := rand.Float64()

	x := float64(st.minStake) * math.Pow(1-r, -1/(st.alpha-1))

	if st.limitStake && x > float64(st.maxStake) {
		x = float64(st.maxStake)
	}
	if st.rounding {
		x = math.Round(x*float64(st.roundBy)) * float64(st.roundBy)
	}

	return int(x)
}
