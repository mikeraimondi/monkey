package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/mikeraimondi/monkey/ast"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
	key   HashKey
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	if i.key.Type == "" {
		i.key = HashKey{Type: i.Type(), Value: uint64(i.Value)}
	}

	return i.key
}

type Boolean struct {
	Value bool
	key   HashKey
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.key.Type == "" {
		var value uint64

		if b.Value {
			value = 1
		} else {
			value = 0
		}
		b.key = HashKey{Type: b.Type(), Value: value}
	}

	return b.key
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	out := ast.StringBuilder{}

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.MustWrite("fn")
	out.MustWrite("(")
	out.MustWrite(strings.Join(params, ", "))
	out.MustWrite(") {\n")
	out.MustWrite(f.Body.String())
	out.MustWrite("\n}")

	return out.String()
}

type String struct {
	Value string
	key   HashKey
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	if s.key.Type == "" {
		// TODO deal with collisions
		// chaining or open addressing
		h := fnv.New64a()
		_, err := h.Write([]byte(s.Value))
		if err != nil {
			panic(err)
		}

		s.key = HashKey{Type: s.Type(), Value: h.Sum64()}
	}

	return s.key
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "built-in function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	out := ast.StringBuilder{}

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.MustWrite("[")
	out.MustWrite(strings.Join(elements, ", "))
	out.MustWrite("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	out := ast.StringBuilder{}

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.MustWrite("{")
	out.MustWrite(strings.Join(pairs, ", "))
	out.MustWrite("}")

	return out.String()
}
