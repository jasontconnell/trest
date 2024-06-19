package process

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jasontconnell/trest/data"
)

var varreg *regexp.Regexp = regexp.MustCompile("[$]([a-zA-Z0-9]+)")

func Run(g *data.Group) []data.Result {
	first := g.Groups[0]
	results := runGroup(first, g, nil)
	return results
}

func runGroup(g *data.Group, root *data.Group, res data.Response) []data.Result {
	var results []data.Result
	if g.Request.Url != "" {
		requrl, err := getUrl(g, root, res)
		method := getMethod(g, root)

		if err != nil {
			log.Fatal(fmt.Errorf("can't get url %s. %w", g.Request.Url, err))
		}

		resp, result := getResponse(g, requrl, method, g.Request.Body, getHeaders(g, root), res)
		results = append(results, result)
		g.Responses = data.GetResponseValues(resp, g.RootElement, g.Variables, res)
	}

	var wg sync.WaitGroup

	for _, c := range g.Groups {
		wg.Add(len(g.Responses))
		for _, r := range g.Responses {
			go func(g1 *data.Group, r1 data.Response) {
				rsub := runGroup(g1, root, r1)
				results = append(results, rsub...)
				wg.Done()
			}(c, r)
		}
	}

	wg.Wait()

	return results
}

func getResponse(s *data.Group, requrl, method, body string, headers []data.Variable, res data.Response) (data.RawResponse, data.Result) {
	result := data.Result{}

	m := make(map[string]string)
	for k, v := range res {
		m[k] = v
	}

	h := make(map[string]string)
	for _, v := range headers {
		h[v.Name] = v.Property
	}

	for k, v := range m {
		requrl = strings.Replace(requrl, "$"+k, v, -1)
		body = strings.Replace(body, "$"+k, v, -1)
	}

	requrl = applyResponseValues(s, res, requrl)
	body = applyResponseValues(s, res, body)

	log.Println(requrl, body)

	buf := bytes.NewBufferString(body)
	r, err := http.NewRequest(method, requrl, buf)
	if err != nil {
		result.Err = fmt.Errorf("couldn't build request. %s %s %s %w", requrl, method, body, err)
		return data.DefaultResponse, result
	}

	for k, v := range h {
		r.Header.Add(k, v)
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		result.Err = fmt.Errorf("failure in request. %s %s %s %w", requrl, method, body, err)
		return data.DefaultResponse, result
	}
	defer resp.Body.Close()

	result.Status = resp.StatusCode
	result.Body = body
	result.Url = requrl

	dresp, err := data.ParseResponse(resp, res)
	if err != nil {
		result.Err = fmt.Errorf("parsing response. %s %s %s %w", requrl, method, body, err)
	}
	result.Duration = time.Since(start)

	return dresp, result
}

func applyResponseValues(s *data.Group, res data.Response, val string) string {
	left := varreg.FindAllStringSubmatch(val, -1)
	if len(left) > 0 {
		for _, mg := range left {
			v := mg[1]
			wd := mg[0]
			varval := getVariable(s, res, v)
			val = strings.ReplaceAll(val, wd, varval)
		}
	}
	return val
}

func getUrl(s *data.Group, root *data.Group, res data.Response) (string, error) {
	u := s.Request.Url
	if u == "" {
		return "", fmt.Errorf("blank url %s", s.Name)
	}
	for k, v := range res {
		u = strings.ReplaceAll(u, "$"+k, v)
	}

	uri, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("parsing url from s. %s %w", u, err)
	}
	if uri.IsAbs() {
		return uri.String(), nil
	}

	abs, err := url.Parse(root.RootUrl + u)
	if err != nil {
		return "", fmt.Errorf("parsing url from root + s. %s %w", root.RootUrl+u, err)
	}

	return abs.String(), nil
}

func getMethod(s *data.Group, root *data.Group) string {
	method := s.Request.Method
	if method == "" {
		method = root.Method
	}
	return method
}

func getHeaders(s *data.Group, root *data.Group) []data.Variable {
	headers := s.Headers
	if headers == nil || len(headers) == 0 {
		headers = root.Headers
	}
	return headers
}

func getVariable(s *data.Group, resp data.Response, name string) string {
	if v, ok := resp[name]; ok {
		return v
	}

	val := ""
	cur := s
	found := false
	for cur != nil && !found {
		for k, v := range cur.Response.RequestData {
			if k == name {
				val = v
				found = true
				break
			}
		}
		cur = cur.Parent
	}
	return val
}

func print(line string, r data.Response) {
	for k, v := range r {
		line = strings.ReplaceAll(line, "$"+k, v)
	}
	log.Println(line)
}
