# php-serialize
![Coverage](https://img.shields.io/badge/Coverage-53.9%25-yellow)

Go-based decoder for php-serialized objects.

## Usage

```go
// php -r 'echo serialize(["foo" => "bar"]);' -> a:1:{s:3:"foo";s:3:"bar";}

import "github.com/popsul/php-serialize/v2/decoder"

// ...

obj, err := decoder.Decode([]byte(`a:1:{s:3:"foo";s:3:"bar";}`))
if err != nil {
    panic(err)
}
fmt.Println(obj.Type) // a
for key := range obj.Array {
    fmt.Println(key.Type) // i for array or s for hashes and sprarced arrays
    fmt.Println(key.Str) // foo
    arrValue = obj.Array[key]
    fmt.Println(arrValue.Type) // s
    fmt.Println(arrValue.Str) // bar
}
```

## Supported types

* `i` -- integer
* `d` -- double (float)
* `b` -- boolean
* `s` -- string
* `a` -- array
* `o` -- object
* `n` -- null

## License

MIT
