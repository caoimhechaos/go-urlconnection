package urlconnection

import (
	"net"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

type tcpConnection struct{}

/*
Connect creates a new TCP connection.
Expected format for the URL is tcp://[::1]:123 or tcp://[::1]/123.
*/
func (tcpConnection) Connect(ctx context.Context, dest *url.URL) (net.Conn, error) {
	var host, port, hostport string

	host = dest.Host
	hostport = host
	if len(dest.Path) > 0 {
		port = dest.Path[1:]
		hostport = net.JoinHostPort(host, port)
	}

	return net.Dial("tcp", hostport)
}

/*
ConnectTimeout creates a new TCP connection with a timeout for establishing.
Expected format for the URL is tcp://[::1]:123 or tcp://[::1]/123.
Abort attempt after "timeout" has expired.
*/
func (tcpConnection) ConnectTimeout(ctx context.Context, dest *url.URL,
	timeout time.Duration) (net.Conn, error) {
	var host, port, hostport string

	host = dest.Host
	hostport = host
	if len(dest.Path) > 0 {
		port = dest.Path[1:]
		hostport = net.JoinHostPort(host, port)
	}

	return net.DialTimeout("tcp", hostport, timeout)
}

/*
Register the connection handler for TCP.
*/
func init() {
	RegisterConnectionHandler("tcp", tcpConnection{})
}
