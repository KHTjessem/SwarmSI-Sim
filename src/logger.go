package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// type logger struct {
// }

type logObject struct {
	// Round information
	Round       int
	TotalPayout float64
	RoundPrice  float64

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
	db.ExecuteScript("simSchema.sql")
	defer db.Close()

	// Add new run:
	runID := db.insertNewRun(nodeCount, description)

	// json file
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)

	name := fmt.Sprintf("%v.json", time.Now().Format("2006-01-02-15.04.05"))
	f, err := os.Create(parent + "/Results/" + name)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	defer f.Close()
	defer fmt.Fprint(f, "\n]")

	_, err = fmt.Fprint(f, "[\n{}")
	if err != nil {
		fmt.Println(err)
		return
	}

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

		// sql
		// TODO: Change payouts and roundprice to just use int instead of float.
		rid := db.InsertnewRound(tx, runID, lgo.Round, int(lgo.RoundPrice), int(lgo.TotalPayout))
		txCounter += 1

		for _, n := range *lgo.Nodes {
			db.InsertNodeRound(tx, rid, &n)
			txCounter += 1
		}

		if txCounter >= 50000 {
			txCounter = -1
			tx.Commit()
		}

		// json
		// js, err := json.Marshal(lgo)
		// // make it json for writing to json file.
		// if err != nil {
		// 	println(err)
		// }

		// _, err = fmt.Fprint(f, ",\n"+string(js))
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}

	if txCounter >= 1 {
		tx.Commit()
	}
}
