package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Rules [][]struct {
	Match   string `json:"match"`
	Content string `json:"content"`
}

type FOFA struct {
	RuleID         string `json:"rule_id"`
	Level          string `json:"level"`
	Softhard       string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Rules          Rules  `json:"rules"`
}

type ArgsInfo struct {
	Host    string
	Hosts   string
	CmsJson string
}

var URL = flag.String("url", "", "输入url")
var Urllist = flag.String("urls", "", "输入urls的文件名")
var wg sync.WaitGroup

func main() {

	/*var Info ArgsInfo
	Info.Flag()
	jsonConfigList := getConfigs()
	unmarshelledConfigs := deserializeJson(jsonConfigList)
	for _, configObj := range unmarshelledConfigs {
		fmt.Printf("Product: %v Rules: %v", configObj.Product, configObj.Rules)
	}*/
	Banner()
	flag.Parse()

	var bar Bar
	if *URL != "" {

		bar.NewOption(0, 100)
		//bar.NewOptionWithGraph(0, 100, "#")
		for i := 0; i <= 100; i++ {
			time.Sleep(100 * time.Millisecond)
			bar.Play(int64(i))
		}
		bar.Finish()
		wg.Add(1)

		run(*URL)
		wg.Wait()

	}

	if *Urllist != "" {

		bar.NewOption(0, 100)
		//bar.NewOptionWithGraph(0, 100, "#")
		for i := 0; i <= 100; i++ {
			time.Sleep(100 * time.Millisecond)
			bar.Play(int64(i))
		}
		bar.Finish()

		file, err := os.Open(*Urllist)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		wg.Add(1)

		for _, line := range lines {

			go run(line)

		}

	}
	if *URL == "" && *Urllist == "" {
		fmt.Println("请使用-url 或 -urls 来指定目标")
		os.Exit(0)
	}
	wg.Wait()
	fmt.Println("done")
}

func get_req(url string, servers *string, headers *string, bodys *string) {
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

func get_server(value string, servers string) (result bool) {
	result = false
	if servers == value {
		result = true
	}
	return result
}
func get_title(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile("<title>(.*?)</title>")
	matchArr := compileRegex.FindStringSubmatch(bodys)
	if matchArr != nil {
		if matchArr[1] == value {
			result = true
		}
	}
	return result
}

func get_body(value string, bodys string) (result bool) {
	result = strings.Contains(bodys, value)
	return result
}

func get_banner(value string, bodys string) (result bool) {
	result = false
	compileRegex := regexp.MustCompile("(?im)<\\s*banner.*>(.*?)<\\s*/\\s*banner>")
	matchArr := compileRegex.FindAllStringSubmatch(bodys, -1)
	var i int
	for i = 0; i < len(matchArr); i++ {
		if value == matchArr[i][1] {
			result = true
		}
	}
	return result
}

func get_port(value string, urls string) (result bool) {

	u, _ := url.Parse(urls)
	result = false
	host := u.Host
	address := net.ParseIP(host)
	if address == nil {
		result = false
	} else {
		ho := strings.Split(host, ":")
		if len(ho) > 1 {
			port := ho[1]
			if port == value {
				result = true
			}
		}
	}
	return result
}

func get_head(value string, headers string) (result bool) {
	what_headers, _ := json.Marshal(headers)
	child_string := string(what_headers)
	what_headers = []byte(strings.Replace(child_string, "\"", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "\\", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "[", "", -1))
	what_headers = []byte(strings.Replace(string(what_headers), "]", "", -1))
	result = strings.Contains(string(what_headers), value)
	return result
}

func check_model(match string, value string, url string, bodys string, heades string, servers string) (result bool) {
	switch {
	case match == "body_contains":
		result = get_body(value, bodys)
	case match == "protocol_contains":
		result = false
	case match == "title_contains":
		result = get_title(value, bodys)
	case match == "banner_contains":
		result = get_banner(value, bodys)
	case match == "header_contains":
		result = get_head(value, heades)
	case match == "port_contains":
		result = get_port(value, url)
	case match == "server":
		result = get_server(value, servers)
	case match == "title":
		result = get_title(value, bodys)
	case match == "cert_contains":
		result = false
	case match == "server_contains":
		result = get_server(value, servers)
	case match == "protocol":
		result = false
	default:
		result = false
	}
	return result
}

func ContainsInSlice(items []bool, item bool) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func resulte_write(url string, value string) {
	month := time.Now().Month()
	months := strconv.Itoa(int(month))
	day := time.Now().Day()
	days := strconv.Itoa(int(day))
	var filepath string
	filepath = months + "-" + days + ".txt"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(url)
	write.WriteString("cms :" + value + "\n")
	write.Flush()
}

func run(url string) {
	var bodys string
	var headers string
	var servers string

	get_req(url, &servers, &headers, &bodys)
	jsonConfigList := getConfigs()
	unmarshelledConfigs := deserializeJson(jsonConfigList)
	for _, configObj := range unmarshelledConfigs {

		i := 0
		for i = 0; i < len(configObj.Rules); i++ {

			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, bodys, headers, servers))
				}
				if !(ContainsInSlice(arr, false)) {
					if ContainsInSlice(arr, true) {
						resulte_write(url, configObj.Product)
					}

				}

			} else {
				if check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, bodys, headers, servers) {

					resulte_write(url, configObj.Product)

				}

			}

		}

	}
	defer wg.Done()
}

func getConfigs() string {

	b, err := ioutil.ReadFile("fofa.json")

	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	return str
}

func deserializeJson(configJson string) []FOFA {

	jsonAsBytes := []byte(configJson)
	configs := make([]FOFA, 0)
	err := json.Unmarshal(jsonAsBytes, &configs)
	if err != nil {
		panic(err)
	}
	return configs
}

func Banner() {
	banner := `
    _____                     __        __      __      _     _
    /    '                  /    )    /    )    / |     /|   /
---/____----_--_----__-----\---------/---------/__|----/-| -/--
  /        / /  )  /   )     \      /         /   |   /  | /
_/____    / /  /  (___/  (____/    (____/    /    |  /   |/___

		                       by Ch1n3mo ver:0.1
	`
	print(banner)
}

/*func (Info *ArgsInfo) Flag() {
	Banner()
	//可以指定的参数
	flag.StringVar(&Info.Host, "host", "", "Test a host,http://xxxxx")
	flag.StringVar(&Info.Hosts, "hosts", "", "Filename with hosts,One host per line")
	flag.StringVar(&Info.CmsJson, "cmsjson", "cms.json", "Cms fingerprint feature json file, The default is cms.json")
	flag.Parse()
	if Info.Host == "" && Info.Hosts == "" {
		log.Fatalln("err:./no host parameter")
		return
	}

	if Info.Host != "" && Info.Hosts != "" {
		log.Fatalln("err:./only one host parameter")
		return
	}

}*/

type Bar struct {
	percent int64  //百分比
	cur     int64  //当前进度位置
	total   int64  //总进度
	rate    string //进度条
	graph   string //显示符号
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //初始化进度条位置
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
