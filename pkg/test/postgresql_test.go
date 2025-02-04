package test_test

import (
	"context"
	"testing"

	// Packages
	test "github.com/djthorpe/go-pg/pkg/test"
	assert "github.com/stretchr/testify/assert"
)

///////////////////////////////////////////////////////////////////////////////
// UNIT TESTS

func Test_Postgresql_001(t *testing.T) {
	assert := assert.New(t)

	// Create a new container with postgresql package
	container, pool, err := test.NewPgxContainer(context.Background(), t.Name(), true)
	assert.NoError(err)
	assert.NotNil(container)
	assert.NotNil(pool)
}
