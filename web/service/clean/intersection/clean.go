package intersection

import (
	"github.com/martin-flower/roboz-web/service"
	"github.com/martin-flower/roboz-web/service/direction"
	"go.uber.org/zap"
)

type Cleaner struct{}

// Clean implementation using intersections of lines
// Load all the lines from the commands
// Count the number of points where the lines crossor touch each other
func (c Cleaner) Clean(start service.Coordinate, commands []service.Command) (cleaned int) {

	length, intersections := getLengthAndIntersectionsForCommands(start, commands)
	cleaned = length - intersections
	zap.S().Infof("cleaned:%d, length:%d, intersections:%d", cleaned, length, intersections)

	return
}

type line struct {
	Start service.Coordinate
	End   service.Coordinate
}

func getLengthAndIntersectionsForCommands(start service.Coordinate, commands []service.Command) (length int, intersections int) {
	var lines = []line{}
	length, lines = getLengthAndLines(start, commands)
	zap.S().Infof("length:%d, lines:%d", length, len(lines))
	intersections = getIntersectionsForLines(lines)
	return
}

type void struct{}

func getIntersectionsForLines(lines []line) (intersections int) {
	intersectionsMap := make(map[int]map[int16]void)
	for lineIndex := int16(0); lineIndex < int16(len(lines)-1); lineIndex++ {
		if lineIndex%100 == 0 {
			zap.S().Infof("checking line %d of %d for intersections", lineIndex+1, len(lines))
		} else {
			zap.S().Debugf("checking line %d of %d for intersections", lineIndex+1, len(lines))
		}
		getIntersection(intersectionsMap, lineIndex, lines[lineIndex], lines[lineIndex+1:])
	}
	for _, lines := range intersectionsMap {
		intersections = intersections + len(lines)
	}
	return
}

func getIntersection(intersectionsMap map[int]map[int16]void, lineIndex int16, line line, lines []line) {
	if len(lines) < 1 {
		return
	}
	for i := 0; i < len(lines); i++ {
		getIntersectionOfTwoLines(intersectionsMap, lineIndex, line, lines[i])
	}
}

func getIntersectionOfTwoLines(intersectionsMap map[int]map[int16]void, lineIndex int16, line1 line, line2 line) {
	zap.S().Debugf("%v,%v", line1, line2)
	if parallel(line1, line2) {
		return
	}
	if perpendicular(line1, line2) {
		return
	}
	if line1.Start.X > line2.End.X && line1.Start.Y > line2.End.Y {
		return
	}
	if line1.Start.X == line1.End.X && line1.Start.X == line2.Start.X && line2.Start.X == line2.End.X {
		sameVertical(intersectionsMap, lineIndex, line1, line2)
		return
	}
	if line1.Start.Y == line1.End.Y && line1.Start.Y == line2.Start.Y && line2.Start.Y == line2.End.Y {
		sameHorizontal(intersectionsMap, lineIndex, line1, line2)
		return
	}
	checkPoints(intersectionsMap, lineIndex, line1, line2)
}

func sameHorizontal(intersectionsMap map[int]map[int16]void, lineIndex int16, line1 line, line2 line) {
	if line1.End.X < line2.Start.X {
		return
	}
	if line2.End.X < line1.Start.X {
		return
	}
	// lines overlap - add the overlapping points to the intersections map
	start := maxx(line1.Start.X, line2.Start.X)
	end := minn(line1.End.X, line2.End.X)

	for x := start; x <= end; x++ {
		addPointToIntersectionsMap(intersectionsMap, lineIndex, x, line1.Start.Y)
	}
}

func sameVertical(intersectionsMap map[int]map[int16]void, lineIndex int16, line1 line, line2 line) {
	if line1.End.Y < line2.Start.Y {
		return
	}
	if line2.End.Y < line1.Start.Y {
		return
	}
	// lines overlap - add the overlapping points to the intersections map
	start := maxx(line1.Start.Y, line2.Start.Y)
	end := minn(line1.End.Y, line2.End.Y)

	for y := start; y <= end; y++ {
		addPointToIntersectionsMap(intersectionsMap, lineIndex, line1.Start.X, y)
	}
}

func maxx(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func minn(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func checkPoints(intersectionsMap map[int]map[int16]void, lineIndex int16, line1 line, line2 line) {

	// 1:horizontal, 2:vertical
	if line1.Start.Y == line1.End.Y && line2.Start.X == line2.End.X {
		if line1.Start.Y <= line2.End.Y && line1.End.X >= line2.Start.X && line1.Start.X <= line2.Start.X {
			zap.S().Debugf(" +1 line1:%v,line2:%v", line1, line2)
			// X from line2, Y from line 1
			addPointToIntersectionsMap(intersectionsMap, lineIndex, line2.Start.X, line1.Start.Y)
			return
		}
		zap.S().Debugf(" +0 line1:%v,line2:%v", line1, line2)
		return
	}

	// 1:vertical, 2:horizontal
	if line1.Start.X == line1.End.X && line2.Start.Y == line2.End.Y {
		if line1.Start.X <= line2.End.X && line1.End.Y >= line2.Start.Y && line1.Start.Y <= line2.Start.Y {
			zap.S().Debugf(" +1 line1:%v,line2:%v", line1, line2)
			// X from line1, Y from line2
			addPointToIntersectionsMap(intersectionsMap, lineIndex, line1.Start.X, line2.Start.Y)
			return
		}
		zap.S().Debugf(" +0 line1:%v,line2:%v", line1, line2)
		return
	}

	zap.S().Fatalf("programming error - shouldn't get here")
}

func addPointToIntersectionsMap(intersectionsMap map[int]map[int16]void, lineIndex int16, x int, y int) {
	point := point(x, y)
	if intersectionsMap[point] == nil {
		intersectionsMap[point] = make(map[int16]void)
	}
	intersectionsMap[point][lineIndex] = void{}
}

func perpendicular(line1, line2 line) bool {
	if line1.Start.X == line1.End.X && line2.Start.Y == line2.End.Y {
		// line1 vertical, line2 horizontal
		if line1.Start.X < line2.Start.X || line1.Start.X > line2.End.X {
			return true
		}
	}
	if line1.Start.Y == line1.End.Y && line2.Start.X == line2.End.X {
		// line1 horizontal, line2 vertical
		if line1.Start.Y < line2.Start.Y || line1.Start.Y > line2.End.Y {
			return true
		}
	}
	return false
}

func parallel(line1 line, line2 line) bool {
	if line1.Start.X == line1.End.X && line2.Start.X == line2.End.X && line1.Start.X != line2.Start.X {
		// parallel vertical
		return true
	}
	if line1.Start.Y == line1.End.Y && line2.Start.Y == line2.End.Y && line1.Start.Y != line2.Start.Y {
		// parallel horizontal
		return true
	}
	return false
}

func getLengthAndLines(position service.Coordinate, commands []service.Command) (length int, lines []line) {
	for _, command := range commands {
		var end service.Coordinate
		start := position
		length, end = getLengthAndEnd(length, position, command)
		switch command.Direction {
		case direction.North, direction.East:
			lines = append(lines, line{Start: start, End: end})
		case direction.South, direction.West:
			// positive line direction
			lines = append(lines, line{Start: end, End: start})
		}
		position = end // where the robot is now
		zap.S().Debugf("lines:%v", lines)
	}
	zap.S().Debugf("lines:%v", lines)
	return
}

const max = 100000
const min = -100000

func getLengthAndEnd(lengthSoFar int, start service.Coordinate, command service.Command) (length int, end service.Coordinate) {

	switch command.Direction {
	case direction.North:
		end.X = start.X
		end.Y = start.Y + command.Steps
		if end.Y > max {
			end.Y = max
		}
		length = lengthSoFar + (end.Y - start.Y) + 1
	case direction.South:
		end.X = start.X
		end.Y = start.Y - command.Steps
		if end.Y < min {
			end.Y = min
		}
		length = lengthSoFar + (start.Y - end.Y) + 1
	case direction.East:
		end.X = start.X + command.Steps
		if end.X > max {
			end.X = max
		}
		end.Y = start.Y
		length = lengthSoFar + (end.X - start.X) + 1
	case direction.West:
		end.X = start.X - command.Steps
		if end.X < min {
			end.X = min
		}
		end.Y = start.Y
		length = lengthSoFar + (start.X - end.X) + 1
	default:
		{
			zap.S().Fatalf("programming error, unsupported direction %", command.Direction)
		}
	}

	return
}

// point converts coordinates to a single 64 bit integer for using as a key in the map
func point(x int, y int) int {

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
