package pg

import (
	"testing"

	// Packages
	"github.com/stretchr/testify/assert"
)

func Test_Opts_001(t *testing.T) {
	assert := assert.New(t)

	// Empty options
	o, err := apply()
	if assert.NoError(err) {
		assert.NotNil(o)
	}
}

func Test_Opts_002(t *testing.T) {
	assert := assert.New(t)

	// Empty options
	o, err := apply()
	if assert.NoError(err) {
		assert.NotNil(o)
		assert.Equal("host=localhost pool_max_conns=10 port=5432", o.Encode())
	}
}

func Test_Opts_003(t *testing.T) {
	assert := assert.New(t)

	// Host/Port
	o, err := apply(
		WithHostPort("host", "999"),
	)
	if assert.NoError(err) {
		assert.NotNil(o)
		assert.Equal("host=host pool_max_conns=10 port=999", o.Encode())
	}
}

func Test_Opts_004(t *testing.T) {
	assert := assert.New(t)

	// Host/Port
	o, err := apply(
		WithCredentials("user", "password"),
	)
	if assert.NoError(err) {
		assert.NotNil(o)
		assert.Equal("dbname=user host=localhost password=password pool_max_conns=10 port=5432 user=user", o.Encode())
	}
}

func Test_Opts_005(t *testing.T) {
	assert := assert.New(t)

	// Host/Port
	o, err := apply(
		WithCredentials("user", "password"),
		WithDatabase("db"),
		WithHostPort("host", "999"),
	)
	if assert.NoError(err) {
		assert.NotNil(o)
		assert.Equal("dbname=db host=host password=password pool_max_conns=10 port=999 user=user", o.Encode())
	}
}

func Test_Opts_006(t *testing.T) {
	assert := assert.New(t)

	// Host/Port
	o, err := apply(
		WithSSLMode("disable"),
	)
	if assert.NoError(err) {
		assert.NotNil(o)
		assert.Equal("host=localhost pool_max_conns=10 port=5432 sslmode=disable", o.Encode())
	}
}
