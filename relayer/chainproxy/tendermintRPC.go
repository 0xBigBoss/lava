package chainproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/lavanet/lava/relayer/parser"
	"github.com/lavanet/lava/relayer/sentry"
	pairingtypes "github.com/lavanet/lava/x/pairing/types"
	spectypes "github.com/lavanet/lava/x/spec/types"
)

type TendemintRpcMessage struct {
	JrpcMessage
	cp *tendermintRpcChainProxy
}

type tendermintRpcChainProxy struct {
	//embedding the jrpc chain proxy because the only diff is on parse message
	JrpcChainProxy
}

func (m TendemintRpcMessage) GetParams() []interface{} {
	return m.msg.Params
}

func (m TendemintRpcMessage) GetResult() json.RawMessage {
	return m.msg.Result
}

func (m TendemintRpcMessage) ParseBlock(inp string) (int64, error) {
	return parser.ParseDefaultBlockParameter(inp)
}

func (cp *tendermintRpcChainProxy) FetchLatestBlockNum(ctx context.Context) (int64, error) {
	serviceApi, ok := cp.GetSentry().GetSpecApiByTag(spectypes.GET_BLOCKNUM)
	if !ok {
		return spectypes.NOT_APPLICABLE, errors.New(spectypes.GET_BLOCKNUM + " tag function not found")
	}

	params := []interface{}{}
	nodeMsg, err := cp.newMessage(&serviceApi, serviceApi.GetName(), spectypes.LATEST_BLOCK, params)
	if err != nil {
		return spectypes.NOT_APPLICABLE, err
	}

	_, err = nodeMsg.Send(ctx)
	if err != nil {
		return spectypes.NOT_APPLICABLE, err
	}

	blocknum, err := parser.ParseBlockFromReply(nodeMsg.GetMsg().(*JsonrpcMessage), serviceApi.Parsing.ResultParsing)
	if err != nil {
		return spectypes.NOT_APPLICABLE, err
	}

	return blocknum, nil
}

func (cp *tendermintRpcChainProxy) FetchBlockHashByNum(ctx context.Context, blockNum int64) (string, error) {
	serviceApi, ok := cp.GetSentry().GetSpecApiByTag(spectypes.GET_BLOCK_BY_NUM)
	if !ok {
		return "", errors.New(spectypes.GET_BLOCK_BY_NUM + " tag function not found")
	}

	var nodeMsg NodeMessage
	var err error
	if serviceApi.GetParsing().FunctionTemplate != "" {
		nodeMsg, err = cp.ParseMsg("", []byte(fmt.Sprintf(serviceApi.Parsing.FunctionTemplate, blockNum)), "")
	} else {
		params := make([]interface{}, 0)
		params = append(params, blockNum)
		nodeMsg, err = cp.newMessage(&serviceApi, serviceApi.GetName(), spectypes.LATEST_BLOCK, params)
	}

	if err != nil {
		return "", err
	}

	_, err = nodeMsg.Send(ctx)
	if err != nil {
		return "", err
	}

	blockData, err := parser.ParseMessageResponse((nodeMsg.GetMsg().(*JsonrpcMessage)), serviceApi.Parsing.ResultParsing)
	if err != nil {
		return "", err
	}

	// blockData is an interface array with the parsed result in index 0.
	// we know to expect a string result for a hash.
	hash, ok := blockData[spectypes.DEFAULT_PARSED_RESULT_INDEX].(string)
	if !ok {
		return "", errors.New("hash not string parseable")
	}

	return hash, nil
}

func NewtendermintRpcChainProxy(nodeUrl string, nConns uint, sentry *sentry.Sentry) ChainProxy {
	return &tendermintRpcChainProxy{
		JrpcChainProxy: JrpcChainProxy{
			nodeUrl: nodeUrl,
			nConns:  nConns,
			sentry:  sentry,
		},
	}
}

func (cp *tendermintRpcChainProxy) newMessage(serviceApi *spectypes.ServiceApi, method string, requestedBlock int64, params []interface{}) (*TendemintRpcMessage, error) {
	nodeMsg := &TendemintRpcMessage{
		JrpcMessage: JrpcMessage{serviceApi: serviceApi,
			msg: &JsonrpcMessage{
				Version: "2.0",
				ID:      []byte("1"), //TODO:: use ids
				Method:  method,
				Params:  params,
			},
			requestedBlock: requestedBlock},
		cp: cp,
	}
	return nodeMsg, nil
}

func (cp *tendermintRpcChainProxy) ParseMsg(path string, data []byte, connectionType string) (NodeMessage, error) {
	// connectionType is currently only used only in rest api
	// Unmarshal request
	var msg JsonrpcMessage
	if string(data) != "" {
		//assuming jsonrpc
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return nil, err
		}

	} else {
		//assuming URI
		var parsedMethod string
		idx := strings.Index(path, "?")
		if idx == -1 {
			parsedMethod = path
		} else {
			parsedMethod = path[0:idx]
		}

		msg = JsonrpcMessage{
			ID:      []byte("1"),
			Version: "2.0",
			Method:  parsedMethod,
		} //other parameters don't matter
		if strings.Contains(path[idx+1:], "=") {
			params_raw := strings.Split(path[idx+1:], "&") //list with structure ['height=0x500',...]
			msg.Params = make([]interface{}, len(params_raw))
			for i := range params_raw {
				msg.Params[i] = params_raw[i]
			}
		} else {
			msg.Params = make([]interface{}, 0)
		}
		//convert the list of strings to a list of interfaces
	}
	//
	// Check api is supported and save it in nodeMsg
	serviceApi, err := cp.getSupportedApi(msg.Method)
	if err != nil {
		return nil, err
	}

	requestedBlock, err := parser.ParseBlockFromParams(msg, serviceApi.BlockParsing)
	if err != nil {
		return nil, err
	}

	nodeMsg := &TendemintRpcMessage{
		JrpcMessage: JrpcMessage{serviceApi: serviceApi,
			msg: &msg, requestedBlock: requestedBlock},
		cp: cp,
	}
	return nodeMsg, nil
}

func (cp *tendermintRpcChainProxy) PortalStart(ctx context.Context, privKey *btcec.PrivateKey, listenAddr string) {
	//
	// Setup HTTP Server
	app := fiber.New(fiber.Config{})

	app.Use("/ws/:dappId", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:dappId", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				c.WriteMessage(mt, []byte("Error Received: "+err.Error()))
				break
			}
			log.Println("ws: in <<< ", string(msg))

			reply, err := SendRelay(ctx, cp, privKey, "", string(msg), "")
			if err != nil {
				log.Println(err)
				c.WriteMessage(mt, []byte("Error Received: "+err.Error()))
				break
			}

			if err = c.WriteMessage(mt, reply.Data); err != nil {
				log.Println("write:", err)
				c.WriteMessage(mt, []byte("Error Received: "+err.Error()))
				break
			}
			log.Println("out >>> ", string(reply.Data))
		}
	}))

	app.Post("/:dappId/*", func(c *fiber.Ctx) error {
		log.Println("jsonrpc in <<< ", string(c.Body()))
		reply, err := SendRelay(ctx, cp, privKey, "", string(c.Body()), "")
		if err != nil {
			log.Println(err)
			return c.SendString(fmt.Sprintf(`{"error": "unsupported api","more_information" %s}`, err))
		}

		log.Println("out >>> ", string(reply.Data))
		return c.SendString(string(reply.Data))
	})

	app.Get("/:dappId/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		log.Println("urirpc in <<< ", path)
		reply, err := SendRelay(ctx, cp, privKey, path, "", "")
		if err != nil {
			log.Println(err)
			if string(c.Body()) != "" {
				return c.SendString(fmt.Sprintf(`{"error": "unsupported api", "recommendation": "For jsonRPC use POST", "more_information": "%s"}`, err))
			}
			return c.SendString(fmt.Sprintf(`{"error": "unsupported api","more_information" %s}`, err))
		}
		log.Println("out >>> ", string(reply.Data))
		return c.SendString(string(reply.Data))
	})
	//
	// Go
	err := app.Listen(listenAddr)
	if err != nil {
		log.Println(err)
	}
}

func (nm *TendemintRpcMessage) Send(ctx context.Context) (*pairingtypes.RelayReply, error) {
	// Get node
	log.Println("Sending started")

	rpc, err := nm.cp.conn.GetRpc(true)
	if err != nil {
		return nil, err
	}
	defer nm.cp.conn.ReturnRpc(rpc)

	//
	// Call our node
	var result json.RawMessage
	var responseUnmarsheled interface{}
	connectCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	log.Println("Sending Call")
	responseUnmarsheled, err = rpc.Call(connectCtx, &result, nm.msg.Method, nm.msg.Params)

	if err != nil {
		// error from call
		nm.msg.Result = []byte(fmt.Sprintf("%s", err))
		return nil, err
	}

	var data []byte
	switch rpc.ClientType() {
	case "HTTPClient":
		log.Println("HTTPClient results")
		data, err = json.Marshal(responseUnmarsheled)
		if err != nil {
			// if marshaling the response interface failed, try to marshal the result directly
			data, err = json.Marshal(result)
		}
	case "WebSocketClient":
		log.Println("WebSocketClient results")
		data, err = rpc.GetResponse()
	case "URIClient":
		log.Println("URIClient results")
	}
	log.Println("res:", string(data))

	if err != nil {
		// error from parsing response
		nm.msg.Result = []byte(fmt.Sprintf("%s", err))
		return nil, err
	}
	nm.msg.Result = (json.RawMessage)(data)
	reply := &pairingtypes.RelayReply{
		Data: data,
	}
	return reply, nil
}
