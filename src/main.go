package main

import (
	"fmt"
	"time"
)

const NODECOUNT = 2048
const DBNAME = "simRes2.db"
const DESCRIPTION = "Simulation of 2048 nodes," +
	"static network - same stake, using weighted randomness"

var SETUPSEED int64 = 123123
var SIMSEED int64 = 123123 // For random: time.Now().Unix()

func main() {
	print("Hello Swarm!")

	// Logger channels
	logchan := make(chan *logObject, 100000)
	logStopped := make(chan bool)
	// Start logger, TODO: maybe get description from command
	// line argument/user input
	go logger(logchan, logStopped, NODECOUNT, DESCRIPTION)

	s := &simulator{
		totalNodeCount: NODECOUNT,
		swarmnetwork: &FixedIdealSwarmNetwork{
			networkNodeCount:  NODECOUNT,
			stakeDistribution: EqualStake{amount: 100}},
		rentoracle: &FixedRentOracle{fixedPrice: 199},
		postage:    &simpleFixedPostage{},
		logChan:    logchan,

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
	str += "\nWating for all results to be written to db"
	print(str)

	closemsg := logObject{Round: -1}
	s.logChan <- &closemsg
	<-logStopped // Wait for DB to finish writing

}
