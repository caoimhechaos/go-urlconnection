/*-
 * Copyright (c) 2012 Tonnerre Lombard <tonnerre@ancient-solutions.com>,
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
	"net"
	"net/url"
	"time"
)

type tcpConnection struct {}

/**
 * TCP connection creator.
 * Make a TCP connection to the given host (in format
 * tcp://[::1]:123 or tcp://[::1]/123.
 */
func (self tcpConnection) Connect(dest *url.URL) (net.Conn, error) {
	var host, port, hostport string

	host = dest.Host
	hostport = host
	if len(dest.Path) > 0 {
		port = dest.Path[1:]
		hostport = net.JoinHostPort(host, port)
	}

	return net.Dial("tcp", hostport)
}

/**
 * TCP connection creator.
 * Make a TCP connection to the given host (in format
 * tcp://[::1]:123 or tcp://[::1]/123. Abort attempt after
 * "timeout" has expired.
 */
func (self tcpConnection) ConnectTimeout(dest *url.URL,
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

/**
 * Register the connection handler for TCP.
 */
func init() {
	RegisterConnectionHandler("tcp", tcpConnection{})
}
