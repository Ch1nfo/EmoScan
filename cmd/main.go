package main

import (
	core "EmoScan/core"
	"flag"
	"fmt"
	"os"
)

func main() {

	Banner()
	flag.Parse()

	if *core.URL != "" {
		fmt.Println("扫描中")
		core.Wg.Add(1)
		data := core.Get_req(*core.URL)
		core.Run(*core.URL, data.Bodys, data.Header, data.Server)
	}

	if *core.Urllist != "" {
		fmt.Println("扫描中")
		lines := core.Get_urllist()
		for _, line := range lines {
			core.Wg.Add(1)
			data := core.Get_req(line)
			go core.Run(line, data.Bodys, data.Header, data.Server)
		}

	}

	if *core.URL == "" && *core.Urllist == "" {
		fmt.Println("请输入-url 或 -file 来指定目标")
		os.Exit(0)
	}
	core.Wg.Wait()

	fmt.Println("done")
}

func Banner() {
	banner := `
    _____                     __        __      __      _     _
    /    '                  /    )    /    )    / |     /|   /
---/____----_--_----__-----\---------/---------/__|----/-| -/--
  /        / /  )  /   )     \      /         /   |   /  | /
_/____    / /  /  (___/  (____/    (____/    /    |  /   |/___

		                       by Ch1nfo ver:1.0
	`
	print(banner)
}
