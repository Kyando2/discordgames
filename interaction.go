package main

func (in *Interact) Answer(v int) Answer {
	return Answer{
		To:    in,
		Value: v,
	}
}

func (ch *Chain) Push(in *InThread) {
	ch.Q <- in
}

// Probably require a discordgo session
func (ch *Chain) Work() {
	go func(ch chan *InThread) {
		// tohandle := <-ch
		// handle `tohandle`
		// using discord
	}(ch.Q)
}
