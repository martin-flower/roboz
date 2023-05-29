package simplest

import (
	"github.com/martin-flower/roboz-web-go/service"
	"github.com/martin-flower/roboz-web-go/service/direction"
	"go.uber.org/zap"
)

type Cleaner struct{}

// Clean simplest implementation with no optimization. Does not scale.
func (c Cleaner) Clean(start service.Coordinate, commands []service.Command) (cleaned int) {

	position := start
	cleanedSoFar := []service.Coordinate{}

	for _, command := range commands {
		var cleanedThisCommand = []service.Coordinate{}
		position, cleanedThisCommand = move(position, command)
		cleanedSoFar = consolidate(cleanedSoFar, cleanedThisCommand)
	}
	zap.S().Infof("finishing at %v", position)
	cleaned = len(cleanedSoFar)
	return
}

func move(from service.Coordinate, command service.Command) (to service.Coordinate, cleaned []service.Coordinate) {
	zap.S().Debugf("moving from %v", from)

	cleaned = append(cleaned, from)

	switch command.Direction {
	case direction.North:
		to.Y = from.Y + command.Steps
		to.X = from.X
		for step := 1; step <= command.Steps; step++ {
			cleaned = append(cleaned, service.Coordinate{X: from.X, Y: from.Y + step})
		}
	case direction.South:
		to.Y = from.Y - command.Steps
		to.X = from.X
		for step := 1; step <= command.Steps; step++ {
			cleaned = append(cleaned, service.Coordinate{X: from.X, Y: from.Y - step})
		}
	case direction.East:
		to.Y = from.Y
		to.X = from.X + command.Steps
		for step := 1; step <= command.Steps; step++ {
			cleaned = append(cleaned, service.Coordinate{X: from.X + step, Y: from.Y})
		}
	case direction.West:
		to.Y = from.Y
		to.X = from.X - command.Steps
		for step := 1; step <= command.Steps; step++ {
			cleaned = append(cleaned, service.Coordinate{X: from.X - step, Y: from.Y})
		}
	default:
		{
			zap.S().Errorf("programming error, unsupported direction %", command.Direction)
		}
	}

	zap.S().Debugf("moved to %v", to)
	return
}

// add c2 service.Coordinates to c1 provided the c2 service.Coordinate does not exist in c1
// return the consolidated service.Coordinates
func consolidate(c1s []service.Coordinate, c2s []service.Coordinate) (consolidated []service.Coordinate) {
	consolidated = c1s
	for _, c2 := range c2s {
		if !contains(c1s, c2) {
			consolidated = append(consolidated, c2)
		}
	}
	return
}

func contains(coordinates []service.Coordinate, item service.Coordinate) bool {
	for _, coordinate := range coordinates {
		if coordinate == item {
			zap.S().Debugf("c1s contains c2 - %v, %v", coordinates, item)
			return true
		}
	}
	zap.S().Debugf("c1s does not contain c2 - %v, %v", coordinates, item)
	return false
}
