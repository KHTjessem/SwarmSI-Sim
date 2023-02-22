package main

import "math/rand"

type SwarmNetwork interface {
	CreateSwarmNetwork()
	UpdateNetwork()
	SelectNeighbourhood() *neighbourhood
	SelectWinner() *node
}

// type neighbourhood interface {

// }

type FixedIdealSwarmNetwork struct {
	// Amount of nodes needs to be divisible by four
	networkNodeCount int
	neighbourhoods   []neighbourhood
	nodes            map[uint64]*node
}

func (sn *FixedIdealSwarmNetwork) CreateSwarmNetwork() {
	//sn.neighbourhoods := make([]neighbourhood, 0, sn.networkNodeCount/4)
	sn.nodes = make(map[uint64]*node)

	newhood := neighbourhood{}
	for i := 0; i < sn.networkNodeCount; i++ {
		nn := node{id: uint64(i), earnings: 0}
		sn.nodes[nn.id] = &nn
		newhood.nodes = append(newhood.nodes, &nn)

		if (i+1)%4 == 0 {
			sn.neighbourhoods = append(sn.neighbourhoods, newhood)

			newhood = neighbourhood{}
		}
	}
}

func (sn *FixedIdealSwarmNetwork) SelectNeighbourhood() *neighbourhood {
	// Anchor is selected at random. Here it is assumed the chance is 1/len(neighbourhoods)
	ind := rand.Intn(len(sn.neighbourhoods))
	return &sn.neighbourhoods[ind]

}
func (sn *FixedIdealSwarmNetwork) SelectWinner() *node {
	// In this it's assumed that every node has the same stake. final winner is therefor 1/len(neighbourhood)
	hoodind := rand.Intn(len(sn.neighbourhoods))
	winind := rand.Intn(len(sn.neighbourhoods[hoodind].nodes))
	return sn.neighbourhoods[hoodind].nodes[winind]
}
func (sn *FixedIdealSwarmNetwork) UpdateNetwork() {
	// Fixed network, no change
}
