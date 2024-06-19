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
	parts := strings.Fields(line)
	return parts[0]
}

func parseElement(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func parsePrint(line string) string {
	return line
}

func parseVariable(line string) data.Variable {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		panic("invalid variable format " + line)
	}
	name := parts[0]
	val := parts[1]
	tp := parts[2]
	return data.Variable{
		Name:     name,
		Property: val,
		Type:     tp,
	}
}
