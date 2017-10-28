package framework

import (
	"bytes"
	"strings"
)

type processor struct {
	Selector func([]string) bool
	Process  func([]string)
}

type BaseContext struct {
	Framework

	Commands     []Command
	HelpCommands []string
}

func NewBaseContext(fw Framework) *BaseContext {
	return &BaseContext{
		Framework:    fw,
		HelpCommands: []string{"help", "h"},
	}
}

func (c *BaseContext) AddCommand(command ...Command) {

	c.Commands = append(c.Commands, command...)
}

func (c *BaseContext) Process(cmdArgs []string) {
	defer func() {
		if r := recover(); r != nil {
			c.Error(r)
		}
	}()

	if len(cmdArgs) == 0 {
		return
	}

	if c.matchHelpCommand(cmdArgs[0]) {
		c.usage()
		return
	}

	var command *Command
	for i := range c.Commands {
		if c.Commands[i].matchCommand(cmdArgs[0]) {
			command = &c.Commands[i]
			break
		}
	}

	if command == nil {
		c.badCommandInformer()
		return
	}

	args := []string{}
	if len(cmdArgs) > 1 {
		args = cmdArgs[1:]
	}

	if !command.Args(command, args) {
		c.badCommandInformer()
		return
	}

	if err := command.Action(args); err != nil {
		c.Error(err)
	}
}

func (c *BaseContext) usage() {
	tw := c.TabWriter().Init()
	defer tw.Flush()
	for _, cmd := range c.Commands {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(strings.Join(append([]string{cmd.Name}, cmd.Aliases...), "/"))
		if len(cmd.Usage) > 0 {
			buf.WriteString(" ")
			buf.WriteString(cmd.Usage)
		}
		buf.WriteString(":")

		tw.AppendRow(buf.String(), cmd.Description)
	}
}

func (c *BaseContext) matchHelpCommand(cmd string) bool {
	cmd = strings.ToLower(cmd)
	for i := range c.HelpCommands {
		if cmd == strings.ToLower(c.HelpCommands[i]) {
			return true
		}
	}
	return false
}

func (c *BaseContext) badCommandInformer() {
	c.Error("Bad command, see help")
}
