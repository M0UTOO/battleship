package main

import (
	"battleship/player"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Entre your name : ")
	var name string
	var port int
	fmt.Scanln(&name)
	for {
		fmt.Println("Enter your port (between 8000 and 9000) : ")
		fmt.Scanln(&port)
		if port < 8000 || port > 9000 {
			fmt.Println("Invalid port, please enter a port between 8000 and 9000")
		} else {
			break
		}
	}
	listBoats := createBoats()
	player := player.Player{name, port, listBoats}

	go startServer(port, &player)

	fmt.Println("Do you want to send a request ? (y/n)")
	var answer string
	fmt.Scanln(&answer)
	if answer == "y" {
		sendRequest("http://localhost:8000/hit", "POST")
	}
}

func board(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		player.Hello()
	}
}

func boats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		player.Hello()
	}
}

func hit(player *player.Player) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			req, err := json.Marshal(player)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(req)
		}
	}
}

func sendRequest(url string, method string) {
	req, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()

	var player player.Player
	err = json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		fmt.Println(err)
	}
	for _, boat := range player.Boats {
		fmt.Println(boat.Name)
		fmt.Println(boat.Size)
		fmt.Println(boat.Direction)
		fmt.Println(boat.StartingCoordinates)
		fmt.Println(boat.BoatParts)
	}
}

func startServer(port int, player *player.Player) {
	http.HandleFunc("/board", board)
	http.HandleFunc("boats", boats)
	http.HandleFunc("/hit", hit(player))
	address := strings.Join([]string{":", strconv.Itoa(port)}, "")
	http.ListenAndServe(address, nil)
}

func createBoats() []player.Boat {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var boats []player.Boat
	var isOccupied []player.Coordinates

	// Carrier
	carrier := createBoat("Carrier", 5, r, &isOccupied)
	boats = append(boats, carrier)

	// Battleship
	battleship := createBoat("Battleship", 4, r, &isOccupied)
	boats = append(boats, battleship)

	// Cruiser
	cruiser := createBoat("Cruiser", 3, r, &isOccupied)
	boats = append(boats, cruiser)

	// Submarine
	submarine := createBoat("Submarine", 3, r, &isOccupied)
	boats = append(boats, submarine)

	// Destroyer
	destroyer := createBoat("Destroyer", 2, r, &isOccupied)
	boats = append(boats, destroyer)

	return boats
}

func createBoat(name string, size int, r *rand.Rand, isOccupied *[]player.Coordinates) player.Boat {
	var boat player.Boat
	boat.Name = name
	boat.Size = size
	boat.Direction = r.Intn(4)
	boat.StartingCoordinates.X, boat.StartingCoordinates.Y = getposition(boat.Direction, boat.Size, isOccupied)
	boat.BoatParts = make([]int, boat.Size)
	for i := 0; i < boat.Size; i++ {
		boat.BoatParts[i] = 0
	}
	return boat
}

func getposition(direction, size int, isOccupied *[]player.Coordinates) (int, int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var x int
	var y int
	var coord player.Coordinates
	var isGoodPositon bool
	max := 0
	if size == 5 {
		max = 6
	} else if size == 4 {
		max = 7
	} else if size == 3 {
		max = 8
	} else {
		max = 9
	}
	for !isGoodPositon {
		if direction == 0 {
			x = r.Intn(10) + 1
			y = r.Intn(max) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied)
		} else if direction == 1 {
			x = r.Intn(max) + size
			y = r.Intn(10) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied)
		} else if direction == 2 {
			x = r.Intn(10) + 1
			y = r.Intn(max) + size
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied)
		} else {
			x = r.Intn(max) + 1
			y = r.Intn(10) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied)
		}
	}
	return x, y
}

func isAlreadyOccupied(coord player.Coordinates, isOccupied *[]player.Coordinates) bool {
	for _, c := range *isOccupied {
		if c.X == coord.X && c.Y == coord.Y {
			return false
		}
	}
	*isOccupied = append(*isOccupied, coord)
	return true
}
