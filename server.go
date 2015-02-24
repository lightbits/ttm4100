package main

import (
    "log"
    "fmt"
    "net"
)

const SV_LISTEN_ADDRESS = "127.0.0.1:12345"

func listenForIncomingConnections(incoming_connection chan *net.TCPConn) {
    local, err := net.ResolveTCPAddr("tcp", SV_LISTEN_ADDRESS)
    if err != nil {
        log.Fatal(err)
    }

    listener, err := net.ListenTCP("tcp", local)
    if err != nil {
        log.Fatal(err)
    }

    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            log.Fatal(err)
        }
        incoming_connection <- conn
    }
}

type ClientPacket struct {
    Connection *net.TCPConn
    Content     string
}

func listenToClient(incoming_cl_packet chan ClientPacket, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 32)
        _, err := conn.Read(buffer)
        if err != nil {
            log.Fatal(err)
        }
        incoming_cl_packet <- ClientPacket{conn, string(buffer)}
    }
}

func sendToClient(content string, conn *net.TCPConn) {
    // TODO: Null-delimiting and/or fixed message length
    // Currently just sending the content and hoping that
    // the recipient gets the data in a single packet
    _, err := conn.Write([]byte(content))
    if err != nil {
        log.Fatal(err)
    }
}

func prettyPrintClientMessage(username string, content string) {
    // This prints to the console with colored usernames!!
    // (if your console supports colors, that is)
    fmt.Printf("\x1b[35m%s\x1b[0m said: %s\n", username, content)
}

func main() {
    connections           := make(map[*net.TCPConn]string)
    incoming_connection   := make(chan *net.TCPConn)
    incoming_cl_packet    := make(chan ClientPacket)

    prettyPrintClientMessage("John doe", "Hey guys!")

    go listenForIncomingConnections(incoming_connection)

    for {
        select {
            case conn := <- incoming_connection:
                connections[conn] = "NoName"
                go listenToClient(incoming_cl_packet, conn)
                fmt.Println(conn.RemoteAddr(), "connected")

            case ClientPacket := <- incoming_cl_packet:
                who := ClientPacket.Connection.RemoteAddr()
                prettyPrintClientMessage(who.String(), ClientPacket.Content)
        }
    }

    for conn := range(connections) {
        conn.Close()
    }
}
