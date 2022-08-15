package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type HTTPClient struct {
	*client.Client
}

func (h *HTTPClient) Close() { // nothing to do in http
}

func (h *HTTPClient) GetResponse() (json.RawMessage, error) { // shouldnt be called in HTTP
	return nil, nil
}

func (h *HTTPClient) ClientType() string {
	return "HTTPClient"
}

func (h *HTTPClient) Call(ctx context.Context, result *json.RawMessage, method string, params interface{}) (interface{}, error) {
	var paramsFinal map[string]interface{}
	switch p := params.(type) {
	case []interface{}:
		log.Println("got http interface list:", p)
		paramsFinal = map[string]interface{}{
			"arg": p, // couldnt find other way to trigger []interface{} params
		}
	case map[string]interface{}:
		log.Println("got http map:", p)
		paramsFinal = p
	default:
		return nil, fmt.Errorf("unknown type %v", p)
	}
	return h.Client.Call(ctx, method, paramsFinal, result)
}
