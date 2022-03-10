package main

import (
	Bar "EmoScan/bar"
	Run "EmoScan/run"
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var URL = flag.String("url", "", "输入url")
var Urllist = flag.String("urls", "", "输入urls的文件名")
var wg sync.WaitGroup

func main() {

	Banner()
	flag.Parse()

	var bar Bar.Bar
	if *URL != "" {

		bar.NewOption(0, 100)
		//bar.NewOptionWithGraph(0, 100, "#")
		for i := 0; i <= 100; i++ {
			time.Sleep(100 * time.Millisecond)
			bar.Play(int64(i))
		}
		bar.Finish()

		Run.Run(*URL)

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

		wg.Add(3)

		for _, line := range lines {

			go Run.Run(line)

		}

	}
	if *URL == "" && *Urllist == "" {
		fmt.Println("请使用-url 或 -urls 来指定目标")
		os.Exit(0)
	}
	wg.Wait()
	fmt.Println("Done")
}

func Banner() {
	banner := `
    _____                     __        __      __      _     _
    /    '                  /    )    /    )    / |     /|   /
---/____----_--_----__-----\---------/---------/__|----/-| -/--
  /        / /  )  /   )     \      /         /   |   /  | /
_/____    / /  /  (___/  (____/    (____/    /    |  /   |/___

		                       by Ch1nfo ver:0.2
	`
	print(banner)
}
