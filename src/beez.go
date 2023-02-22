package main

type neighbourhood struct {
	nodeCount int
	nodes     []*node
}

type node struct {
	id       uint64
	earnings float64
}
