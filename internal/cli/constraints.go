package cli

const (
	TagConstraintNone TagConstraint = iota
	TagConstraintUntagged
	TagConstraintInvalid
)

const (
	PrintConstraintNone PrintConstraint = iota
	PrintConstraintValid
	PrintConstraintInvalid
	PrintConstraintUntagged
)

const (
	UntagConstraintNone UntagConstraint = iota
	UntagConstraintAll
	UntagConstraintInvalid
)

type TagConstraint int
type UntagConstraint int
type PrintConstraint int
