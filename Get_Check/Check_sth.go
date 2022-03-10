package Check

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var Value string
var Result bool
var Url string
var Server *string
var Headers *string
var Bodys *string

func Get_req(url string, servers *string, headers *string, bodys *string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	*servers = r.Header.Get("Server")
	defer func() { _ = r.Body.Close() }()
	a, _ := ioutil.ReadAll(r.Body)
	*bodys = string(a)
	dataType, _ := json.Marshal(r.Header)
	*headers = string(dataType)
}

func Get_server(value string, servers string) (result bool) {
	result = false
	if servers == value {
		Result = true
	}
	return Result
}
func Get_title(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile("<title>(.*?)</title>")
	matchArr := compileRegex.FindStringSubmatch(bodys)
	if matchArr != nil {
		if matchArr[1] == value {
			Result = true
		}
	}
	return Result
}

func Get_body(value string, bodys string) (result bool) {
	Result = strings.Contains(bodys, value)
	return Result
}

func Get_banner(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile(`(?im)<\\s*banner.*>(.*?)<\\s*/\\s*banner>`)
	matchArr := compileRegex.FindAllStringSubmatch(bodys, -1)
	var i int
	for i = 0; i < len(matchArr); i++ {
		if value == matchArr[i][1] {
			Result = true
		}
	}
	return Result
}

func Get_port(value string, urls string) (result bool) {

	u, _ := url.Parse(urls)
	result = false
	host := u.Host
	address := net.ParseIP(host)
	if address == nil {
		Result = false
	} else {
		ho := strings.Split(host, ":")
		if len(ho) > 1 {
			port := ho[1]
			if port == value {
				Result = true
			}
		}
	}
	return result
}

func Get_head(value string, headers string) (result bool) {
	what_headers, _ := json.Marshal(headers)
	child_string := string(what_headers)
	what_headers = []byte(strings.Replace(child_string, "\"", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "\\", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "[", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "]", "", -1))
	Result = strings.Contains(string(what_headers), value)
	return Result
}

func Check_model(match string, value string, url string, bodys string, heades string, servers string) (result bool) {
	switch {
	case match == "body_contains":
		Result = Get_body(value, bodys)
	case match == "protocol_contains":
		Result = false
	case match == "title_contains":
		Result = Get_title(value, bodys)
	case match == "banner_contains":
		Result = Get_banner(value, bodys)
	case match == "header_contains":
		Result = Get_head(value, heades)
	case match == "port_contains":
		Result = Get_port(value, url)
	case match == "server":
		Result = Get_server(value, servers)
	case match == "title":
		Result = Get_title(value, bodys)
	case match == "cert_contains":
		result = false
	case match == "server_contains":
		Result = Get_server(value, servers)
	case match == "protocol":
		Result = false
	default:
		Result = false
	}
	return Result
}
