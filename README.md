# EmoScan

0.1：这是go语言萌新的试水程序。逻辑混乱，效率低下

0.2更新：使用项目结构重构，使用emded将json嵌入二进制文件内部

0.3更新：大幅修改原来的shi山，处理了发包特征，打包了exe

0.35更新：新增输出选项

# Example

单个url：go run EmoScan.go -url http://127.0.0.1 或 EmoScan.exe -url http://127.0.0.1

多个url：go run EmoScan.go -file filename.txt 或 EmoScan.exe -file filename.txt

.txt文件请放入放入运行目录内
默认输出至控制台，若想保存至文件，请增加输入-f T
使用exe时，文件保存位置为运行目录

