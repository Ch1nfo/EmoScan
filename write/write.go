package write

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func Resulte_write(url string, value string) {
	month := time.Now().Month()
	months := strconv.Itoa(int(month))
	day := time.Now().Day()
	days := strconv.Itoa(int(day))
	var filepath = months + "-" + days + ".txt"
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
