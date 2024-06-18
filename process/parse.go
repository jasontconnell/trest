package process

import (
	"strings"

	"github.com/jasontconnell/trest/data"
)

func parseMethod(line string) string {
	return line
}

func parseUrl(line string) data.Request {
	parts := strings.Split(line, " ")
	url := parts[0]
	body := strings.Join(parts[1:], " ")
	return data.Request{
		Url:  url,
		Body: body,
	}
}

func parseName(line string) string {
	return line
}

func parseVariable(line string) data.Variable {
	parts := strings.Split(line, " ")
	name := parts[0]
	val := parts[1]
	return data.Variable{
		Name:  name,
		Value: val,
	}
}
