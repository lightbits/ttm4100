# chat
TTM4100 project

Mockup code
	
	func ListenForIncomingConnections() {
	
		local    = net.ResolveTCPAddr("tcp", ":12345") <- define server listen port
		listener = net.ListenTCP("tcp", local)
		
		connection = listener.AcceptTCP() <- the incoming client's tcp connection
		
		go ListenToClient(connection)
	}
	
	type RequestType int
	const (
		RequestLogin RequestType = iota
		RequestLogout
		RequestMsg
		RequestNames
		RequestHelp
	)
	
	type Data struct {
		request RequestType
		content string
	}
	
	
	message {
		timestamp string
		content string
	}
		
	
	MessageHistory map[username][]message
	
	
	
	func ListenToClient(connection) {
		username = ""
	
		for {
			select {
			case bytes = <- connection.ReadFrom():
				data = ParseBytes(byte)
				
				switch (data.request) {
					case RequestLogin:
					
						if (!ValidateUserName(data.content))
							SendError()
							
						string
						for char in string:
						
						username = data.content
						
						// send message history
						for username, msg in MessageHistory {
							r = Response{
								timestamp: msg.timestamp,
								sender:    username,
								response:  'history',
								content:   msg.content}
								
							bytes = EncodeResponsePacket(r)
							network.SendBytes(bytes)
						}
						
					case RequestLogout:
					
						// close connection
						
					case RequestMsg:
					
						if !username
							SendError(connection)
						else
							MessageHistory[username].append({time.Now(), data.content})
					
					case RequestNames:
					case RequestHelp:
					default:
						SendError(connection)
			}
		}
	}
	
	
	
	
	{‘timestamp’: <timestamp>,‘sender’: <username>,‘response’: <response>,‘content’: <content>}
	
	
	
	
	
	main() {
		connections = map[TCPConn]username
		for {
			select {			
				for c in connections {
					case bytes = <- c.ReadFrom():
						HandleMessage(c, bytes)
				}
			}
		}
	}
	
