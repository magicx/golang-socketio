package gosocketio

import (
	"github.com/graarh/golang-socketio/transport"
	"strconv"
	"net/http"
)

const (
	webSocketProtocol = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl       = "/socket.io/?EIO=3&transport=websocket"
)

/**
Socket.io client representation
*/
type Client struct {
	methods
	Channel
}

/**
Get ws/wss url by host and port
 */
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + host + ":" + strconv.Itoa(port) + socketioUrl
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, tr transport.Transport) (*Client, error) {
	c := &Client{}
	c.initChannel()
	c.initMethods()

	var err error
	c.conn, err = tr.Connect(url)
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return c, nil
}

// +
func DialwithConnAndHeader(url string, tr transport.Transport, f interface{}, header http.Header) (*Client, error) {
	c := &Client{}
	c.initChannel()
	c.initMethods()

	var err error
	//+
	if (nil != f) {
		err = c.On(OnConnection, f)
		if err != nil {
			return nil, err
		}
	}
	
	c.conn, err = tr.ConnectWithHeader(url, header)
	if err != nil {
		return nil, err
	}
	//-
	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return c, nil
}
// -

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
