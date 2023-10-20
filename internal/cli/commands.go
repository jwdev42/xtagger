package cli

const (
	CommandInvalid    Command = ""
	CommandPrint              = "print"
	CommandTag                = "tag"
	CommandUntag              = "untag"
	CommandInvalidate         = "invalidate"
)

type Command string
