package client

import (
	"fmt"
	"io"

	"github.com/gorilla/websocket"
)

var (
	// DefaultReadBufferSize is the default maximum read buffer size.
	DefaultReadBufferSize = 25 * 1024 * 1024

	// DefaultWriteBufferSize is the default maximum write buffer size.
	DefaultWriteBufferSize = 10 * 1024 * 1024
)

// Transport is the common interface to send/receive messages to a target.
type Transport interface {
	Read() ([]byte, error)
	Write([]byte) error
	io.Closer
}

// Conn wraps a gorilla/websocket.Conn connection.
type Conn struct {
	*websocket.Conn
}

// Read reads the next websocket message.
func (c *Conn) Read() ([]byte, error) {
	_, buf, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Write writes a websocket message.
func (c *Conn) Write(buf []byte) error {
	fmt.Println(string(buf))
	return c.WriteMessage(websocket.TextMessage, buf)
}

// Dial dials the specified target's websocket URL.
//
// Note: uses gorilla/websocket.
func Dial(urlstr string, opts ...DialOption) (Transport, error) {
	d := &websocket.Dialer{
		ReadBufferSize:  DefaultReadBufferSize,
		WriteBufferSize: DefaultWriteBufferSize,
	}

	// apply opts
	for _, o := range opts {
		o(d)
	}

	// connect
	conn, _, err := d.Dial(urlstr, nil)
	if err != nil {
		return nil, err
	}

	return &Conn{conn}, nil
}

// DialOption is a dial option.
type DialOption func(*websocket.Dialer)
