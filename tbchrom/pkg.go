// Package tbchrom implements a very simplified version of tbon
// for use in exercise generation.
package tbchrom

import (
	"strings"

	"github.com/Michael-F-Ellis/pipe"
)

// zero indexed chromatic pitch numbers.
var pitchTokens = map[rune]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'X': 10, // mnemonic: Roman numeral X = 10
	'Y': 11,
}

// midiNotes structs get passed down the pipeline to be filled in.
type midiNote struct {
	barText        string
	beatText       string
	subbeatText    string
	barBeats       int
	beatNum        int
	subBeats       int
	subBeatNum     int
	midiPitch      int
	priorMidiPitch int
	velocity       int
	startTicks     int
	duration       int
}

func Parse(fullText string) (notes []midiNote, err error) {
	// For now, treat fullText as one bar.
	var inlet pipe.Inlet = func(chout pipe.PipeChan, c *pipe.Cancellation) {
		defer close(chout)
		var beatTexts []string
		var beatCount int
		for i, beat := range strings.Fields(fullText) {
			if c.Cancelled() {
				return
			}
			beatTexts = append(beatTexts, beat)
			beatCount = i
		}
		for i, s := range beatTexts {
			chout <- midiNote{
				barText:  fullText,
				beatText: s,
				barBeats: beatCount + 1,
				beatNum:  i,
			}
		}
	}
	var outlet pipe.Outlet = func(chin pipe.PipeChan, done chan (struct{}), c *pipe.Cancellation) {
		for item := range chin {
			if c.Cancelled() {
				return
			}
			note := item.(midiNote)
			notes = append(notes, note)
		}
	}
	err = pipe.PipeLine(inlet, outlet)
	return
}
