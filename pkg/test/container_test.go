package test_test

import (
	"context"
	"testing"

	// Packages
	test "github.com/mutablelogic/go-pg/pkg/test"
	assert "github.com/stretchr/testify/assert"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	TEST_HELLOWORLD = "hello-world" // Official "hello world" container
)

///////////////////////////////////////////////////////////////////////////////
// UNIT TESTS

func Test_Container_001(t *testing.T) {
	assert := assert.New(t)

	// Create a new container with hello-world package
	container, err := test.NewContainer(context.Background(), t.Name(), TEST_HELLOWORLD)
	if !assert.NoError(err) {
		t.FailNow()
	}
	defer container.Close(context.Background())
	t.Log(container)
}
