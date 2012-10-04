/*-
 * Copyright (c) 2012 Caoimhe Chaos <caoimhechaos@protonmail.com>,
 *                    Ancient Solutions. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions  of source code must retain  the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions  in   binary  form  must   reproduce  the  above
 *    copyright  notice, this  list  of conditions  and the  following
 *    disclaimer in the  documentation and/or other materials provided
 *    with the distribution.
 *
 * THIS  SOFTWARE IS  PROVIDED BY  ANCIENT SOLUTIONS  AND CONTRIBUTORS
 * ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO,  THE IMPLIED WARRANTIES OF  MERCHANTABILITY AND FITNESS
 * FOR A  PARTICULAR PURPOSE  ARE DISCLAIMED.  IN  NO EVENT  SHALL THE
 * FOUNDATION  OR CONTRIBUTORS  BE  LIABLE FOR  ANY DIRECT,  INDIRECT,
 * INCIDENTAL,   SPECIAL,    EXEMPLARY,   OR   CONSEQUENTIAL   DAMAGES
 * (INCLUDING, BUT NOT LIMITED  TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE,  DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT  LIABILITY,  OR  TORT  (INCLUDING NEGLIGENCE  OR  OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package urlconnection

import (
	"errors"
	"net"
	"net/url"
	"time"
)

type ConnectionHandler interface {
	Connect(dest *url.URL) (net.Conn, error)
	ConnectTimeout(dest *url.URL, timeout time.Duration) (net.Conn, error)
}

var handlers map[string]ConnectionHandler = make(map[string]ConnectionHandler)

/**
 * Used by clients to register their handlers. The URL schema "name"
 * (the part before the dot, i.e. name://...) will be associated with
 * the "handler", which creates connections.
 */
func RegisterConnectionHandler(name string, handler ConnectionHandler) {
	handlers[name] = handler
}

/**
 * Establish a connection to the service specified in "tourl".
 * For example, tcp://[::1]:8080 will establish a TCP connection
 * to port 8080 on localhost, as will tcp://[::1]/8080.
 */
func Connect(tourl string) (net.Conn, error) {
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
		return handler.Connect(u)
	}

	return nil, errors.New("No handler found for " + u.Scheme)
}

/**
 * Establish a connection to the service specified in "tourl".
 * For example, tcp://[::1]:8080 will establish a TCP connection
 * to port 8080 on localhost, as will tcp://[::1]/8080.
 * Abort the attempt after "timeout" has expired.
 */
func ConnectTimeout(tourl string, timeout time.Duration) (net.Conn, error) {
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
		return handler.ConnectTimeout(u, timeout)
	}

	return nil, errors.New("No handler found for " + u.Scheme)
}
