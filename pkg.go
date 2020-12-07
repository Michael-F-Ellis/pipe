// Package pipe provides funcs for building processing pipelines.
package pipe

// PipeLineFuncs must read from inchan until it closes
// and write to outchan, closing it when done. PipeLineFuncs
// should call c.Cancel if an unrecoverable error occurs
type Inlet func(chout PipeChan, c *Cancellation)
type Section func(chin, chout PipeChan, c *Cancellation)
type Outlet func(chin PipeChan, done chan (struct{}), c *Cancellation)

type Cancellation struct {
	cancelled bool
	err       error
}
type PipeChan chan (interface{})

func (c *Cancellation) Cancel(err error) {
	c.cancelled = true
	c.err = err
}

func (c *Cancellation) Cancelled() bool {
	return c.cancelled
}

// PipeLine runs a pipeline.
func PipeLine(inlet Inlet, outlet Outlet, sections ...Section) (err error) {
	var c Cancellation
	// Provide the inlet function with an output channel
	var chout = make(chan (interface{}), 1)
	// Launch it.
	go inlet(chout, &c)
	// Launch the internal pipeline sections
	for _, section := range sections {
		// Use the prior output channel as the input
		chin := chout
		// Make a new output channel
		chout = make(chan (interface{}), 1)
		// Launch the current section
		go section(chin, chout, &c)
	}
	// Make a "done" channel for the outlet
	done := make(chan (struct{}), 1)
	// Launch it with the preceding output channel as its input.
	go outlet(chout, done, &c)

	// Block on the done channel (which the outlet function
	// must close when all processing is complete.)
	<-done
	err = c.err
	return
}
