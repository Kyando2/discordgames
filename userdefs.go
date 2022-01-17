package main

// Wrapper object around a userdata and id
type User struct {
	data *userData
	id   string
}

type Item struct {
	ID       int
	Quantity int
}

type Inventory struct {
	Data []Item
}

type userData struct {
	pos  position
	char character
	emp  empire
	inv  Inventory
	rank int
}
