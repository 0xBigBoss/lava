package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	spectypes "github.com/lavanet/lava/x/spec/types"
)

const (
	PARSE_PARAMS = 0
	PARSE_RESULT = 1
)

type RPCInput interface {
	GetParams() []interface{}
	GetResult() json.RawMessage
	ParseBlock(block string) (int64, error)
}

func ParseDefaultBlockParameter(block string) (int64, error) {
	switch block {
	case "latest":
		return spectypes.LATEST_BLOCK, nil
	case "earliest":
		return spectypes.EARLIEST_BLOCK, nil
	case "pending":
		return spectypes.PENDING_BLOCK, nil
	default:
		//try to parse a number
	}
	blockNum, err := strconv.ParseInt(block, 0, 64)
	if err != nil {
		return spectypes.NOT_APPLICABLE, fmt.Errorf("invalid block value, could not parse block %s, error: %s", block, err)
	}
	if blockNum < 0 {
		return spectypes.NOT_APPLICABLE, fmt.Errorf("invalid block value, block value was negative %d", blockNum)
	}
	return blockNum, nil
}

//this function returns the block that was requested,
func Parse(rpcInput RPCInput, blockParser spectypes.BlockParser, dataSource int) ([]interface{}, error) {
	var retval []interface{}
	var err error

	switch blockParser.ParserFunc {
	case spectypes.PARSER_FUNC_EMPTY:
		return nil, nil
	case spectypes.PARSER_FUNC_PARSE_BY_ARG:
		retval, err = ParseByArg(rpcInput, blockParser.ParserArg, dataSource)
	case spectypes.PARSER_FUNC_PARSE_CANONICAL:
		retval, err = ParseCanonical(rpcInput, blockParser.ParserArg, dataSource)
	case spectypes.PARSER_FUNC_PARSE_DICTIONARY:
		retval, err = ParseDictionary(rpcInput, blockParser.ParserArg, dataSource)
	case spectypes.PARSER_FUNC_PARSE_DICTIONARY_OR_ORDERED:
		retval, err = ParseDictionaryOrOrdered(rpcInput, blockParser.ParserArg, dataSource)
	case spectypes.PARSER_FUNC_DEFAULT:
		retval = ParseDefault(rpcInput, blockParser.ParserArg, dataSource)
		// case default -> retval, err = ParseByArg(rpcInput, blockParser.ParserArg, dataSource) with default.
	default:
		return nil, fmt.Errorf("unsupported block parser parserFunc")
	}

	if err != nil {
		return nil, err
	}

	return retval, nil

}

func ParseDefault(rpcInput RPCInput, input []string, dataSource int) []interface{} {
	retArr := make([]interface{}, 0)
	retArr = append(retArr, fmt.Sprintf("%s", input[0]))
	return retArr
}

//this function returns the block that was requested,
func ParseBlockFromParams(rpcInput RPCInput, blockParser spectypes.BlockParser) (int64, error) {
	result, err := Parse(rpcInput, blockParser, PARSE_PARAMS)
	if err != nil || result == nil {
		return spectypes.NOT_APPLICABLE, err
	}
	return rpcInput.ParseBlock(result[0].(string))
}

//this function returns the block that was requested,
func ParseBlockFromReply(rpcInput RPCInput, blockParser spectypes.BlockParser) (int64, error) {
	result, err := Parse(rpcInput, blockParser, PARSE_RESULT)
	if err != nil || result == nil {
		return spectypes.NOT_APPLICABLE, err
	}

	blockstr, ok := result[0].(string)
	if !ok {
		return spectypes.NOT_APPLICABLE, errors.New("block number is not string parseable")
	}

	if strings.Contains(blockstr, "\"") {
		blockstr, err = strconv.Unquote(blockstr)
		if err != nil {
			return spectypes.NOT_APPLICABLE, err
		}

	}
	return rpcInput.ParseBlock(blockstr)
}

//this function returns the block that was requested,
func ParseMessageResponse(rpcInput RPCInput, resultParser spectypes.BlockParser) ([]interface{}, error) {
	return Parse(rpcInput, resultParser, PARSE_RESULT)
}

// Move to RPCInput
func GetDataToParse(rpcInput RPCInput, dataSource int) ([]interface{}, error) {
	switch dataSource {
	case PARSE_PARAMS:
		return rpcInput.GetParams(), nil
	case PARSE_RESULT:
		log.Println("Parse Result")
		interfaceArr := []interface{}{}
		var data map[string]interface{}
		unmarshalled := rpcInput.GetResult()
		log.Println("unmarshaled:", string(unmarshalled))

		// Try to unmarshal and if the data is unmarshalable then return the data itself
		err := json.Unmarshal(unmarshalled, &data)
		if err != nil {
			interfaceArr = append(interfaceArr, unmarshalled)
		} else {
			interfaceArr = append(interfaceArr, data)
		}

		return interfaceArr, nil
	default:
		return nil, fmt.Errorf("unsupported block parser parserFunc")
	}
}

func ParseByArg(rpcInput RPCInput, input []string, dataSource int) ([]interface{}, error) {
	//specified block is one of the direct parameters, input should be one string defining the location of the block
	if len(input) != 1 {
		return nil, fmt.Errorf("invalid input format, input length: %d", len(input))
	}
	inp := input[0]
	param_index, err := strconv.ParseUint(inp, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, input isn't an unsigned index: %s, error: %s", inp, err)
	}

	unmarshalledData, err := GetDataToParse(rpcInput, dataSource)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, data is not json: %s, error: %s", unmarshalledData, err)
	}

	if uint64(len(unmarshalledData)) < param_index {
		return nil, fmt.Errorf("invalid rpc input and input index: wanted param: %d params: %s", param_index, unmarshalledData)
	}
	block := unmarshalledData[param_index]
	//TODO: turn this into type assertion instead

	retArr := make([]interface{}, 0)
	retArr = append(retArr, fmt.Sprintf("%s", block))
	return retArr, nil
}

//expect input to be keys[a,b,c] and a canonical object such as
// {
//   "a": {
//       "b": {
//          "c": "wanted result"
//        }
//    }
// }
// should output an interface array with "wanted result" in first index 0
func ParseCanonical(rpcInput RPCInput, input []string, dataSource int) ([]interface{}, error) {
	inp := input[0]
	param_index, err := strconv.ParseUint(inp, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, input isn't an unsigned index: %s, error: %s", inp, err)
	}

	unmarshalledData, err := GetDataToParse(rpcInput, dataSource)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, data is not json: %s, error: %s", unmarshalledData, err)
	}

	if uint64(len(unmarshalledData)) < param_index {
		return nil, fmt.Errorf("invalid rpc input and input index: wanted param: %d params: %s", param_index, unmarshalledData)
	}

	blockContainer := unmarshalledData[param_index]
	for _, key := range input[1:] {
		// type assertion for blockcontainer
		if blockContainer, ok := blockContainer.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("invalid parser input format, blockContainer is %v and not map[string]interface{} and tried to get a field inside: %s", blockContainer, key)
		}

		// assertion for key
		if container, ok := blockContainer.(map[string]interface{})[key]; ok {
			blockContainer = container
		} else {
			return nil, fmt.Errorf("invalid input format, blockContainer %s does not have field inside: %s", blockContainer, key)
		}

	}

	retArr := make([]interface{}, 0)
	retArr = append(retArr, fmt.Sprintf("%s", blockContainer))
	return retArr, nil
}

func ParseDictionary(rpcInput RPCInput, input []string, dataSource int) ([]interface{}, error) {
	if len(input) != 2 {
		return nil, fmt.Errorf("invalid input format, input length: %d and needs to be 2", len(input))
	}
	prop_name := input[0]
	inner_separator := input[1]

	unmarshalledData, err := GetDataToParse(rpcInput, dataSource)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, data is not json: %s, error: %s", unmarshalledData, err)
	}

	for _, val := range unmarshalledData {
		if prop, ok := val.(string); ok {
			splitted := strings.SplitN(prop, inner_separator, 2)
			if splitted[0] != prop_name {
				continue
			} else {
				retArr := make([]interface{}, 0)
				retArr = append(retArr, splitted[1])
				return retArr, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid input format, did not find prop name %s on params: %s", prop_name, unmarshalledData)
}

func ParseDictionaryOrOrdered(rpcInput RPCInput, input []string, dataSource int) ([]interface{}, error) {
	if len(input) != 3 {
		return nil, fmt.Errorf("invalid input format, input length: %d and needs to be 3", len(input))
	}
	prop_name := input[0]
	inner_separator := input[1]
	inp := input[2]
	param_index, err := strconv.ParseUint(inp, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, input isn't an unsigned index: %s, error: %s", inp, err)
	}

	unmarshalledData, err := GetDataToParse(rpcInput, dataSource)
	if err != nil {
		return nil, fmt.Errorf("invalid input format, data is not json: %s, error: %s", unmarshalledData, err)
	}

	for _, val := range unmarshalledData {
		if prop, ok := val.(string); ok {
			splitted := strings.SplitN(prop, inner_separator, 2)
			if splitted[0] != prop_name || len(splitted) < 2 {
				continue
			} else {
				retArr := make([]interface{}, 0)
				retArr = append(retArr, splitted[1])
				return retArr, nil
			}
		}
	}
	//did not find a named property
	if uint64(len(unmarshalledData)) < param_index {
		return nil, fmt.Errorf("invalid rpc input and input index: wanted param idx: %d params: %s", param_index, unmarshalledData)
	}
	block := unmarshalledData[param_index]
	//TODO: turn this into type assertion instead

	retArr := make([]interface{}, 0)
	retArr = append(retArr, fmt.Sprintf("%s", block))
	return retArr, nil
}
