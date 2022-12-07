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
)

func (g *Game) client() {
	var err error
	g.conn, err = net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	log.Println("I'm connected")

	reader := bufio.NewReader(g.conn)
	for {

		message, _ := reader.ReadString('\n')
		fmt.Print("Server answer : (received)" + message + "\n")

		if strings.Contains(message, "4 players are connected") {
			g.done = true
		} else if strings.Contains(message, "you are the player") {
			g.myRunner, _ = strconv.Atoi(message[len(message)-2 : len(message)-1])
		} else if strings.Contains(message, "All the players are ready") {
			g.done = true
		} else if strings.Contains(message, ":r") {
			message = message[2 : len(message)-1]
			times := strings.Split(message, ",")
			for nb, _ := range g.runners {
				timePlayer, _ := strconv.Atoi(times[nb])
				g.runners[nb].runTime = time.Duration(timePlayer)
				//we resolved the case of lost message
				g.runners[nb].arrived = true
			}
			g.done = true
		} else if strings.Contains(message, ":c") {
			g.nbPlayer, _ = strconv.Atoi(message[2 : len(message)-1])
		} else if strings.Contains(message, ":nbplayer") {
			g.nbPlayer++
		} else if strings.Contains(message, ":keys") {
			bools := strings.Split(message, ",")
			fmt.Println(len(bools))
			fmt.Println("start")
			for _, values := range bools {
				fmt.Println(values)
			}
			fmt.Println("end")
			fmt.Println(bools)
			index, _ := strconv.Atoi(bools[1])
			left, _ := strconv.ParseBool(bools[2])
			right, _ := strconv.ParseBool(bools[3])
			space, _ := strconv.ParseBool(bools[4])
			fmt.Println("space", space, bools[4])
			g.keys_bool[index] = [3]bool{left, right, space}
		} else if strings.Contains(message, ":space") {
			index, _ := strconv.Atoi(message[6 : len(message)-1])
			g.counter_space[index] = true
		}

	}
}
