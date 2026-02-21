package object

type ObjectType string

const (
	INTEGER_TYPE = "INTEGER"
	FLOAT_TYPE   = "FLOAT"
	STRING_TYPE  = "STRING"
	BOOLEAN_TYPE = "BOOLEAN"
	NULL_TYPE    = "NULL"
)

type Object interface {
	Type() ObjectType
	GetValue() any
}

type IntegerType struct {
	Value int64
}

func (i *IntegerType) Type() ObjectType {
	return INTEGER_TYPE
}

func (i *IntegerType) GetValue() int64 {
	return i.Value
}
