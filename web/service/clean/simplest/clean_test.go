package simplest

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

// takes about 3 seconds
func Test3Command(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 3)
}
