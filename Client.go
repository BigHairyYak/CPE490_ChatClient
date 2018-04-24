package main

import (
	"log"
	"net"
	"bufio"
	"os"
	"encoding/json"
	"strings"
)

type Message struct {
	Username	string`json:"username"`
	Command		string	`json:"command"`
	Message		string	`json:"message"`
}

var serverIP, serverPort, username string
var msgCommand string			//commands to be sent in message
var inputReader bufio.Reader
var cmdReader strings.Reader

func check(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func main(){
	if len(os.Args) != 4{	//Program name, IP, port, username
		log.Fatal("Improper number of arguments!")
	}

	serverIP = os.Args[1]
	serverPort = os.Args[2]
	username = os.Args[3]

	println("Connecting to " + serverIP + ":" + serverPort)
	conn, err := net.Dial("tcp", serverIP + ":" + serverPort)
	check(err)

	joinMSg, _ := json.Marshal(Message{username, "join", ""})
	conn.Write(joinMSg)		//Send server join command
/*	print("local: ")
	println(conn.LocalAddr())
	print("remote: ")
	println(conn.RemoteAddr())*/
	reader := bufio.NewReader(conn)

	go sendMessages(conn)

	for {
		serverMsg, _ := reader.ReadBytes('\n') //Reads until a newline escape appears
		print(string(serverMsg))				  	 //No need for println because \n forces new line
	}
	defer conn.Close()
}

func sendMessages(connection net.Conn){
	var msg Message
	inputReader = *bufio.NewReader(os.Stdin)
	for {
		msgBody, _ := inputReader.ReadString('\r')
		msgCommand := "say"			//Unless otherwise stated the client is trying to send a message

		//Isolate any /command in the message, to be compared to a list in the server
		if msgBody[0] == byte('/'){	//If the first character is '/' then it is a command
			msgCommand = strings.SplitAfter(msgBody, " ")[0]	//Isolate command (with / and end space)
			msgBody = strings.TrimLeft(msgBody, msgCommand)			//Remove command from message contents
			msgCommand = strings.ToLower(msgCommand)				//Make lowercase for use in enum
			msgCommand = strings.TrimLeft(msgCommand, "/")	//Eliminate /
			msgCommand = strings.TrimRight(msgCommand, " ")	//Eliminate trailing space
			msgCommand = strings.TrimRight(msgCommand, "\n")	//Eliminate endline
		}

		msg = Message{username, msgCommand, msgBody}
		encodedMsg, _ := json.Marshal(msg)

		connection.Write(encodedMsg)
	}
}