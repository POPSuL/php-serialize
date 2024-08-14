package decoder

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/popsul/php-serialize/v2/types"
)

var (
	errInvalidInteger      = errors.New("invalid integer")
	errInvalidFloat        = errors.New("invalid float")
	errInvalidBool         = errors.New("invalid bool")
	errInvalidArrayKey     = errors.New("invalid array key")
	errInvalidPropertyName = errors.New("invalid property name type")
	errUnexpected          = errors.New("unexpected case")
)

func Decode(in []byte) (v *types.Value, e error) {
	v, _, e = parseNext(in, 0)
	return
}

func peekInteger(in []byte, i int) (int, int, error) {
	var out = []byte{}
	for i < len(in) {
		if in[i] == ';' || in[i] == ':' {
			break
		}
		out = append(out, in[i])
		i++
	}
	if len(out) == 0 {
		return 0, 0, errInvalidInteger
	}
	val, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, 0, err
	}

	return val, i, nil
}

func peekFloat(in []byte, i int) (float64, int, error) {
	var out = []byte{}
	for i < len(in) {
		if in[i] == ';' {
			break
		}
		out = append(out, in[i])
		i++
	}
	if len(out) == 0 {
		return 0, 0, errInvalidFloat
	}
	val, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		return 0, 0, err
	}
	return val, i, nil
}

func peekBool(in []byte, i int) (bool, error) {
	switch in[i] {
	case '1':
		return true, nil
	case '0':
		return false, nil
	default:
		return false, errInvalidBool
	}
}

func parseNext(in []byte, i int) (*types.Value, int, error) {
	for i < len(in) {
		b := in[i]
		i++ // skip type
		i++ // skip delimiter
		switch b {
		case 'i':
			val, offset, err := peekInteger(in, i)
			if err != nil {
				return nil, i, err
			}
			return &types.Value{
				Type:    types.TypeInt,
				Integer: val,
			}, offset + 1, nil
		case 'd':
			val, offset, err := peekFloat(in, i)
			if err != nil {
				return nil, i, err
			}
			return &types.Value{
				Type:  types.TypeFloat,
				Float: val,
			}, offset + 1, nil
		case 's', 'E':
			val, offset, err := peekInteger(in, i)
			if err != nil {
				return nil, i, err
			}
			i = offset + 2
			if b == 's' {
				return types.NewString(string(in[i : i+val])), i + val + 2, nil
			} else {
				return types.NewEnum(string(in[i : i+val])), i + val + 2, nil
			}
		case 'N':
			return &types.Value{
				Type: types.TypeNull,
			}, i, nil
		case 'b':
			v, err := peekBool(in, i)
			if err != nil {
				return nil, i, err
			}
			return &types.Value{
				Type: types.TypeBoolean,
				Bool: v,
			}, i + 2, nil
		case 'a':
			arrLen, offset, err := peekInteger(in, i)
			if err != nil {
				return nil, i, err
			}
			if arrLen == 0 {
				return &types.Value{
					Type: types.TypeArray,
				}, i + 4, nil
			}
			i = offset + 2
			m := make(map[*types.Value]*types.Value)
			for ; arrLen > 0; arrLen-- {
				key, off, err := parseNext(in, i)
				if err != nil {
					return nil, i, err
				}
				if key.Type != types.TypeInt && key.Type != types.TypeString {
					return nil, i, errInvalidArrayKey
				}
				i = off
				val, off, err := parseNext(in, i)
				if err != nil {
					return nil, i, err
				}
				i = off
				m[key] = val
			}
			return &types.Value{
				Type:  types.TypeArray,
				Array: m,
			}, i + 1, nil
		case 'O':
			nameLen, offset, err := peekInteger(in, i)
			if err != nil {
				return nil, i, err
			}
			i = offset + 2
			name := in[i : i+nameLen]
			i = i + nameLen + 2
			propsCount, offset, err := peekInteger(in, i)
			if err != nil {
				return nil, i, err
			}
			i = offset + 2
			obj := types.NewObject(string(name))
			for ; propsCount > 0; propsCount-- {
				propName, offset, err := parseNext(in, i)
				if err != nil {
					return nil, i, err
				}
				if propName.Type != types.TypeString {
					return nil, i, errInvalidPropertyName
				}
				i = offset
				propValue, offset, err := parseNext(in, i)
				if err != nil {
					return nil, i, err
				}
				i = offset
				nameParts := bytes.Split([]byte(propName.Str), []byte{0x00})
				if len(nameParts) == 1 {
					obj.Public[propName.Str] = propValue
				} else if len(nameParts) == 3 {
					owner := string(nameParts[1])
					prop := string(nameParts[2])
					if owner[0] == '*' {
						obj.Protected[prop] = propValue
						continue
					}
					if owner == string(name) {
						obj.Private[prop] = propValue
						continue
					}
					if _, exists := obj.Parents[owner]; !exists {
						obj.Parents[owner] = types.NewObject(owner)
					}
					obj.Parents[owner].Private[prop] = propValue
				} else {
					return nil, i, errInvalidPropertyName
				}
			}
			return &types.Value{
				Type:   types.TypeObject,
				Object: obj,
			}, i + 1, nil
		default:
			return nil, i, fmt.Errorf("unknown type at %d", i)
		}
	}
	return nil, i, errUnexpected
}
