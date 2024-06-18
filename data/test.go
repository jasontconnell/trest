package data

import (
	"fmt"
	"strings"
)

type Group struct {
	Name      string
	Method    string
	Request   Request
	Groups    []*Group
	Variables []Variable
	Parent    *Group
	Repeat    bool
	Depth     int
}

func (g Group) String() string {
	space := strings.Repeat(" ", g.Depth)
	s := fmt.Sprintf("%[3]sName: %s\n%[3]sMethod: %s\n%[3]sRequest: %v\n", g.Name, g.Method, g.Request, space)
	if g.Repeat {
		s += space + fmt.Sprintf("Repeats\n")
	}

	for _, gs := range g.Groups {
		s += fmt.Sprintf("%v", gs)
	}
	return s
}

type Variable struct {
	Name  string
	Value string
}

func (v Variable) String() string {
	return fmt.Sprintf("Name: %s Value: %s", v.Name, v.Value)
}

type Request struct {
	Url  string
	Body string
}

func (req Request) String() string {
	return fmt.Sprintf("Url: %s Body: %s", req.Url, req.Body)
}
