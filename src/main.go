package main

import (
	"fmt"
	"time"
)

func main() {
	print("Hello Swarm!")

	// Database for storing simulation steps
	// db, err := NewDatabase()
	// if err != nil {
	// 	panic(err)
	// }

	// db.ExecuteScript("simSchema.sql")

	s := &simulator{
		totalNodeCount: 16,
		swarmnetwork:   &FixedIdealSwarmNetwork{networkNodeCount: 2048},
		rentoracle:     &FixedRentOracle{fixedPrice: 199},
		postage:        &simpleFixedPostage{},
		logChan:        make(chan *logObject, 100000),
		//db:             db,
		round:        0,
		maxRounds:    350666,
		mathRandSeed: 123123,
		// with 15 minutes pr round, 350666 rounds is approx 10 years
	}

	s.Setup()
	start := time.Now()
	s.MainLoop()
	end := time.Now()

	ela := end.Sub(start)

	str := fmt.Sprintf("Done! Took %v", ela)
	print(str)
}
