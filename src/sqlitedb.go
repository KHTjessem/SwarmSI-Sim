package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	con *sql.DB
}

func NewDatabase() (*Database, error) {
	dbpath := filepath.Join(".", "data", DBNAME)
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		fmt.Println("No database file was found, creating a new one")
	}

	conn, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	return &Database{conn}, nil
}

func (db *Database) Close() error {
	return db.con.Close()
}

// Executes "filename" sql script located in "sql-script" folder
func (db *Database) ExecuteScript(filename string) error {
	file := filepath.Join(".", "sql-script", filename)
	sqlScript, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	tx, err := db.con.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(sqlScript)); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Inserts a new run into database.
// Returns runID for the current simulation
func (db *Database) insertNewRun(nodeCount int, runDescription string, setupSeed int64, simseed int64) int {
	q, err := db.con.Prepare("INSERT INTO run(nodeCount, runDesc, setupSeed, simulationSeed) VALUES(?,?,?,?)")
	if err != nil {
		panic(err)
	}

	res, err := q.Exec(nodeCount, runDescription, setupSeed, simseed)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(id)
}

// Inserts
func (db *Database) InsertnewRound(tx *sql.Tx, runID int, roundNumber int, roundPayout int, totalPayout int) int {
	q := "INSERT INTO rounds(runID, round, roundPayout, totalPayout) VALUES(?,?,?,?)"

	res, err := tx.Exec(q, runID, roundNumber, roundPayout, totalPayout)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(id)
}

func (db *Database) InsertNodeRound(tx *sql.Tx, roundID int, n *node) {
	q := "INSERT INTO nround(nodeID, roundID, earnings, stake) VALUES(?,?,?,?)"
	_, err := tx.Exec(q, n.Id, roundID, n.Earnings, n.stake)
	if err != nil {
		panic(err)
	}
}

func (db *Database) InsertNodeRoundBulk(tx *sql.Tx, bulk *string) {
	q := "INSERT INTO nround(nodeID, roundID, earnings, stake) VALUES "
	q = q + *bulk
	_, err := tx.Exec(q)
	if err != nil {
		panic(err)
	}
}

func (db *Database) InsertNodes(nodecount int) {
	existing := new(int)
	ncq := "SELECT count(nodeID) as nodeC FROM node"
	res, err := db.con.Query(ncq)
	if err != nil {
		print(err.Error())
	}
	res.Next()
	err = res.Scan(existing)
	if err != nil {
		print(err.Error())
	}
	print(*existing)

	newNodes := nodecount - *existing

	if newNodes <= 0 {
		return
	}
	tx, err := db.con.Begin()
	if err != nil {
		println(err.Error())
	}
	for i := *existing; i < newNodes+*existing; i++ {
		q := "INSERT INTO node(nodeID) VALUES(?)"
		tx.Exec(q, i)
	}
	tx.Commit()
}
