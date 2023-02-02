package main

func main() {
	print("Hello Swarm!")

	s := &simulator{
		totalNodeCount: 1200,
		swarmnetwork:   &SimpleSwarmNetwork{},
		rentoracle:     &FixedRentOracle{},
		round:          0,
	}

	s.MainLoop()
}
