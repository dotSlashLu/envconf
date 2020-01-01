/*
 package envconf fills struct with system enviroment variables
 api:
	type S struct {
		A string `env:"STR" envdefault:"str"`
	}
	s := S{}
	err := envconf.Prefix("XAE_").Fill(&s)
	err = envconf.Fill(&s)

 supported types:
	string
	[]string //TODO
	number
	[]number //TODO
	bool: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
*/
package main

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
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		fv := stv.Field(i)
		k, hasEnv := f.Tag.Lookup("env")
		if !hasEnv {
			continue
		}
		if !fv.CanSet() {
			return ErrUnsettable
		}
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

func Prefix(p string) *EC {
	return &EC{cfg{prefix: p}}
}

func Fill(in interface{}) error {
	return (&EC{}).Fill(in)
}
