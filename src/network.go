package main

import "math/rand"

type SwarmNetwork interface {
	CreateSwarmNetwork() *[]*neighbourhood
}

// type neighbourhood interface {

// }

type SimpleSwarmNetwork struct {
	networkNodeCount int
	neighbourhoods   []*neighbourhood
	maxRedundancy    int
}

func (sn *SimpleSwarmNetwork) CreateSwarmNetwork() *[]*neighbourhood {
	nbh := make([]*neighbourhood, 0, sn.networkNodeCount/2)

	nc := 1
	for i := 0; i < sn.networkNodeCount; i += nc {
		nc = rand.Intn(sn.maxRedundancy) + 1
		nbh = append(nbh, &neighbourhood{nodeCount: nc})
	}
	sn.neighbourhoods = nbh
	return &sn.neighbourhoods
}
