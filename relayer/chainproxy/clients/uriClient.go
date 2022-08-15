package clients

import (
	"context"
	"encoding/json"

	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type URIClient struct {
	*client.URIClient
}

func (h *URIClient) Close() { // nothing to do in http
}

func (h *URIClient) ClientType() string {
	return "URIClient"
}

func (h *URIClient) GetResponse() (json.RawMessage, error) { // shouldnt be called in URI
	return nil, nil
}

func (h *URIClient) Call(ctx context.Context, result *json.RawMessage, method string, params interface{}) (interface{}, error) {
	// requestBytes, err := json.Marshal(request)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal request: %w", err)
	// }

	// requestBuf := bytes.NewBuffer(requestBytes)
	// httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.address, requestBuf)
	// if err != nil {
	// 	return nil, fmt.Errorf("request setup failed: %w", err)
	// }

	// httpRequest.Header.Set("Content-Type", "application/json")

	// if c.username != "" || c.password != "" {
	// 	httpRequest.SetBasicAuth(c.username, c.password)
	// }

	// httpResponse, err := c.client.Do(httpRequest)
	// if err != nil {
	// 	return nil, err
	// }

	// defer httpResponse.Body.Close()

	// responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read response body: %w", err)
	// }

	// return unmarshalResponseBytes(responseBytes, id, result)
	return nil, nil
}
