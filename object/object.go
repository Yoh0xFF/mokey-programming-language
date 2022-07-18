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
type Int struct {
	Value int64
}

func (i *Int) Type() ObjectType { return INT_OBJ }
func (i *Int) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

/* Boolean object */
type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL_OBJ }
func (b *Bool) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

/* Null obect */
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

/* Return value object */
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

/* Error object */
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return fmt.Sprintf("Error: %s", e.Message) }

/* Function object */
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatementNode
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FN_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
