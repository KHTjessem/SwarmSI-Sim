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

	// logChan chan *logObject
	saver storer

	// Stat tracking (generated form running simulation)
	round       int
	roundPrice  int
	roundWinner *node
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
		s.roundPrice = s.rentoracle.GetRentPrice()

		//select winner
		s.roundWinner = s.swarmnetwork.SelectWinner()

		// collection
		s.postage.CollectWinnings(s.roundPrice, s.roundWinner)

		// Log changes
		s.saver.save(s)

		// Simulate change
		s.swarmnetwork.UpdateNetwork()
	}

	fmt.Printf("Done with %v rounds\n", s.maxRounds)
}

func (s *simulator) createRoundStat() *logObject {
	lgo := logObject{Round: s.round,
		TotalPayout: s.postage.GetTotalPayout(),
		RoundPrice:  s.roundPrice,
		Nodes:       s.swarmnetwork.GetNodeArray()}

	return &lgo
}
