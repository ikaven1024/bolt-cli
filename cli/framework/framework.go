package framework

import (
	"fmt"
	"io"

	"github.com/ikaven1024/bolt-cli/cli/util/color_str"
	"github.com/ikaven1024/bolt-cli/db"

	"github.com/chzyer/readline"
)

type frameworkImpl struct {
	contextStack []Context
	quitCh       chan struct{}
	readline     *readline.Instance
	writer       io.Writer
	errWriter    io.Writer
	tabWriter    TabWriter
	db           db.Interface
}

var _ Framework = &frameworkImpl{}

func New(db db.Interface, stdout, stderr io.Writer) Framework {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          color_str.Red("»"),
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}

	return &frameworkImpl{
		db:        db,
		readline:  rl,
		writer:    stdout,
		errWriter: stderr,
		tabWriter: newSimpleTabWriter(stdout, 10, 4, 2, ' ', 10),

		quitCh: make(chan struct{}),
	}
}

func (f *frameworkImpl) Run() {
	defer f.readline.Close()
	defer f.tabWriter.Flush()

	for f.loop() {
		select {
		case <-f.quitCh:
			return
		default:
			break
		}
	}
}

func (f *frameworkImpl) DB() db.Interface {
	return f.db
}

func (f *frameworkImpl) TabWriter() TabWriter {
	return f.tabWriter
}

func (f *frameworkImpl) ReadLine() *readline.Instance {
	return f.readline
}

// call before process.
func (f *frameworkImpl) EnterContext(context Context) {
	f.contextStack = append(f.contextStack, context)
}

func (f *frameworkImpl) ExitContext() {
	f.currContext().Exit()
	f.contextStack = f.contextStack[:len(f.contextStack)-1]
	if len(f.contextStack) == 0 {
		f.quit()
	}
}

func (f *frameworkImpl) Echof(format string, a ...interface{}) {
	fmt.Fprintf(f.writer, format, a...)
}

func (f *frameworkImpl) EchoLine(a ...interface{}) {
	fmt.Fprintln(f.writer, a...)
}

func (f *frameworkImpl) Error(a ...interface{}) {
	fmt.Fprintln(f.errWriter, a...)
}

func (f *frameworkImpl) loop() bool {
	line, err := f.readline.Readline()
	if err == readline.ErrInterrupt {
		if len(line) == 0 {
			return false
		} else {
			return true
		}
	} else if err == io.EOF {
		return false
	}

	cmd := CommandAndArgs([]byte(line))
	if len(cmd) > 0 {
		f.currContext().Process(cmd)
	}
	return true
}

func (f *frameworkImpl) currContext() Context {
	if len(f.contextStack) == 0 {
		panic("no context is set")
	}
	return f.contextStack[len(f.contextStack)-1]
}

func (f *frameworkImpl) quit() {
	close(f.quitCh)
}
