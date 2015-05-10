/*-
 * Copyright (c) 2015 Caoimhe Chaos <caoimhechaos@protonmail.com>,
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

	"github.com/coreos/go-etcd/etcd"
)

type etcdConnection struct {
	etcd_conn *etcd.Client
}

/*
Set the etcd connection parameters to be used.

cert, key and ca are optional; if empty strings are passed, an unencrypted
connection will be used.
*/
func SetupEtcd(servers []string, cert, key, ca string) error {
	var conn *etcd.Client
	var err error

	if len(cert) == 0 || len(key) == 0 {
		conn = etcd.NewClient(servers)
	} else {
		conn, err = etcd.NewTLSClient(servers, cert, key, ca)
		if err != nil {
			return err
		}
	}

	RegisterConnectionHandler("etcd", &etcdConnection{
		etcd_conn: conn,
	})
	return nil
}

/*
Use an existing etcd client to pick backends.
*/
func UseExistingEtcd(client *etcd.Client) {
	RegisterConnectionHandler("etcd", &etcdConnection{
		etcd_conn: client,
	})
}

/*
Queries etcd for host:port pairs for the given URL. Returns the
corresponding pairs as a string.
*/
func (conn etcdConnection) lookup(dest *url.URL) (ret []string, err error) {
	var node *etcd.Node
	var resp *etcd.Response

	if conn.etcd_conn == nil {
		err = errors.New("Please use SetupEtcd first")
		return
	}

	resp, err = conn.etcd_conn.Get(dest.Path, false, false)
	if err != nil {
		return
	}

	if resp.Node.Dir {
		for _, node = range resp.Node.Nodes {
			ret = append(ret, node.Value)
		}
	} else {
		ret = append(ret, resp.Node.Value)
	}

	return
}

/*
Connect to a host:port pair given in an etcd file.
Makes a TCP connection to the given host:port pair.
*/
func (conn etcdConnection) Connect(dest *url.URL) (net.Conn, error) {
	var candidates []string
	var candidate string
	var err error

	candidates, err = conn.lookup(dest)
	if err != nil {
		return nil, err
	}
	err = errors.New("No connection candidates have been found")
	for _, candidate = range candidates {
		var c net.Conn

		c, err = net.Dial("tcp", candidate)
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}

/*
Connect to a host:port pair given in an etcd file.
Makes a TCP connection to the given host:port pair.
The attempt is aborted after "timeout".
*/
func (conn etcdConnection) ConnectTimeout(dest *url.URL,
	timeout time.Duration) (net.Conn, error) {
	var candidates []string
	var candidate string
	var err error

	candidates, err = conn.lookup(dest)
	if err != nil {
		return nil, err
	}
	err = errors.New("No connection candidates have been found")
	for _, candidate = range candidates {
		var c net.Conn
		var err error

		c, err = net.DialTimeout("tcp", candidate, timeout)
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}
