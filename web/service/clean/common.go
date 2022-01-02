package clean

import (
	"math/rand"
	"testing"
	"time"

	"github.com/martin-flower/roboz-web/service"
	"github.com/martin-flower/roboz-web/service/direction"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Init() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func CleanTest(t *testing.T, cleaner service.Cleaner) {

	start := service.Coordinate{X: 1, Y: 1}
	commands := []service.Command{{Direction: direction.North, Steps: 3}, {Direction: direction.South, Steps: 1}}

	// function under test
	cleaned := cleaner.Clean(start, commands)
	assert.Equal(t, 4, cleaned)

	commands = []service.Command{{Direction: direction.North, Steps: 2}, {Direction: direction.South, Steps: 4}}

	// function under test
	cleaned = cleaner.Clean(start, commands)
	assert.Equal(t, 5, cleaned)
}

func CleanSquarePathTest(t *testing.T, cleaner service.Cleaner) {
	start := service.Coordinate{X: 0, Y: 0}
	commands := []service.Command{{Direction: direction.North, Steps: 100}, {Direction: direction.East, Steps: 100}, {Direction: direction.South, Steps: 100}, {Direction: direction.West, Steps: 100}}

	// function under test
	cleaned := cleaner.Clean(start, commands)
	assert.Equal(t, 400, cleaned)

	start = service.Coordinate{X: 100, Y: 100}

	// function under test
	cleaned = cleaner.Clean(start, commands)
	assert.Equal(t, 400, cleaned)
}

func CommandsTest(t *testing.T, cleaner service.Cleaner, total int) {
	start := GetRandomStart()
	commands := GetCommandsWithRandomSteps(total)

	// function under test
	cleaned := cleaner.Clean(start, commands)

	// approximate assertion for random data
	assert.Greater(t, cleaned, 10000*total)
}

// GetRandomStart returns random coordinate in range (100000,100000) - (-100000,-100000)
func GetRandomStart() service.Coordinate {
	rand.Seed(time.Now().UnixNano())
	return service.Coordinate{X: rand.Intn(200001) - 100000, Y: rand.Intn(200001) - 100000}
}

// GetCommandsWithRandomSteps returns the requested number of commands
// each containing between 0 and 100000 steps
func GetCommandsWithRandomSteps(total int) (commands []service.Command) {
	for i := 0; i < total; i++ {
		command := service.Command{}
		command.Direction = direction.GetDirections()[rand.Intn(3)]
		command.Steps = rand.Intn(100000)
		commands = append(commands, command)
	}
	return commands
}
