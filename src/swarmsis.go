package main

import (
	"fmt"
	"math/rand"
)

type simulator struct {
	// configs
	totalNodeCount int
	maxRounds      int
	SetupSeed      int64
	simulationSeed int64

	// parts
	swarmnetwork SwarmNetwork
	rentoracle   RentOracle
	postage      postageContract

	logChan chan *logObject

	// Stat tracking (generated form running simulation)
	round int
}

func (s *simulator) Setup() {
	// set the seed
	rand.Seed(s.SetupSeed)

	s.swarmnetwork.CreateSwarmNetwork()
}

func (s *simulator) MainLoop() {
	rand.Seed(s.simulationSeed)
	// The main loop of the simulator
	print("Staring simulation\n")
	for s.round = 0; s.round < s.maxRounds; s.round++ {
		roundPrice := s.rentoracle.GetRentPrice()

		//select winner
		winner := s.swarmnetwork.SelectWinner()

		// collection
		s.postage.CollectWinnings(roundPrice, winner)

		// Simulate change
		s.swarmnetwork.UpdateNetwork()

		// Log changes
		lgo := s.createRoundStat(roundPrice)
		s.logChan <- lgo
	}

	fmt.Printf("Done with %v rounds\n", s.maxRounds)
}

func (s *simulator) createRoundStat(roundPrice int) *logObject {
	lgo := logObject{Round: s.round,
		TotalPayout: s.postage.GetTotalPayout(),
		RoundPrice:  roundPrice,
		Nodes:       s.swarmnetwork.GetNodeArray()}

	return &lgo
}
