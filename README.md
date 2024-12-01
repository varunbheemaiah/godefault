# godefault

`godefault` is a Go package that allows you to specify default values for struct fields using a `default` tag. When a struct is serialized and fields are empty (zero value), `godefault` sets the fields to the specified default values unless the `omitempty` option is used.

## Features

- **Set Default Values**: Automatically assign default values to empty struct fields.
- **Supports Basic Types**: Works with strings, numbers, booleans, and `time.Duration`.
- **Nested Structs**: Recursively sets defaults in nested structs.
- **Respects `omitempty`**: Does not set default values for fields marked with `omitempty`.

## Installation

```bash
go get github.com/varunbheemaiah/godefault
```

## Documentation

For detailed documentation, please refer to the [godoc](https://pkg.go.dev/github.com/varunbheemaiah/godefault).

## Usage

```go
package main

import "github.com/varunbheemaiah/godefault"

type Person struct {
    Name  string `json:"name" default:"NoName"`
    Age   int    `json:"age" default:"18"`
    Email string `json:"email,omitempty" default:"noemail@example.com"`
}

func main() {
    p := Person{}
    err := godefault.SetDefaults(&p)
    if err != nil {
        // Handle error
    }
    // Use p as needed
}
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License.