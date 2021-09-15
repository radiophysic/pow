package lib

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Open stands for client-to-server connection.
func Open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

type HandleFunc func(*bufio.ReadWriter)

// EndPoint is a thread-safe structure.
type EndPoint struct {
	m        sync.RWMutex
	listener net.Listener
	handler  map[string]HandleFunc
	log      *log.Logger
}

func NewEndPoint(logger *log.Logger) *EndPoint {
	return &EndPoint{
		handler: map[string]HandleFunc{},
		log:     logger,
	}
}

// AddHandleFunc Add data type processing method.
func (e *EndPoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// handleMessage Verify the requested data type and send it to the corresponding processing function.
func (e *EndPoint) handleMessage(conn net.Conn) {
	rw := bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)
	defer func() {
		if err := conn.Close(); err != nil {
			e.log.Println("Closing connection error" + err.Error())
		}
	}()

	for {
		cmd, err := rw.ReadString('\n')
		switch {
		case errors.Is(err, io.EOF):
			e.log.Println("Client disconnected.")
			return
		case err != nil:
			e.log.Println("read error")
			return
		}

		cmd = strings.Trim(cmd, "\n ")
		e.m.RLock()
		handleCmd, ok := e.handler[cmd]
		if !ok {
			e.log.Println("Unregistered request data type.")
			return
		}
		handleCmd(rw)
	}
}

// Listen stands for server listen for client connections.
func (e *EndPoint) Listen(port string) (err error) {
	var conn net.Conn
	e.listener, err = net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "Service cannot be bound on port"+port)
	}

	e.log.Println("Service live: ", e.listener.Addr().String())

	for {
		conn, err = e.listener.Accept()
		if err != nil {
			e.log.Println("Heart request monitoring failed!")
			continue
		}
		go e.handleMessage(conn)
	}
}
