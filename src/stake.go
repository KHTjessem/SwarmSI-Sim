package main

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
