package intersection

import (
	"fmt"
	"os"
	"testing"

	"github.com/martin-flower/roboz-web-go/service"
	"github.com/martin-flower/roboz-web-go/service/clean"
	"github.com/martin-flower/roboz-web-go/service/direction"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	clean.Init()
	code := m.Run()
	os.Exit(code)
}

// go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof -http=localhost:8080 cpu.out
// go tool pprof -http=localhost:8080 mem.out

func Benchmark(b *testing.B) {
	start := clean.GetRandomStart()
	commands := clean.GetCommandsWithRandomSteps(200)

	// function under test
	Cleaner{}.Clean(start, commands)
}

func TestClean(t *testing.T) {
	clean.CleanTest(t, Cleaner{})
}

// about 1 second ..
func Test50Commands(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 50)
}

// takes about 3 seconds
func Test150Commands(t *testing.T) {
	clean.CommandsTest(t, Cleaner{}, 150)
}

// maximum number of commands according to the specification
// each with random number of steps up to 100000
// about 30s
func Test10000Commands(t *testing.T) {
	t.Skipf("need to run manually with go test -run Test10000Commands -timeout 20h -v")
	start := clean.GetRandomStart()
	commands := clean.GetCommandsWithRandomSteps(10000)

	// function under test
	cleaned := Cleaner{}.Clean(start, commands)

	// no useful assertion with random numbers
	assert.Greater(t, cleaned, 100000000)
}

// -- -- --

// about 3s
func TestCleanHorizontal5000(t *testing.T) {
	commands := []service.Command{}
	for i := 0; i < 5000; i++ {
		command := service.Command{}
		command.Direction = direction.East
		command.Steps = 11
		commands = append(commands, command)
		command = service.Command{}
		command.Direction = direction.West
		command.Steps = 10
		commands = append(commands, command)
	}
	assert.Equal(t, 5011, Cleaner{}.Clean(service.Coordinate{X: 0, Y: 0}, commands))
}

func TestCleanLimits(t *testing.T) {

	commands := []service.Command{}
	command := service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 100001, Cleaner{}.Clean(service.Coordinate{X: 0, Y: 0}, commands))

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 2, len(commands))
	assert.Equal(t, 100001, Cleaner{}.Clean(service.Coordinate{X: 0, Y: 0}, commands))

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 3, len(commands))
	assert.Equal(t, 100001, Cleaner{}.Clean(service.Coordinate{X: 0, Y: 0}, commands))

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 4, len(commands))
	assert.Equal(t, 100001, Cleaner{}.Clean(service.Coordinate{X: 0, Y: 0}, commands))
}

func TestGetLengthAndIntersectionsForCommandsLimits(t *testing.T) {

	commands := []service.Command{}
	command := service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	l, i := getLengthAndIntersectionsForCommands(service.Coordinate{X: 0, Y: 0}, commands)
	assert.Equal(t, 100001, l)
	assert.Equal(t, 0, i)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	l, i = getLengthAndIntersectionsForCommands(service.Coordinate{X: 0, Y: 0}, commands)
	assert.Equal(t, 100002, l)
	assert.Equal(t, 1, i)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 3, len(commands))
	l, i = getLengthAndIntersectionsForCommands(service.Coordinate{X: 0, Y: 0}, commands)
	assert.Equal(t, 100003, l)
	assert.Equal(t, 2, i)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)
	assert.Equal(t, 4, len(commands))
	l, i = getLengthAndIntersectionsForCommands(service.Coordinate{X: 0, Y: 0}, commands)
	assert.Equal(t, 100004, l)
	assert.Equal(t, 3, i)
}

func TestGetIntersection(t *testing.T) {

	intersectionsMap := map[int]map[int16]void{}
	getIntersection(
		intersectionsMap, 0,
		line{service.Coordinate{X: 0, Y: 0}, service.Coordinate{X: 0, Y: 0}},
		[]line{
			{service.Coordinate{X: 0, Y: 0}, service.Coordinate{X: 0, Y: 0}},
			{service.Coordinate{X: 0, Y: 0}, service.Coordinate{X: 0, Y: 0}},
			{service.Coordinate{X: 0, Y: 0}, service.Coordinate{X: 0, Y: 0}},
			{service.Coordinate{X: 0, Y: 0}, service.Coordinate{X: 0, Y: 0}},
		},
	)
	assert.Len(t, intersectionsMap, 1)
	numberOfLines := 0
	for _, value := range intersectionsMap {
		numberOfLines = numberOfLines + len(value)
	}
	assert.Equal(t, 1, numberOfLines)
}

func TestIntersectionOfTwoLines(t *testing.T) {

	type test struct {
		expectedPoints int
		expectedLines  int
		separator1     string
		line1StartX    int
		line1StartY    int
		line1EndX      int
		line1EndY      int
		separator2     string
		line2StartX    int
		line2StartY    int
		line2EndX      int
		line2EndY      int
	}

	tests := []test{
		{1, 1, "", 2, 0, 2, 5, "", 0, 1, 3, 1},  // vertical, horizontal
		{0, 0, "", 0, 0, 2, 0, "", 0, 1, 2, 1},  // parallel
		{1, 1, "", 0, 0, 0, 0, "", 0, 0, 0, 0},  // one point, length one
		{1, 1, "", 0, 0, 1, 0, "", 1, 0, 2, 0},  // one point
		{2, 2, "", 0, 0, 1, 0, "", 0, 0, 2, 0},  // same horizontal
		{2, 2, "", 0, 0, 1, 0, "", 0, 0, 3, 0},  // same horizontal
		{2, 2, "", 0, -1, 0, 1, "", 0, 0, 0, 2}, // same vertical
		{1, 1, "", 0, -1, 0, 1, "", 0, 1, 0, 9}, // same vertical
	}

	for index, tc := range tests {
		intersectionsMap := map[int]map[int16]void{}
		getIntersectionOfTwoLines(
			intersectionsMap, 0,
			line{service.Coordinate{X: tc.line1StartX, Y: tc.line1StartY}, service.Coordinate{X: tc.line1EndX, Y: tc.line1EndY}},
			line{service.Coordinate{X: tc.line2StartX, Y: tc.line2StartY}, service.Coordinate{X: tc.line2EndX, Y: tc.line2EndY}})

		assert.Len(t, intersectionsMap, tc.expectedPoints, fmt.Sprintf("failed test case index:%d", index))
		numberOfLines := 0
		for _, value := range intersectionsMap {
			numberOfLines = numberOfLines + len(value)
		}
		assert.Equal(t, tc.expectedLines, numberOfLines, fmt.Sprintf("failed test case index:%d", index))
	}
}

func TestGetIntersectionForLines(t *testing.T) {

	type test struct {
		expected    int
		separator1  string
		line1StartX int
		line1StartY int
		line1EndX   int
		line1EndY   int
		separator2  string
		line2StartX int
		line2StartY int
		line2EndX   int
		line2EndY   int
		separator3  string
		line3StartX int
		line3StartY int
		line3EndX   int
		line3EndY   int
		separator4  string
		line4StartX int
		line4StartY int
		line4EndX   int
		line4EndY   int
		separator5  string
		line5StartX int
		line5StartY int
		line5EndX   int
		line5EndY   int
	}

	tests := []test{
		{4, "", 1, 1, 4, 1, "", 2, -1, 2, 4, "", 1, 3, 5, 3, "", 4, 3, 4, 3, "", 5, 2, 5, 4},
	}

	for index, tc := range tests {
		assert.Equal(t, tc.expected, getIntersectionsForLines(
			[]line{
				{service.Coordinate{X: tc.line1StartX, Y: tc.line1StartY}, service.Coordinate{X: tc.line1EndX, Y: tc.line1EndY}},
				{service.Coordinate{X: tc.line2StartX, Y: tc.line2StartY}, service.Coordinate{X: tc.line2EndX, Y: tc.line2EndY}},
				{service.Coordinate{X: tc.line3StartX, Y: tc.line3StartY}, service.Coordinate{X: tc.line3EndX, Y: tc.line3EndY}},
				{service.Coordinate{X: tc.line4StartX, Y: tc.line4StartY}, service.Coordinate{X: tc.line4EndX, Y: tc.line4EndY}},
				{service.Coordinate{X: tc.line5StartX, Y: tc.line5StartY}, service.Coordinate{X: tc.line5EndX, Y: tc.line5EndY}},
			},
		),
			fmt.Sprintf("failed test case index:%d", index),
		)
	}
}

func TestLengthAndIntersectionsForCommands(t *testing.T) {

	type test struct {
		length        int
		intersections int
		start         service.Coordinate
		commands      []service.Command
	}

	tests := []test{
		{7, 1, service.Coordinate{X: 2, Y: 5}, []service.Command{{Direction: direction.North, Steps: 3}, {Direction: direction.East, Steps: 2}}},
		{14, 3, service.Coordinate{X: 2, Y: 5}, []service.Command{
			{Direction: direction.North, Steps: 1},
			{Direction: direction.East, Steps: 2},
			{Direction: direction.South, Steps: 3},
			{Direction: direction.West, Steps: 4}},
		},
		{12, 4, service.Coordinate{X: 0, Y: 0}, []service.Command{
			{Direction: direction.North, Steps: 2},
			{Direction: direction.East, Steps: 2},
			{Direction: direction.South, Steps: 2},
			{Direction: direction.West, Steps: 2}},
		},
	}

	for index, tc := range tests {
		l, i := getLengthAndIntersectionsForCommands(tc.start, tc.commands)
		assert.Equal(t, tc.length, l, fmt.Sprintf("failed length, test case index: %d", index))
		assert.Equal(t, tc.intersections, i, fmt.Sprintf("failed intersections, test case index: %d", index))
	}
}

func TestLengthAndEnd(t *testing.T) {
	length, end := getLengthAndEnd(0, service.Coordinate{X: 0, Y: 0}, service.Command{Direction: direction.North, Steps: 100000})
	assert.Equal(t, 100001, length)
	assert.Equal(t, service.Coordinate{X: 0, Y: 100000}, end)

	length, end = getLengthAndEnd(0, service.Coordinate{X: 100000, Y: 100000}, service.Command{Direction: direction.North, Steps: 100000})
	assert.Equal(t, 1, length)
	assert.Equal(t, service.Coordinate{X: 100000, Y: 100000}, end)

	length, end = getLengthAndEnd(0, service.Coordinate{X: 100000, Y: 100000}, service.Command{Direction: direction.East, Steps: 100000})
	assert.Equal(t, 1, length)
	assert.Equal(t, service.Coordinate{X: 100000, Y: 100000}, end)
}

func TestLengthAndLines(t *testing.T) {

	commands := []service.Command{}
	command := service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)

	command = service.Command{}
	command.Direction = direction.North
	command.Steps = 100000
	commands = append(commands, command)

	length, lines := getLengthAndLines(service.Coordinate{X: 0, Y: 0}, commands)

	assert.Len(t, lines, len(commands))
	assert.Equal(t, 0, lines[0].Start.X)
	assert.Equal(t, 0, lines[0].End.X)
	assert.Equal(t, 0, lines[0].Start.Y)
	assert.Equal(t, 100000, lines[0].End.Y)
	for i := 1; i < len(lines); i++ {
		assert.Equal(t, 0, lines[i].Start.X)
		assert.Equal(t, 0, lines[i].End.X)
		assert.Equal(t, 100000, lines[i].Start.Y)
		assert.Equal(t, 100000, lines[i].End.Y)
	}
	assert.Equal(t, 100000+len(commands), length)
}
