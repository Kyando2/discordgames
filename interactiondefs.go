package main

// Compendium type object that
// contains a list of choices, which are
// a string value and an ID for it to be
// identified. It also contains
//  a 'prompt' to be sent to the user,
// and a potential image to add to the prompt.
type Interact struct {
	Choices []struct {
		Value string
	}
	Value string
	Loc   string
}

// Represents a thread of interactions
// helps track when multiple interactions
// are going on concurrently.
type InThread struct {
	In *Interact
	ID int
}

// Represents an answer to a thread of interactions
// The Value is the index of the selected Choice of the
// original Interact Choices slice.
type Answer struct {
	To    *InThread
	Value int
	ID    int
}

// Array of interactions to process.
type Chain struct {
	Q chan *InThread
}

// Compendium type object of
// all possible interactions.
var possibilities map[int]Interact
var threadcounter int

func init() {
	threadcounter = 0
	// handle generating possibilities from file
}
