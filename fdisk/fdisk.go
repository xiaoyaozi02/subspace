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
func GetSubspaceTotalCapacity() (int64, error) {
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
	var totalCapacity int64
	for _, line := range lines {
		if strings.Contains(line, "subspace") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				capacityStr := fields[1] // 获取容量字符串,包含单位
				capacity, err := strconv.ParseFloat(capacityStr[:len(capacityStr)-1], 64) // 去掉单位后转换为float64
				if err == nil {
					switch capacityStr[len(capacityStr)-1] { // 检查单位
					case 'T':
						totalCapacity += int64(capacity * 1024 * 1024 * 1024 * 1024) // 转换为bytes
					case 'G':
						totalCapacity += int64(capacity * 1024 * 1024 * 1024) // 转换为bytes
					case 'M':
						totalCapacity += int64(capacity * 1024 * 1024) // 转换为bytes
					case 'K':
						totalCapacity += int64(capacity * 1024) // 转换为bytes
					}
				}
			}
		}
	}
	return totalCapacity, nil
}

func ConvertBytesToTB(bytes int64) float64 {
    return float64(bytes) / (1024 * 1024 * 1024 * 1024)
}