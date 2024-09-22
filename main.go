package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	inputFilePath  = flag.String("input", "", "Path to the input file")
	outputFilePath = flag.String("output", "", "Path to the output file including the file name")
	portLimit      = flag.Int("l", 100, "Limit on the number of open ports per IP")
	mode           = flag.String("mode", "", "Mode of operation: 'filter', 'clean', 'quote', 'cidr','cidr-quote'")
)

func main() {
	flag.Parse()

	if *inputFilePath == "" || *outputFilePath == "" {
		fmt.Println("Please provide both input and output paths")
		return
	}

	switch *mode {
	case "filter":
		err := filterIPPorts(*inputFilePath, *outputFilePath, *portLimit)
		if err != nil {
			fmt.Printf("Error filtering IP ports: %v\n", err)
		}
	//	清除一些内网地址、DSN服务器地址
	case "clean":
		err := cleanIPs(*inputFilePath, *outputFilePath)
		if err != nil {
			fmt.Printf("Error cleaning IPs: %v\n", err)
		}
	//	处理ip资产为fofa语法支持的格式：ip="xx.xx.xx.x"
	case "quote":
		err := addQuotesToIPs(*inputFilePath, *outputFilePath)
		if err != nil {
			fmt.Printf("Error adding quotes to IPs: %v\n", err)
		}
	//	提取资产列表中的C段（大于5条
	case "cidr":
		err := extractAndFilterCIDRs(*inputFilePath, *outputFilePath)
		if err != nil {
			fmt.Printf("Error extracting and filtering CIDRs: %v\n", err)
		}
	default:
		fmt.Println("Invalid mode. Please use 'filter', 'clean', or 'quote'.")
	}
}

func filterIPPorts(inputFile, outputFile string, portLimit int) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

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
		return err
	}

	writer := bufio.NewWriter(output)

	for ip, ports := range ipPortsMap {
		if len(ports) < portLimit {
			for _, port := range ports {
				_, err := writer.WriteString(fmt.Sprintf("%s:%s\n", ip, port))
				if err != nil {
					return err
				}
			}
		}
	}

	return writer.Flush()
}
func cleanIPs(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	excludeIPs := []string{
		"114.114.114.114",
		"8.8.8.8",
		"0.0.0.1",
		"0.0.0.0",
		"127.0.0.1",
		"1.1.1.1",
		"114.114.114.114",
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(output)

	for scanner.Scan() {
		line := scanner.Text()
		ip := net.ParseIP(line)
		if ip == nil || !isPrivateIP(ip) {
			if isExcludedIP(ip.String(), excludeIPs) {
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					return err
				}
			}
		}
	}

	return writer.Flush()
}

func isPrivateIP(ip net.IP) bool {
	return ip.IsPrivate()
}

func isExcludedIP(ip string, excludeIPs []string) bool {
	for _, excludedIP := range excludeIPs {
		if ip == excludedIP {
			return false
		}
	}
	return true
}
func addQuotesToIPs(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(output)

	for scanner.Scan() {
		line := scanner.Text()
		// 添加引号
		quotedLine := fmt.Sprintf("ip=\"%s\"", line)
		_, err := writer.WriteString(quotedLine + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

/*新增C段提取功能，把出现超过4次的IP提取C段出来*/
func extractAndFilterCIDRs(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	cidrMap := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipStr := scanner.Text()
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Printf("Invalid IP address: %s\n", ipStr)
			continue
		}
		cidr := getCIDR(ipStr)
		cidrMap[cidr]++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	writer := bufio.NewWriter(output)
	for cidr, count := range cidrMap {
		if count >= 5 {
			_, err := writer.WriteString(cidr + "\n")
			if err != nil {
				return err
			}
		}
	}

	return writer.Flush()
}

func getCIDR(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	ip = ip.To4()
	if ip == nil {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d.0/24", ip[0], ip[1], ip[2])
}
