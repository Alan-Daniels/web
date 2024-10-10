package database

import (
	"encoding/json"
	"fmt"
)

type Map map[string]interface{}

type RawResponse Map

func toMapArr(data interface{}) []Map {
	out := make([]Map, 0)

	dataAsArray := data.([]interface{})
	for _, dta := range dataAsArray {
		out = append(out, toMap(dta))
	}

	return out
}

func toMap(data interface{}) Map {
	out := make(Map, 0)

	dataAsMap := data.(map[interface{}]interface{})
	for k, v := range dataAsMap {
		key := k.(string)
		val := v
		switch v.(type) {
		case map[interface{}]interface{}:
			val = toMap(v)
			break
		case []interface{}:
			val = toMapArr(v)
			break
		default:
		}

		out[key] = val
	}

	return out
}

func toRawResponses(data interface{}, err error) ([]RawResponse, error) {
	if err != nil {
		return nil, err
	}
	responses := make([]RawResponse, 0)

	dataAsArray := data.([]interface{})
	for _, dta := range dataAsArray {
		responses = append(responses, (RawResponse)(toMap(dta)))
	}

	return responses, nil
}

func Unmarshal[T any](response any, err error) (*T, error) {
	Resp := new(T)

	if err != nil {
		return nil, err
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(response)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialise response '%+v' to object: %w", response, err)
	}

	err = json.Unmarshal(jsonBytes, Resp)

	if err != nil {
		return nil, err
	}

	return Resp, nil
}

func UnmarshalResponse[T any](response RawResponse, err error) (*T, error) {
	return Unmarshal[T](response["result"], err)
}
