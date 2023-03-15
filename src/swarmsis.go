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

	logChan    chan *logObject
	logStopped chan bool

	// Stat tracking (generated form running simulation)
	round int
}

func (s *simulator) Setup() {
	// set the seed
	rand.Seed(s.mathRandSeed)

	s.swarmnetwork.CreateSwarmNetwork()

	// Start logger
	s.logStopped = make(chan bool)
	go logger(s.logChan, s.logStopped, s.totalNodeCount, "Simulation of 2048 nodes, static network - same stake")
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

		// Log changes
		lgo := s.createRoundStat(roundPrice)
		s.logChan <- lgo
	}

	fmt.Printf("Rounds: %v", s.maxRounds)

	// End logger. TODO: Consider move this to main.go
	closemsg := logObject{Round: -1}
	s.logChan <- &closemsg
	<-s.logStopped // File has been written
}

func (s *simulator) createRoundStat(roundPrice float64) *logObject {
	lgo := logObject{Round: s.round,
		TotalPayout: s.postage.GetTotalPayout(),
		RoundPrice:  roundPrice,
		Nodes:       s.swarmnetwork.GetNodeArray()}

	return &lgo
}
