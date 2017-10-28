package framework

import (
	"bytes"
	"sort"
	"strings"
)

type PositionalArgs func(cmd *Command, args []string) bool

type Command struct {
	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// Expected arguments
	Args PositionalArgs
	// Description of command
	Description string

	Action func([]string) error
}

func (c *Command) matchCommand(cmd string) bool {
	cmd = strings.ToLower(cmd)
	if strings.ToLower(c.Name) == cmd {
		return true
	}
	for _, aliase := range c.Aliases {
		if strings.ToLower(aliase) == cmd {
			return true
		}
	}
	return false
}

type commandsByName []Command

var _ sort.Interface = &commandsByName{}

func (b commandsByName) Len() int           { return len(b) }
func (b commandsByName) Less(i, j int) bool { return b[i].Name < b[j].Name }
func (b commandsByName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func CommandAndArgs(line []byte) (parts []string) {
	var startQuote byte
	var index int
	var splits [][]byte

	for {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			return
		}

		switch line[0] {
		case '"', '\'':
			startQuote = line[0]
			line = line[1:]
			index = bytes.IndexByte(line, startQuote)
			if index < 0 {
				// lack of end quote, we tolerate this mistake and regard the rest part as one.
				parts = append(parts, string(line))
				return
			} else if index == len(line)-1 {
				parts = append(parts, string(line[:index]))
				return
			}
			parts = append(parts, string(line[:index]))
			line = line[index+1:]
		default:
			splits = bytes.SplitN(line, []byte(" "), 2)
			parts = append(parts, string(splits[0]))
			if len(splits) == 1 {
				return
			}
			line = splits[1]
		}
	}
}
