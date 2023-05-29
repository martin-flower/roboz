package sortedset

import (
	"go.uber.org/zap"

	"github.com/martin-flower/roboz-web-go/service"
	"github.com/martin-flower/roboz-web-go/service/direction"
	"github.com/wangjia184/sortedset"
)

type Cleaner struct{}

// Clean sortedset implementation performs much better than simplest,
// but does not scale to the edges of the test data.
func (c Cleaner) Clean(start service.Coordinate, commands []service.Command) (cleaned int) {

	position := start
	cleanedSoFar := sortedset.New()

	for _, command := range commands {
		var cleanedThisCommand = []service.Coordinate{}
		position, cleanedThisCommand = move(position, command)
		cleanedSoFar = consolidate(cleanedSoFar, cleanedThisCommand)
	}
	zap.S().Infof("finishing at %+v", position)
	cleaned = cleanedSoFar.GetCount()
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
func consolidate(c1s *sortedset.SortedSet, c2s []service.Coordinate) (consolidated *sortedset.SortedSet) {
	consolidated = c1s
	for _, c2 := range c2s {
		consolidated.AddOrUpdate(c2.Key(), sortedset.SCORE(c2.X+c2.Y), c2)
	}
	return
}
