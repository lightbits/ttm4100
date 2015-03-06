package main

import (
    "net"
    "log"
    "fmt"
    "bufio"
    "os"
    "time"
    "encoding/json"
    "strings"
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

type userInputPack struct {
    Command string
    Payload string
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

func listenForUserInput(user_input chan userInputPack, userInputTrigger chan int) {
    reader := bufio.NewReader(os.Stdin)
    for{
        select{
        case <- userInputTrigger:
                fmt.Printf("Enter text: ")
                line, _, err := reader.ReadLine()
                if err != nil {
                    log.Println(err)
                }
                streng := strings.TrimSpace(string(line))
                splitIndex := strings.Index(streng," ")
                var c, p string
                if splitIndex != -1 {
                    c = streng[:splitIndex]
                    if splitIndex != len(streng) {
                        p = streng[splitIndex+1:]
                    }else{
                        p = ""
                    }                   
                }else{
                    c = streng
                    p = ""
                }
                log.Println("Command: ", c)
                log.Println("Payload: ", p)
                user_input <- userInputPack{c,p}
        }
    }
}

func sendToServer(request, content string, conn *net.TCPConn) {
    _, err := conn.Write(clientPackToNetworkPackage( ClientPackage{request, content} ))
    if err != nil {
        log.Fatal(err)
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

    user_input := make(chan userInputPack)
    userInputTrigger := make(chan int)
    go listenForUserInput(user_input, userInputTrigger)

   
    var waitingOnServer bool = false
    userInputTrigger <- 1

    for {
        select {
            case msg := <- incoming_message:
                log.Println("Received message from server")
                printServerPackage(msg)
                waitingOnServer = false
           
            case input := <- user_input:
                switch input.Command {
                    case "login":
                        sendToServer("login","",conn)
                        waitingOnServer = true
                    case "logout":
                        sendToServer("logout","",conn)
                        waitingOnServer = true
                    case "msg":
                        sendToServer("msg",input.Payload,conn)
                        waitingOnServer = true
                    case "names":
                        sendToServer("names","",conn)
                        waitingOnServer = true
                    case "help":
                        sendToServer("help", "",conn)
                        waitingOnServer = true
                    default:
                        fmt.Println("Ugyldig kommando")
                        waitingOnServer = false
                        userInputTrigger <- 1
                }
            case <- time.After(3 * time.Second):
                if(waitingOnServer){
                    waitingOnServer = false
                    fmt.Println("waitingOnServer timed out")
                    userInputTrigger <- 1
                }
        }
    }
}
