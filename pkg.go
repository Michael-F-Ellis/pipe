// Package pipe provides funcs for building processing pipelines.
package pipe

// PipeLineFuncs must read from inchan until it closes
// and write to outchan, closing it when done. PipeLineFuncs
// should panic on errors that can't be fixed within
// the func.
type Inlet func(chout chan (interface{}))
type Section func(chin, chout chan (interface{}))
type Outlet func(chin chan (interface{}), done chan (struct{}))

// PipeLine runs a pipeline.
func PipeLine(inlet Inlet, outlet Outlet, sections ...Section) {
	var chout = make(chan (interface{}), 1)
	go inlet(chout)
	for _, itemFunc := range sections {
		chin := chout
		chout = make(chan (interface{}), 1)
		go itemFunc(chin, chout)
	}
	done := make(chan (struct{}), 1)
	go outlet(chout, done)
	<-done
}
