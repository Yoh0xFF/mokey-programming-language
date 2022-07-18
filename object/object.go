package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	INT_OBJ  = "INT"
	BOOL_OBJ = "BOOL"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	FN_OBJ = "FN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

/* Integer object */
type IntObject struct {
	Value int64
}

func (i *IntObject) Type() ObjectType { return INT_OBJ }
func (i *IntObject) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

/* Boolean object */
type BoolObject struct {
	Value bool
}

func (b *BoolObject) Type() ObjectType { return BOOL_OBJ }
func (b *BoolObject) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

/* NullObject obect */
type NullObject struct{}

func (n *NullObject) Type() ObjectType { return NULL_OBJ }
func (n *NullObject) Inspect() string  { return "null" }

/* Return value object */
type ReturnValueObject struct {
	ValueObject Object
}

func (rv *ReturnValueObject) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValueObject) Inspect() string  { return rv.ValueObject.Inspect() }

/* ErrorObject object */
type ErrorObject struct {
	Message string
}

func (e *ErrorObject) Type() ObjectType { return ERROR_OBJ }
func (e *ErrorObject) Inspect() string  { return fmt.Sprintf("Error: %s", e.Message) }

/* FunctionObject object */
type FunctionObject struct {
	ParamNodes []*ast.IdentifierNode
	BodyNode   *ast.BlockStatementNode
	Env        *Environment
}

func (f *FunctionObject) Type() ObjectType { return FN_OBJ }
func (f *FunctionObject) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.ParamNodes {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.BodyNode.String())
	out.WriteString("\n}")

	return out.String()
}
