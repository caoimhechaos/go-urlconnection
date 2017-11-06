package urlconnection

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type etcdConnection struct {
	etcdConn *clientv3.Client
}

/*
SetupEtcd instantiates an etcd connection handler with a new etcd client
configured from the parameters passed.

cert, key and ca are optional; if empty strings are passed, an unencrypted
connection will be used.
*/
func SetupEtcd(servers []string, cert, key, ca string) error {
	var conn *clientv3.Client
	var config clientv3.Config = clientv3.Config{
		Endpoints:   servers,
		DialTimeout: 30 * time.Second,
	}
	var tc = new(tls.Config)
	var err error

	tc.RootCAs, err = x509.SystemCertPool()
	if err != nil {
		return err
	}

	if len(ca) > 0 {
		var x509cert *x509.Certificate
		var certPEMBlock []byte
		var certDERBlock *pem.Block

		certPEMBlock, err = ioutil.ReadFile(ca)
		if err != nil {
			return err
		}

		certDERBlock, _ = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			return errors.New("Error decoding certificate " + ca)
		}

		x509cert, err = x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			return err
		}

		tc.RootCAs.AddCert(x509cert)
	}

	if len(cert) > 0 && len(key) > 0 {
		var x509cert tls.Certificate

		x509cert, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return err
		}

		tc.Certificates = append(tc.Certificates, x509cert)
	}
	config.TLS = tc

	conn, err = clientv3.New(config)
	if err != nil {
		return err
	}

	RegisterConnectionHandler("etcd", &etcdConnection{
		etcdConn: conn,
	})
	return nil
}

/*
UseExistingEtcd instantiates an etcd connection handler pointing at a
preconfigured etcd client.
*/
func UseExistingEtcd(client *clientv3.Client) {
	RegisterConnectionHandler("etcd", &etcdConnection{
		etcdConn: client,
	})
}

/*
lookup queries etcd for host:port pairs for the given URL and returns the
corresponding pairs as a string.
*/
func (conn etcdConnection) lookup(ctx context.Context, dest *url.URL) (ret []string, err error) {
	var kv *mvccpb.KeyValue
	var resp *clientv3.GetResponse

	if conn.etcdConn == nil {
		err = errors.New("Please use SetupEtcd first")
		return
	}

	resp, err = conn.etcdConn.Get(ctx, dest.Path)
	if err != nil {
		return
	}

	for _, kv = range resp.Kvs {
		ret = append(ret, string(kv.Value))
	}

	return
}

/*
Connect connects to a host:port pair given in an etcd file.
Makes a TCP connection to the given host:port pair.
*/
func (conn etcdConnection) Connect(ctx context.Context, dest *url.URL) (net.Conn, error) {
	var candidates []string
	var candidate string
	var err error

	candidates, err = conn.lookup(ctx, dest)
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
ConnectTimeout connects to a host:port pair given in an etcd file.
Makes a TCP connection to the given host:port pair.
The attempt is aborted after "timeout".
*/
func (conn etcdConnection) ConnectTimeout(ctx context.Context, dest *url.URL,
	timeout time.Duration) (net.Conn, error) {
	var candidates []string
	var candidate string
	var err error

	candidates, err = conn.lookup(ctx, dest)
	if err != nil {
		return nil, err
	}
	err = errors.New("No connection candidates have been found")
	for _, candidate = range candidates {
		var c net.Conn

		c, err = net.DialTimeout("tcp", candidate, timeout)
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}
