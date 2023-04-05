package main

import (
	"fmt"
	"time"
)

// TODO: Consider moving constants to their own file.
const NODECOUNT = 2080
const ADDRESSLENGTH = 128
const DBNAME = "simRes2.db"
const DESCRIPTION = "Simulation of 2080 nodes," +
	"static Kademlia network:128 bit address - Bucket stake distribution, stake is 43300"

var SETUPSEED int64 = 123123
var SIMSEED int64 = 123123 // For random: time.Now().Unix()

func main() {
	print("Hello Swarm!\n")

	// Logger channels
	logchan := make(chan *logObject, 100000)
	logStopped := make(chan bool)
	// Start logger, TODO: maybe get description from command
	// line argument/user input
	go logger(logchan, logStopped, NODECOUNT, DESCRIPTION)

	// stake := PowerDistStake{
	// 	alpha:      2.5,
	// 	minStake:   100,
	// 	rounding:   true,
	// 	roundBy:    100,
	// 	limitStake: false,
	// }

	stake := bucketSumStake{
		stake:  43300,
		nc:     new(int),
		bucket: new(int),
	}
	*stake.bucket = 1
	*stake.nc = 0

	//
	// NETWORKS
	//

	// swnet := &FixedIdealSwarmNetwork{
	// 	networkNodeCount:  NODECOUNT,
	// 	stakeDistribution: stake,
	// }

	swnet := &KademSwarmNet{
		addressLength:     ADDRESSLENGTH,
		nodeCount:         NODECOUNT,
		stakeDistribution: stake,
		addressBook:       make(map[uint64]*node),
		kademAddress:      make(map[string]*node),
		kadem2Indx:        make(map[string]int),
		nodes:             make([]*node, 0, NODECOUNT),
	}

	//
	// SIMULATOR
	//
	s := &simulator{
		totalNodeCount: NODECOUNT,
		swarmnetwork:   swnet,
		rentoracle:     &FixedRentOracle{fixedPrice: 199},
		postage:        &simpleFixedPostage{},
		logChan:        logchan,

		round:          0,
		maxRounds:      120000,
		SetupSeed:      SETUPSEED,
		simulationSeed: SIMSEED,
		// with 15 minutes pr round, 350666 rounds is approx 10 years
	}

	s.Setup()
	start := time.Now()
	s.MainLoop()
	end := time.Now()

	ela := end.Sub(start)
	str := fmt.Sprintf("Done with simulation! Took %v", ela)
	str += "\nWating for all results to be written to db\n"
	print(str)

	closemsg := logObject{Round: -1}
	s.logChan <- &closemsg
	<-logStopped // Wait for DB to finish writing
	end = time.Now()

	ela = end.Sub(start)
	str = fmt.Sprintf("Done with everything, total runtime: %v", ela)
	print(str)
}
