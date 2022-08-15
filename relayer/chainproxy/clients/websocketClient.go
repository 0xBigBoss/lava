package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type WebSocketClient struct {
	*client.WSClient
}

func (c *WebSocketClient) Close() {
	c.Stop() // maybe other way to close the connection
}

func (c *WebSocketClient) ClientType() string {
	return "WebSocketClient"
}

func (h *WebSocketClient) GetResponse() (json.RawMessage, error) { // shouldnt be called in HTTP
	msg := <-h.WSClient.ResponsesCh
	if msg.Error != nil {
		return nil, fmt.Errorf("Code: %d, Message: %s, Data: %s", msg.Error.Code, msg.Error.Message, msg.Error.Data)
	}
	return msg.Result, nil
}

func (h *WebSocketClient) Call(ctx context.Context, result *json.RawMessage, method string, params interface{}) (interface{}, error) {
	switch p := params.(type) {
	case []interface{}:
		log.Println("got http interface list:", p)
		return nil, h.WSClient.CallWithArrayParams(ctx, method, p)
	case map[string]interface{}:
		log.Println("got http map:", p)
		return nil, h.WSClient.Call(ctx, method, p)
	default:
		return nil, fmt.Errorf("unknown type %v", p)
	}
}
