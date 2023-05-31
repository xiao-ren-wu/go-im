package websocket

import (
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
	name.ServerRegistration
	dal.ChannelMap
	dal.Acceptor
	dal.MessageListener
	dal.StateListener
	once    sync.Once
	options ServerOptions
}

func NewServer(listen string, service nameing.ServiceRegistration) dal.Server {
	return &Server{
		listen:             listen,
		ServerRegistration: service,
		options: ServerOptions{
			loginwait: constants.DefaultLoginWait,
			readwait:  constants.DefauReadWait,
			writewait: constants.DefaultWriteWait,
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	if s.Acceptor == nil {
		s.Acceptor = new(defaultAcceptor)
	}
}
