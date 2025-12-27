package testkit

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

// EnvSeed can be used to set a seed for this run via an environment variable.
// This is useful for reproducing failed tests that used testkit.NewRand.
// The seed is printed to the error log for each failed test:
//
//	testkit_seed=1234567
//
// To use the failure in the next run simply prepend this line to the go test command:
//
//	testkit_seed=1234567 go test ./...
const EnvSeed = "testkit_seed"

func init() {
	if v, ok := os.LookupEnv(EnvSeed); ok {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Fatalf("invalid %s=%s, must be in64", EnvSeed, v)
		}
		seed = i
	} else {
		seed = Now.UnixNano()
	}
}

// seed for this test run
var seed int64

// PRNG is a wrapper around math/rand.Rand.
type PRNG struct {
	*rand.Rand
}

// NewPRNG returns a new pseudo-random number generator.
// By default PRNG is seeded to the current unix timestamp.
// The seed is constant during a single test run.
// Set a custom seed via an envvar before running the tests:
//
//	testkit_seed=1234567 go test ./...
func NewPRNG(t *testing.T) *PRNG {
	t.Cleanup(printSeed(t))
	return &PRNG{
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// printSeed prints the seed to test's error log if the test fails.
func printSeed(t *testing.T) func() {
	return func() {
		if t.Failed() {
			t.Logf("%s=%d", EnvSeed, seed)
		}
	}
}
