package main

import (
	"fmt"
	"math/rand"
)

type simulator struct {
	// configs
	totalNodeCount int
	maxRounds      int
	mathRandSeed   int64

	// parts
	swarmnetwork SwarmNetwork
	rentoracle   RentOracle
	postage      postageContract

	// Stat tracking (generated form running simulation)
	round int
}

func (s *simulator) Setup() {
	// set the seed
	rand.Seed(s.mathRandSeed)

	s.swarmnetwork.CreateSwarmNetwork()
}

func (s *simulator) MainLoop() {
	// The main loop of the simulator
	print("Staring simulation")
	for s.round = 0; s.round < s.maxRounds; s.round++ {
		roundPrice := s.rentoracle.GetRentPrice()

		//select winner
		winner := s.swarmnetwork.SelectWinner()

		// collection
		s.postage.CollectWinnings(roundPrice, winner)

		// Simulate change
		s.swarmnetwork.UpdateNetwork()

	}

	fmt.Printf("Rounds: %v", s.maxRounds)

}
