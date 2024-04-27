package fdisk

import (
	"fmt"
	"os/exec"
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