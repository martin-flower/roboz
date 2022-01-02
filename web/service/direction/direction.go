package direction

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func GetDirections() []Direction {
	return []Direction{North, South, East, West}
}

func Valid(s string) bool {
	switch s {
	case "north", "south", "east", "west":
		return true
	}
	return false
}

func FromString(s string) Direction {
	switch s {
	case "north":
		return North
	case "south":
		return South
	case "east":
		return East
	case "west":
		return West
	default:
		panic("invalid direction %s")
	}
}

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	}
	return "unknown"
}
