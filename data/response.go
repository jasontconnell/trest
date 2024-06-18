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

func GetResponseValues(resp RawResponse, root string, vars []*Variable) []Response {
	var list []Response
	fmt.Println(len(resp.Data))
	for k, v := range resp.Data {
		fmt.Println(k, root)
		fmt.Printf("%T\n", v)
		if k != root {
			continue
		}
		switch val := v.(type) {
		case []interface{}:
			list = extractArray(val, vars)
		default:
			log.Printf("unhandled %T\n", val)
		}
	}
	return list
}

func extractArray(ary []interface{}, vars []*Variable) []Response {
	resp := []Response{}
	for _, obj := range ary {
		switch val := obj.(type) {
		case map[string]interface{}:
			r := Response{}
			for _, v := range vars {
				r[v.Name] = extractValue(val, v.Property)
			}
			resp = append(resp, r)
		case float64:
			r := Response{}
			log.Println(val)
			for _, v := range vars {
				if v.Name == "@index" {
					r[v.Name] = fmt.Sprintf("%d", int(val))
				}
			}
			resp = append(resp, r)
		case int:
			r := Response{}
			for _, v := range vars {
				if v.Name == "@index" {
					r[v.Name] = fmt.Sprintf("%d", val)
				}
			}
			resp = append(resp, r)
		default:
			log.Printf("unsupported type %T\n", val)
		}
	}
	return resp
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
