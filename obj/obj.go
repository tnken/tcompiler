package obj

import (
	"fmt"

	"github.com/takeru56/tcompiler/code"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	FUNCTION_OBJ = "FUNCTION"
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

type Function struct {
	Instructions code.Instructions
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string  { return fmt.Sprintf("function%p", f) }

// TODO: のちほど32bitに対応する
// ひとまず2byte(16bit)で表現
func (f *Function) Size() int { return len(f.Instructions) }
