package Run

import (
	"sync"

	Check "EmoScan/Get_Check"
	Read "EmoScan/Read"
	"EmoScan/assets"
	Write "EmoScan/write"
)

var wg sync.WaitGroup

func ContainsInSlice(items []bool, item bool) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Run(url string) {
	var bodys string
	var headers string
	var servers string

	Check.Get_req(url, &servers, &headers, &bodys)
	JsonConfigList := assets.JsonAsBytes
	unmarshelledConfigs := Read.DeserializeJson(JsonConfigList)
	for _, configObj := range unmarshelledConfigs {

		i := 0
		for i = 0; i < len(configObj.Rules); i++ {

			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, Check.Check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, bodys, headers, servers))
				}
				if !(ContainsInSlice(arr, false)) {
					if ContainsInSlice(arr, true) {
						Write.Resulte_write(url, configObj.Product)
					}

				}

			} else {
				if Check.Check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, bodys, headers, servers) {

					Write.Resulte_write(url, configObj.Product)

				}

			}

		}

	}
	defer wg.Done()
}
