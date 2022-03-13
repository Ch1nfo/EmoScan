package main

import (
	Check "EmoScan/GetCheck"
	Run "EmoScan/Run"
	"flag"
	"fmt"
	"os"
)

func main() {

	Banner()
	flag.Parse()

	if *Check.URL != "" {
		fmt.Println("扫描中")
		Check.Wg.Add(1)
		data := Check.Get_req(*Check.URL)
		Run.Run(*Check.URL, data.Bodys, data.Header, data.Server)
	}

	if *Check.Urllist != "" {
		fmt.Println("扫描中")
		lines := Check.Get_urllist()
		for _, line := range lines {
			Check.Wg.Add(1)
			data := Check.Get_req(line)
			go Run.Run(line, data.Bodys, data.Header, data.Server)
		}

	}

	if *Check.URL == "" && *Check.Urllist == "" {
		fmt.Println("请使用-url 或 -file 来指定目标")
		os.Exit(0)
	}
	Check.Wg.Wait()

	fmt.Println("done")
}

func Banner() {
	banner := `
    _____                     __        __      __      _     _
    /    '                  /    )    /    )    / |     /|   /
---/____----_--_----__-----\---------/---------/__|----/-| -/--
  /        / /  )  /   )     \      /         /   |   /  | /
_/____    / /  /  (___/  (____/    (____/    /    |  /   |/___

		                       by Ch1nfo ver:0.35
	`
	print(banner)
}
