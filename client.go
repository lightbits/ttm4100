package main

import (
    "net"
    "log"
    "fmt"
    "bufio"
    "os"
    "time"
    "strings"
    "encoding/json"
)

//----ENCODING START
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

func networkPackageToServerPack(byteArr []byte) ServerPackage {
    var ServerPack ServerPackage
    err := json.Unmarshal(byteArr[:], &ServerPack)
    if (err != nil) {log.Println(err)}
    return ServerPack
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
        fmt.Printf(">> ")
        line, _, err := reader.ReadLine()
        if err != nil {
            fmt.Println(err)
        }
        user_input <- string(line)
    }
}

func sendToServer(payload ClientPackage, conn *net.TCPConn) {
    _, err := conn.Write(clientPackToNetworkPackage(payload))
    if err != nil {
        log.Fatal(err)
    }
}

func parseUserInput(input string) ClientPackage {
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

    return ClientPackage{request, content}
}

func prettyPrintClientMessage(username string, content string) {
    fmt.Printf("\x1b[35m%s\x1b[0m said: %s\n", username, content)
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

    prettyPrintClientMessage("John doe", "Hey guys!")

    incoming_message := make(chan ServerPackage)
    go listenForMessages(incoming_message, conn)

    user_input := make(chan string)
    go listenForUserInput(user_input)

    for {
        select {
            case msg := <- incoming_message:
                printServerPackage(msg)

            case input := <- user_input:
                payload := parseUserInput(input)
                sendToServer(payload, conn)
        }
    }
}
