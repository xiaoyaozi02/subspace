package fdisk

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetSubspaceMountInfo() (int, string) {
	//执行df -h 并找出sub硬盘
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing df command: ", err)
		return 0, ""
	}

	// 解析输出并获取挂载信息
	lines := strings.Split(string(output),"\n")
	var mountCount int
	var totalSize string
	for _, line := range lines {
		//跳过标题行和空行
		if strings.HasPrefix(line, "Filesystem") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		// 检查是否包含 "subspace"
		if strings.Contains(line, "subspace") {
			mountCount++
			fields := strings.Fields(line)
            // 第一个字段是文件系统,第二个字段是已用容量
            totalSize = fields[1]
		}
	}

	return mountCount, totalSize
}

//计算出所有subspace硬盘的总容量
func GetSubspaceTotalCapacity() (int, error) {
	// 执行 df -h 命令获取硬盘信息
	cmd := exec.Command("df", "-h")
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// 将命令输出转换为字符串
	output := string(out)

	// 按行分割输出
	lines := strings.Split(output, "\n")

	// 计算总容量
	totalCapacity := 0
	for _, line := range lines {
		if strings.Contains(line, "subspace") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				capacityStr := strings.TrimRight(fields[1], "T") // 去掉容量中的单位"T"
				capacity, err := strconv.Atoi(capacityStr)
				if err == nil {
					totalCapacity += capacity
				}
			}
		}
	}
	
	return totalCapacity, nil
}