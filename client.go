package main

import (
    "net"
    "log"
    "fmt"
    "bufio"
    "os"
    "time"
    "encoding/json"
)

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

func getTime() string{
    const layout = "Jan 2, 2006 kl 02:00"
    return time.Now().Format(layout)
}
//-----ENCODING END

func listenForMessages(incoming_message chan ServerPackage, conn *net.TCPConn) {
    for {
        buffer := make([]byte, 1024)
        bytes_read, err := conn.Read(buffer)
        if err != nil {
            log.Fatal(err)
        }
        incoming_message <- networkPackageToServerPack(buffer[:bytes_read])
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

    incoming_message := make(chan ServerPackage)
    go listenForMessages(incoming_message, conn)

    user_input := make(chan string)
    go listenForUserInput(user_input)

    for {
        select {
            case msg := <- incoming_message:
                log.Println("Received message from server")
                printServerPackage(msg)

            case input := <- user_input:
                _, err := conn.Write(clientPackToNetworkPackage(ClientPackage{"msg",input}))
                if err != nil {
                    log.Fatal(err)
                }
        }
    }
}
