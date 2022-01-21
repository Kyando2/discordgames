package main

func (in *InThread) Answer(v int) Answer {
	return Answer{
		To:    in,
		Value: v,
	}
}

func (ch *Chain) Push(in *InThread) {
	ch.Q <- in
}
