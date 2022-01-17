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
	t.Log(wrld.userdata)
}

var atbenchresult int

func BenchmarkAt(b *testing.B) {
	wrld, err := NewWorld()
	b.ResetTimer()
	if err != nil {
		b.Errorf(err.Error())
	}
	var r int
	for n := 0; n < b.N; n++ {
		r = wrld.At(rand.Intn(2000), rand.Intn(2000))
	}
	atbenchresult = r
}
