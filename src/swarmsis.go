package main

import "math/rand"

type simulator struct {
	// configs
	totalNodeCount int

	// parts
	swarmnetwork SwarmNetwork
	rentoracle   RentOracle

	// structs
	neighbourhoods []*neighbourhood

	// Stat tracking (generated form running simulation)
	round int
}

func (s *simulator) Setup() {
	s.neighbourhoods = *s.swarmnetwork.CreateSwarmNetwork()
}

func (s *simulator) MainLoop() {
	// The main loop of the simulator
	print("Staring new round")
	s.round += 1

	// Move the selection process to the network
	anch := rand.Intn(len(s.neighbourhoods))
	actNeigbourhood := s.neighbourhoods[anch]

	pot := s.rentoracle.GetRentPrice()
	actNeigbourhood.SelectWinner().AddEarnings(pot)
}
