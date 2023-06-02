package websocket

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/xiao-ren-wu/go-im/dal/nameing"
	"github.com/xiao-ren-wu/go-im/middleware"
	"net/http"
	"sync"
	"time"

	"github.com/xiao-ren-wu/go-im/dal"
	"github.com/xiao-ren-wu/go-im/dal/constants"
)

type ServerOptions struct {
	loginwait time.Duration
	readwait  time.Duration
	writewait time.Duration
}

type Server struct {
	listen string
	nameing.ServerRegistration
	dal.ChannelMap
	dal.Acceptor
	dal.MessageListener
	dal.StateListener
	once    sync.Once
	options ServerOptions
}

func NewServer(listen string, service nameing.ServerRegistration) dal.Server {
	return &Server{
		listen:             listen,
		ServerRegistration: service,
		options: ServerOptions{
			loginwait: constants.DefaultLoginWait,
			readwait:  constants.DefaultReadWait,
			writewait: constants.DefaultWriteWait,
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	if s.Acceptor == nil {
		s.Acceptor = new(dal.DefaultAcceptor)
	}
	if s.StateListener == nil {
		return fmt.Errorf("StateListener is nil")
	}
	if s.ChannelMap == nil {
		s.ChannelMap = dal.NewChannelMap(100)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rawConn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			resp(w, http.StatusBadRequest, err.Error())
			return
		}
		conn := NewWsConn(rawConn)
		id, _, err := s.Accept(conn, s.options.loginwait)
		if err != nil {
			middleware.L.Error("accept error, err info: %v", err)
			_ = conn.WriteFrame(constants.OpClose, []byte(err.Error()))
			_ = conn.Close()
			return
		}
		if _, ok := s.Get(id); ok {
			middleware.L.Warn("channel %s existed", id)
			_ = conn.WriteFrame(constants.OpClose, []byte("channelId is repeated"))
			_ = conn.Close()
			return
		}
		channel := dal.NewChannel(id, conn)
		channel.SetWriteWait(s.options.writewait)
		channel.SetReadWait(s.options.readwait)
		s.Add(channel)

		go func(ch dal.Channel) {
			err := ch.Readloop(s.MessageListener)
			if err != nil {
				middleware.L.Warn(err.Error())
			}
			s.Remove(id)
			if err = s.Disconnect(id); err != nil {
				middleware.L.Warn(err.Error())
			}
			if err = ch.Close(); err != nil {
				middleware.L.Warn(err.Error())
			}
		}(channel)
	})

	return http.ListenAndServe(s.listen, mux)
}

func resp(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errMsg))
}

func (s *Server) SetAcceptor(acceptor dal.Acceptor) {
	s.Acceptor = acceptor
}

func (s *Server) SetMessageListener(listener dal.MessageListener) {
	s.MessageListener = listener
}

func (s *Server) SetStateListener(listener dal.StateListener) {
	s.StateListener = listener
}

func (s *Server) SetReadWait(duration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) SetChannelMap(channelMap dal.ChannelMap) {
	s.ChannelMap = channelMap
}

func (s *Server) Push(s2 string, bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Shutdown(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
