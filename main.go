package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// 从环境变量读取配置
	url := os.Getenv("TARGET_URL")
	if url == "" {
		log.Fatal("TARGET_URL environment variable is required")
	}

	intervalStr := os.Getenv("INTERVAL")
	if intervalStr == "" {
		intervalStr = "30" // 默认 30 秒
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		log.Fatal("INTERVAL must be a positive integer (seconds)")
	}

	showResponse := os.Getenv("SHOW_RESPONSE")
	if showResponse == "" {
		showResponse = "true" // 默认显示响应
	}
	printResponse := showResponse == "true" || showResponse == "1" || showResponse == "yes"

	log.Printf("Starting heartbeat service")
	log.Printf("Target URL: %s", url)
	log.Printf("Interval: %d seconds", interval)
	log.Printf("Show Response: %v", printResponse)
	log.Println("---")

	// 循环执行：请求完成后等待间隔时间再执行下一次
	for {
		makeRequest(url, printResponse)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func makeRequest(url string, printResponse bool) {
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return
	}

	if printResponse && len(body) > 0 {
		log.Printf("Status: %d | Duration: %v | Response: %s", resp.StatusCode, duration, string(body))
	} else {
		log.Printf("Status: %d | Duration: %v", resp.StatusCode, duration)
	}
}
