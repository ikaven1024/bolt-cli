package framework

func NoArgs(cmd *Command, args []string) bool {
	return len(args) == 0
}

func RequiresMinArgs(min int) PositionalArgs {
	return func(_ *Command, args []string) bool {
		return len(args) >= min
	}
}

func RequiresMaxArgs(max int) PositionalArgs {
	return func(_ *Command, args []string) bool {
		return len(args) <= max
	}
}

func RequiresRangeArgs(min int, max int) PositionalArgs {
	return func(_ *Command, args []string) bool {
		l := len(args)
		return l >= min && l <= max
	}
}

func ExactArgs(number int) PositionalArgs {
	return func(_ *Command, args []string) bool {
		return len(args) == number
	}
}
