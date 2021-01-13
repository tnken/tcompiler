package obj

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Size() int
}

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// TODO: のちほど32bitに対応する
// ひとまず2byte(16bit)で表現
func (i *Integer) Size() int { return 2 }
