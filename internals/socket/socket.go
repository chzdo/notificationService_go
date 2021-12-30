package socket

import (
	"fmt"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type Socket struct {
	Socket socketio.Server
}

type SocketData struct {
	NameSpace string
	Room      string
	Event     string
	Data      interface{}
}

var Websocket = &Socket{
	Socket: initSocket(),
}

func initSocket() socketio.Server {

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(c socketio.Conn) error {
		orgId := getQuery(c.URL().RawQuery)["orgId"]
		c.Join(fmt.Sprintf(`room-%s`, orgId))
		fmt.Println("connected")
		return nil
	})

	server.OnConnect("/social", func(c socketio.Conn) error {
		orgId := getQuery(c.URL().RawQuery)["orgId"]
		c.Join(fmt.Sprintf(`room-%s`, orgId))
		fmt.Println("connected 2")
		return nil
	})

	server.OnEvent("/", "INCOMING:FACEBOOK", func(s socketio.Conn, msg interface{}) interface{} {
		s.Emit("NEW:FACEBOOK", msg)
		return msg
	})

	server.OnEvent("/", "INCOMING:TWEET", func(s socketio.Conn, msg interface{}) interface{} {
		s.Emit("NEW:TWEET", msg)
		return msg
	})
	server.OnEvent("/", "INCOMING:MAIL", func(s socketio.Conn, msg interface{}) interface{} {
		s.Emit("NEW:MAIL", msg)
		return msg
	})
	server.OnEvent("/", "INCOMING:WHATSAPP", func(s socketio.Conn, msg interface{}) interface{} {
		s.Emit("NEW:WHATSAPP", msg)
		return msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	return *server
}

func getQuery(s string) map[string]interface{} {

	firstSplit := strings.Split(s, "&")
	mapToReturn := make(map[string]interface{})

	for _, value := range firstSplit {

		temp := strings.Split(value, "=")
		mapToReturn[temp[0]] = temp[1]

	}

	return mapToReturn
}

func (socket *Socket) BroadCastNotification(namespace string, room string, event string, data interface{}) bool {
	return socket.Socket.BroadcastToRoom(namespace, room, event, data)
}
