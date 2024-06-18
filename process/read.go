package process

import (
	"bufio"
	"fmt"
	"log"
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

	log.Println(len(lines))
	g, err := parseFile(lines)
	return g, nil
}

var spaces = regexp.MustCompile("^( +)(.+)$")

func parseFile(lines []string) (*data.Group, error) {
	g := &data.Group{
		Groups: []*data.Group{},
	}
	depth := 0

	cur := g

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
			log.Println("parse method", cur.Method)
		case data.Url:
			cur.Request = parseUrl(rest)
			log.Println("parse request", cur.Request)
		case data.Test, data.Each:
			repeat := dir == data.Each
			var n *data.Group
			if cur.Depth == depth {
				cur = cur.Parent
			}

			n = &data.Group{
				Name:   parseName(rest),
				Repeat: repeat,
				Depth:  depth,
			}
			log.Println("parse", dir, n.Name)

			cur.Groups = append(cur.Groups, n)
			n.Parent = cur
			cur = n
		case data.Var:
			v := parseVariable(rest)
			cur.Variables = append(cur.Variables, v)
			log.Println("parse variable", v)
		}
	}

	return g, nil
}
