package player

func Hello() {
	println("hello")
}

type Player struct {
	Pseudo            string
	Port              int
	Boats             []Boat
	OGrenade          bool
	XGrenade          bool
	Nuke              bool
	Heal              bool
	Combo             int
	BoatPartDestroyed int
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

type HitReq struct {
	Boats    []Boat        `json:"boats"`
	BoatsMap []Coordinates `json:"boatsMap"`
}
