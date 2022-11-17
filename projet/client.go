package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
)

func (g *Game)client() {
	var err  error
	g.conn, err = net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	//defer conn.Close()
	log.Println("Je suis connecté")
	/*
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Reponse du serveur : (reçu)" + message + "\n")

*/
	reader := bufio.NewReader(g.conn)
	for {
		message, _ := reader.ReadString('\n')
		fmt.Print("Reponse du serveur : (reçu)" + message + "\n")
		if strings.Contains(message,"4 joueurs sont connectés") {
			g.done=true
		}

		if strings.Contains(message,"tu est le joueur ") {
			g.myRunner,_=strconv.Atoi(message[len(message)-2:len(message)-1])
		}

		if strings.Contains(message,"tous les joueurs sont pret") {
			g.done=true
		}
		/*reader1 := bufio.NewReader(os.Stdin)
		fmt.Print("message à envoyé :")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(g.conn, text+"\n")*/
	}

}

