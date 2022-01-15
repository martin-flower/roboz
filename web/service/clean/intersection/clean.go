package intersection

import (
	"github.com/martin-flower/roboz-web/service"
	"github.com/martin-flower/roboz-web/service/direction"
	"go.uber.org/zap"
)

type Cleaner struct{}

var parallelc int
var perpendicularc int
var checkpoints1 int
var checkpoints0 int
var quadrantc int
var sameverticalc int
var samehorizontalc int

// Clean implementation using intersections of lines
// Load all the lines from the commands
// Count the number of points where the lines crossor touch each other
func (c Cleaner) Clean(start service.Coordinate, commands []service.Command) (cleaned int) {

	length, intersections := getLengthAndIntersectionsForCommands(start, commands)

	zap.S().Infof("parallel:%d, perpendicular:%d, quadrantc:%d, sameverticalc:%d, samehorizontalc:%d, checkpoints0:%d, checkpoints1:%d", parallelc, perpendicularc, quadrantc, sameverticalc, samehorizontalc, checkpoints0, checkpoints1)

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

func getIntersectionsForLines(lines []line) (intersections int) {
	intersectionsMap := make(map[int][]int, 10000*len(lines))
	for lineIndex := 0; lineIndex < len(lines)-1; lineIndex++ {
		zap.S().Debugf("checking line %d of %d for intersections", lineIndex+1, len(lines))
		getIntersection(intersectionsMap, lineIndex+1, lines[lineIndex], lines[lineIndex+1:])
	}
	for _, lines := range intersectionsMap {
		intersections = intersections + len(lines)
	}
	return
}

func getIntersection(intersectionsMap map[int][]int, lineIndex int, line line, lines []line) {
	if len(lines) < 1 {
		return
	}
	for i := 0; i < len(lines); i++ {
		getIntersectionOfTwoLines(intersectionsMap, lineIndex, line, lines[i])
	}
}

func getIntersectionOfTwoLines(intersectionsMap map[int][]int, lineIndex int, line1 line, line2 line) {
	zap.S().Debugf("%v,%v", line1, line2)
	if parallel(line1, line2) {
		parallelc++
		return
	}
	if perpendicular(line1, line2) {
		perpendicularc++
		return
	}
	if line1.Start.X > line2.End.X && line1.Start.Y > line2.End.Y {
		quadrantc++
		return
	}
	if line1.Start.X == line1.End.X && line1.Start.X == line2.Start.X && line2.Start.X == line2.End.X {
		sameverticalc++
		sameVertical(intersectionsMap, lineIndex, line1, line2)
		return
	}
	if line1.Start.Y == line1.End.Y && line1.Start.Y == line2.Start.Y && line2.Start.Y == line2.End.Y {
		samehorizontalc++
		sameHorizontal(intersectionsMap, lineIndex, line1, line2)
		return
	}
	checkPoints(intersectionsMap, lineIndex, line1, line2)
}

func sameHorizontal(intersectionsMap map[int][]int, lineIndex int, line1 line, line2 line) {
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
		point := point(x, line1.Start.Y)
		lineIndexes := intersectionsMap[point]
		if !contains(lineIndexes, lineIndex) {
			intersectionsMap[point] = append(lineIndexes, lineIndex)
		}
	}
}

func sameVertical(intersectionsMap map[int][]int, lineIndex int, line1 line, line2 line) {
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
		point := point(line1.Start.X, y)
		lineIndexes := intersectionsMap[point]
		if !contains(lineIndexes, lineIndex) {
			update(intersectionsMap, point, lineIndexes, lineIndex)
		}
	}
}

func update(intersectionsMap map[int][]int, point int, lineIndexes []int, lineIndex int) {
	intersectionsMap[point] = append(lineIndexes, lineIndex)
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

func checkPoints(intersectionsMap map[int][]int, lineIndex int, line1 line, line2 line) {

	// 1:horizontal, 2:vertical
	if line1.Start.Y == line1.End.Y && line2.Start.X == line2.End.X {
		if line1.Start.Y <= line2.End.Y && line1.End.X >= line2.Start.X && line1.Start.X <= line2.Start.X {
			checkpoints1++
			zap.S().Debugf(" +1 line1:%v,line2:%v", line1, line2)
			// X from line2, Y from line 1
			point := point(line2.Start.X, line1.Start.Y)
			lineIndexes := intersectionsMap[point]
			if contains(lineIndexes, lineIndex) {
				return
			}
			lineIndexes = append(lineIndexes, lineIndex)
			intersectionsMap[point] = lineIndexes
			return
		}
		checkpoints0++
		zap.S().Debugf(" +0 line1:%v,line2:%v", line1, line2)
		return
	}

	// 1:vertical, 2:horizontal
	if line1.Start.X == line1.End.X && line2.Start.Y == line2.End.Y {
		if line1.Start.X <= line2.End.X && line1.End.Y >= line2.Start.Y && line1.Start.Y <= line2.Start.Y {
			checkpoints1++
			zap.S().Debugf(" +1 line1:%v,line2:%v", line1, line2)
			// X from line1, Y from line2
			point := point(line1.Start.X, line2.Start.Y)
			lineIndexes := intersectionsMap[point]
			if contains(lineIndexes, lineIndex) {
				return
			}
			lineIndexes = append(lineIndexes, lineIndex)
			intersectionsMap[point] = lineIndexes
			return
		}
		checkpoints0++
		zap.S().Debugf(" +0 line1:%v,line2:%v", line1, line2)
		return
	}

	zap.S().Fatalf("programming error - shouldn't get here")
}

func contains(numbers []int, number int) bool {
	for _, n := range numbers {
		if n == number {
			return true
		}
	}
	return false
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

// point converts coordinates to a 64 bit integer for using as a key in the map.
// max int64  : -9223372036854775808 to 9223372036854775807
// squareroot of 9223372036854775807 is 3037000500
// x, y between -100000 and 100000
// small enough for int64, but too big for int32
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
