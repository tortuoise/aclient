package aclient

import (
	"testing"
	"unicode"
)

func TestDisplay(t *testing.T) {

	//want := unicode.Tamil
	//Display(0, unicode.MaxASCII)

}

func TestCreate(t *testing.T) {

	err := Create(0, unicode.MaxASCII)

	if err != nil {
		t.Errorf("%v", err)
	}

}
