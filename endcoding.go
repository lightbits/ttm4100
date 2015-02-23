package main

import(
	"fmt"
	"time"
	"encoding/json"
	"log"
	"bytes"
)

type clientPackage struct {
	Request	string
	Content 	string
}

type serverPackage struct {
	Timestamp	time.Time
	Sender 	string
	Response	string
	Content	string
}

func clientPackToNetworkPackage(pack clientPackage) []byte {
	var csvBuffer bytes.Buffer
	csvBuffer.WriteString("{")
	csvBuffer.WriteString("´request´:")
	csvBuffer.WriteString(pack.Request)
	csvBuffer.WriteString(",´content´:")
	csvBuffer.WriteString(pack.Content)
	csvBuffer.WriteString("}")
	byteArr, err := json.Marshal(csvBuffer.String())
	if err != nil {log.Println(err)}
	return byteArr
}

func serverPackToNetworkPackage(pack serverPackage) []byte {
	var csvBuffer bytes.Buffer
	csvBuffer.WriteString("{")
	csvBuffer.WriteString("´timestamp´:")
	csvBuffer.WriteString(pack.Timestamp.String())
	csvBuffer.WriteString(",´sender´:")
	csvBuffer.WriteString(pack.Sender)
	csvBuffer.WriteString(",´response´:")
	csvBuffer.WriteString(pack.Response)
	csvBuffer.WriteString(",´content´:")
	csvBuffer.WriteString(pack.Content)
	csvBuffer.WriteString("}")
	byteArr, err := json.Marshal(csvBuffer.String())
	if err != nil {log.Println(err)}
	return byteArr
}

func networkPackageToClientPack(b []byte) clientPackage{
	var csvString string
	err := json.Unmarshal(b[:], &csvString)
	if (err != nil) {log.Println(err)}
	var clientPack clientPackage
	fmt.Println(csvString)
	return clientPack
}



func main() {
	clientTestpackage	 := clientPackage{"login" , "eriklil"}
	serverTestpackage := serverPackage{time.Now(), "gunnar", "message", "Hei, hvordan står det til. Hilsen Gunnar"}

	fmt.Println("Client Test Package:")
	printClientPackage(clientTestpackage)
	fmt.Println()
	fmt.Println("Server Test Package")
	printServerPackage(serverTestpackage)
	fmt.Println()
	fmt.Println("Ferdig!!!")

	b1 := clientPackToNetworkPackage(clientTestpackage)
	if (b1 != nil) {
		fmt.Println("Alt gikk bra med clientTestpackage")
	}

	b2 := serverPackToNetworkPackage(serverTestpackage)
	if (b2 != nil) {
		fmt.Println("Alt gikk bra med serverTestpackage")
	}

	_=networkPackageToClientPack(b1)
}

func printClientPackage(pack clientPackage){
	fmt.Println("Request = ",pack.Request)
	fmt.Println("Content = ",pack.Content)
}

func printServerPackage(pack serverPackage){
	fmt.Println("Timestamp = ",pack.Timestamp)
	fmt.Println("Sender = ", pack.Sender)
	fmt.Println("Resonse = ", pack.Response)
	fmt.Println("Content = ",pack.Content)
}