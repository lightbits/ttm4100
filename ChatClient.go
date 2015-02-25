package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
	"strings"
	"encoding"
)

const addres string = "127.0.0.1:6000"

const (
	
	COMMAND_LOGIN = iota
	COMMAND_LOGOUT
	COMMAND_MSG
	COMMAND_NAMES
	COMMAND_HELP
	COMMAND_INVALID
)


func main() {

    conn, err := net.Dial("tcp", addres)
    if err != nil {
		fmt.Println("Did not connect to server")        
		fmt.Println(err)
        os.Exit(1)
    }
	fmt.Println("Connected to server")

	reader := bufio.NewReader(os.Stdin)


	// enter username	
	fmt.Print("Enter a username: ")
	username, _ := reader.ReadString('\n')
	    

		
	// choose message or command
	for {
		command := getcommand(reader.ReadString('\n')[0])	
		pack ClientPackage
		switch (command) {
			case COMMAND_LOGIN:
				
			case COMMAND_LOGOUT:
			case COMMAND_MSG:
			case COMMAND_NAMES:
			case COMMAND_HELP:
				
			case COMMAND_INVALID:
				fmt.Println("Please enter a valid command")
		}		

	}
}

func getcommand(input string) (command int){
	input = strings.ToLower(input)
	switch (input) {
		case "login": 	
			return COMMAND_LOGIN
		case "logout": 	
			return COMMAND_LOGOUT
		case "msg":		
			return COMMAND_MSG
		case "names":	
			return COMMAND_NAMES
		case "help":	
			return COMMAND_HELP
	}		
	return COMMAND_INVALID
}


func doLogin() {

}

func doLogout() {

}

func doSendMessage() {

}

func doNames() {

}

func doHelp() {

}
