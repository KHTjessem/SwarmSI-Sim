package main

type neighbourhood struct {
	nodeCount int
	nodes     []*node
}

type node struct {
	Id       uint64
	Earnings float64
	stake    float64
}
