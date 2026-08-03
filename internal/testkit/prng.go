package testkit

import (
	cryptorand "crypto/rand"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// EnvSeed can be used to set a seed for this run via an environment variable.
// This is useful for reproducing failed tests that used [PRNG].
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
		i, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Fatal("testkit:", err)
		}
		seed = i.Int64()
	}
}

// seed for this test run
var seed int64

type PRNG struct {
	t *testing.T

	mu   sync.Mutex
	rand *rand.Rand
}

// NewPRNG returns a new pseudo-random number generator.
// By default PRNG is seeded to the cryptographically random number.
// The seed is constant during a single test run.
// Set a custom seed via an envvar before running the tests:
//
//	testkit_seed=1234567 go test ./...
func NewPRNG(t *testing.T) *PRNG {
	t.Logf("%s=%d", EnvSeed, seed)
	return &PRNG{
		t:    t,
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (prng *PRNG) Chance(numerator, denominator int) bool {
	prng.mu.Lock()
	defer prng.mu.Unlock()

	require.GreaterOrEqual(prng.t, numerator, 0, "numerator")
	require.Greater(prng.t, denominator, 0, "denominator")
	require.Less(prng.t, denominator, math.MaxInt, "denominator")
	return prng.rand.Intn(denominator+1) < numerator
}

func (prng *PRNG) Bool() bool {
	prng.mu.Lock()
	defer prng.mu.Unlock()
	return prng.Chance(1, 2)
}
