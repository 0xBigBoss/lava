package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type WebSocketClient struct {
	*client.WSClient
}

func (c *WebSocketClient) Close() {
	c.Stop() // maybe other way to close the connection
}

func (c *WebSocketClient) CreateDailer(remoteAddr string) (func(string, string) (net.Conn, error), error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		u.Scheme = "tcp"
	}
	dialFn := func(proto, addr string) (net.Conn, error) {
		var timeout = 10 * time.Second
		return net.DialTimeout(proto, u.Host+u.EscapedPath(), timeout)
	}
	return dialFn, nil
}

func (c *WebSocketClient) ClientType() string {
	return "WebSocketClient"
}

func (c *WebSocketClient) GetResponse() (json.RawMessage, error) { // shouldnt be called in HTTP
	msg := <-c.WSClient.ResponsesCh
	if msg.Error != nil {
		return nil, fmt.Errorf("Code: %d, Message: %s, Data: %s", msg.Error.Code, msg.Error.Message, msg.Error.Data)
	}
	return msg.Result, nil
}

func (c *WebSocketClient) Call(ctx context.Context, result *json.RawMessage, method string, params interface{}) (interface{}, error) {
	log.Println("params:", params, "method:", method)
	log.Println("endpoint", c.Endpoint)
	log.Println("address", c.Address)
	// err := c.WSClient.Start()
	// if err != nil {
	// 	return nil, err
	// }
	// defer c.Stop()
	var res error
	switch p := params.(type) {
	case []interface{}:
		log.Println("got http interface list:", p)
		res = c.WSClient.CallWithArrayParams(ctx, method, p)
	case map[string]interface{}:
		log.Println("got http map:", p)
		res = c.WSClient.Call(ctx, method, p)
	default:
		return nil, fmt.Errorf("unknown type %v", p)
	}

	return nil, res
}
