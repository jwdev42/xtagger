package cli

const (
	CommandInvalid     Command = ""
	CommandPrint               = "print"
	CommandTag                 = "tag"
	CommandUntag               = "untag"
	CommandRecalculate         = "recalc"
)

type Command string
