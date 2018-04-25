package main

import (
	"log"
	"net"
	"bufio"
	"os"
	"encoding/json"
	"strings"
	"os/signal"
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

	println("deferring close connection")
	//defer quit(conn)
	//defer conn.Close()

	c := make(chan os.Signal, 1)		//This code taken from StackOverflow
	signal.Notify(c, os.Interrupt)		//Used for clean disconnect w/o breaking server on ^C interrupt
	go func(){
		for sig := range c {
			// sig is a ^C, handle it
			log.Printf("%v", sig)
			quit(conn)
			//conn.Close()
			os.Exit(1)
		}
	}()

	go sendMessages(conn)

	for {
		serverMsg, _ := reader.ReadBytes('\n') //Reads until a newline escape appears
		print(string(serverMsg))               //No need for println because \n forces new line
	}
	println("end of main")
}

func quit(connection net.Conn){
	var msg = Message{username, "quit", " "}
	encodedMsg, _ := json.Marshal(msg)
	println("Sending quit message")
	connection.Write(encodedMsg)
	connection.Close()
	os.Exit(1)
}

func sendMessages(connection net.Conn){
	var msg Message
	inputReader = *bufio.NewReader(os.Stdin)
	for {
		msgBody, _ := inputReader.ReadString('\n')
		msgCommand := "say"			//Unless otherwise stated the client is trying to send a message

		//Isolate any /command in the message, to be compared to a list in the server
		if len(msgBody) != 0 && msgBody[0] == byte('/'){	//If the first character is '/' then it is a command
			msgCommand = strings.SplitAfter(msgBody, " ")[0]	//Isolate command (with / and end space)
			msgBody = strings.TrimLeft(msgBody, msgCommand)			//Remove command from message contents
			msgCommand = strings.ToLower(msgCommand)				//Make lowercase for use in enum
			msgCommand = strings.TrimLeft(msgCommand, "/")	//Eliminate /
			msgCommand = strings.TrimRight(msgCommand, " ")	//Eliminate trailing space
			msgCommand = strings.TrimRight(msgCommand, "\n")	//Eliminate endline
		}

		if msgCommand == "quit"{
			quit(connection)
		}
		msg = Message{username, msgCommand, msgBody}
		encodedMsg, _ := json.Marshal(msg)

		connection.Write(encodedMsg)
	}
}