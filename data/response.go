package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RawResponse struct {
	Data        map[string]interface{}
	RequestData Response
}

type Response map[string]string

var DefaultResponse RawResponse = RawResponse{}

func ParseResponse(resp *http.Response, presp Response) (RawResponse, error) {
	var raw map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&raw)
	return RawResponse{Data: raw, RequestData: presp}, err
}

func GetResponseValues(resp RawResponse, root string, vars []*Variable, res Response) []Response {
	var list []Response
	for k, v := range resp.Data {
		if k != root {
			continue
		}
		switch val := v.(type) {
		case []interface{}:
			list = extractArray(val, vars, res)
		}
	}
	return list
}

func extractArray(ary []interface{}, vars []*Variable, res Response) []Response {
	resp := []Response{}
	for _, obj := range ary {
		switch val := obj.(type) {
		case map[string]interface{}:
			r := Response{}
			for _, v := range vars {
				r[v.Name] = extractValue(val, v.Property)
			}
			mergeResponses(r, res)
			resp = append(resp, r)
		case float64:
			r := Response{}
			for _, v := range vars {
				if v.Property == "@index" {
					r[v.Name] = fmt.Sprintf("%d", int(val))
				}
			}
			mergeResponses(r, res)
			resp = append(resp, r)
		case int:
			r := Response{}
			for _, v := range vars {
				if v.Property == "@index" {
					r[v.Name] = fmt.Sprintf("%d", val)
				}
			}
			mergeResponses(r, res)
			resp = append(resp, r)
		default:
			log.Printf("unsupported type %T\n", val)
		}
	}

	return resp
}

func mergeResponses(r Response, res Response) {
	for k, v := range res {
		r[k] = v
	}
}

func extractValue(m map[string]interface{}, prop string) string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}

	if prop == "@index" {
		log.Println("getting @index", m)
	}

	var value string
	parts := strings.Split(prop, ".")
	for i, p := range parts {
		if val, ok := m[p]; ok {
			switch tval := val.(type) {
			case map[string]interface{}:
				value = extractValue(tval, strings.Join(parts[i+1:], "."))
			case string:
				value = tval
			case int:
				value = fmt.Sprintf("%d", tval)
			case float64:
				value = fmt.Sprintf("%d", int(tval))
			case []interface{}:
				log.Println("array", val)
			}
		}
	}
	return value
}
