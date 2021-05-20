package pipe

import (
	"fmt"
	"testing"
)

type Item struct {
	s string
}

func TestPipeLine(t *testing.T) {
	foo := Foo{2, true}
	err := PipeLine(foo, inlet, outlet, section1, section2)
	if err != nil {
		t.Errorf("%v", err)
	}
}

type Foo struct {
	x int
	y bool
}

func TestPipeLineError(t *testing.T) {
	var streamingInlet Inlet = func(chout PipeChan, c *Cancellation) {
		defer close(chout)
		for {
			if c.Cancelled() {
				return
			}
			chout <- Item{s: "blah"}
		}
	}
	var badSection Section = func(chin, chout PipeChan, c *Cancellation) {
		c.Cancel(fmt.Errorf("I bailed."))
		close(chout)
	}
	foo := Foo{2, true}
	err := PipeLine(foo, streamingInlet, outlet, section1, badSection, section2)
	if err == nil {
		t.Error("Expected failure but the error was nil")
	}
	exp := "I bailed."
	got := err.Error()
	if got != exp {
		t.Errorf(`exp: "%v", got "%v"`, exp, got)
	}
}

var inlet Inlet = func(outchan PipeChan, c *Cancellation) {
	defer close(outchan)
	for _, s := range []string{"foo", "bar"} {
		if c.Cancelled() {
			return
		}
		outchan <- Item{s: s}
	}
}
var section1 Section = func(inchan, outchan PipeChan, c *Cancellation) {
	defer close(outchan)
	ctx := c.Context()
	for item := range inchan {
		if c.Cancelled() {
			return
		}
		it := item.(Item)
		it.s += " stage1" + fmt.Sprintf("%v", ctx.(Foo))
		outchan <- it
		outchan <- it

	}
}
var section2 Section = func(inchan, outchan PipeChan, c *Cancellation) {
	defer close(outchan)
	for item := range inchan {
		if c.Cancelled() {
			return
		}
		it := item.(Item)
		it.s += " stage2"
		outchan <- it
	}
}

var outlet Outlet = func(inchan PipeChan, done chan (struct{}), c *Cancellation) {
	for item := range inchan {
		if c.Cancelled() {
			return
		}
		it := item.(Item)
		fmt.Println(it.s)
	}
	done <- struct{}{}
}
