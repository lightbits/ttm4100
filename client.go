package main

import (
    "net"
    "log"
    "fmt"
    "bufio"
    "os"
)

func listenForMessages(incoming_message chan string, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 1024)
        _, err := conn.Read(buffer)
        if err != nil {
            log.Fatal(err)
        }
        incoming_message <- string(buffer)
    }
}

func listenForUserInput(user_input chan string) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("Enter text: ")
        line, _, err := reader.ReadLine()
        if err != nil {
            fmt.Println(err)
        }
        user_input <- string(line)
    }
}

func main() {
    // TODO: Take server address as user input?

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

    incoming_message := make(chan string)
    go listenForMessages(incoming_message, conn)

    user_input := make(chan string)
    go listenForUserInput(user_input)

    for {
        select {
            case msg := <- incoming_message:
                log.Println("Received", len(msg))

            case input := <- user_input:
                _, err := conn.Write([]byte(input))
                if err != nil {
                    log.Fatal(err)
                }
        }
    }
}
