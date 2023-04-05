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
		x = math.Round(x/float64(st.roundBy)) * float64(st.roundBy)
	}

	return int(x)
}

// 1st node is in bucket 1, it gets the full stake of X.
// Secound and third are in bucket two, and each gets X/2 stake.
// Fourth, fifth and sixth nodes in bucket 3 gets X/3 stake.
// So the boucket index is the amount of nodes, stake is divided by
// amount of nodes.
type bucketSumStake struct {
	stake int

	nc     *int
	bucket *int
}

// Keeps track of bukects internally. Bucket B_i has i nodes.
// B_1 has 1 nodes, B_2 has 2 etc.
func (bss bucketSumStake) GetStake(nodeID int) int {
	// st := bss.stake / bss.bucket

	st := bss.stake / *bss.bucket

	*bss.nc++
	if *bss.nc == *bss.bucket {
		*bss.bucket++
		*bss.nc = 0
	}

	return st
}
