package bosh

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

func (b *Bosh) AssertEquals(actual interface{}, expected interface{}) bool {
	if actual != expected {
		return fail(fmt.Sprintf("Value expected to be %v but was %v", expected, actual))
	}

	return true
}

func (b *Bosh) AssertNotEquals(actual interface{}, expected interface{}) bool {
	if actual == expected {
		return fail(fmt.Sprintf("Value expected to be %v but was %v", expected, actual))
	}

	return true
}

func (b *Bosh) AssertTrue(value bool) bool {
	if value != true {
		return fail(fmt.Sprintf("Value expected to be true but was %v", value))
	}

	return true
}

func (b *Bosh) AssertFalse(value bool) bool {
	if value == true {
		return fail(fmt.Sprintf("Value expected to be false but was %v", value))
	}

	return true
}


func (b *Bosh) AssertNil(value interface{}) bool {
	if value != nil {
		return fail(fmt.Sprintf("Value expected to be nil but was %v", value))
	}

	return true
}

func (b *Bosh) AssertNotNil(value interface{}) bool {
	if value != nil {
		return fail(fmt.Sprintf("Value expected to be %v but was nil", value))
	}

	return true
}

func fail(message string) bool {
	log.Printf(color.RedString("[ASSERTION ERROR] " + message))

	return false
}