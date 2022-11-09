package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

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
	log.Println("4 personnes sont connectées")
	for i:=0;i<4;i++{
		fmt.Fprintf(connexion[i],"4 joueurs sont connectés"+"\n")
	}
	/*
	var i = 0
	for {
		message, _ := bufio.NewReader(connexion[i%4]).ReadString('\n')
		fmt.Print("Message Reçu", message)
		fmt.Fprintf(connexion[i%4],message+"\n")
		i++
	}
	*/
	
}

