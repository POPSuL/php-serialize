package decoder_test

import (
	"fmt"
	"testing"

	. "github.com/popsul/php-serialize/v2/decoder"
	. "github.com/popsul/php-serialize/v2/types"
	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	in := "i:1;"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeInt, out.Type)
	assert.Equal(t, 1, out.Integer)
}

func TestFloat(t *testing.T) {
	in := "d:1.22;"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeFloat, out.Type)
	assert.Equal(t, 1.22, out.Float)
}

func TestString(t *testing.T) {
	in := "s:1:\"s\";"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type)
	assert.Equal(t, "s", out.Str)
}

func TestEnum(t *testing.T) {
	in := "E:6:\"Test:A\";"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeEnum, out.Type)
	assert.Equal(t, "Test:A", out.Enum)
}

func TestNull(t *testing.T) {
	in := "N;"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeNull, out.Type)
}

func TestBoolFalse(t *testing.T) {
	in := "b:0;"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeBoolean, out.Type)
	assert.Equal(t, false, out.Bool)
}

func TestBoolTrue(t *testing.T) {
	in := "b:1;"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeBoolean, out.Type)
	assert.Equal(t, true, out.Bool)
}

func TestBoolInvalid(t *testing.T) {
	in := "b:2;"
	_, err := Decode([]byte(in))
	assert.Error(t, err)
}

func TestArrayEmpty(t *testing.T) {
	in := "a:0:{}"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeArray, out.Type)
	assert.Equal(t, 0, len(out.Array))
}

func TestArraySimple(t *testing.T) {
	in := "a:2:{i:0;i:1;i:1;i:2;}"
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeArray, out.Type)
	assert.Equal(t, 2, len(out.Array))
	for k := range out.Array {
		v := out.Array[k]
		assert.Equal(t, TypeInt, k.Type)
		assert.Equal(t, TypeInt, v.Type)
		assert.Equal(t, k.Integer+1, v.Integer)
	}
}

func getByKey(v *Value, k interface{}) *Value {
	if v == nil || v.Type != TypeArray {
		return nil
	}
	for key := range v.Array {
		switch k.(type) {
		case int:
			if key.Type == TypeInt && k == key.Integer {
				return v.Array[key]
			}
		case string:
			if key.Type == TypeString && k == key.Str {
				return v.Array[key]
			}
		default:
			fmt.Printf("I don't know about type %T!\n", k)
		}
	}
	return nil
}

func TestArrayNested(t *testing.T) {
	in := `a:2:{i:0;i:1;s:1:"x";a:1:{i:0;i:-1;}}`
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeArray, out.Type)
	assert.Equal(t, 2, len(out.Array))

	first := getByKey(out, 0)
	assert.NotNil(t, first)
	assert.Equal(t, TypeInt, first.Type)
	assert.Equal(t, 1, first.Integer)

	second := getByKey(out, "x")
	assert.NotNil(t, second)
	assert.Equal(t, TypeArray, second.Type)
	assert.Equal(t, 1, len(second.Array))

	third := getByKey(second, 0)
	assert.NotNil(t, third)
	assert.Equal(t, TypeInt, third.Type)
	assert.Equal(t, -1, third.Integer)
}

func TestObjectSimple(t *testing.T) {
	in := `O:8:"stdClass":1:{s:1:"0";i:1;}`
	out, err := Decode([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, TypeObject, out.Type)
	assert.Equal(t, out.Object.Name, "stdClass")
	assert.Empty(t, out.Object.Parents)
	assert.Empty(t, out.Object.Private)
	assert.Empty(t, out.Object.Protected)
	assert.Len(t, out.Object.Public, 1)
	assert.Equal(t, 1, out.Object.Public["0"].Integer)
}

func TestObjectComplex(t *testing.T) {
	//in := "O:1:\"B\":4:{s:4:\" A b\";i:1;s:4:\" * c\";i:3;s:1:\"a\";i:1;s:4:\" B b\";i:2;}"
	in := []byte{
		79, 58, 49, 58, 34, 66, 34, 58, 52, 58, 123, 115, 58, 52, 58, 34, 0, 65, 0, 98, 34, 59, 105,
		58, 49, 59, 115, 58, 52, 58, 34, 0, 42, 0, 99, 34, 59, 105, 58, 51, 59, 115, 58, 49, 58, 34,
		97, 34, 59, 105, 58, 49, 59, 115, 58, 52, 58, 34, 0, 66, 0, 98, 34, 59, 105, 58, 50, 59, 125,
	}
	out, err := Decode(in)
	assert.NoError(t, err)
	assert.Equal(t, TypeObject, out.Type)
	assert.Equal(t, out.Object.Name, "B")
	assert.Len(t, out.Object.Parents, 1)
	assert.Equal(t, 1, out.Object.Public["a"].Integer)
	assert.Equal(t, 2, out.Object.Private["b"].Integer)
	assert.Equal(t, 3, out.Object.Protected["c"].Integer)
	assert.Equal(t, "A", out.Object.Parents["A"].Name)
	assert.Equal(t, 1, out.Object.Parents["A"].Private["b"].Integer)
}
