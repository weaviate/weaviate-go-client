package dev_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func TestAssert(t *testing.T) {
	require.NotPanics(t, func() { dev.Assert(true, "ok") })
	require.Panics(t, func() { dev.Assert(false, "not ok") })
}

func TestAssertNotNil(t *testing.T) {
	require.NotPanics(t, func() { dev.AssertNotNil(new(string), "not nil") })
	require.Panics(t, func() { dev.AssertNotNil(nil, "nil") })
}
