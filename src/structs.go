package main

type snode struct {
	id       int
	earnings float64
}

func (s *snode) AddEarnings(pot float64) {
	s.earnings += pot
}

// func (n *neighbourhood) SelectWinner() *snode {
// 	winner := rand.Intn(n.nodeCount)

// 	return n.swNodes[winner]
// }
