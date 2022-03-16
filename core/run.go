package Core

import (
	"EmoScan/assets"
	Read "EmoScan/struct"
	utils "EmoScan/utils"
	"flag"
	"fmt"
)

var FILE = flag.String("f", "", "save to file")

func Run(url string, Bodys string, Headers string, Servers string) {

	unmarshelledConfigs := Read.DeserializeJson(assets.JsonAsBytes)
	for _, configObj := range unmarshelledConfigs {

		i := 0
		for i = 0; i < len(configObj.Rules); i++ {
			if len(configObj.Rules[i]) >= 2 {
				arr := make([]bool, len(configObj.Rules))
				j := 0
				for j = 0; j < len(configObj.Rules[i]); j++ {
					arr = append(arr, Check_model(configObj.Rules[i][j].Match, configObj.Rules[i][j].Content, url, Bodys, Headers, Servers))
				}
				if !(utils.ContainsInSlice(arr, false)) {
					if utils.ContainsInSlice(arr, true) {
						if *FILE != "" {
							utils.Resulte_write(url, configObj.Product)
						} else {
							fmt.Println(url, configObj.Product)
						}

					}
				}

			} else {
				if Check_model(configObj.Rules[i][0].Match, configObj.Rules[i][0].Content, url, Bodys, Headers, Servers) {
					if *FILE != "" {
						utils.Resulte_write(url, configObj.Product)
					} else {
						fmt.Println(url, configObj.Product)
					}

				}
			}
		}
	}
	Wg.Done()
}
