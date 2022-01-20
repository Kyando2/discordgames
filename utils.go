package main

func contains(l []int, s int) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}

type HasID interface {
	getID() int
}

type IDArr interface {
	getIDs() []int
}
