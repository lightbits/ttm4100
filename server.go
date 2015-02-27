package main

import (
    "log"
    "fmt"
    "net"
    "encoding/json"
    "time"
)

// TODO: Should this be configurable on startup?
const SV_LISTEN_ADDRESS = "127.0.0.1:12345"

//----ENCODING START
type InternalClientPackage struct {
    Connection *net.TCPConn
    Payload ClientPackage
}

type ClientPackage struct {
    Request string
    Content string
}

type ServerPackage struct {
    Timestamp string
    Sender string
    Response string
    Content string
}

func clientPackToNetworkPackage(pack ClientPackage) []byte {
    byteArr, err := json.Marshal(pack)
    if err != nil {log.Println(err)}
    return byteArr
}

func serverPackToNetworkPackage(pack ServerPackage) []byte {
    byteArr, err := json.Marshal(pack) 
    if err != nil {log.Println(err)}
    return byteArr
}

func networkPackageToClientPack(byteArr []byte) ClientPackage {
    var ClientPack ClientPackage
    err := json.Unmarshal(byteArr[:], &ClientPack)
    if(err != nil) {log.Println(err)}
    return ClientPack
}

func networkPackageToServerPack(byteArr []byte) ServerPackage {
    var ServerPack ServerPackage
    err := json.Unmarshal(byteArr[:], &ServerPack)
    if (err != nil) {log.Println(err)}
    return ServerPack
}

func validClientPack(p ClientPackage) bool {
    if (p.Request != "login" && p.Request != "logout" && p.Request != "msg" && p.Request != "names" && p.Request != "help"){
        return false
    }
    return true
}

func validServerPack(p ServerPackage) bool {
    if (p.Response != "error" && p.Response!="info" && p.Response!="history" && p.Response!="message"){
        return false
    }
    return true
}

func printClientPackage(pack ClientPackage){
    fmt.Println("Request = ",pack.Request)
    fmt.Println("Content = ",pack.Content)
}

func printServerPackage(pack ServerPackage){
    fmt.Println("Timestamp = ",pack.Timestamp)
    fmt.Println("Sender = ", pack.Sender)
    fmt.Println("Resonse = ", pack.Response)
    fmt.Println("Content = ",pack.Content)
}
//-----ENCODING END

//-----SUBFUNCTIONS START
func getTime() string{
    const layout = "Jan 2, 2006 kl 02:00"
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

func listenToClient(incoming_cl_packet chan InternalClientPackage, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 1024)
        bytes_read, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(conn.RemoteAddr(), "lost connection")
            return
        }
        incoming_cl_packet <- InternalClientPackage{conn, networkPackageToClientPack(buffer[:bytes_read])}
    }
}

func sendToClient(sender, response, content string, conn *net.TCPConn) {
    _, err := conn.Write(serverPackToNetworkPackage( ServerPackage{getTime(), sender, response, content} ))
    if err != nil {
        log.Fatal(err)
    }
}

// TODO: This actually belongs in the client, not the server
func prettyPrintClientMessage(username string, content string) {
    // This prints to the console with colored usernames
    fmt.Printf("\x1b[35m%s\x1b[0m said: %s\n", username, content)
}

//----SUBFUNCTIONS END

func main() {
    connections                 := make(map[*net.TCPConn]string)
    incoming_connection   := make(chan *net.TCPConn)
    incoming_cl_packet     := make(chan InternalClientPackage)

    prettyPrintClientMessage("John doe", "Hey guys!")

    go listenForIncomingConnections(incoming_connection)

    for {
        select {
            case conn := <- incoming_connection:
                connections[conn] = "" // Initial username
                go listenToClient(incoming_cl_packet, conn)
                fmt.Println(conn.RemoteAddr(), "connected")

            case ClientPacket := <- incoming_cl_packet:
                who := ClientPacket.Connection.RemoteAddr() //ip-adress
                fmt.Println("Recived a packet from: ", who.String())
                printClientPackage(ClientPacket.Payload)
        }
    }

    for conn := range(connections) {
        conn.Close()
    }
    fmt.Println("All connections closed. I´l take the day off")
}
