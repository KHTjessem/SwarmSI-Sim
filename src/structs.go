package main

import "math/rand"

type snode struct {
	id       int
	earnings float64
}

func (s *snode) AddEarnings(pot float64) {
	s.earnings += pot
}

type neighbourhood struct {
	nodeCount int
	swNodes   []*snode
}

func (n *neighbourhood) SelectWinner() *snode {
	winner := rand.Intn(n.nodeCount)

	return n.swNodes[winner]
}
