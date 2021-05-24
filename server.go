package main
import (
	"net"
	"strings"
	"log"
)
const (
	PORT = ":8080"
	END_BYTES = "\000\001\002\003\004\005"
)
var (
	Connections = make(map[net.Conn]bool)
)
func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		panic("server error")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil { break }
		go handleConnect(conn)//многопоточность
	}
}
func handleConnect (conn net.Conn) {
	Connections[conn] = true

	var (
		message string 
		buffer = make([]byte, 512)
	)
	close: for {
		message = ""
		for {
			length, err := conn.Read(buffer)
			if err != nil { break close }
			message += string(buffer[:length])
			if strings.HasSuffix(message, END_BYTES) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		log.Println(message);
		for c := range Connections {
			if c == conn { continue }
			c.Write([]byte(strings.ToUpper(message) + END_BYTES))
		
		}
		// delete(Connections, conn)
	}
}