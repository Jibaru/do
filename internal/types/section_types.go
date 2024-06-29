package types

type String string
type Int int
type Float float64
type Bool bool
type Map map[string]interface{}
type File struct {
	Path string
}

type ReferenceToVariable struct {
	Value string
}

func NewReferenceToVariable(value string) ReferenceToVariable {
	return ReferenceToVariable{Value: value}
}

// HasBasicTypesValues returns true if all values in the map are basic types (String, Int, Float, Bool)
func (m Map) HasBasicTypesValues() bool {
	for _, v := range m {
		switch v.(type) {
		case String, Bool, Int, Float:
			continue
		default:
			return false
		}
	}
	return true
}

// HasStringValues returns true if all values in the map are strings
func (m Map) HasStringValues() bool {
	for _, v := range m {
		switch v.(type) {
		case String:
			continue
		default:
			return false
		}
	}
	return true
}
