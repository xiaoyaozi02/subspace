package main

import (
	"bufio"
	"fmt"
	"io"
	"jk_hash/ddding"
	"jk_hash/fdisk"
	"jk_hash/ip"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
)

//更改为你日志文件的所在目录
var logDir = "/run/user/1001"
var keyword = "Successfully"
var startTime time.Time
var previousCount int
var currentCount int

func main() {
	// 初始化监控
	c := make(chan notify.EventInfo)
	if err := notify.Watch(logDir, c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	fmt.Println("本机IP: ", ip.GetLoacalIPAddresses())

	//检查硬盘挂载情况
	mountCount, totalSize := fdisk.GetSubspaceMountInfo()
	//计算所有硬盘挂载总和
	totalCapacity, err := fdisk.GetSubspaceTotalCapacity()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	totalCapacityInTB := fdisk.ConvertBytesToTB(totalCapacity)

	fmt.Printf("subspace挂载硬盘数量: %d\n", mountCount)
	fmt.Printf("subspace单个硬盘容量: %siB\n", totalSize)
	fmt.Printf("subspace硬盘总容量为: %.2fTiB\n", totalCapacityInTB)

	// 设置初始开始时间为当前时间
	startTime = time.Now()
	// 定时任务，每隔1分钟检查一次
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		// 检查是否已经过去了24小时
		if time.Since(startTime) >= 24*time.Hour {
			ipAddress := ip.GetLoacalIPAddresses()
			message := fmt.Sprintf("本机IP:%v\n 24小时内爆块: %v 次 \n Local IP: %s",ipAddress, currentCount-previousCount, ipAddress)
			ddding.SendToDingTalkGroup(message)

			// 保存当前周期的统计结果作为上一个周期的统计结果
			previousCount = currentCount

			// 重置开始时间为当前时间，开始新的24小时周期
			startTime = time.Now()
			currentCount = 0
		}

		// 统计关键词出现次数
		currentCount = countKeywordOccurrences()
	}
}

var lastPosition map[string]int64

func init() {
	lastPosition = make(map[string]int64)
}

func countKeywordOccurrences() int {
	count := 0
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "sub") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// 恢复上次统计的位置
			lastPos, ok := lastPosition[path]
			if ok {
				file.Seek(lastPos, 0)
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				count += strings.Count(line, keyword)
			}

			// 记录当前位置
			lastPosition[path], _ = file.Seek(0, io.SeekCurrent)
			if err := scanner.Err(); err != nil {
				log.Printf("Error scanning file: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking path: %v", err)
	}
	return count
}

/*
func countKeywordOccurrences() int {
	count := 0
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "sub") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			buf := make([]byte, 1024)
			for {
				n, err := file.Read(buf)
				if n == 0 || err != nil {
					break
				}
				count += strings.Count(string(buf[:n]), keyword)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error counting keyword occurrences: %v", err)
	}
	return count
}
*/

//这里注释的代码会出现当日志文件重置会发送负数消息

/*
func countKeywordOccurrences() int {
	count := 0
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "sub") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				count += strings.Count(line, keyword)
			}

			if err := scanner.Err(); err != nil {
				log.Printf("Error scanning file: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking path: %v", err)
	}
	return count
}
*/