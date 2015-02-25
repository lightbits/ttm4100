# chat
TTM4100 project

How to use client/server
------------------------
1. go run server.go in one console
2. go run client.go in another
3. yayifications!

Todo
----
Misc
* Make diagrams

Server and client
* Send messages as json-encoded byte arrays
* Parse incoming json-encoded byte arrays into structs

Client
* Take server ip address as user input in chat client
* Only accept messages of form
	- /login <username>
	- /logout
	- /msg <message> (this is default if no /* is specified)
	- /names
	- /help
* Parse server response
	- names should print a list
	- info does what?
	- history iterates over all previous messages and prints them?

Server
* Take "listen" port as input on startup?
* How is "history" and "names" transferred?
* Parse client request
	- Validate request field
	- Handle request based on type
* If the sender does not yet have a username, i.e.
	``if (connections[ClientPacket.Connection] == "")``
	we should send back an error, saying they need to /login
