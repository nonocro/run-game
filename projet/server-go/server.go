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
			w.Add(4)
			for i:=0;i<4;i++{
				go message_choix(readers[i])
			}
			compt =4
			w.Wait()
			log.Println("tous les joueur sont pret !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connexion[i],"tous les joueurs sont pret"+"\n")		
			}
		}
		
		for compt<8{
			var result []string = make([]string,4)
			w2.Add(4)
			for i:=0;i<4;i++{
				go message_resultat(readers[i],result)
			}
			compt =8
			w2.Wait()
			log.Println(result)
			log.Println("tous les joueur sont arrivés !!!!")
			for i:=0;i<4;i++{
				fmt.Fprintf(connexion[i],":r"+strings.Join(result,",")+"\n")		
			}
		}
	}
}

func message_choix(reader *bufio.Reader){
	message, _ := reader.ReadString('\n')
	log.Println(message)
	w.Done()
}


func message_resultat(reader *bufio.Reader, result []string){
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
