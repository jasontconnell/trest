package data

import (
	"fmt"
	"strings"
)

type Group struct {
	Name      string
	Method    string
	Request   Request
	RootUrl   string
	Groups    []*Group
	Variables []*Variable
	Headers   []Variable
	Parent    *Group
	Repeat    bool
	Depth     int
	Print     string
	Response  RawResponse
	Responses []Response
}

func (g Group) String() string {
	space := strings.Repeat(" ", g.Depth)
	s := fmt.Sprintf("%[4]sName: %[1]s\n%[4]sMethod: %[2]s\n%[4]sRequest: %[3]v\n", g.Name, g.Method, g.Request, space)
	if g.Repeat {
		s += space + fmt.Sprintf("Repeats\n")
	}

	for _, gs := range g.Groups {
		s += fmt.Sprintf("%v", gs)
	}
	return s
}

type Variable struct {
	Name     string
	Property string
	Type     string
}

func (v Variable) String() string {
	return fmt.Sprintf("Name: %s Value: %s Type: %s", v.Name, v.Property, v.Type)
}

type Request struct {
	Method string
	Url    string
	Body   string
}

func (req Request) String() string {
	return fmt.Sprintf("Method: %s Url: %s Body: %s", req.Method, req.Url, req.Body)
}
