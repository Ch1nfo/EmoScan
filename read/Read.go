package Read

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func GetConfigs() string {

	b, err := ioutil.ReadFile("fofa.json")

	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	return str
}

func DeserializeJson(configJson []byte) []FOFA {

	jsonAsBytes := configJson
	configs := make([]FOFA, 0)
	err := json.Unmarshal(jsonAsBytes, &configs)
	if err != nil {
		panic(err)
	}
	return configs
}
