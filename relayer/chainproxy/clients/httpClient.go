package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
	log.Println("params:", params, "method:", method)
	switch p := params.(type) {
	case []interface{}:
		log.Println("got http interface list:", p)
		var paramsFinal = make(map[string]interface{}, len(p))
		for idx, v := range p {
			paramsFinal[strconv.Itoa(idx)] = v
		}
	case map[string]interface{}:
		log.Println("got http map:", p)
		paramsFinal = p
	default:
		return nil, fmt.Errorf("unknown type %v", p)
	}
	log.Println("call started")
	ret, err := h.Client.Call(ctx, method, paramsFinal, result)
	log.Println(fmt.Sprintf("res: %s, err: %s", ret, err))
	return ret, err
}
