package main

import (
	"bytes"
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
	Nodes []*node
}

func logger(logchan chan *[]byte, done chan bool) {
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
	defer func() { done <- true }() // Tell main loop that writing is done.
	defer f.Close()
	defer fmt.Fprint(f, "\n]")
	closemsg := []byte("CLOSE")

	_, err = fmt.Fprint(f, "[\n{}")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		str := <-logchan
		if bytes.Equal(*str, closemsg) {
			break
		}
		_, err = fmt.Fprint(f, ",\n"+string(*str))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
