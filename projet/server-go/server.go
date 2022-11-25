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

var w sync.WaitGroup
var w2 sync.WaitGroup
var w3 sync.WaitGroup

func main()  {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	var connection []net.Conn // array of connection
	var readers  []*bufio.Reader // array of reader for each connection

	var count int = 0

	for count<4{
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return 
		}
		defer conn.Close()
		connection = append(connection,conn)
		readers = append(readers,bufio.NewReader(conn))
		log.Println("A player is connected ")
		fmt.Fprintf(conn,"you are the player "+strconv.Itoa(count)+"\n")
		for i:=0;i<=count;i++{
			fmt.Fprintf(connection[i],":c"+strconv.Itoa(count+1)+"\n")
		}
		count++
	}
	count = 0
	
	log.Println("4 players are connected")

	for i:=0;i<4;i++{
		fmt.Fprintf(connection[i],"4 players are connected"+"\n")

	}
	
	for {
		for count<4{
			w.Add(4)
			for i:=0;i<4;i++{
				go choice_message(readers[i])
			}
			count =4
			w.Wait()
			log.Println("All the players are ready !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connection[i],"All the players are ready"+"\n")		
			}
		}
		

		for count<8{
			var result []string = make([]string,4)
			w2.Add(4)
			for i:=0;i<4;i++{
				go result_message(readers[i],result)
			}
			count =8
			w2.Wait()
			log.Println(result)
			log.Println("All the players have arrived !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connection[i],":r"+strings.Join(result,",")+"\n")		
			}
		}

		for count<12{
			w3.Add(4)
			for i:=0;i<4;i++{
				go reset_message(readers[i],connection)
			}
			count =12
			w3.Wait()
			log.Println("All the players want to restart !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connection[i],":again"+"\n")		
			}
		}
		count = 4
	}
}

func choice_message(reader *bufio.Reader){
	message, _ := reader.ReadString('\n')
	log.Println(message)
	w.Done()
}

func reset_message(reader *bufio.Reader, connection []net.Conn){
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message,"restart"){
		message,_ =reader.ReadString('\n')	
	}
	for _,conn := range connection{
		fmt.Fprintf(conn,":nbplayer++"+"\n")
	}
	log.Println(message)
	w3.Done()
}

func result_message(reader *bufio.Reader, result []string){
	message, _ := reader.ReadString('\n')
	for !strings.Contains(message,":r"){
		message,_ =reader.ReadString('\n')	
	}
	
	var indice int
	indice,_ = strconv.Atoi(message[:1])
	result[indice] = message[3:len(message)-1]
	log.Println(message)
	w2.Done()
}



