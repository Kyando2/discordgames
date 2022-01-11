package main

import (
	"math/rand"
	"testing"
)

func TestNewWorld(t *testing.T) {
	wrld, err := NewWorld()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(wrld.comp.items)
}

func BenchmarkAt(b *testing.B) {
	wrld, err := NewWorld()
	if err != nil {
		b.Errorf(err.Error())
	}
	for n := 0; n < b.N; n++ {
		wrld.At(rand.Intn(2000), rand.Intn(2000))
	}
}
