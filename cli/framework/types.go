package framework

import (
	"github.com/ikaven1024/bolt-cli/db"

	"github.com/chzyer/readline"
)

type Framework interface {
	Run()
	DB() db.Interface

	TabWriter() TabWriter
	ReadLine() *readline.Instance

	EnterContext(context Context)
	ExitContext()

	Echof(string, ...interface{})
	EchoLine(...interface{})
	Error(...interface{})
}

type Context interface {
	Process([]string)
	Exit()
}

type TabWriter interface {
	Init() TabWriter
	AppendRow(colms ...string)
	Flush() error
}
