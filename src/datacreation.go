package main

import "fmt"

type logObject struct {
	// Round information
	Round       int
	TotalPayout int
	RoundPrice  int

	// Node information
	Nodes *[]node
}

type lgoCreator func(roundPrice int) *logObject

type storer interface {
	save(*simulator)
}

// Saves no results.
type saveNothing struct {
	finished chan bool
}

func (sn saveNothing) init() {
	sn.finished = make(chan bool, 1)
	sn.finished <- true
}
func (sn saveNothing) save(lgo *logObject) {
	// empty
}
func (sn saveNothing) mainloop() {
	// empty
}

// Saves all or X rounds to sqlite database.
// This can create a big database since if storing all round states
// of all nodes.
// The more nodes and rounds, the more space is required.
// set 'everyXRound' to 1 to store all rounds.
type saveFullStateSql struct {
	db          *Database
	dataChannel chan *logObject
	finished    chan bool

	runID       int
	everyXRound *int
}

// Does some initial config and starts the mainloop.
func (sar *saveFullStateSql) init() {
	sar.dataChannel = make(chan *logObject, 10000)
	sar.finished = make(chan bool, 1)

	// Start main loop
	go sar.mainloop()
}

// Sends logObject to the dataChannel where it is proccessed
// in the main loop.
func (sar saveFullStateSql) save(s *simulator) {
	if s.round%*sar.everyXRound != 0 {
		return
	}
	lgo := s.createRoundStat()
	sar.dataChannel <- lgo
}
func (sar saveFullStateSql) close() {
	closemsg := logObject{Round: -1}
	sar.dataChannel <- &closemsg
	<-sar.finished // Wait for DB to finish writing
}
func (sar *saveFullStateSql) mainloop() {
	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	// sar.db = db

	db.ExecuteScript("simSchema.sql")
	db.InsertNodes(NODECOUNT)

	// Add new run:
	sar.runID = db.insertNewRun(NODECOUNT, DESCRIPTION, SETUPSEED, SIMSEED)

	txCounter := 0
	tx, err := db.con.Begin()
	if err != nil {
		panic(err)
	}

	for {
		lgo := <-sar.dataChannel
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

		rid := sar.db.InsertnewRound(tx, sar.runID, lgo.Round, lgo.RoundPrice, lgo.TotalPayout)
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
		sar.db.InsertNodeRoundBulk(tx, &bulkstr)

		if txCounter >= 20000 {
			txCounter = -1
			tx.Commit()
		}

	}

	if txCounter >= 1 {
		tx.Commit()
	}

	// Done writing, close db
	db.Close()

	sar.finished <- true
}

// This only stores the gini for everyX round.
// Set to 1 to store all rounds, 10 to store every 10 rounds,
// 100 for every 100 rounds etc.
type saveGiniJson struct {
	dataChannel chan *logObject
	finished    chan bool
	everyX      int
}

func (sgj *saveGiniJson) init() {
	sgj.dataChannel = make(chan *logObject, 10000)
	sgj.finished = make(chan bool, 1)

	go sgj.mainloop()

}
func (sgj saveGiniJson) save(s *simulator) {
	lgo := s.createRoundStat()
	sgj.dataChannel <- lgo
}
func (sgj *saveGiniJson) mainloop() {

}

// This saves only the change each round.
// You will need to count from 0 up to the round you
// want the state of, like to a blockchain.
type saveOnlyChange struct {
	dataChannel chan *logObject
	finished    chan bool
}

func (soc *saveOnlyChange) init() {
	soc.dataChannel = make(chan *logObject, 10000)
	soc.finished = make(chan bool, 1)

	go soc.mainloop()

}
func (soc saveOnlyChange) save(s *simulator) {
	lgo := s.createRoundStat()
	soc.dataChannel <- lgo
}

func (soc saveOnlyChange) close() {
	closemsg := logObject{Round: -1}
	soc.dataChannel <- &closemsg
}

func (soc *saveOnlyChange) mainloop() {

}
