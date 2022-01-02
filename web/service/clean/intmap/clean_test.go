package intmap

import (
	"os"
	"testing"

	"github.com/martin-flower/roboz-web/service/clean"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	clean.Init()
	code := m.Run()
	os.Exit(code)
}

func TestClean(t *testing.T) {
	clean.CleanTest(t, Cleaner{})
}

func TestCleanSquarePath(t *testing.T) {
	clean.CleanSquarePathTest(t, Cleaner{})
}

// about 2 seconds ..
func Test500Commands(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 500)
}

// about 3 seconds ..
func Test1000Commands(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 1000)
}

// maximum number of commands according to the specification
// each with random number of steps up to 100000
// about 30s
func Test10000Commands(t *testing.T) {
	t.Skipf("need to run manually with go test -run Test10000Commands -timeout 70s -v")
	start := clean.GetRandomStart()
	commands := clean.GetCommandsWithRandomSteps(10000)

	// function under test
	cleaned := Cleaner{}.Clean(start, commands)

	// no useful assertion with random numbers
	assert.Greater(t, cleaned, 100000000)
}

// demonstrate that 64 bit integer will accommodate test data
func TestPair(t *testing.T) {
	assert.Equal(t, 40000400000, pair(100000, 100000))
	assert.Equal(t, 39999999999, pair(-100000, -100000))
	assert.Equal(t, 4000004000000, pair(1000000, 1000000))
	assert.Equal(t, 3999999999999, pair(-1000000, -1000000))
}
