package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	//"sync"
	"strconv"
	"strings"
	"sync"
)

const (
	StateWelcomeScreen int = iota // Title screen (iota = 0 et les autre constante sont incrémenté automatiquement de 1)
	StateChooseRunner             // Player selection screen
	StateRun                      // Run
	StateResult                   // Results announcement
)

var w sync.WaitGroup // use to synchronise all go-routines 

func main() {

	var state int = StateWelcomeScreen
	var connection []net.Conn   // array of connection
	var readers []*bufio.Reader // array of reader for each connection
	var count int = 0

	//connection reception
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	//loop for player's connection
	for count < 4 {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		defer conn.Close()
		connection = append(connection,conn)
		readers = append(readers,bufio.NewReader(conn))
		log.Println("player "+strconv.Itoa(count)+" is connected ")
		fmt.Fprintf(conn,"you are the player "+strconv.Itoa(count)+"\n")
		count++
		messageToAll(connection,":c"+strconv.Itoa(count))
	}

	log.Println("4 players are connected")
	messageToAll(connection, "4 players are connected")
	state++

	//loop of listening to all clients
	//each loops will call a go-routine function to listen to each client in particular
	for {
		//loop for the choosing state until all player have choosing a character
		for state == StateChooseRunner {
			w.Add(4)
			for i := 0; i < 4; i++ {
				go choice_message(readers[i], connection, i)
			}

			w.Wait()
			log.Println("All the players are ready !!!!")
			messageToAll(connection, "All the players are ready")
			state++
		}
		//loop for the running state until all player have arrived
		for state == StateRun {
			var result []string = make([]string, 4)
			w.Add(4)
			for i := 0; i < 4; i++ {
				go result_message(readers[i], result, connection)
			}

			w.Wait()
			log.Println(result)
			log.Println("All the players arrived !!!!")
			messageToAll(connection, ":r"+strings.Join(result, ","))
			state++
		}
		//loop for the result state until all player want to restart
		for state == StateResult {
			w.Add(4)
			for i := 0; i < 4; i++ {
				go reset_message(readers[i], connection)
			}
			w.Wait()
			log.Println("All the players want to restart !!!!")
			state = StateRun
		}

	}
}

//send a message to all clients
func messageToAll(connection []net.Conn, msg string) {
	for i := 0; i < len(connection); i++ {
		fmt.Fprintf(connection[i], msg+"\n")
	}
}

//manages the choice of characters of one client, it receive and send back to all the mouvement of each clients
//end when receive ":skins" message from the clients
func choice_message(reader *bufio.Reader, connection []net.Conn, nbPlayer int) {
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message, ":skins"){
		if strings.Contains(message,"true"){
			bools := strings.Split(message,",")
			messageToAll(connection,":key"+","+bools[1]+","+bools[2]+","+bools[3]+","+bools[4]+",")
		}
		message, _ = reader.ReadString('\n')
	}
	w.Done()
}

//manage the restart state, receive the reset message and increment the counter of each player
func reset_message(reader *bufio.Reader, connection []net.Conn) {
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message, "restart") {
		message, _ = reader.ReadString('\n')
	}
	messageToAll(connection, ":nbplayer++")
	log.Println(message)
	w.Done()
}

//Manage the run and the results, receive the space command and send it back to all clients
func result_message(reader *bufio.Reader, result []string, connection []net.Conn) {
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message, ":r") {
		if strings.Contains(message, ":space") {
			var numRunner int
			for _, conn := range connection {
				numRunner, _ = strconv.Atoi(message[6 : len(message)-1])
				fmt.Fprintf(conn, ":space"+strconv.Itoa(numRunner)+"\n")
			}
		}
		message, _ = reader.ReadString('\n')
	}
	var index int
	index, _ = strconv.Atoi(message[:1])
	result[index] = message[3 : len(message)-1]
	w.Done()
}
