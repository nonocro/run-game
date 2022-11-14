package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	//"sync"
)

func main()  {


	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	var connexion []net.Conn
	compt := 0

	for compt<4{
		conn, err := listener.Accept()
		connexion = append(connexion,conn)
		if err != nil {
			log.Println("accept error:", err)
			return 
		}
		defer conn.Close()
		log.Println("Un client s'est connecté")
		compt++
	}
	log.Println("ok")
	
	
	log.Println("4 personnes sont connectées")
	for i:=0;i<4;i++{
		fmt.Fprintf(connexion[i],"4 joueurs sont connectés"+"\n")
	}

	for {
		for i:=0;i<4;i++{
			message, _ := bufio.NewReader(connexion[i]).ReadString('\n')
			if len(message)==1{
				break
			}
			fmt.Print("Message Reçu", message)
			fmt.Fprintf(connexion[i], message+"\n")
		}
	}
}

