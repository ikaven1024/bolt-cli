package cli

import (
	"fmt"
	"os"

	"github.com/ikaven1024/bolt-cli/cli/command"
	"github.com/ikaven1024/bolt-cli/cli/framework"
	"github.com/ikaven1024/bolt-cli/db"
)

const info = `Welcome to the boltDB monitor.

Type 'help;' or 'h' for help.
Type 'ctrl+C' to clear the current input statement.
Type 'ctrl+C' to exit when empty input statement.
`

func Run(db db.Interface) {
	stdout, stderr := os.Stdout, os.Stderr
	fmt.Fprintln(stdout, info)

	fw := framework.New(db, stdout, stderr)
	fw.EnterContext(command.NewRootContext(fw))
	fw.Run()
}
