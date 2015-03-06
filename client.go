package main

import (
    "net"
    "log"
    "fmt"
    "time"
    "bufio"
    "os"
    "strings"
    "./coding"
)

func listenForMessages(incoming_message chan coding.ServerPackage, connection_terminated chan bool, conn *net.TCPConn) {
    buffer := make([]byte, 2048)
    for {
        bytes_read, err := conn.Read(buffer)
        if err != nil {
            conn.Close()
            fmt.Println("Lost connection with server")
            connection_terminated <- true
            return
        }
        packages := coding.NetworkPacketToServerPackages(buffer[:bytes_read])
        for _, p := range(packages) {
            incoming_message <- p
        }
    }
}

func listenForUserInput(user_input chan string) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf(">>")
        line, _, err := reader.ReadLine()
        if err != nil {
            fmt.Println(err)
        }
        user_input <- string(line)
    }
}

func sendToServer(payload coding.ClientPackage, conn *net.TCPConn) {
    _, err := conn.Write(coding.ClientPackageToNetworkPacket(payload))
    if err != nil {
        log.Println("Could not send message")
        log.Fatal(err)
    }
}

func parseUserInput(input string) (coding.ClientPackage) {

    request := ""
    content := ""
    req_begin := strings.Index(input, "/")
    if req_begin >= 0 {
        req_end := strings.Index(input[req_begin:], " ")
        if req_end > req_begin {
            request = input[req_begin + 1 : req_end]
            content = input[req_end + 1 :]
        } else {
            request = input[req_begin + 1 :]
            content = ""
        }
    } else {
        request = "msg"
        content = input
    }

    return coding.ClientPackage{request, content}
}

func prettyPrint(when, username, content string) {
    fmt.Printf("At \x1b[30;1m%s \x1b[35m%s\x1b[0m said: %s\n>>", when, username, content)
}

// This will block until we actually connect!
func connectToServer(addr string) (*net.TCPConn){
    remote, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }

    for {
        conn, err := net.DialTCP("tcp", nil, remote)
        if err == nil {
            return conn
        }
        log.Println("Could not connect to server. Retrying...")
        time.Sleep(1 * time.Second)
    }
}

func main() {
    incoming_message := make(chan coding.ServerPackage)
    connection_terminated := make(chan bool)
    user_input := make(chan string)

    // TODO: Take server address as user input?
    server_addr := "127.0.0.1:12345"

    conn := connectToServer(server_addr)

    go listenForMessages(incoming_message, connection_terminated, conn)
    go listenForUserInput(user_input)

    //------Starting client------
    fmt.Println("Welcome to BabySeal chat client")
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
                    log.Fatal("Unknown server response")
            }

        case input := <- user_input:
            payload := parseUserInput(input)
            sendToServer(payload, conn)

        case <- connection_terminated:
            conn = connectToServer(server_addr)
        }
    }
}
