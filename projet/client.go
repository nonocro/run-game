package main
/*
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	//"strconv"
)

func main() {
	conn, err := net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("Je suis connecté")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("message à envoyé :")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Reponse du serveur : (reçu)" + message + "\n")
	}

}

*/