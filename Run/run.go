package Run

import (
	Check "EmoScan/GetCheck"
	Read "EmoScan/Read"
	"EmoScan/assets"
	Write "EmoScan/write"
)

func Run(url string, Bodys string, Headers string, Servers string) {

	unmarshelledConfigs := Read.DeserializeJson(assets.JsonAsBytes)
	for _, configObj := range unmarshelledConfigs {

		i := 0
		for i = 0; i < len(configObj.Rules); i++ {
			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, Check.Check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, Bodys, Headers, Servers))
				}
				if !(ContainsInSlice(arr, false)) {
					if ContainsInSlice(arr, true) {
						Write.Resulte_write(url, configObj.Product)
					}
				}

			} else {
				if Check.Check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, Bodys, Headers, Servers) {
					Write.Resulte_write(url, configObj.Product)
				}
			}
		}
	}
	Check.Wg.Done()
}

func ContainsInSlice(items []bool, item bool) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
