package main

import (
    "net"
    "log"
    "fmt"
    "bufio"
    "os"
    "strings"
    "./coding"
)

func listenForMessages(incoming_message chan coding.ServerPackage, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 1024)
        bytes_read, err := conn.Read(buffer)
        if err != nil {
            log.Fatal(err)
        }
        incoming_message <- coding.NetworkPacketToServerPackage(buffer[:bytes_read])
    }
}

func listenForUserInput(user_input chan string) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf(">> ")
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
        log.Fatal(err)
    }
}

func parseUserInput(input string) coding.ClientPackage {
    request := ""
    content := ""

    /*
    User input is of the form
        /request content
    If the /request field is not given,
    we interpret the input as a <msg> payload
    with content = input.
    */
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
    fmt.Printf("At \x1b[30;1m%s \x1b[35m%s\x1b[0m said: %s\n", when, username, content)
}

func main() {
    // TODO: Take server address as user input?
    // For now, use localhost
    remote, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
    if err != nil {
        log.Fatal(err)
    }

    conn, err := net.DialTCP("tcp", nil, remote)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    log.Println("Connected to", conn.RemoteAddr())

    incoming_message := make(chan coding.ServerPackage)
    go listenForMessages(incoming_message, conn)

    user_input := make(chan string)
    go listenForUserInput(user_input)

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
                case "error":
                    prettyPrint(timestamp, sender, content)
                default:
                    fmt.Println("Unknown server response")
                }
                fmt.Printf(">> ")

            case input := <- user_input:
                payload := parseUserInput(input)
                sendToServer(payload, conn)
        }
    }
}
