/*
Functions for abstracting target resolution for establishing new network
connections by using URLs.
*/
package urlconnection

import (
	"errors"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

/*
ConnectionHandler specifies the basic functionality which must be provided by
any connection setup handlers in order to be registered.
*/
type ConnectionHandler interface {
	Connect(ctx context.Context, dest *url.URL) (net.Conn, error)
	ConnectTimeout(ctx context.Context, dest *url.URL, timeout time.Duration) (net.Conn, error)
}

var handlers map[string]ConnectionHandler = make(map[string]ConnectionHandler)

/*
RegisterConnectionHandler is used by connection APIs to register their handlers.
The URL schema "name" (the part before the dot, i.e. name://...) will be
associated with the "handler", which creates connections.
*/
func RegisterConnectionHandler(name string, handler ConnectionHandler) {
	handlers[name] = handler
}

/*
Connect establishes a new connection to the URL, determining the correct handler
to use for the schema.
For example, tcp://[::1]:8080 will establish a TCP connection to port 8080 on
localhost, as will tcp://[::1]/8080.
*/
func Connect(tourl string) (net.Conn, error) {
	var ctx context.Context = context.Background()
	var u *url.URL
	var err error
	var ok bool

	u, err = url.Parse(tourl)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		return nil, errors.New("Absolute URL required")
	}

	if _, ok = handlers[u.Scheme]; ok {
		var handler ConnectionHandler = handlers[u.Scheme]
		return handler.Connect(ctx, u)
	}

	return nil, errors.New("No handler found for " + u.Scheme)
}

/*
ConnectContext establishes a new connection to the URL, passing the context down
to the handler and determining the correct handler to use for the schema.
For example, tcp://[::1]:8080 will establish a TCP connection to port 8080 on
localhost, as will tcp://[::1]/8080.
*/
func ConnectContext(ctx context.Context, tourl string) (net.Conn, error) {
	var u *url.URL
	var err error
	var ok bool

	u, err = url.Parse(tourl)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		return nil, errors.New("Absolute URL required")
	}

	if _, ok = handlers[u.Scheme]; ok {
		var handler ConnectionHandler = handlers[u.Scheme]
		return handler.Connect(ctx, u)
	}

	return nil, errors.New("No handler found for " + u.Scheme)
}

/*
ConnectTimeout establishes a new connection to the URL, waiting at most the
specified timeout for the connection to be established, determining the correct
handler to use for the schema.
For example, tcp://[::1]:8080 will establish a TCP connection to port 8080 on
localhost, as will tcp://[::1]/8080.
*/
func ConnectTimeout(tourl string, timeout time.Duration) (net.Conn, error) {
	var ctx context.Context
	var u *url.URL
	var err error
	var ok bool

	ctx, _ = context.WithTimeout(context.Background(), timeout)

	u, err = url.Parse(tourl)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		return nil, errors.New("Absolute URL required")
	}

	if _, ok = handlers[u.Scheme]; ok {
		var handler ConnectionHandler = handlers[u.Scheme]
		return handler.ConnectTimeout(ctx, u, timeout)
	}

	return nil, errors.New("No handler found for " + u.Scheme)
}
