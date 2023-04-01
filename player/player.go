package player

func Hello() {
	println("hello")
}

type Player struct {
	Pseudo string
	Port   int
	Boats  []Boat
}

type Boat struct {
	Name                string
	Size                int
	Direction           int
	StartingCoordinates Coordinates
	BoatParts           []int
}

type Coordinates struct {
	X        int
	Y        int
	BoatName string
	BoatPart int
}
