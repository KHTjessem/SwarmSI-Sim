package main

import "math/rand"

type SwarmNetwork interface {
	CreateSwarmNetwork()
	UpdateNetwork()
	SelectNeighbourhood() *neighbourhood
	SelectWinner() *node
	GetNodeArray() *[]node
	GetNodeAdressMap() map[uint64]*node
}

// type neighbourhood interface {

// }

type FixedIdealSwarmNetwork struct {
	// Amount of nodes needs to be divisible by four
	networkNodeCount  int
	neighbourhoods    []neighbourhood
	nodeAddressMap    map[uint64]*node
	nodes             []*node
	stakeDistribution StakeCreator
}

func (sn *FixedIdealSwarmNetwork) CreateSwarmNetwork() {
	//sn.neighbourhoods := make([]neighbourhood, 0, sn.networkNodeCount/4)
	sn.nodeAddressMap = make(map[uint64]*node)

	newhood := neighbourhood{}
	for i := 0; i < sn.networkNodeCount; i++ {
		nn := node{
			Id:       uint64(i),
			Earnings: 0,
			stake:    sn.stakeDistribution.GetStake(i),
		}

		sn.nodeAddressMap[nn.Id] = &nn
		newhood.nodes = append(newhood.nodes, &nn)
		sn.nodes = append(sn.nodes, &nn)

		if (i+1)%4 == 0 {
			newhood.nodeCount = len(newhood.nodes)
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
	nbhood := sn.SelectNeighbourhood()

	// It's weigthed by the stake of the nodes.
	weigthSum := 0
	for i := 0; i < nbhood.nodeCount; i++ {
		weigthSum += nbhood.nodes[i].stake
	}
	num := rand.Intn(weigthSum)

	// Should always return a winner.
	// Since num should be less than total
	// weighted sum.
	for i := 0; i < nbhood.nodeCount; i++ {
		num -= nbhood.nodes[i].stake
		if num <= 0 {
			return nbhood.nodes[i]
		}
	}

	// If it gets here, something is wrong
	panic("Found no winning node")
}
func (sn *FixedIdealSwarmNetwork) UpdateNetwork() {
	// Fixed network, no change
}

// Creates an array of nodes at their current state.
// Used for storing nodes data at each round.
func (sn *FixedIdealSwarmNetwork) GetNodeArray() *[]node {
	var nodes []node
	for _, v := range sn.nodes {
		nodes = append(nodes, *v)
	}
	return &nodes
}
func (sn *FixedIdealSwarmNetwork) GetNodeAdressMap() map[uint64]*node {
	return sn.nodeAddressMap
}
