package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
此代码用户做处理ip:port的资产
有的时候主机开放端口过多、可能为蜜罐，若进行漏洞探测等操作会导致资源浪费
*/
// 用作处理ip:port列表文件中开放大量端口信息-默认超过100条的不保存
var (
	inputFilePath  = flag.String("input", "", "Path to the input file")
	outputFilePath = flag.String("output", "", "Path to the output file including the file name")
	portLimit      = flag.Int("l", 100, "Limit on the number of open ports per IP")
)

func main() {
	flag.Parse()

	if *inputFilePath == "" || *outputFilePath == "" {
		fmt.Println("Please provide both input and output paths")
		return
	}

	file, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ipPortsMap := make(map[string][]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		ip := parts[0]
		port := parts[1]

		ipPortsMap[ip] = append(ipPortsMap[ip], port)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = os.MkdirAll(filepath.Dir(*outputFilePath), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output folder:", err)
		return
	}

	file, err = os.Create(*outputFilePath)
	if err != nil {
		fmt.Println("Error creating filtered file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for ip, ports := range ipPortsMap {
		if len(ports) < *portLimit {
			for _, port := range ports {
				_, err := writer.WriteString(fmt.Sprintf("%s:%s\n", ip, port))
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Filtered file created successfully at:", *outputFilePath)
}
