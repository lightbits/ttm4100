package main

import (
    "log"
    "fmt"
    "net"
    "strings"
    "./coding"
    "time"
)

// TODO: Should this be configurable on startup?
const SV_LISTEN_ADDRESS = "127.0.0.1:12345"

func getTime() string{
    const layout = "Jan 2 15:04"
    return time.Now().Format(layout)
}

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

type IncomingClientRequest struct {
    Socket  *net.TCPConn
    Payload coding.ClientPackage
}

func listenToClient(incoming_request chan IncomingClientRequest, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 1024)
        bytes_read, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(conn.RemoteAddr(), "lost connection")
            return
        }
        payload := coding.NetworkPacketToClientPackage(buffer[:bytes_read])
        request := IncomingClientRequest{conn, payload}
        incoming_request <- request
    }
}

func sendToClient(sender, response, content string, conn *net.TCPConn) {
    srv_struct := coding.ServerPackage{getTime(), sender, response, content}
    net_packet := coding.ServerPackageToNetworkPacket(srv_struct)
    _, err := conn.Write(net_packet)
    if err != nil {
        log.Fatal(err)
    }
}

func sendHistoryToClient(history []coding.ServerPackage, conn *net.TCPConn) {
    net_packet := coding.ServerPackagesToNetworkPacket(history)
    _, err := conn.Write(net_packet)
    if err != nil {
        log.Fatal(err)
    }
}

type Connection struct {
    Socket   *net.TCPConn
    Username string
}

func main() {
    connections         := make(map[string]Connection)
    incoming_connection := make(chan *net.TCPConn)
    incoming_requests   := make(chan IncomingClientRequest)
    message_history     := make([]coding.ServerPackage, 0)

    go listenForIncomingConnections(incoming_connection)
    fmt.Println("Ready for incoming connections. Bring it on!")
    for {
        select {
            case socket := <- incoming_connection:
                address := socket.RemoteAddr().String()
                connections[address] = Connection{socket, ""}
                go listenToClient(incoming_requests, socket)
                fmt.Println(address, "connected")

            case client_request := <- incoming_requests:
                socket  := client_request.Socket
                address := socket.RemoteAddr().String()
                request := client_request.Payload.Request
                content := client_request.Payload.Content
                switch (request) {
                    case "login":

                        connections[address] = Connection{socket, content}
                        sendToClient("server", "info", fmt.Sprintf("Your username is now %s", content), socket)

                        sendHistoryToClient(message_history, socket)

                    case "logout":
                        sendToClient("server", "info", "Goodbye!", socket)
                        connections[address].Socket.Close()
                        delete(connections, address)

                    case "msg":

                        username := connections[address].Username
                        if username == "" {
                            sendToClient("server", "error", "You must /login first.", socket)
                        } else {
                            message_history = append(message_history, coding.ServerPackage{getTime(), username, "message", content})

                            for dst_address, connection := range(connections) {
                                dst_socket := connection.Socket
                                if dst_address != address {
                                    sendToClient(username, "message", content, dst_socket)
                                }
                            }
                        }


                    case "names":

                        name_list := make([]string, len(connections))
                        i := 0
                        for _, connection := range(connections) {
                            username := connection.Username
                            if username == "" {
                                name_list[i] = "Noname"
                            } else {
                                name_list[i] = username
                            }
                            i++
                        }

                        names := strings.Join(name_list, ", ")

                        sendToClient("server", "info", names, socket)

                    case "help":

                        help := "The following commands are recognized by the server:\n/login <username>: Login with given username\n/logout: Disconnect from server\n/msg <message>: Send a message to everyone else\n/names: Get a list of people connected\n/help: See this list again."
                        sendToClient("server", "info", help, socket)

                    default:
                        sendToClient("server", "error", "Unknown command.", socket)
                }
                fmt.Println("Recived a packet from:", address)
                fmt.Println("request:", request)
                fmt.Println("content:", content)
                fmt.Println()
        }
    }

    for _, connection := range(connections) {
        connection.Socket.Close()
    }
    fmt.Println("All connections closed. I´ll take the day off")
}
