package main

import (
	"fmt"
	"time"
)

// TODO: Consider moving constants to their own file.
const NODECOUNT = 32896
const ADDRESSLENGTH = 128

// with 15 minutes pr round, 350666 rounds is approx 10 years
const ROUNDS = 350000
const DBNAME = "nss.db"
const DESCRIPTION = "" +
	"Depth 16 - 128 add - 20 runs - Shifting node placement - Stake:50k"

var SETUPSEED int64 = 123123
var SIMSEED int64 = time.Now().Unix() // For random: time.Now().Unix()
var STAKESEED = 210210210

func main() {
	print("Hello Swarm!\n")

	// Storing results
	// saver := &saveNothing{}

	saver := saveFullStateSql{
		everyXRound: new(int),
	}
	*saver.everyXRound = 10000

	saver.init()

	// stake := EqualStake{
	// 	amount: 10,
	// }

	// stake := PowerDistStake{
	// 	alpha:      2.5,
	// 	minStake:   10,
	// 	rounding:   true,
	// 	roundBy:    10,
	// 	limitStake: false,
	// }

	stake := spreadStake{
		stake:    50000,
		nc:       new(int),
		operator: new(int),
	}
	*stake.operator = 1
	*stake.nc = 0

	//
	// NETWORKS
	//

	// swnet := &FixedIdealSwarmNetwork{
	// 	networkNodeCount:  NODECOUNT,
	// 	stakeDistribution: stake,
	// }

	// swnet := &KademSwarmTree{
	// 	addressLength:     ADDRESSLENGTH,
	// 	nodeCount:         NODECOUNT,
	// 	stakeDistribution: stake,
	// 	kademTree:         bintree{root: &binNode{prefix: ""}},
	// 	fullySaturate:     false,
	// 	addressBook:       make(map[uint64]*node),
	// 	kademAddress:      make(map[string]*node),
	// 	nodes:             make([]*node, 0, NODECOUNT),
	// }

	swnet := &KademSwarmTreeStorageDepth{
		addressLength:     ADDRESSLENGTH,
		nodeCount:         NODECOUNT,
		stakeDistribution: stake,
		kademTree:         bintree{root: &binNode{prefix: ""}},
		fullySaturate:     false,
		storageDepth:      8,
		addressBook:       make(map[uint64]*node),
		kademAddress:      make(map[string]*node),
		nodes:             make([]*node, 0, NODECOUNT),
	}

	//
	// SIMULATOR
	//
	s := &simulator{
		totalNodeCount: NODECOUNT,
		swarmnetwork:   swnet,
		rentoracle:     &FixedRentOracle{fixedPrice: 1},
		postage:        &simpleFixedPostage{},
		saver:          saver,

		round:          0,
		maxRounds:      ROUNDS,
		SetupSeed:      SETUPSEED,
		simulationSeed: SIMSEED,
	}

	s.Setup()
	start := time.Now()
	s.MainLoop()
	end := time.Now()

	ela := end.Sub(start)
	str := fmt.Sprintf("Done with simulation! Took %v", ela)
	str += "\nWating for all results to be written to db\n"
	print(str)

	saver.close()
	end = time.Now()

	ela = end.Sub(start)
	str = fmt.Sprintf("Done with everything, total runtime: %v", ela)
	print(str)
}
