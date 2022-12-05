package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
)

func (g *Game)client() {
	var err  error
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

		if strings.Contains(message,"4 players are connected") {
			g.done=true
		}else if strings.Contains(message,"you are the player") {
			g.myRunner,_=strconv.Atoi(message[len(message)-2:len(message)-1])
		}else if strings.Contains(message,"All the players are ready") {
			g.done=true
		}else if strings.Contains(message,":r") {
			message = message[2:len(message)-1]
			times := strings.Split(message,",")
			for nb,_ := range g.runners {
				timePlayer,_ :=strconv.Atoi(times[nb])
				g.runners[nb].runTime=time.Duration(timePlayer)
				//we resolved the case of lost message
				g.runners[nb].arrived = true
			}
			g.done=true
		}else if strings.Contains(message,":c") {
			g.nbPlayer,_ = strconv.Atoi(message[2:len(message)-1])
		}else if strings.Contains(message,":nbplayer") {
			g.nbPlayer++
		}else if strings.Contains(message,":space") {
			index,_ := strconv.Atoi(message[6:len(message)-1])
			g.counter_space[index] = true
		}else if strings.Contains(message,":key") {
			mouv := strings.Split(message,",")
			player,_:= strconv.Atoi(mouv[1])
			if player != g.myRunner{
				if mouv[2] == "2" {
					g.runners[player].ServerChoose(false,false,true)
				}else if mouv[2] == "0"{
					g.runners[player].ServerChoose(true,false,false)
				}else if mouv[2] == "1" {
					g.runners[player].ServerChoose(false,true,false)
	
				}
			}
		}

	
	} 
}

