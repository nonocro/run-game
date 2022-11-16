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
	conn, err := net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	log.Println("Je suis connecté")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	for {
		fmt.Print("Reponse du serveur : (reçu) je suis le joueur" + message + "\n")
		g.MyRunner,_ = strconv.Atoi(message)
		fmt.Print(g.MyRunner)
	}
	


	for {
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("message à envoyé :")
		//text, _ := reader.ReadString('\n')
		//fmt.Fprintf(conn, text+"\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Reponse du serveur : (reçu)" + message + "\n")
		if strings.Contains(message,"4 joueurs sont connectés") {
			g.done=true
		}
	}

}

