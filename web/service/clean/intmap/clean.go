package intmap

import (
	"github.com/martin-flower/roboz-web/service"
	"github.com/martin-flower/roboz-web/service/direction"
	"go.uber.org/zap"
)

type void struct{}

type Cleaner struct{}

// Clean implementation using a map containing only 64 bit integer keys.
// Can scale to the edges of the test data, although it requires
// a computer with at least 16GB of memory.
// This map of only keys is an example of the Go idiom of implementing a set.
func (c Cleaner) Clean(start service.Coordinate, commands []service.Command) (cleaned int) {

	position := start
	// allocate space for up to 400 million cleaned coordinates
	// (if testing using maximum number of commands, 10000
	// each having a random number of steps from 0 to 100000
	// keys occupy 3,2 billion bytes (8 bytes per integer)
	cleanedCoordinates := make(map[int]void, 40000*len(commands))

	for index, command := range commands {
		if (index+1)%100 == 0 {
			zap.S().Infof("processing %d of %d - (%d coordinates in map)", index+1, len(commands), len(cleanedCoordinates))
		}
		cleanedCoordinates, position = move(cleanedCoordinates, position, command)
	}
	zap.S().Infof("finishing at %v", position)
	cleaned = len(cleanedCoordinates)
	return
}

const max = 100000
const min = -1000000

func move(cleanedCoordinates map[int]void, from service.Coordinate, command service.Command) (map[int]void, service.Coordinate) {
	zap.S().Debugf("moving from %v", from)

	to := service.Coordinate{}

	cleanedCoordinates[pair(from.X, from.Y)] = void{}

	switch command.Direction {
	case direction.North:
		to.X = from.X
		to.Y = from.Y + command.Steps
		if to.Y > max {
			to.Y = max
		}
		for step := 1; step <= command.Steps; step++ {
			cleanedCoordinates[pair(from.X, from.Y+step)] = void{}
		}
	case direction.South:
		to.X = from.X
		to.Y = from.Y - command.Steps
		if to.Y < min {
			to.Y = min
		}
		for step := 1; step <= command.Steps; step++ {
			cleanedCoordinates[pair(from.X, from.Y-step)] = void{}
		}
	case direction.East:
		to.X = from.X + command.Steps
		if to.X > max {
			to.X = max
		}
		to.Y = from.Y
		for step := 1; step <= command.Steps; step++ {
			cleanedCoordinates[pair(from.X+step, from.Y)] = void{}
		}
	case direction.West:
		to.X = from.X - command.Steps
		if to.X < min {
			to.X = min
		}
		to.Y = from.Y
		for step := 1; step <= command.Steps; step++ {
			cleanedCoordinates[pair(from.X-step, from.Y)] = void{}
		}
	default:
		{
			zap.S().Fatalf("programming error, unsupported direction %", command.Direction)
		}
	}

	cleanedCoordinates[pair(to.X, to.Y)] = void{}

	zap.S().Debugf("moved to %v", to)
	return cleanedCoordinates, to
}

// pair converts coordinates to a 64 bit integer for using as a key in the map.
// max int64  : -9223372036854775808 to 9223372036854775807
// squareroot of 9223372036854775807 is 3037000500
// x, y between -100000 and 100000
// small enough for int64, but too big for int32
func pair(x int, y int) int {

	// with inspiration from
	// @see https://www.vertexfragment.com/ramblings/cantor-szudzik-pairing-functions/

	xx := x * 2
	if x < 0 {
		xx = (x * -2) - 1
	}

	yy := y * 2
	if y < 0 {
		yy = (y * -2) - 1
	}

	if xx >= yy {
		return (xx*xx + xx + yy)
	}
	return yy*yy + xx
}
