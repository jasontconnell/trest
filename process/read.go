package process

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jasontconnell/trest/data"
)

func ReadTests(filename string) (*data.Group, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file %s. %w", filename, err)
	}
	defer f.Close()

	lines := []string{}
	scn := bufio.NewScanner(f)
	for scn.Scan() {
		lines = append(lines, scn.Text())
	}

	g, err := parseFile(lines)
	return g, err
}

var spaces = regexp.MustCompile("^( +)(.+)$")

func parseFile(lines []string) (*data.Group, error) {
	g := &data.Group{
		Groups: []*data.Group{},
	}
	depth := 0

	cur := g
	headers := false

	for _, line := range lines {
		m := spaces.FindStringSubmatch(line)
		var instr string
		if len(m) == 3 {
			depth = len(m[1])
			instr = strings.TrimLeft(m[0], " ")
		} else {
			depth = 0
			instr = line
		}

		if len(instr) == 0 {
			continue
		}

		kw := strings.Split(instr, " ")
		dir := data.Directive(kw[0])
		rest := strings.Join(kw[1:], " ")

		switch dir {
		case data.Method:
			cur.Method = parseMethod(rest)
		case data.Url:
			cur.Request = parseUrl(rest)
		case data.Print:
			cur.Print = parsePrint(rest)
		case data.Test:
			var n *data.Group
			if cur.Depth == depth && cur.Parent != nil {
				cur = cur.Parent
			}

			n = &data.Group{
				Name:        parseName(rest),
				RootElement: parseElement(rest),
				Depth:       depth,
			}
			cur.Groups = append(cur.Groups, n)
			n.Parent = cur
			cur = n
			headers = false
		case data.Headers:
			headers = true
		case data.Variables:
			headers = false
		case data.Var:
			v := parseVariable(rest)
			if !headers {
				cur.Variables = append(cur.Variables, &v)
			} else {
				cur.Headers = append(cur.Headers, v)
			}
		}
	}

	var root *data.Group
	if len(g.Groups) == 1 {
		root = g.Groups[0]
	}

	return root, nil
}
