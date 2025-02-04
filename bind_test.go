package pg_test

import (
	"testing"

	// Packages
	"github.com/djthorpe/go-pg"
	"github.com/stretchr/testify/assert"
)

func Test_Bind_001(t *testing.T) {
	assert := assert.New(t)

	t.Run("1", func(t *testing.T) {
		bind := pg.NewBind("a", "b")
		assert.NotNil(bind)
		assert.True(bind.Has("a"))
		assert.Equal("b", bind.Get("a"))
	})

	t.Run("2", func(t *testing.T) {
		bind := pg.NewBind("a", "b", "c")
		assert.Nil(bind)
	})

	t.Run("3", func(t *testing.T) {
		bind := pg.NewBind("a", 100)
		assert.NotNil(bind)
		assert.True(bind.Has("a"))
		assert.Equal(100, bind.Get("a"))
	})

	t.Run("4", func(t *testing.T) {
		bind := pg.NewBind()
		assert.NotNil(bind)
		assert.Equal("@a", bind.Set("a", "b"))
		assert.True(bind.Has("a"))
		assert.Equal("b", bind.Get("a"))
	})

	t.Run("5", func(t *testing.T) {
		bind := pg.NewBind("", "b")
		assert.Nil(bind)
	})

	t.Run("6", func(t *testing.T) {
		bind := pg.NewBind()
		assert.NotNil(bind)
		assert.Equal("", bind.Set("", "b"))
	})

}

func Test_Bind_002(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In  string
		Out string
	}{
		{In: `$schema`, Out: "schema"},
		{In: `${'schema'}`, Out: "'schema'"},
		{In: `${"schema"}`, Out: `"schema"`},
		{In: `$1`, Out: `$1`},
		{In: `${1}`, Out: `$1`},
		{In: `$$`, Out: `$$`},
		{In: `${'single'}`, Out: `'''single'''`},
		{In: `${"single"}`, Out: `"'single'"`},
		{In: `${'double'}`, Out: `'"double"'`},
		{In: `${"double"}`, Out: `"""double"""`},
	}

	bind := pg.NewBind(
		"schema", "schema",
		"single", "'single'",
		"double", "\"double\"",
	)

	for _, test := range tests {
		t.Run(test.In, func(t *testing.T) {
			assert.Equal(test.Out, bind.Replace(test.In))
		})
	}
}

func Test_Bind_003(t *testing.T) {
	assert := assert.New(t)

	bind := pg.NewBind(
		"list", []string{"a", "b", "c"},
	)
	assert.Equal("IN ('a','b','c')", bind.Replace("IN (${'list'})"))
}
