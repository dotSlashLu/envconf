// package envconf fills struct with system enviroment variables
//
// api:
//	type St struct {
//	    A string `env:"PATH" envdefault:"/usr/bin"`
//	    B int32  `env:"I32"`
//
//		// embedded structs are supported
//	    C struct {
//	        Term string `env:"TERM" envdefault:"screen"`
//	    }
//
//		// you can set a common prefix for fields of
//		// an embedded struct by the envprefix tag
//		D struct {
//			Client string `env:"CLIENT"`
//			TTY    string `env:"TTY"`
//		} `envprefix:"SSH_"`
//	}
//	s := St{}
//
//	// if all environment variable keys share a common prefix
//	err := envconf.Prefix("XAE_").Fill(&s)
//
//	// or just call Fill
//	err = envconf.Fill(&s)
//
// supported types:
//	- string
//	- number
//	- bool: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
//	- struct with fields of supported type
//	- TODO
//		- []string
//		- []number
package envconf

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrNotPtr          = errors.New("input is not a pointer")
	ErrNotSt           = errors.New("input is not a pointer to a struct")
	ErrTypNotSupported = errors.New("field type not supported")
	ErrUnsettable      = errors.New("field not settable")
	ErrCast            = errors.New("can't cast string to field type")
)

type cfg struct {
	prefix string
}

type EC struct {
	config cfg
}

func (ec *EC) env(k string) string {
	if ec.config.prefix != "" {
		k = ec.config.prefix + k
	}
	return os.Getenv(k)
}

// fillFields fills each field of a struct
// fPrefix: field specific prefix
func (ec *EC) fillFields(st reflect.Type, stv reflect.Value, fPrefix string) (err error) {
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		fv := stv.Field(i)

		// process embedded struct
		if fv.Kind() == reflect.Struct && fv.CanSet() {
			p := f.Tag.Get("envprefix")
			err = ec.fillFields(fv.Type(), fv, p)
			if err != nil {
				return
			}
			continue
		}

		k, hasEnv := f.Tag.Lookup("env")
		if !hasEnv {
			continue
		}
		if !fv.CanSet() {
			return ErrUnsettable
		}
		k = fPrefix + k
		v := ec.env(k)
		defaultv, hasDefault := f.Tag.Lookup("envdefault")
		if v == "" && hasDefault {
			v = defaultv
		}
		// TODO: support slice types
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			iv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return ErrCast
			}
			fv.SetInt(iv)
		case reflect.Bool:
			bv, err := strconv.ParseBool(v)
			if err != nil {
				return ErrCast
			}
			fv.SetBool(bv)
		default:
			err = ErrTypNotSupported
			return
		}
	}
	return
}

// Fill does the actual filling process
func (ec *EC) Fill(in interface{}) (err error) {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Ptr {
		return ErrNotPtr
		// requires 1.13
		// return fmt.Errorf("%w: %s", ErrNotPtr, v.Kind())
	}
	stv := reflect.Indirect(v)
	if stv.Kind() != reflect.Struct {
		return ErrNotSt
		// requires 1.13
		// return fmt.Errorf("%w: %s", ErrNotSt, stv.Kind())
	}
	st := stv.Type()
	err = ec.fillFields(st, stv, "")
	return
}

// Prefix configures the EC instance then returns it.
// Usually a Fill call is chained after it.
func Prefix(p string) *EC {
	return &EC{cfg{prefix: p}}
}

// Fill fills a struct without configuration.
// It accepts a pointer to a struct and sets input in place.
// If an invalid type is passed in, or a field cannot be processed, it
// returns an error.
func Fill(in interface{}) error {
	return (&EC{}).Fill(in)
}
