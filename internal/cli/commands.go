package cli

const (
	CommandInvalid Command = ""
	CommandPrint           = "print"
	CommandTag             = "tag"
	CommandUntag           = "untag"
)

type Command string
