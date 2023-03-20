package main

import "fmt"

// type logger struct {
// }

type logObject struct {
	// Round information
	Round       int
	TotalPayout int
	RoundPrice  int

	// Node information
	Nodes *[]node
}

func logger(logchan chan *logObject, done chan bool, nodeCount int, description string) {
	defer func() { done <- true }() // Tell main loop that writing is done.

	// sqlite
	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.ExecuteScript("simSchema.sql")
	db.InsertNodes(NODECOUNT)

	// Add new run:
	runID := db.insertNewRun(nodeCount, description, SETUPSEED, SIMSEED)

	txCounter := 0
	tx, err := db.con.Begin()
	if err != nil {
		panic(err)
	}

	for {
		lgo := <-logchan
		if lgo.Round == -1 {
			// -1 means it is done
			break
		}

		if txCounter == -1 {
			txCounter += 1
			tx, err = db.con.Begin()
			if err != nil {
				panic(err)
			}
		}

		rid := db.InsertnewRound(tx, runID, lgo.Round, lgo.RoundPrice, lgo.TotalPayout)
		txCounter += 1

		nodes := *lgo.Nodes

		actNode := &nodes[0]
		bulkstr := fmt.Sprintf("(%v, %v, %v, %v)",
			actNode.Id, rid, actNode.Earnings, actNode.stake)
		for i := 1; i < len(nodes); i++ {
			actNode = &nodes[i]
			bulkstr += fmt.Sprintf(",(%v, %v, %v, %v)",
				actNode.Id, rid, actNode.Earnings, actNode.stake)

			txCounter += 1
		}
		// db.InsertNodeRound(tx, rid, &n)
		db.InsertNodeRoundBulk(tx, &bulkstr)

		if txCounter >= 20000 {
			txCounter = -1
			tx.Commit()
		}

	}

	if txCounter >= 1 {
		tx.Commit()
	}

	// print("Applying index")
	// db.ExecuteScript("indexes.sql")
}
