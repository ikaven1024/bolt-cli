package app

import (
	"flag"
	"fmt"
)

// Flag is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recommended that
// this interface be implemented.
type Flag interface {
	fmt.Stringer
	// Apply Flag settings to the given flag set
	Apply(*flag.FlagSet)
	GetName() string
}

type baseFlag struct {
	Name        string
	Usage       string
	Value       string
	Destination *string
}

type StringFlag struct {
}
