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
	createBoats()
	player := player.Player{name, port}

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
	fmt.Println(player)
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

	// Carrier
	var carrier player.Boat
	carrier.Name = "Carrier"
	carrier.Size = 5
	carrier.Direction = r.Intn(4)
	carrier.StartingCoordinates.X, carrier.StartingCoordinates.Y = getposition(carrier.Direction, carrier.Size)
	if carrier.Direction == 0 {
		carrier.StartingCoordinates.X = r.Intn(10) + 1
		carrier.StartingCoordinates.Y = r.Intn(6) + 5
	} else if carrier.Direction == 1 {
		carrier.StartingCoordinates.X = r.Intn(6) + 1
		carrier.StartingCoordinates.Y = r.Intn(10) + 1
	} else if carrier.Direction == 2 {
		carrier.StartingCoordinates.X = r.Intn(10) + 1
		carrier.StartingCoordinates.Y = r.Intn(6) + 1
	} else {
		carrier.StartingCoordinates.X = r.Intn(6) + 5
		carrier.StartingCoordinates.Y = r.Intn(10) + 1
	}
	carrier.BoatParts = make([]int, carrier.Size)
	for i := 0; i < carrier.Size; i++ {
		carrier.BoatParts[i] = 0
	}
	boats = append(boats, carrier)

}

func getposition(direction, size int) (int, int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// max := 10 - size
	if size == 5 {

		if direction == 0 {
			return r.Intn(10) + 1, r.Intn(6) + 5
		} else if direction == 1 {
			return r.Intn(6) + 1, r.Intn(10) + 1
		} else if direction == 2 {
			return r.Intn(10) + 1, r.Intn(6) + 1
		} else {
			return r.Intn(6) + 5, r.Intn(10) + 1
		}
	}
}
