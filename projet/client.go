package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"sync"
)
var m sync.Mutex //use to avoid race-condition on game's attributes

//manage all communication from the server and act consequently
func (g *Game) client() {
	var err error
	//connect itself to the server with IP give as a parameter
	g.conn, err = net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer g.conn.Close()
	log.Println("I'm connected")

	reader := bufio.NewReader(g.conn)

	//loop of listening to the server
	for {
		message, _ := reader.ReadString('\n')
		fmt.Print("Server answer : (received)" + message + "\n")
		// test the content of the message from server
		if strings.Contains(message, "4 players are connected") { // received the validation when all the player are connected
			m.Lock()
			g.done = true
			m.Unlock()
		} else if strings.Contains(message, "you are the player") { // tell you your player number
			m.Lock()
			g.myRunner, _ = strconv.Atoi(message[len(message)-2 : len(message)-1])
			m.Unlock()
		} else if strings.Contains(message, "All the players are ready") { // received the validation when all the player choose their color 
			m.Lock()
			g.done = true
			m.Unlock()
		} else if strings.Contains(message, ":r") { // reveived all the result and mange to affect them to each runners
			message = message[2 : len(message)-1]
			times := strings.Split(message, ",")
			m.Lock()
			for nb, _ := range g.runners {
				timePlayer, _ := strconv.Atoi(times[nb])
				g.runners[nb].runTime = time.Duration(timePlayer)
				//we resolved the case of lost message
				g.runners[nb].arrived = true
			}
			g.done = true
			m.Unlock()
		} else if strings.Contains(message, ":c") { // received the number of player connected
			m.Lock()
			g.nbPlayer, _ = strconv.Atoi(message[2 : len(message)-1])
			m.Unlock()
		} else if strings.Contains(message, ":nbplayer") { // received message to update the counter of players
			m.Lock()
			g.nbPlayer++
			m.Unlock()
		} else if strings.Contains(message, ":space") { // received the input (ex :space1)) to adapt speed of each runner
			index, _ := strconv.Atoi(message[6 : len(message)-1])
			m.Lock()
			g.counter_space[index] = true
			m.Unlock()
		}else if strings.Contains(message,":key") { // received the input to manage the synchronisation of the color selection
			mouv := strings.Split(message,",")
			player,_:= strconv.Atoi(mouv[1])
			m.Lock()
			left, _ := strconv.ParseBool(mouv[2])
			right, _ := strconv.ParseBool(mouv[3])
			space, _ := strconv.ParseBool(mouv[4])
			g.keys_bool[player] = [3]bool{left, right, space}
			m.Unlock()
		}
	}
}
