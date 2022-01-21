package main

import "log"

// Starts a goroutine to do work.
// The goroutine will first receive
// from the chain. Then, it will call
// ''f'' with the value received from
// the chain. Then, it will check if
// the value returned by ''f'' warrant
// the execution of the goroutine
// stopping, and if not it will
// resume execution.
func (ch *Chain) Work(f func(*InThread, ...interface{}) (bool, error), args ...interface{}) {
	go func(ch chan *InThread) {
		done := false
		var v *InThread
		var err error
		for !done {
			v = <-ch
			done, err = f(v, args...)
			if err != nil {
				log.Fatalf("Error running worker: %s", err)
			}
		}
	}(ch.Q)
}

// Will create a ThreadIn with a reference
// to the Interact found at the provided
// ''id'' index in the possibilities map.
// Increments the threadcounter and uses
// the incremented value as an InThread ID
// and return the ID.
func (ch *Chain) Event(id int) int {
	// Generate InThread
	threadcounter++
	s := possibilities[id]
	th := InThread{
		In: &s,
		ID: threadcounter,
	}
	ch.Push(&th)
	return threadcounter
}

// Creates a new Chain
// ''Work'' may be called on it
// when it is ready to ~~go~~!
func NewChain() Chain {
	return Chain{
		Q: make(chan *InThread),
	}
}
