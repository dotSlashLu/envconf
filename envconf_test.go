package main

import (
	"github.com/stretchr/testify/assert"
	"os"
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
	os.Setenv("I32", "123123")
	a := struct {
		A string `env:"PATH"`
		B int32  `env:"I32"`
	}{}
	t.Run("string and int fields", func(t *testing.T) {
		err := Fill(&a)
		assert.NoError(t, err, "should not error")
		t.Log(a)
	})
}
