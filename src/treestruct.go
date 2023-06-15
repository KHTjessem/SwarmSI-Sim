package main

import "errors"

// Right is 0, left is 1.
type bintree struct {
	root *binNode
}
type binNode struct {
	depth  int
	prefix string //

	parent *binNode
	right  *binNode
	left   *binNode

	address string
	leaf    *node

	allNodesBelow   int
	allNodeBelowArr []*node
}

func (b *bintree) InsertNode(nod *node, address string) error {
	const zero = "0"

	actNode := b.root
	for i := range address {
		bin := address[i]

		// Move down tree, add node if nil
		if bin == zero[0] { // right
			if actNode.right == nil {
				actNode.right = &binNode{parent: actNode,
					prefix: address[:i+1],
					depth:  i + 1}
				actNode = actNode.right
			} else {
				actNode = actNode.right
			}
		} else { // left
			if actNode.left == nil {
				actNode.left = &binNode{parent: actNode,
					prefix: address[:i+1],
					depth:  i + 1}
				actNode = actNode.left
			} else {
				actNode = actNode.left
			}
		}

	}

	if actNode.leaf != nil {
		return errors.New("address occupied")
	}

	actNode.leaf = nod
	actNode.address = address
	actNode.allNodeBelowArr = append(actNode.allNodeBelowArr, nod)
	actNode.allNodesBelow++

	// Update parents
	parent := actNode.parent
	for parent != nil {
		parent.allNodesBelow++
		parent.allNodeBelowArr = append(parent.allNodeBelowArr, nod)
		parent = parent.parent
	}
	return nil
}

// Navigates as far as possible to given prefix/address
func (b *bintree) navigate(pref string) *binNode {
	const zero = "0"
	actNode := b.root
	for i := range pref {
		bin := pref[i]

		if bin == zero[0] { // right
			if actNode.right == nil {
				return actNode
			} else {
				actNode = actNode.right
			}
		} else { // left
			if actNode.left == nil {
				return actNode
			} else {
				actNode = actNode.left
			}
		}
	}
	return actNode

}

// Navigates as far as possible to given prefix/address
func (b *bintree) navigateWithStop(pref string, stopdepth int) *binNode {
	const zero = "0"
	actNode := b.root
	for i := range pref {
		bin := pref[i]

		if bin == zero[0] { // right
			if actNode.right == nil {
				return actNode
			} else {
				actNode = actNode.right
			}
		} else { // left
			if actNode.left == nil {
				return actNode
			} else {
				actNode = actNode.left
			}
		}
		if actNode.depth == stopdepth {
			return actNode
		}
	}
	return actNode

}

func (b *bintree) FindClosestNodes(address string) []*node {
	const zero = "0"
	// First navigate as far as we can with given address.
	actNode := b.navigate(address)

	nodes := make([]*node, 0, 10)
	if actNode.depth == len(address) { // have leaf node
		nodes = append(nodes, actNode.leaf)
	}

	// Move up tree and add nodes
	for len(nodes) < 4 {
		actNode = actNode.parent
		if address[actNode.depth] == zero[0] {
			if actNode.left != nil {
				nodes = append(nodes, actNode.left.allNodeBelowArr...)
			}
		} else {
			if actNode.right != nil {
				nodes = append(nodes, actNode.right.allNodeBelowArr...)
			}
		}
	}

	return nodes
}
