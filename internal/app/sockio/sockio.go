package sockio

import (
	"log"
	"time"

	"github.com/google/wire"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// ProviderSet 注入
var ProviderSet = wire.NewSet(New)

func New() (*socketio.Server, func(), error) {

	// cfg := config.C

	//log.Println(" --== cfg:", cfg)
	//server := socketio.NewServer(nil)
	transports := []transport.Transport{
		polling.Default,
		websocket.Default,
	}
	server := socketio.NewServer(&engineio.Options{
		Transports:   transports,
		PingInterval: time.Second * 30,
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

		log.Println(" --== connected:", s.ID())
		log.Println(" --== connected URL:", s.URL())
		log.Println(" --== connected LocalAddr:", s.LocalAddr())
		log.Println(" --== connected RemoteAddr:", s.RemoteAddr())
		log.Println(" --== connected RemoteHeader:", s.RemoteHeader())

		//s.Emit("ok", "welcome")

		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println(" ====000-- OnError error:", e)
		log.Println(" ====000-- OnError error:", s.ID())
		log.Println(" ====000-- OnError error:", s.RemoteAddr())
		log.Println(" ====000-- OnError error:", s.URL())
	})

	server.OnEvent("/", "start", func(s socketio.Conn, msg string) {
		log.Println(" ------ ==== msg: ", msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println(" ===== 99999 = closed ", reason)
	})

	// server.Adapter(&socketio.RedisAdapterOptions{
	// 	Addr:    cfg.Redis.DSN(),
	// 	Prefix:  "socket.io",
	// 	Network: "tcp",
	// })

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	cleanFunc := func() {
		server.Close()
		log.Println(" clean in sockio ")
	}

	return server, cleanFunc, nil
}
