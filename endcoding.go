package main

import(
	"fmt"
	"time"
)

type clientPackage struct {
	request string
	content string
}

type serverPackage struct {
	timestamp time.Time
	sender string
	response string
	content string
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

}

func printClientPackage(pack clientPackage){
	fmt.Println("Request = ",pack.request)
	fmt.Println("Content = ",pack.content)
}

func printServerPackage(pack serverPackage){
	fmt.Println("Timestamp = ",pack.timestamp)
	fmt.Println("Sender = ", pack.sender)
	fmt.Println("Resonse = ", pack.response)
	fmt.Println("Content = ",pack.content)
}