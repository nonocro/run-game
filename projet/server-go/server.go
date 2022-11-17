package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	//"sync"
	"strconv"
)

func main()  {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	var connexion []net.Conn
	var readers  []*bufio.Reader
	var compt int = 0

	for compt<4{
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return 
		}
		defer conn.Close()
		connexion = append(connexion,conn)
		readers = append(readers,bufio.NewReader(conn))
		log.Println("Un client s'est connecté")
		fmt.Fprintf(conn,"tu est le joueur "+strconv.Itoa(compt)+"\n")
		compt++
	}
	compt = 0
	
	log.Println("4 personnes sont connectées")

	for i:=0;i<4;i++{
		fmt.Fprintf(connexion[i],"4 joueurs sont connectés"+"\n")

	}
	
	for {
		for compt<4{
			for i:=0;i<4;i++{
				message, _ := readers[i].ReadString('\n')
				log.Println(message)
				compt+=1
			}
			log.Println("tous les joueur sont pret !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connexion[i],"tous les joueurs sont pret"+"\n")		
			}
		}
		

		
	}
}

