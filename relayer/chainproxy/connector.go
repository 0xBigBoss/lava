package chainproxy

//
// Right now this is only for Ethereum
// TODO: make this into a proper connection pool that supports
// the chainproxy interface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/lavanet/lava/relayer/chainproxy/clients"
	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Client interface {
	Close()
	Call(ctx context.Context, result *json.RawMessage, method string, params interface{}) (interface{}, error)
	ClientType() string
	GetResponse() (json.RawMessage, error)
}

type Connector struct {
	lock        sync.Mutex
	freeClients []Client
	usedClients int
}

func addPortAndParse(originalAddress string) string {
	port := "80" // set default port
	// search where to put the port. some ips has complex paths like: http://user:password@1.1.1.1/path/more_path/etc...
	// understand if there is anything after the ip. or its empty
	urlSplitted := strings.Split(originalAddress, `.`)                                                         // split dots
	lastIpAddressPart := urlSplitted[len(urlSplitted)-1]                                                       // get last dot and path
	lastPartPaths := strings.SplitN(lastIpAddressPart, "/", 2)                                                 // split paths and last ip number
	lastPartPaths[0] = strings.Replace(lastPartPaths[0], lastPartPaths[0], (lastPartPaths[0] + ":" + port), 1) // replace last ip number with ip + port
	finalURL := strings.Join(urlSplitted[0:len(urlSplitted)-1], `.`) + `.` + strings.Join(lastPartPaths, "/")  // join all parts
	println("FINAL URL:", finalURL)
	return finalURL
}

func getClient(ctx context.Context, addr string) (Client, error) {
	log.Println("getting Client Address:", addr)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Port() == "" {
		// we need to add a default port 80.
		addr = addPortAndParse(addr)
	}

	switch u.Scheme {
	case "http", "https":
		log.Println("Http Chosen")
		c, err := client.New(addr)
		if err != nil {
			return nil, err
		}
		return &clients.HTTPClient{Client: c}, nil
		// TODO support URI client from address.
	case "ws", "wss":
		log.Println("WS chosen")
		c, err := client.NewWS(addr, "")
		if err != nil {
			return nil, err
		}
		wsc := &clients.WebSocketClient{WSClient: c}
		wsc.Dialer, err = wsc.CreateDailer(addr)
		if err != nil {
			return nil, err
		}
		return wsc, nil
	case "stdio":
		return nil, fmt.Errorf("unsupported scheme: %q", u.Scheme)
	case "":
		return nil, fmt.Errorf("unsupported scheme: %q", u.Scheme)
	default:
		return nil, fmt.Errorf("no known transport for URL scheme %q", u.Scheme)
	}

}

func NewConnector(ctx context.Context, nConns uint, addr string) *Connector {
	log.Println("Creating New Connector")
	connector := &Connector{
		freeClients: make([]Client, 0, nConns),
	}

	for i := uint(0); i < nConns; i++ {
		var rpcClient Client
		var err error
		for {
			if ctx.Err() != nil {
				connector.Close()
				return nil
			}
			log.Println("A new client")
			nctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			rpcClient, err = getClient(nctx, addr)
			log.Println("client created: ", rpcClient.ClientType())
			if err != nil {
				log.Println("retrying", err)
				cancel()
				continue
			}
			cancel()
			break
		}
		connector.freeClients = append(connector.freeClients, rpcClient)
	}

	go connector.connectorLoop(ctx)
	return connector
}

func (connector *Connector) connectorLoop(ctx context.Context) {
	<-ctx.Done()
	log.Println("connectorLoop ctx.Done")
	connector.Close()
}

func (connector *Connector) Close() {
	for {
		connector.lock.Lock()
		log.Println("Connector closing", len(connector.freeClients))
		for i := 0; i < len(connector.freeClients); i++ {
			connector.freeClients[i].Close()
		}
		connector.freeClients = make([]Client, 0) // removing all clients.

		if connector.usedClients > 0 {
			log.Println("Connector closing, waiting for in use clients", connector.usedClients)
			connector.lock.Unlock()
			time.Sleep(100 * time.Millisecond)
		} else {
			connector.lock.Unlock()
			break
		}
	}
}

func (connector *Connector) GetRpc(block bool) (Client, error) {
	log.Println("GetRpc ")
	connector.lock.Lock()
	defer connector.lock.Unlock()
	countPrint := 0

	if len(connector.freeClients) == 0 {
		if !block {
			return nil, errors.New("out of clients")
		} else {
			for {
				if countPrint < 3 {
					countPrint++
				}
				connector.lock.Unlock()
				time.Sleep(50 * time.Millisecond)
				connector.lock.Lock()
				if len(connector.freeClients) != 0 {
					break
				}
			}
		}
	}

	ret := connector.freeClients[len(connector.freeClients)-1]
	connector.freeClients = connector.freeClients[:len(connector.freeClients)-1]
	connector.usedClients++

	return ret, nil
}

func (connector *Connector) ReturnRpc(rpc Client) {
	connector.lock.Lock()
	defer connector.lock.Unlock()

	connector.usedClients--
	connector.freeClients = append(connector.freeClients, rpc)
}
