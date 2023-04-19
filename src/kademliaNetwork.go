package main

import (
	"math/big"
	"math/rand"
	"sort"
)

func randomBitString(amount int) string {
	const bits = "01"
	bin := make([]byte, amount)
	for i := range bin {
		bin[i] = bits[rand.Intn(len(bits))]
	}
	return string(bin)
}

type KademSwarmNetArr struct {
	addressLength     int
	nodeCount         int
	stakeDistribution StakeCreator
	// kademTree         bintree

	// addressBook is mapping nodeID to the node
	addressBook map[uint64]*node // TODO: is it needed for this imp=?
	// Need extra book for kademlia address
	kademAddress map[string]*node
	kadem2Indx   map[string]int
	nodes        []*node
	flip         int // Index where first bit is 1
}

func (ksn *KademSwarmNetArr) CreateSwarmNetwork() {
	ksn.flip = -1

	// Create nodes and insert into data structures.
	for i := 0; i < ksn.nodeCount; i++ {
		n := &node{Id: uint64(i), stake: ksn.stakeDistribution.GetStake(i)}
		ksn.addressBook[uint64(i)] = n // TODO: might not be used, del.
		ksn.nodes = append(ksn.nodes, n)

		// Kademlia address.
		nAdd := randomBitString(ksn.addressLength)
		for j := 0; j < 100; j++ {
			//for { // Avoid infinite loop
			if _, ok := ksn.kademAddress[nAdd]; !ok {
				break
			}
			nAdd = randomBitString(ksn.addressLength)
		}
		n.address = nAdd
		ksn.kademAddress[nAdd] = n
	}

	// Sort "nodes" array based on kademlia address.
	// Done by converting address to bigint
	sort.Slice(ksn.nodes, func(i, j int) bool {
		n1 := new(big.Int)
		n1.SetString(ksn.nodes[i].address, 2)
		n2 := new(big.Int)
		n2.SetString(ksn.nodes[j].address, 2)
		return n1.Cmp(n2) < 0
	})

	// Map of kadem address to index in nodes.
	for i, v := range ksn.nodes {
		ksn.kadem2Indx[v.address] = i
		if v.address[0] == 49 { // checks if address starts with 1 (49 is ASCII for 1).
			if ksn.flip == -1 {
				ksn.flip = i // This is where the first bit goes from 0 to 1 in address.
			}
		}
	}

}

// Need to sort ksn.nodes if any changes to nodes.
func (ksn *KademSwarmNetArr) UpdateNetwork() {
	// TODO: implement
}

func (ksn *KademSwarmNetArr) SelectNeighbourhood() *neighbourhood {
	anc := randomBitString(ksn.addressLength)

	indx := 0
	// find the closest nodes to anchor that forms a neighbourhood of ATLEAST 4 nodes
	if anc[0] == 49 && ksn.flip != -1 { // is first bit 1? (49 ASCII is 1)
		indx = ksn.flip
	}

	counter := make([][]*node, ksn.addressLength+1)
	// for i := 0; i < len(counter); i++ {
	// 	counter[i] = make([]*node, 0, 4)
	// }

	for i := indx; i < len(ksn.nodes); i++ {
		nod := ksn.nodes[i]

		prefLen := 0
		for j := 0; j < len(nod.address); j++ {
			if nod.address[j] == anc[j] {
				prefLen++
			} else {
				break
			}
		}

		counter[prefLen] = append(counter[prefLen], nod)
	}

	ne := neighbourhood{nodes: make([]*node, 0, 8)}

	for i := len(counter) - 1; i >= 0; i-- {
		if counter[i] == nil {
			continue
		}
		ne.nodes = append(ne.nodes, counter[i]...)
		if len(ne.nodes) >= 4 {
			ne.nodeCount = len(ne.nodes)
			break
		}
	}

	return &ne
}

func (ksn *KademSwarmNetArr) SelectWinner() *node {
	nbhood := ksn.SelectNeighbourhood()

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

// creates a "copy" of the nodes at the time the method is ran.
// Used for storing the data for later use.
func (ksn *KademSwarmNetArr) GetNodeArray() *[]node {
	var nodes []node
	for _, v := range ksn.nodes {
		nodes = append(nodes, *v)
	}
	return &nodes
}

func (ksn *KademSwarmNetArr) GetNodeAdressMap() map[uint64]*node {
	return ksn.addressBook
}

type KademSwarmTree struct {
	addressLength     int
	nodeCount         int
	stakeDistribution StakeCreator
	kademTree         bintree

	// addressBook is mapping nodeID to the node
	addressBook map[uint64]*node // TODO: is it needed for this imp=?
	// Need extra book for kademlia address
	kademAddress map[string]*node
	nodes        []*node
}

func (kdst *KademSwarmTree) CreateSwarmNetwork() {
	for i := 0; i < kdst.nodeCount; i++ {
		// Create node
		n := &node{Id: uint64(i), stake: kdst.stakeDistribution.GetStake(i)}

		//Create Kademlia address.
		nAdd := randomBitString(kdst.addressLength)
		for j := 0; j < 100; j++ {
			//for { // Avoid infinite loop
			if _, ok := kdst.kademAddress[nAdd]; !ok {
				break
			}
			nAdd = randomBitString(kdst.addressLength)
		}
		n.address = nAdd

		// Add node to data structures
		kdst.nodes = append(kdst.nodes, n)
		kdst.addressBook[uint64(i)] = n // TODO: might not be used, del.

		kdst.kademAddress[nAdd] = n
		kdst.kademTree.InsertNode(n, n.address)
	}
}
func (kdst KademSwarmTree) UpdateNetwork() {
	return
}

func (kdst KademSwarmTree) SelectNeighbourhood() *neighbourhood {
	anch := randomBitString(kdst.addressLength)
	nodes := kdst.kademTree.FindClosestNodes(anch)
	nei := neighbourhood{nodes: nodes, nodeCount: len(nodes)}

	return &nei
}

func (kdst KademSwarmTree) SelectWinner() *node {
	nbhood := kdst.SelectNeighbourhood()

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
