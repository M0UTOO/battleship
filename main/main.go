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

	"github.com/enescakir/emoji"
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
	isFree := "true"
	for isFree != "false" {
		fmt.Println("Enter your port (between 8000 and 9000) : ")
		fmt.Scanln(&port)
		if port < 8000 || port > 9001 {
			fmt.Println("Invalid port, please enter a port between 8000 and 9000")
		} else {
			isFree = checkIfPortIsFree(port)
			if isFree == "true" {
				fmt.Println("This port is already used, please enter another port")
			}
		}
	}
	CallClear()
	listBoats := createBoats(&isOccupied)
	player := player.Player{name, port, listBoats, false, false, false, false, 0, 0}

	go startServer(port, &player, &isOccupied)

	for check != "quit" {

		playerList = nil

		getPlayers(&playerList, port)

		time.Sleep(1 * time.Second)

		if len(playerList) == 0 {
			waitingForPlayers(&playerList, port)
			CallClear()
		}

		for check != "menu" {
			fmt.Println("Menu :")
			fmt.Println("\"play\" - Interact with others players")
			fmt.Println("\"weapons\" - See your available weapons")
			fmt.Println("\"heal\" - Select a boat to heal")
			fmt.Println("\"board\" - See your current board")
			fmt.Println("\"quit\" - Quit the game")
			fmt.Scanln(&check)

			CallClear()

			if check == "weapons" {
				if player.XGrenade == true {
					fmt.Println("X Grenade (1) - Can be activated")
				} else {
					fmt.Println("X Grenade (0) - Need a combo of 3 to be activated")
				}
				if player.OGrenade == true {
					fmt.Println("x9 Grenade (1) - Can be activated")
				} else {
					fmt.Println("x9 Grenade (0) - Need a combo of 5 to be activated")
				}
				if player.Nuke == true {
					fmt.Println("Nuke Grenade (1) - Can be activated")
				} else {
					fmt.Println("Nuke Grenade (0) - Where do I find this...?")
				}
				if player.Heal == true {
					fmt.Println("Heal (1) - Can be activated")
				} else {
					fmt.Println("Heal (0) - Need to destory " + strconv.Itoa(10-player.Combo) + " more boats parts")
				}
				fmt.Println("")
			}

			if check == "heal" {
				if player.Heal == true {
					for check != "back" {
						isAlive := false
						count := 0
						for _, boat := range player.Boats {
							isAlive = false
							for _, part := range boat.BoatParts {
								if part == 0 {
									isAlive = true
								}
							}
							if !isAlive {
								fmt.Println(boat.Name + ": Destroyed\n")
							} else {
								fmt.Println(boat.Name + ": Alive\n")
								count++
							}
						}
						if count == 0 {
							fmt.Println("You can't heal your boats, they are all destroyed\n")
							check = "back"
						} else if count == 5 {
							fmt.Println("All your boats are alive, you don't need to heal them\n")
							check = "back"
						} else {
							fmt.Println("You still have " + strconv.Itoa(count) + " alive which one do you want to heal ?\n")
							fmt.Println("\"Carrier\" - Heal your carrier")
							fmt.Println("\"Battleship\" - Heal your battleship")
							fmt.Println("\"Cruiser\" - Heal your cruiser")
							fmt.Println("\"Destroyer\" - Heal your destroyer")
							fmt.Println("\"Submarine\" - Heal your submarine")
							fmt.Println("\"back\" - Go back to the menu")
							fmt.Scanln(&check)
							CallClear()
							if check == "Carrier" || check == "Battleship" || check == "Cruiser" || check == "Destroyer" || check == "Submarine" {
								checkIfHealable := true
								for _, boat := range player.Boats {
									if check == boat.Name {
										for _, part := range boat.BoatParts {
											if part == 0 {
												checkIfHealable = false
											}
										}
										if checkIfHealable == true {
											for _, part := range boat.BoatParts {
												boat.BoatParts[part] = 0
											}
											fmt.Println("Your " + check + " have been healed\n")
											player.Heal = false
											player.BoatPartDestroyed = 0
											check = "back"
										} else {
											fmt.Println("Your " + check + " is not destroyed\n")
										}
									}
								}
							}
						}
					}
				} else {
					fmt.Println("You don't have an available heal for now, go destroy " + strconv.Itoa(10-player.BoatPartDestroyed) + " more boat parts !\n")
				}
			}

			if check == "???" {
				CallClear()
				fmt.Println("Welcome in the secret zone. Here you can type some controller input to get some secret stuff !\n")
				for check != "back" {
					fmt.Println("This is how it works, you have all this inputs :\n")
					fmt.Println("\"U\" - Up " + emoji.UpArrow)
					fmt.Println("\"L\" - Left " + emoji.LeftArrow)
					fmt.Println("\"R\" - Right " + emoji.RightArrow)
					fmt.Println("\"D\" - Down " + emoji.DownArrow)
					fmt.Println("\"A\" - A button " + emoji.AButtonBloodType)
					fmt.Println("\"B\" - B button " + emoji.BButtonBloodType)
					fmt.Println("\"S\" - Start " + emoji.PlayButton)
					fmt.Println("")
					fmt.Println("\"back\" - Go back to the menu\n")
					fmt.Println("You can combine them to get some secret stuff !\n")
					fmt.Println("For example, if you type \"SBABA\" you will get a secret message !\n")
					fmt.Scanln(&check)

					if check == "UUDDLRLRBAS" {
						if player.Nuke == false {
							CallClear()
							fmt.Println("You found the Konami code !\n")
							fmt.Println("You can now use the nuke to make some big damage !\n")
							player.Nuke = true
						} else {
							CallClear()
							fmt.Println("You already found the Konami code !\n")
							fmt.Println("You can now use the nuke to make some big damage !\n")
						}

					} else if check == "SBABA" {
						CallClear()
						fmt.Println("You found the secret message !\n")
						fmt.Println("I hear that if u type a famous gaming code, you can get a nuke !\n")

					} else if check == "back" {
						CallClear()

					} else {
						CallClear()
						fmt.Println("This is not a valid code try again :)\n")
					}
				}
			}

			if check == "play" {

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

				fmt.Println("Which player do you want to interact with ?")
				fmt.Println("\"number of the player\" - Interact with this player")
				fmt.Println("\"menu\" - Go back to the menu")
				fmt.Println("\"quit\" - Quit the game")
				var answer string
				fmt.Scanln(&answer)

				if answer == "quit" {
					return
				}

				i, _ := strconv.Atoi(answer)

				check = ""

				for s, player := range playerList {
					if s+1 == i {
						var isConnected bool = true
						for check != "return" {
							fmt.Println(player.BoatPartDestroyed)
							fmt.Println("What do you want to do with " + player.Pseudo + ":")
							fmt.Println("\"board\" - Show the board of " + player.Pseudo)
							fmt.Println("\"boats\" - Show the boats of " + player.Pseudo)
							fmt.Println("\"hit\" - Attack " + player.Pseudo + " (you will have to enter the coordinates of the boat you want to attack)")
							fmt.Println("\"return\" - Go back to the menu")
							fmt.Println("\"quit\" - Quit the game")
							fmt.Scanln(&check)
							CallClear()
							if check == "board" || check == "Board" {
								url := "http://localhost:" + strconv.Itoa(player.Port) + "/board"
								isConnected = getRouteInfo(url, "board", player.Pseudo, nil, nil)
								if isConnected == false {
									check = "return"
								}
							} else if check == "boats" || check == "Boats" {
								url := "http://localhost:" + strconv.Itoa(player.Port) + "/boats"
								isConnected = getRouteInfo(url, "boats", player.Pseudo, nil, nil)
								if isConnected == false {
									check = "return"
								}
							} else if check == "hit" || check == "Hit" {
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
								if countDestroyedPlayer == len(player.Boats) {
									fmt.Println("All your boats are destroyed, you lost the game, you cannot attack anymore")
								} else if countDestroyed == len(player.Boats) {
									fmt.Println("All the boats of " + player.Pseudo + " are destroyed, you won the game, you cannot attack anymore")
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
									isConnected = getRouteInfo(url, "hit", player.Pseudo, body, &player)
									if isConnected == false {
										check = "return"
									}
								}
							} else if check == "whereAreMyBoats?" {

								for _, coordinate := range isOccupied {
									println("Boat name : ", coordinate.BoatName, " : ", coordinate.X, coordinate.Y)
								}

							} else {
								if check != "return" && check != "quit" {
									fmt.Println("Invalid answer, please enter 'board', 'boats', 'hit' or 'quit'")
								}
							}
						}
					}
				}
			}
		}
	}
}

func getRouteInfo(url string, route string, pseudo string, body []byte, user *player.Player) bool {
	resp := &http.Response{}
	var err error
	if route == "hit" {
		resp, err = http.Post(url, "application/json", strings.NewReader(string(body)))
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		fmt.Println(pseudo + " is not connected anymore (I think he/she rage quit...)\n")
		return false
	}
	defer resp.Body.Close()
	var res string
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(res, "You hit a") {
		user.BoatPartDestroyed += 1
		if user.BoatPartDestroyed == 10 {
			user.Heal = true
		}
	}
	fmt.Println(res)
	return true
}

func checkIfPortIsFree(port int) string {
	url := "http://localhost:" + strconv.Itoa(port) + "/isFree"
	resp, err := http.Get(url)
	if err != nil {
		return "false"
	}
	defer resp.Body.Close()
	return "true"
}

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
								isDestroyed := true
								for _, part := range boat.BoatParts {
									if part != 2 {
										isDestroyed = false
									}
								}
								if isDestroyed {
									txt += ("You destroyed a " + c.BoatName + " !\n")
								}
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

func isFree(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

	}
}

func startServer(port int, player *player.Player, isOccupied *[]player.Coordinates) {
	http.HandleFunc("/IsFree", isFree)
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
