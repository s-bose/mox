package eval

const (
	INT_TYPE    = "INT"
	FLOAT_TYPE  = "FLOAT"
	STRING_TYPE = "STRING"
	BOOL_TYPE   = "BOOL"
	NULL_TYPE   = "NULL"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INT_TYPE
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_TYPE
}

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType {
	return BOOL_TYPE
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_TYPE
}

type NullType struct{}

func (n *NullType) Type() ObjectType {
	return NULL_TYPE
}
