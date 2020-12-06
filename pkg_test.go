package pipe

import (
	"fmt"
	"testing"
)

type Item struct {
	s string
}

func TestPipeLine(t *testing.T) {
	PipeLine(inlet, outlet, section1, section2)
}

var inlet Inlet = func(outchan chan (interface{})) {
	for _, s := range []string{"foo", "bar"} {
		outchan <- Item{s: s}
	}
	close(outchan)
}
var section1 Section = func(inchan, outchan chan (interface{})) {
	for item := range inchan {
		it := item.(Item)
		it.s += " stage1"
		outchan <- it
		outchan <- it

	}
	close(outchan)
}
var section2 Section = func(inchan, outchan chan (interface{})) {
	for item := range inchan {
		it := item.(Item)
		it.s += " stage2"
		outchan <- it
	}
	close(outchan)
}

var outlet Outlet = func(inchan chan (interface{}), done chan (struct{})) {
	for item := range inchan {
		it := item.(Item)
		fmt.Println(it.s)
	}
	done <- struct{}{}
}
