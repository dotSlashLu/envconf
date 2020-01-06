package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestPrefix(t *testing.T) {
	p := "XAE_"
	ec := Prefix(p)
	assert.Equal(t, ec.config.prefix, p, "ec prefix should equal to param")
}

func TestECFillNonPtr(t *testing.T) {
	p := "XAE_"
	s := struct {
		A string
	}{"a"}
	err := Prefix(p).Fill(s)
	assert.EqualError(t, err, ErrNotPtr.Error())
}

func TestECFill(t *testing.T) {
	envval := "xae path"
	os.Setenv("XAE_PATH", envval)
	p := "XAE_"
	s := struct {
		A string `env:"PATH" envdefault:"abc"`
	}{"a"}
	// s := "abc"
	err := Prefix(p).Fill(&s)
	t.Log(s)
	assert.NoError(t, err, "fill should not return an error")
	assert.Equal(t, s.A, envval)
}

func TestFill(t *testing.T) {
	i32v := int32(123123)
	sv := "string value"
	os.Setenv("I32", strconv.FormatInt(int64(i32v), 10))
	os.Setenv("SV", sv)

	a := struct {
		A string `env:"SV"`
		B int32  `env:"I32"`
	}{}
	t.Run("string and int fields", func(t *testing.T) {
		err := Fill(&a)
		assert.NoError(t, err, "should not error")
		assert.Equal(t, a.A, sv)
		assert.Equal(t, a.B, i32v)
		t.Log(a)
	})

	b := struct {
		A string
		B struct {
			A string `env:"SV"`
		}
	}{}
	t.Run("embeded struct", func(t *testing.T) {
		err := Fill(&b)
		assert.NoError(t, err, "should not error")
		assert.Equal(t, b.B.A, sv)
		t.Log(fmt.Sprintf("%+v", b))
	})

	c := struct {
		A struct {
			a string `env:"SV"`
		}
	}{}
	t.Run("embeded unsettable struct field", func(t *testing.T) {
		err := Fill(&c)
		assert.EqualError(t, err, ErrUnsettable.Error())
		t.Log(fmt.Sprintf("%+v", c))
	})
}
