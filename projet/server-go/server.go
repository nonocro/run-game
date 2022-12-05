package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	//"sync"
	"strconv"
	"sync"
	"strings"
)

const (
	StateWelcomeScreen int = iota // Title screen (iota = 0 et les autre constante sont incrémenté automatiquement de 1)
	StateChooseRunner             // Player selection screen
	StateRun                      // Run
	StateResult                   // Results announcement
)

var w sync.WaitGroup

func main()  {

	var state int = StateWelcomeScreen
	var connection []net.Conn // array of connection
	var readers  []*bufio.Reader // array of reader for each connection
	var count int = 0

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	for count<4{
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
		messageToAll(connection,":c"+strconv.Itoa(count+1))
		count++
	}
	
	log.Println("4 players are connected")
	messageToAll(connection,"4 players are connected")
	state++

	for {
		for state == StateChooseRunner{
			w.Add(4)
			for i:=0;i<4;i++{
				go choice_message(readers[i],connection,i)
			}
			
			w.Wait()
			log.Println("All the players are ready !!!!")
			messageToAll(connection,"All the players are ready")
			state++
		}
		

		for state == StateRun{
			var result []string = make([]string,4)
			w.Add(4)
			for i:=0;i<4;i++{
				go result_message(readers[i],result,connection)
			}
			
			w.Wait()
			log.Println(result)
			log.Println("All the players arrived !!!!")
			messageToAll(connection,":r"+strings.Join(result,","))
			state++
		}

		for state==StateResult{
			w.Add(4)
			for i:=0;i<4;i++{
				go reset_message(readers[i],connection)
			}
			w.Wait()
			log.Println("All the players want to restart !!!!")
			state = StateRun
		}
		
	}
}


func messageToAll(connection []net.Conn,msg string){
	for i:=0;i<len(connection);i++{
		fmt.Fprintf(connection[i],msg+"\n")		
	}
}

func choice_message(reader *bufio.Reader,connection []net.Conn,nbPlayer int){
	message, _ := reader.ReadString('\n')
	for {
		bools := strings.Split(message,",")
		if bools[4]=="true" {
			messageToAll(connection,":key"+","+bools[1]+","+"2"+",")
			w.Done()
			break
		}else if bools[2]=="true"{
			messageToAll(connection,":key"+","+bools[1]+","+"0"+",")
		}else if bools[3]=="true"{
			messageToAll(connection,":key"+","+bools[1]+","+"1"+",")
		}
		message, _ = reader.ReadString('\n')
		log.Println(message)
	}
	
}

func reset_message(reader *bufio.Reader, connection []net.Conn){
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message,"restart"){
		message,_ =reader.ReadString('\n')
	}
	messageToAll(connection,":nbplayer++")
	log.Println(message)
	w.Done()
}

func result_message(reader *bufio.Reader, result []string, connection []net.Conn){
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message,":r"){
		if strings.Contains(message,":space"){
			var numRunner int
			for _,conn := range connection{
				fmt.Println(message)
				numRunner,_ = strconv.Atoi(message[6:len(message)-1])
				fmt.Println(numRunner)
				fmt.Fprintf(conn,":space"+strconv.Itoa(numRunner)+"\n")
			}
		}
		message,_ =reader.ReadString('\n')	
	}
	
	var indice int
	indice,_ = strconv.Atoi(message[:1])
	result[indice] = message[3:len(message)-1]
	log.Println(message)
	w.Done()
}



