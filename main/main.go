package main

import (
	"battleship/player"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Entre your name : ")
	var name string
	var port int
	var isOccupied []player.Coordinates
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
	listBoats := createBoats(&isOccupied)
	player := player.Player{name, port, listBoats}

	go startServer(port, &player, &isOccupied)

	fmt.Println("Do you want to send a request ? (y/n)")
	var answer string
	fmt.Scanln(&answer)
	if answer == "y" {
		for i := 8000; i < 9001; i++ {
			if i != port {
				go raw_connect("localhost", strconv.Itoa(i))
			}
		}
	}
	time.Sleep(100 * time.Second)
}

func raw_connect(host string, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return
	}
	if conn != nil {
		defer conn.Close()
		url := "http://" + host + ":" + port + "/get-player"
		req, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer req.Body.Close()
	}
}

func getPlayer(user *player.Player) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Printf("ok")
			res, _ := json.Marshal(user)
		
		}
	}
}

func board(user *player.Player, isOccupied *[]player.Coordinates) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprintln(w, "(O = Not Discovered | W = Water | X = Boat Hit)")
			fmt.Fprintln(w, "Board of ", user.Pseudo, ":")
			var isBoat bool
			for i := 1; i < 11; i++ {
				fmt.Fprint(w, "|")
				for j := 1; j < 11; j++ {
					isBoat = false
					for _, c := range *isOccupied {
						if c.X == j && c.Y == i {
							if c.BoatName == "Water" {
								fmt.Fprint(w, "W")
								isBoat = true
							} else {
								for _, boat := range user.Boats {
									if boat.Name == c.BoatName {
										if boat.BoatParts[c.BoatPart] == 0 {
											fmt.Fprint(w, "O")
											isBoat = true
										} else if boat.BoatParts[c.BoatPart] == 2 {
											fmt.Fprint(w, "X")
											isBoat = true
										}
									}
								}
							}
						}
					}
					if !isBoat {
						fmt.Fprint(w, "O")
					}
					fmt.Fprint(w, "|")
				}
				fmt.Fprintln(w, "")
			}
		}
	}
}

func boats(user *player.Player) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var isAlive bool
			var count int
			fmt.Fprintln(w, "Boats of", user.Pseudo, ":")
			for _, boat := range user.Boats {
				isAlive = false
				for _, part := range boat.BoatParts {
					if part == 0 {
						isAlive = true
					}
				}
				if !isAlive {
					fmt.Fprintln(w, boat.Name, ": Destroyed")
				} else {
					fmt.Fprintln(w, boat.Name, ": Alive")
					count++
				}
			}
			if count == 0 {
				fmt.Fprintln(w, "The player have lost all his boats !")
			} else {
				fmt.Fprintln(w, "The player have", count, "boats alive !")
			}
		}
	}
}

func checkIsAlive(user *player.Player) bool {
	var isAlive bool
	for _, boat := range user.Boats {
		isAlive = false
		for _, part := range boat.BoatParts {
			if part == 0 {
				isAlive = true
			}
		}
		if isAlive {
			return true
		}
	}
	return false
}

func hit(user *player.Player, isOccupied *[]player.Coordinates) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var hit player.Coordinates
			err := json.NewDecoder(r.Body).Decode(&hit)
			if err != nil {
				fmt.Println(err)
			}
			if hit.X < 1 || hit.X > 10 || hit.Y < 1 || hit.Y > 10 {
				fmt.Fprintln(w, "Invalid coordinates")
				return
			}
			for _, c := range *isOccupied {
				if c.X == hit.X && c.Y == hit.Y {
					for _, boat := range user.Boats {
						if boat.Name == c.BoatName {
							if boat.BoatParts[c.BoatPart] == 0 {
								boat.BoatParts[c.BoatPart] = 2
								fmt.Fprintln(w, "You hit a", c.BoatName)
								return
							} else if boat.BoatParts[c.BoatPart] == 2 {
								fmt.Fprintln(w, "You already hit this part of the boat")
								return
							}
						} else if c.BoatName == "Water" {
							fmt.Fprintln(w, "You shot in the water")
							return
						}
					}
				}
			}
			fmt.Fprintln(w, "You shot in the water")
			hit.BoatName = "Water"
			hit.BoatPart = 0
			*isOccupied = append(*isOccupied, hit)
		}
	}
}

// func sendRequest(url string, method string) {
// 	req, err := http.Post(url, "application/json", nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer req.Body.Close()

// 	// var player player.Player
// 	var hitReq player.HitReq
// 	err = json.NewDecoder(req.Body).Decode(&hitReq)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func startServer(port int, player *player.Player, isOccupied *[]player.Coordinates) {
	http.HandleFunc("/get-player", getPlayer(player))
	http.HandleFunc("/board", board(player, isOccupied))
	http.HandleFunc("/boats", boats(player))
	http.HandleFunc("/hit", hit(player, isOccupied))
	address := strings.Join([]string{":", strconv.Itoa(port)}, "")
	http.ListenAndServe(address, nil)
}

func createBoats(isOccupied *[]player.Coordinates) []player.Boat {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var boats []player.Boat

	// Carrier
	carrier := createBoat("Carrier", 5, r, isOccupied)
	boats = append(boats, carrier)

	// Battleship
	battleship := createBoat("Battleship", 4, r, isOccupied)
	boats = append(boats, battleship)

	// Cruiser
	cruiser := createBoat("Cruiser", 3, r, isOccupied)
	boats = append(boats, cruiser)

	// Submarine
	submarine := createBoat("Submarine", 3, r, isOccupied)
	boats = append(boats, submarine)

	// Destroyer
	destroyer := createBoat("Destroyer", 2, r, isOccupied)
	boats = append(boats, destroyer)

	return boats
}

func createBoat(name string, size int, r *rand.Rand, isOccupied *[]player.Coordinates) player.Boat {
	var boat player.Boat
	boat.Name = name
	boat.Size = size
	boat.Direction = r.Intn(4)
	boat.StartingCoordinates.X, boat.StartingCoordinates.Y = getposition(boat.Direction, boat.Size, isOccupied, name)
	boat.BoatParts = make([]int, boat.Size)
	for i := 0; i < boat.Size; i++ {
		boat.BoatParts[i] = 0
	}
	return boat
}

func getposition(direction, size int, isOccupied *[]player.Coordinates, name string) (int, int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var x int
	var y int
	var coord player.Coordinates
	var isGoodPositon bool = false
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
			y = r.Intn(max) + size
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied, direction, size, name)
		} else if direction == 1 {
			x = r.Intn(max) + 1
			y = r.Intn(10) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied, direction, size, name)
		} else if direction == 2 {
			x = r.Intn(10) + 1
			y = r.Intn(max) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied, direction, size, name)
		} else {
			x = r.Intn(max) + size
			y = r.Intn(10) + 1
			coord.X = x
			coord.Y = y
			isGoodPositon = isAlreadyOccupied(coord, isOccupied, direction, size, name)
		}
	}
	return x, y
}

func isAlreadyOccupied(coord player.Coordinates, isOccupied *[]player.Coordinates, direction, size int, name string) bool {
	tmp := player.Coordinates{coord.X, coord.Y, "", 0}
	if direction == 0 {
		for i := 0; i < size; i++ {
			for _, c := range *isOccupied {
				if c.X == tmp.X && c.Y == tmp.Y {
					return false
				}
			}
			tmp.Y = tmp.Y - 1
		}
		for i := 0; i < size; i++ {
			coord.BoatName = name
			coord.BoatPart = i
			*isOccupied = append(*isOccupied, coord)
			coord.Y = coord.Y - 1
		}
	} else if direction == 1 {
		for i := 0; i < size; i++ {
			for _, c := range *isOccupied {
				if c.X == tmp.X && c.Y == tmp.Y {
					return false
				}
			}
			tmp.X = tmp.X + 1
		}
		for i := 0; i < size; i++ {
			coord.BoatName = name
			coord.BoatPart = i
			*isOccupied = append(*isOccupied, coord)
			coord.X = coord.X + 1
		}
	} else if direction == 2 {
		for i := 0; i < size; i++ {
			for _, c := range *isOccupied {
				if c.X == tmp.X && c.Y == tmp.Y {
					return false
				}
			}
			tmp.Y = tmp.Y + 1
		}
		for i := 0; i < size; i++ {
			coord.BoatName = name
			coord.BoatPart = i
			*isOccupied = append(*isOccupied, coord)
			coord.Y = coord.Y + 1
		}
	} else {
		for i := 0; i < size; i++ {
			for _, c := range *isOccupied {
				if c.X == tmp.X && c.Y == tmp.Y {
					return false
				}
			}
			tmp.X = tmp.X - 1
		}
		for i := 0; i < size; i++ {
			coord.BoatName = name
			coord.BoatPart = i
			*isOccupied = append(*isOccupied, coord)
			coord.X = coord.X - 1
		}
	}
	return true
}
