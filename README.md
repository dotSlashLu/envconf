# envconf

[![GoDoc](https://godoc.org/github.com/dotSlashLu/envconf?status.svg)](https://godoc.org/github.com/dotSlashLu/envconf)

Package envconf fills struct with system enviroment variables using struct
tags, default values are supported.

## API

```go
type St struct {
    A string `env:"PATH" envdefault:"/usr/bin"`
    B int32  `env:"I32"`
}
s := St{}

// if all environment variable keys share a common prefix
err := envconf.Prefix("XAE_").Fill(&s)

// or just call Fill
err = envconf.Fill(&s)
```

## Supported Types

- string
- number
- bool: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
- TODO
    - []string
    - []number

