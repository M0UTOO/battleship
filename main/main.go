package main

import (
	"battleship/player"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var clear map[string]func()

func main() {
	// Clear the screen
	CallClear()
	fmt.Println("Enter your name : ")
	var name string
	var port int
	var check string
	var isOccupied []player.Coordinates
	var playerList []player.Player
	fmt.Scanln(&name)
	for {
		fmt.Println("Enter your port (between 8000 and 9000) : ")
		fmt.Scanln(&port)
		if port < 8000 || port > 9000 {
			fmt.Println("Invalid port, please enter a port between 8000 and 9000")
		}
		// checkIfPortIsFree(port)
	}
	CallClear()
	listBoats := createBoats(&isOccupied)
	player := player.Player{name, port, listBoats}

	go startServer(port, &player, &isOccupied)

	for check != "quit" {

		playerList = nil

		getPlayers(&playerList, port)

		time.Sleep(1 * time.Second)

		if len(playerList) == 0 {
			waitingForPlayers(&playerList, port)
			CallClear()
		}

		fmt.Println("This is the list of the players : ")

		for s, player := range playerList {
			s++
			fmt.Println(s, ":", player.Pseudo)
		}

		fmt.Println("Which player do you want to interact with ? (Enter the number of the player or enter 'quit' to quit the game)")
		var answer string
		fmt.Scanln(&answer)

		if answer == "quit" {
			return
		}

		i, _ := strconv.Atoi(answer)

		for s, player := range playerList {
			if s+1 == i {
				fmt.Println("What do you want to do with", player.Pseudo, "? (Enter 'board' to see his board, 'boats' to see his boats, 'hit' to attack or 'quit' to quit the game)")
				fmt.Scanln(&check)
				CallClear()
				if check == "board" {
					url := "http://localhost:" + strconv.Itoa(player.Port) + "/board"
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println(err)
					}
					defer resp.Body.Close()
					var board string
					err = json.NewDecoder(resp.Body).Decode(&board)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(board)
				} else if check == "boats" {
					url := "http://localhost:" + strconv.Itoa(player.Port) + "/boats"
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println(err)
					}
					defer resp.Body.Close()
					var boats string
					err = json.NewDecoder(resp.Body).Decode(&boats)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(boats)
				} else if check == "hit" {
					countDestroyed := 0
					for boats := range listBoats {
						isAlive := false
						for boatParts := range listBoats[boats].BoatParts {
							if listBoats[boats].BoatParts[boatParts] == 0 {
								isAlive = true
							}
						}
						if isAlive == false {
							countDestroyed++
						}
					}
					countDestroyedPlayer := 0
					for boats := range listBoats {
						isAlive := false
						for boatParts := range player.Boats[boats].BoatParts {
							if player.Boats[boats].BoatParts[boatParts] == 0 {
								isAlive = true
							}
						}
						if isAlive == false {
							countDestroyedPlayer++
						}
					}
					fmt.Println(strconv.Itoa(countDestroyed) + "For ennemy")
					fmt.Println(strconv.Itoa(countDestroyed) + "For u")
					if countDestroyedPlayer == len(player.Boats) {
						fmt.Println("All your boats are destroyed, you lost the game, you cannot attack anymore")
					} else {
						var x int = 0
						var y int = 0
						for x < 1 || x > 10 {
							fmt.Println("Enter the coordinate X of the attack (between 1 and 10) : ")
							fmt.Scanln(&x)
							if x < 1 || x > 10 {
								fmt.Println("Invalid coordinate, please enter a coordinate between 1 and 10")
							}
						}
						for y < 1 || y > 10 {
							fmt.Println("Enter the coordinate Y of the attack (between 1 and 10) : ")
							fmt.Scanln(&y)
							if y < 1 || y > 10 {
								fmt.Println("Invalid coordinate, please enter a coordinate between 1 and 10")
							}
						}
						url := "http://localhost:" + strconv.Itoa(player.Port) + "/hit"
						body := []byte(`{"x":` + strconv.Itoa(x) + `,"y":` + strconv.Itoa(y) + `}`)
						resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
						if err != nil {
							fmt.Println(err)
						}
						defer resp.Body.Close()
						var hit string
						err = json.NewDecoder(resp.Body).Decode(&hit)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(hit)
					}
				} else {
					fmt.Println("Invalid answer, please enter 'board', 'boats', 'hit' or 'quit'")
				}
			}
		}
	}
}

// func checkIfPortIsFree(port int) {
// 	url := "http://localhost:" + strconv.Itoa(port) + "/isFree"

// }

func waitingForPlayers(playerList *[]player.Player, port int) {
	for len(*playerList) == 0 {
		CallClear()
		fmt.Println("Waiting for players...")
		getPlayers(playerList, port)
		time.Sleep(10 * time.Second)
	}
}

func getPlayers(playerList *[]player.Player, port int) {
	for i := 8000; i < 9001; i++ {
		if i != port {
			go raw_connect("localhost", strconv.Itoa(i), playerList)
		}
	}
}

func raw_connect(host string, port string, playerList *[]player.Player) error {
	url := "http://" + host + ":" + port + "/get-player"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var player player.Player
	err = json.NewDecoder(resp.Body).Decode(&player)
	if err != nil {
		return err
	}
	if player.Pseudo != "" {
		*playerList = append(*playerList, player)
	}
	return nil
}

func getPlayer(user *player.Player) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			jsonResp, _ := json.Marshal(user)
			w.Write(jsonResp)
		}
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func board(user *player.Player, isOccupied *[]player.Coordinates) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			txt := ("\n(O = Not Discovered | W = Water | X = Boat Hit)\n")
			txt += ("Board of " + user.Pseudo + ":\n")
			var isBoat bool
			for i := 1; i < 11; i++ {
				txt += ("|")
				for j := 1; j < 11; j++ {
					isBoat = false
					for _, c := range *isOccupied {
						if c.X == j && c.Y == i {
							if c.BoatName == "Water" {
								txt += ("W")
								isBoat = true
							} else {
								for _, boat := range user.Boats {
									if boat.Name == c.BoatName {
										if boat.BoatParts[c.BoatPart] == 0 {
											txt += ("O")
											isBoat = true
										} else if boat.BoatParts[c.BoatPart] == 2 {
											txt += ("X")
											isBoat = true
										}
									}
								}
							}
						}
					}
					if !isBoat {
						txt += ("O")
					}
					txt += ("|")
				}
				txt += ("\n")
			}
			jsonResp, _ := json.Marshal(txt)
			w.Write(jsonResp)
		}
	}
}

func boats(user *player.Player) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var isAlive bool
			var count int
			txt := ("\nBoats of " + user.Pseudo + ":\n")
			for _, boat := range user.Boats {
				isAlive = false
				for _, part := range boat.BoatParts {
					if part == 0 {
						isAlive = true
					}
				}
				if !isAlive {
					txt += (boat.Name + ": Destroyed\n")
				} else {
					txt += (boat.Name + ": Alive\n")
					count++
				}
			}
			if count == 0 {
				txt += ("The player have lost all his boats !\n")
			} else {
				txt += ("The player still have " + strconv.Itoa(count) + " boats alive !\n")
			}
			jsonResp, _ := json.Marshal(txt)
			w.Write(jsonResp)
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
			txt := ("")
			for _, c := range *isOccupied {
				if c.X == hit.X && c.Y == hit.Y {
					for _, boat := range user.Boats {
						if boat.Name == c.BoatName {
							if boat.BoatParts[c.BoatPart] == 0 {
								boat.BoatParts[c.BoatPart] = 2
								txt += ("You hit a " + c.BoatName + "\n")
							} else if boat.BoatParts[c.BoatPart] == 2 {
								txt += ("You already hit this part of the boat\n")
							}
						}
					}
				}
				if c.BoatName == "Water" && c.X == hit.X && c.Y == hit.Y {
					txt += ("You, missed your shot landed in the water\n")
				}
			}
			if txt == "" {
				txt += ("You missed, your shot landed in the water\n")
				hit.BoatName = "Water"
				hit.BoatPart = 0
			}
			*isOccupied = append(*isOccupied, hit)
			jsonResp, _ := json.Marshal(txt)
			w.Write(jsonResp)
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
