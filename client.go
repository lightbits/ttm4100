package main

import (
    "net"
    "log"
    "fmt"
    "time"
    "bufio"
    "os"
    "strings"
    "runtime"
    "./coding"
)
const DEBUG = false

func listenForMessages(incoming_message chan coding.ServerPackage, connection_terminated chan bool, conn *net.TCPConn) {
    buffer := make([]byte, 2048)
    for {
        bytes_read, err := conn.Read(buffer)
        if err != nil {
             if DEBUG {
                log.Println(err)
            }
            conn.Close()
            fmt.Println("Lost connection with server")
            connection_terminated <- true
            break
        }else{
            packages := coding.NetworkPacketToServerPackages(buffer[:bytes_read])
            for _, p := range(packages) {
                incoming_message <- p
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func listenForUserInput(user_input chan string, userInputTrigger chan int) {
    reader := bufio.NewReader(os.Stdin)
    for {
        select{
            case <- userInputTrigger:
                fmt.Printf(">> ")
                line, _, err := reader.ReadLine()
                if err != nil {
                    log.Fatal(err)
                }
                user_input <- strings.TrimSpace(string(line))
        }
    }
}

func sendToServer(payload coding.ClientPackage, conn *net.TCPConn) {
    _, err := conn.Write(coding.ClientPackageToNetworkPacket(payload))
    if err != nil {
        log.Println("Could not send message")
        log.Fatal(err)
    }
}

func parseUserInput(input string) (coding.ClientPackage, bool) {
    request := ""
    content := ""
    waitingOnServer := true
    splitIndex := strings.Index(input," ")
    if splitIndex == -1 { //one word
        switch strings.ToLower(input){
        case "msg":
            waitingOnServer = false
            fallthrough
        case "login","logout","names","help":
            request = strings.ToLower(input)
        default:
            request = "msg"
            content = input
            waitingOnServer = false
        }
    }else { //more than one word
        switch strings.ToLower(input[:splitIndex]){
            case "msg":
                waitingOnServer = false
                fallthrough
            case "login","logout","names","help":
                request = strings.ToLower(input[:splitIndex])
                if splitIndex != len(input) { //actually a part two
                    content = input[splitIndex+1:]
                }
            default:
                request = "msg"
                content = input
                waitingOnServer = false
        }
    }
    if DEBUG {
        log.Println("Command: ", request)
        log.Println("Payload: ", content)
    }
    return coding.ClientPackage{request, content}, waitingOnServer
}

func prettyPrint(when, username, content string) {
    fmt.Printf("At \x1b[30;1m%s \x1b[35m%s\x1b[0m said: %s\n", when, username, content)
}

func connectToServer(addr string) (*net.TCPConn, error){
    remote, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        if DEBUG {
            log.Println("Not a valid server adress")
        }
        return nil, err
    }else{
        conn, err := net.DialTCP("tcp", nil, remote)
        if err != nil {
            if DEBUG {
                log.Println("Could not connect to server")
                log.Println(err)
            }
            return nil, err
        }else{
            if DEBUG {
                log.Println("Connected to", conn.RemoteAddr())
            }
            return conn, nil
        }
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    incoming_message := make(chan coding.ServerPackage)
    connection_terminated := make(chan bool)
    user_input := make(chan string, 5)
    userInputTrigger := make(chan int, 5)

    // TODO: Take server address as user input?
    serverAdr := "127.0.0.1:12345"

    var conn *net.TCPConn
    var err error
    var countAttempts int = 0
    
    for{
        conn, err = connectToServer(serverAdr)
        if err != nil {
            countAttempts++
            if countAttempts == 1 {
                fmt.Printf("Connecting to server")
            }
            fmt.Printf(".")
            if countAttempts > 10 {
                fmt.Printf("\n")
                fmt.Println("Closing down BabySeal client")
                os.Exit(0) 
            }
            time.Sleep(1 * time.Second)
        }else{
            defer conn.Close()
            if countAttempts > 0 {
                fmt.Printf("\n")
            }
            break
        }
    }

    go listenForMessages(incoming_message,connection_terminated, conn)    
    go listenForUserInput(user_input, userInputTrigger)

    //------Starting client------  
    fmt.Println("Welcome to BabySeal chat client")
    var waitingOnServer bool = false
    userInputTrigger <- 1
    for {
        select {
        case server_response := <- incoming_message:
            response  := server_response.Response
            content   := server_response.Content
            sender    := server_response.Sender
            timestamp := server_response.Timestamp
            switch (response) {
                case "history":
                    prettyPrint(timestamp, sender, content)
                case "message":
                    prettyPrint(timestamp, sender, content)
                case "info":
                    prettyPrint(timestamp, sender, content)
                    if content == "Goodbye!" {
                        os.Exit(0)
                    }
                case "error":
                    prettyPrint(timestamp, sender, content)
                default:
                    fmt.Println("Unknown server response")
            }
            waitingOnServer = false
            userInputTrigger <- 1

        case input := <- user_input:
            payload, needToWait := parseUserInput(input)
            sendToServer(payload, conn)
            waitingOnServer = needToWait

        case <- time.After(2 * time.Second):
            if(waitingOnServer){
                waitingOnServer = false
                fmt.Println("waiting on server timed out")
                userInputTrigger <- 1
            }
        case <- connection_terminated:
            fmt.Printf("Trying to reconnect to server")
            var numOfAttempst int
            for{
                conn, err = connectToServer(serverAdr)
                if err != nil {
                    if DEBUG {
                        log.Println(err)
                    }else{
                        fmt.Printf(".")
                    }
                    numOfAttempst++
                    if(numOfAttempst > 60){
                        fmt.Printf("\n")
                        fmt.Println("Closing down BabySeal client")
                        os.Exit(0)
                    }
                }else{
                    fmt.Printf("\n")
                    fmt.Println("Reconnected to server :)")
                    go listenForMessages(incoming_message,connection_terminated, conn)
                    userInputTrigger <- 1
                    break
                }
                time.Sleep(1000 * time.Millisecond)
            }
        }
    }
}