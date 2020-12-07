package tbchrom

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	notes, err := Parse("0 2 4 5")
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("%+v", notes)
}
