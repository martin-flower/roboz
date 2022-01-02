package sortedset

import (
	"os"
	"testing"

	"github.com/martin-flower/roboz-web/service/clean"
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

// about 3s ..
func Test100Command(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 100)
}

// about 6s ..
func Test200Commands(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 200)
}
