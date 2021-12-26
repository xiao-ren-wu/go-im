package serv

import (
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	once    sync.Once
	id      string
	address string
	sync.Mutex
	users map[string]net.Conn
}

func NewServer(id, address string) *Server {
	return newServer(id, address)
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(
		logrus.Fields{
			"module": "Server",
			"listen": s.address,
			"id":     s.id,
		},
	)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			_ = conn.Close()
			return
		}
		user := r.URL.Query().Get("user")
		if "" == user {
			_ = conn.Close()
			return
		}
		oldConn, ok := s.addUser(user, conn)
		if ok {
			_ = oldConn.Close()
		}
		log.Infof("user %s in", user)
		go func(user string, conn net.Conn) {
			err = s.readLoop(user, conn)
			if err != nil {
				log.Error(err.Error())
			}
			_ = conn.Close()
			s.delUser(user)
			log.Infof("connection of %s closed", user)
		}(user, conn)
	})
	log.Infoln("started")
	return http.ListenAndServe(s.address, mux)
}

func (s *Server) Shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()
		for _, conn := range s.users {
			_ = conn.Close()
		}
	})
}

func newServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   make(map[string]net.Conn),
	}
}

func (s *Server) readLoop(user string, conn net.Conn) error {
	for {
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}
		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote site close the conn")
		}
		if frame.Header.Masked {
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}
		if frame.Header.OpCode == ws.OpText {
			go s.handle(user, string(frame.Payload))
		}
	}
}

func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
}

func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()
	old, ok := s.users[user]
	s.users[user] = conn
	return old, ok
}

func (s *Server) handle(user string, message string) {
	logrus.Infof("recv msg %s from %s", message, user)
	s.Lock()
	defer s.Unlock()
	broadcast := fmt.Sprintf("%s -- FROM %s", message, user)
	for u, conn := range s.users {
		if u == user {
			continue
		}
		logrus.Infof("send to %s : %s", u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed,err info: %s", u, err.Error())
		}
	}
}

func (s *Server) writeText(conn net.Conn, broadcast string) error {
	frame := ws.NewTextFrame([]byte(broadcast))
	return ws.WriteFrame(conn, frame)
}
