package framework

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type simpleTabWriter struct {
	writer *tabwriter.Writer

	output   io.Writer
	minwidth int
	tabwidth int
	padding  int
	padchar  byte
	flags    uint
}

func newSimpleTabWriter(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) TabWriter {
	return &simpleTabWriter{
		writer: &tabwriter.Writer{},

		output:   output,
		minwidth: minwidth,
		tabwidth: tabwidth,
		padding:  padding,
		padchar:  padchar,
		flags:    flags,
	}
}

func (s *simpleTabWriter) Init() TabWriter {
	s.writer.Init(s.output, s.minwidth, s.tabwidth, s.padding, s.padchar, s.flags)
	return s
}

func (s *simpleTabWriter) AppendRow(colms ...string) {
	for _, c := range colms {
		fmt.Fprint(s.writer, c, "\t")
	}
	fmt.Fprint(s.writer, "\n")
}

func (s *simpleTabWriter) Flush() error {
	return s.writer.Flush()
}
