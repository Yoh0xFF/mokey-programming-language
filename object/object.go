package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"monkey/code"
	"strings"
)

type BuiltinFunction func(args ...Object) Object

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	INT_OBJ    = "INT"
	BOOL_OBJ   = "BOOL"
	STRING_OBJ = "STRING"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	FN_OBJ     = "FN"
	BULTIN_OBJ = "BUILTIN"

	ARRAY_OBJ = "ARRAY"
	HASH_OBJ  = "HASH"

	COMPILED_FN_OBJ = "COMPILED_FUNC_OBJ"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

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
func (i *IntObject) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

/* Boolean object */
type BoolObject struct {
	Value bool
}

func (b *BoolObject) Type() ObjectType { return BOOL_OBJ }
func (b *BoolObject) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *BoolObject) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

/* String object */
type StringObject struct {
	Value string
}

func (s *StringObject) Type() ObjectType { return STRING_OBJ }
func (s *StringObject) Inspect() string  { return s.Value }
func (s *StringObject) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

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

/* Builtin obect */
type BuiltinObject struct {
	Fn BuiltinFunction
}

func (b *BuiltinObject) Type() ObjectType { return BULTIN_OBJ }
func (b *BuiltinObject) Inspect() string  { return "Builtin function" }

/* Array object */
type ArrayObject struct {
	Elements []Object
}

func (a *ArrayObject) Type() ObjectType { return ARRAY_OBJ }
func (a *ArrayObject) Inspect() string {
	var out bytes.Buffer

	strElements := []string{}
	for _, el := range a.Elements {
		strElements = append(strElements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(strElements, ", "))
	out.WriteString("]")

	return out.String()
}

/* Hash object */
type HashPair struct {
	Key   Object
	Value Object
}

type HashObject struct {
	Pairs map[HashKey]HashPair
}

func (h *HashObject) Type() ObjectType { return HASH_OBJ }
func (h *HashObject) Inspect() string {
	var out bytes.Buffer

	strPairs := []string{}
	for _, pair := range h.Pairs {
		strPairs = append(
			strPairs,
			fmt.Sprintf("%s : %s", pair.Key.Inspect(), pair.Value.Inspect()),
		)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(strPairs, ", "))
	out.WriteString("}")

	return out.String()
}

type CompiledFnObject struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFnObject) Type() ObjectType { return COMPILED_FN_OBJ }
func (cf *CompiledFnObject) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
