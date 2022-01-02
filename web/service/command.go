package service

import (
	"fmt"

	"github.com/martin-flower/roboz-web/service/direction"
)

type Command struct {
	Direction direction.Direction
	Steps     int
}

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) Key() string {
	return fmt.Sprintf("x:%d,y:%d", c.X, c.Y)
}

type Cleaner interface {
	Clean(start Coordinate, commands []Command) (cleaned int)
}
