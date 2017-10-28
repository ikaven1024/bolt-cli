package app

// Command is a subcommand for a App.
type Command struct {
	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// List of flags to parse
	Flags []Flag
	// List of child commands
	Subcommands []Command
}
