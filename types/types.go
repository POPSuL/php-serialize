package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Type byte

const (
	TypeInt     Type = 'i'
	TypeFloat   Type = 'd'
	TypeString  Type = 's'
	TypeNull    Type = 'n'
	TypeBoolean Type = 'b'
	TypeArray   Type = 'a'
	TypeObject  Type = 'o'
	TypeEnum    Type = 'e'
)

type Object struct {
	Name      string
	Parents   map[string]*Object
	Public    map[string]*Value
	Protected map[string]*Value
	Private   map[string]*Value
}

type Value struct {
	Type    Type
	Integer int
	Float   float64
	Bool    bool
	Str     string
	Array   map[*Value]*Value
	Object  *Object
	Enum    string
}

func (v *Value) String() string {
	switch v.Type {
	case TypeString:
		return v.Str
	case TypeInt:
		return strconv.FormatInt(int64(v.Integer), 10)
	case TypeBoolean:
		return strconv.FormatBool(v.Bool)
	case TypeNull:
		return "null"
	case TypeArray:
		var sb strings.Builder
		sb.WriteByte('[')
		for k := range v.Array {
			val := v.Array[k]
			sb.Write([]byte(k.String()))
			sb.WriteString(" -> ")
			sb.WriteString(val.String())
			sb.WriteString(", ")
		}
		sb.WriteByte(']')
		return sb.String()
	case TypeObject:
		return fmt.Sprintf("<object of %s>", v.Object.Name)
	case TypeEnum:
		return v.Enum
	}
	return "<unknown>"
}

func NewString(s string) *Value {
	return &Value{
		Type: TypeString,
		Str:  s,
	}
}

func NewEnum(s string) *Value {
	return &Value{
		Type: TypeEnum,
		Enum: s,
	}
}

func NewObject(s string) *Object {
	return &Object{
		Name:      s,
		Parents:   make(map[string]*Object),
		Public:    make(map[string]*Value),
		Protected: make(map[string]*Value),
		Private:   make(map[string]*Value),
	}
}
