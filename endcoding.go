package main

import(
	"fmt"
	"time"
	"encoding/json"
	"log"
)

type ClientPackage struct {
	Request	string
	Content 	string
}

type ServerPackage struct {
	Timestamp	time.Time
	Sender 	string
	Response	string
	Content	string
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

func networkPackageToClientPack(byteArr []byte)(ClientPackage, error){
	var ClientPack ClientPackage
	err := json.Unmarshal(byteArr[:], &ClientPack)
	//if (!validClientPack(ClientPack)){return nil, err}<--- Should this check be in this module?
	return ClientPack, err
}

func networkPackageToServerPack(byteArr []byte)(ServerPackage, error){
	var ServerPack ServerPackage
	err := json.Unmarshal(byteArr[:], &ServerPack)
	//if (!validServerPack(ServerPack)){return nil, err} <--- Should this check be in this module?
	return ServerPack, err
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

func main() {
	ClientTestpackage	 := ClientPackage{"login" , "eriklil"}
	ServerTestpackage := ServerPackage{time.Now(), "gunnar", "message", "Hei, hvordan står det til. Hilsen Gunnar"}

	fmt.Println("----Client Test Package----")
	printClientPackage(ClientTestpackage)
	fmt.Println()
	fmt.Println("----Server Test Package----")
	printServerPackage(ServerTestpackage)
	fmt.Println("-------------------------------------------------------")

	b1 := clientPackToNetworkPackage(ClientTestpackage)
	if (b1 != nil) {
		fmt.Println("Alt gikk bra med clientTestpackage encoding:", string(b1))
	}

	b2 := serverPackToNetworkPackage(ServerTestpackage)
	if (b2 != nil) {
		fmt.Println("Alt gikk bra med serverTestpackage encoding:", string(b2))
	}

	fmt.Println()

	restoredClientPackage, err := networkPackageToClientPack(b1)
	if(err != nil){
		log.Println(err)
	}else{
		fmt.Println("Alt gikk bra med ClientTestpackage decoding:")
		printClientPackage(restoredClientPackage)
	}

	fmt.Println()

	restoredServerPackage, err := networkPackageToServerPack(b2)
	if(err != nil){
		log.Println(err)
	}else{
		fmt.Println("Alt gikk bra med ServerTestpackage decoding:")
		printServerPackage(restoredServerPackage)
	}
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