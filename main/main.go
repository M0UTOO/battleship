package main

import (
	"battleship/player"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
