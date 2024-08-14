package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/popsul/php-serialize/v2/decoder"
	"github.com/popsul/php-serialize/v2/types"
)

func main() {
	out, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Cannot read stdin")
		os.Exit(1)
	}

	zval, err := decoder.Decode(out)
	if err != nil {
		fmt.Printf("Cannot unserialize: %+v\n", err)
		os.Exit(1)
	}
	var sb strings.Builder
	printZval(&sb, zval, 1)
	fmt.Print(sb.String())
}

func indent(l int) []byte {
	return bytes.Repeat([]byte{' '}, l*2)
}

func printZval(sb *strings.Builder, v *types.Value, level int) {
	switch v.Type {
	case types.TypeInt:
		sb.WriteString(strconv.FormatInt(int64(v.Integer), 10))
	case types.TypeNull:
		sb.WriteString("null")
	case types.TypeBoolean:
		sb.WriteString(strconv.FormatBool(v.Bool))
	case types.TypeString:
		sb.WriteByte('"')
		sb.WriteString(v.Str)
		sb.WriteByte('"')
	case types.TypeArray:
		sb.WriteString("[\n")
		for k := range v.Array {
			sb.Write(bytes.Repeat([]byte{' '}, level*2))
			printZval(sb, k, level+1)
			sb.WriteString(" => ")
			printZval(sb, v.Array[k], level+1)
			sb.WriteString(",\n")
		}
		sb.Write(bytes.Repeat([]byte{' '}, (level-1)*2))
		sb.WriteString("]")
	case types.TypeObject:
		sb.WriteString("{\n")
		sb.Write(indent(level))
		sb.WriteString("/* ")
		sb.WriteString(v.Object.Name)
		sb.WriteString(" */\n")
		printMembers(sb, v.Object.Public, "public ", level)
		printMembers(sb, v.Object.Protected, "protected ", level)
		printMembers(sb, v.Object.Private, "private ", level)
		for k := range v.Object.Parents {
			printMembers(sb, v.Object.Parents[k].Private, "private "+k+"#", level)
		}
		sb.Write(bytes.Repeat([]byte{' '}, (level-1)*2))
		sb.WriteString("}")
	case types.TypeEnum:
		sb.WriteString(v.Enum)
	}
}

func printMembers(sb *strings.Builder, m map[string]*types.Value, pref string, level int) {
	for k := range m {
		sb.Write(indent(level))
		sb.WriteString(pref)
		sb.WriteString(k)
		sb.WriteString(" = ")
		printZval(sb, m[k], level+1)
		sb.WriteString("\n")
	}
}
