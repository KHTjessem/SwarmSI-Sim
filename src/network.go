package main

import "math/rand"

type SwarmNetwork interface {
	CreateSwarmNetwork()
	UpdateNetwork()
	SelectNeighbourhood() *neighbourhood
	SelectWinner() *node
	GetNodeArray() *[]*node
	GetNodeAdressMap() map[uint64]*node
}

// type neighbourhood interface {

// }

type FixedIdealSwarmNetwork struct {
	// Amount of nodes needs to be divisible by four
	networkNodeCount int
	neighbourhoods   []neighbourhood
	nodeAddressMap   map[uint64]*node
	nodes            []*node
}

func (sn *FixedIdealSwarmNetwork) CreateSwarmNetwork() {
	//sn.neighbourhoods := make([]neighbourhood, 0, sn.networkNodeCount/4)
	sn.nodeAddressMap = make(map[uint64]*node)

	newhood := neighbourhood{}
	for i := 0; i < sn.networkNodeCount; i++ {
		nn := node{Id: uint64(i), Earnings: 0}
		sn.nodeAddressMap[nn.Id] = &nn
		newhood.nodes = append(newhood.nodes, &nn)
		sn.nodes = append(sn.nodes, &nn)

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

func (sn *FixedIdealSwarmNetwork) GetNodeArray() *[]*node {
	return &sn.nodes
}
func (sn *FixedIdealSwarmNetwork) GetNodeAdressMap() map[uint64]*node {
	return sn.nodeAddressMap
}
