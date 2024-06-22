package types

type String string
type Int int
type Float float64
type Bool bool
type Map map[string]interface{}

type ReferenceToVariable struct {
	Value string
}

func NewReferenceToVariable(value string) ReferenceToVariable {
	return ReferenceToVariable{Value: value}
}
